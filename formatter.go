package swag

import (
	"bytes"
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
)

// Check of @Param @Success @Failure @Response @Header
var specialTagForSplit = map[string]bool{
	paramAttr:    true,
	successAttr:  true,
	failureAttr:  true,
	responseAttr: true,
	headerAttr:   true,
}

var skipChar = map[byte]byte{
	'"': '"',
	'(': ')',
	'{': '}',
	'[': ']',
}

// Formatter implements a formatter for Go source files.
type Formatter struct {
	// debugging output goes here
	debug Debugger
}

// NewFormatter create a new formatter instance.
func NewFormatter() *Formatter {
	formatter := &Formatter{
		debug: log.New(os.Stdout, "", log.LstdFlags),
	}
	return formatter
}

// Format formats swag comments in contents. It uses fileName to report errors
// that happen during parsing of contents.
func (f *Formatter) Format(fileName string, contents []byte) ([]byte, error) {
	fileSet := token.NewFileSet()
	ast, err := goparser.ParseFile(fileSet, fileName, contents, goparser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Formatting changes are described as an edit list of byte range
	// replacements. We make these content-level edits directly rather than
	// changing the AST nodes and writing those out (via [go/printer] or
	// [go/format]) so that we only change the formatting of Swag attribute
	// comments. This won't touch the formatting of any other comments, or of
	// functions, etc.
	maxEdits := 0
	for _, comment := range ast.Comments {
		maxEdits += len(comment.List)
	}
	edits := make(edits, 0, maxEdits)

	for _, comment := range ast.Comments {
		formatFuncDoc(fileSet, comment.List, &edits)
	}

	return edits.apply(contents), nil
}

type edit struct {
	begin       int
	end         int
	replacement []byte
}

type edits []edit

func (edits edits) apply(contents []byte) []byte {
	// Apply the edits with the highest offset first, so that earlier edits
	// don't affect the offsets of later edits.
	sort.Slice(edits, func(i, j int) bool {
		return edits[i].begin > edits[j].begin
	})

	// 每次编辑都写入新分配的切片，避免就地改写调用方传入的 `contents`
	// 否则替换不需要扩容时，调用方无法通过比较内容判断文件是否有改动
	for _, edit := range edits {
		result := make([]byte, 0, len(contents)-(edit.end-edit.begin)+len(edit.replacement))
		result = append(result, contents[:edit.begin]...)
		result = append(result, edit.replacement...)
		result = append(result, contents[edit.end:]...)
		contents = result
	}

	return contents
}

// formatFuncDoc reformats the comment lines in commentList, and appends any
// changes to the edit list.
func formatFuncDoc(fileSet *token.FileSet, commentList []*ast.Comment, edits *edits) {
	// Building the edit list to format a comment block is a two-step process.
	// First, we iterate over each comment line looking for Swag attributes. In
	// each one we find, we replace alignment whitespace with a tab character,
	// then write the result into a tab writer.

	linesToComments := make(map[int]int, len(commentList))

	buffer := &bytes.Buffer{}
	w := tabwriter.NewWriter(buffer, 1, 4, 1, '\t', 0)

	for commentIndex, comment := range commentList {
		text := comment.Text
		if attr, body, found := swagComment(text); found {
			formatted := "// " + attr
			if body != "" {
				formatted += "\t" + splitComment2(attr, body)
				formatted = formatGeneric(attr, formatted)
			}
			_, _ = fmt.Fprintln(w, formatted)
			linesToComments[len(linesToComments)] = commentIndex
		}
	}

	// Once we've loaded all of the comment lines to be aligned into the tab
	// writer, flushing it causes the aligned text to be written out to the
	// backing buffer.
	_ = w.Flush()

	// Now the second step: we iterate over the aligned comment lines that were
	// written into the backing buffer, pair each one up to its original
	// comment line, and use the combination to describe the edit that needs to
	// be made to the original input.
	formattedComments := bytes.Split(buffer.Bytes(), []byte("\n"))
	for lineIndex, commentIndex := range linesToComments {
		comment := commentList[commentIndex]
		*edits = append(*edits, edit{
			begin:       fileSet.Position(comment.Pos()).Offset,
			end:         fileSet.Position(comment.End()).Offset,
			replacement: formattedComments[lineIndex],
		})
	}
}

func splitComment2(attr, body string) string {
	if specialTagForSplit[strings.ToLower(attr)] {
		for i := 0; i < len(body); i++ {
			if skipEnd, ok := skipChar[body[i]]; ok {
				skipStart, n := body[i], 1
				for i++; i < len(body); i++ {
					if skipStart != skipEnd && body[i] == skipStart {
						n++
					} else if body[i] == skipEnd {
						n--
						if n == 0 {
							break
						}
					}
				}
			} else if body[i] == ' ' || body[i] == '\t' {
				j := i
				for ; j < len(body) && (body[j] == ' ' || body[j] == '\t'); j++ {
				}
				body = replaceRange(body, i, j, "\t")
			}
		}
	}
	return body
}

func replaceRange(s string, start, end int, new string) string {
	return s[:start] + new + s[end:]
}

var swagCommentLineExpression = regexp.MustCompile(`^\/\/\s+(@[\S.]+)\s*(.*)`)

func swagComment(comment string) (string, string, bool) {
	matches := swagCommentLineExpression.FindStringSubmatch(comment)
	if matches == nil {
		return "", "", false
	}
	return matches[1], matches[2], true
}

type BracketType int

const (
	BracketTypeGeneric BracketType = iota
	BracketTypeMap
	BracketTypeSlice
)

type Bracket struct {
	LeftPos int
	Type    BracketType
}

// formatGeneric adds spaces inside generic type brackets for better readability
// TODO: remove extra spaces
func formatGeneric(attr string, body string) string {
	if !specialTagForSplit[strings.ToLower(attr)] {
		return body
	}

	if strings.IndexByte(body, '[') == -1 && strings.IndexByte(body, ']') == -1 {
		return body
	}

	arr := strings.Split(body, "\t")
	if len(arr) < 4 {
		return body
	}

	str := arr[3]
	n := len(str)

	var brackets []*Bracket

	for i := 0; i < n; i++ {
		ch := str[i]

		switch ch {
		case '[':
			lb := &Bracket{
				LeftPos: i,
			}
			if isMapLeftBracket(str, i) {
				lb.Type = BracketTypeMap
			} else if isSliceLeftBracket(str, i) {
				lb.Type = BracketTypeSlice
			} else {
				lb.Type = BracketTypeGeneric
			}
			brackets = append(brackets, lb)
		case ']':
			// get the last left bracket and remove it from slice
			if len(brackets) == 0 {
				continue
			}

			lastIndex := len(brackets) - 1
			lastLeftBracket := brackets[lastIndex]
			brackets = brackets[:lastIndex]

			if lastLeftBracket.Type == BracketTypeGeneric {
				added := 0

				// if [ right character  is space
				if lastLeftBracket.LeftPos == 0 || str[lastLeftBracket.LeftPos+1] != ' ' {
					// insert space before left bracket
					str = replaceRange(str, lastLeftBracket.LeftPos+1, lastLeftBracket.LeftPos+1, " ")
					added += 1
				}

				// if ] left character is space
				if i > 0 && str[i+added-1] != ' ' {
					// insert space after right bracket
					str = replaceRange(str, i+added, i+added, " ")
					added += 1
				}

				i += added
				n += added
			}
		}
	}

	arr[3] = str
	return strings.Join(arr, "\t")
}

func isNormalChar(ch uint8) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
}

func isMapLeftBracket(str string, i int) bool {
	if i < 3 {
		return false
	}

	prevChar := str[i-1]

	if prevChar == ' ' {
		return isMapLeftBracket(str, i-1)
	}

	if !isNormalChar(prevChar) {
		return false
	}

	if !isNormalChar(prevChar) {
		return false
	}

	return str[i-3:i] == "map"
}

func isSliceLeftBracket(str string, i int) bool {
	if len(str) < i+1 {
		return false
	}

	nextChar := str[i+1]

	if nextChar == ' ' {
		return isSliceLeftBracket(str, i+1)
	}

	return nextChar == ']'
}

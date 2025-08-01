package swag

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/go-openapi/spec"
)

// Schema parsed schema.
type Schema struct {
	*spec.Schema        //
	PkgPath      string // package import path used to rename Name of a definition int case of conflict
	Name         string // Name in definitions
}

// TypeSpecDef the whole information of a typeSpec.
type TypeSpecDef struct {
	// ast file where TypeSpec is
	File *ast.File

	// the TypeSpec of this type definition
	TypeSpec *ast.TypeSpec

	Enums []EnumValue

	// path of package starting from under ${GOPATH}/src or from module path in go.mod
	PkgPath    string
	ParentSpec ast.Decl

	SchemaName string

	NotUnique bool
}

// Name the name of the typeSpec.
func (t *TypeSpecDef) Name() string {
	if t.TypeSpec != nil && t.TypeSpec.Name != nil {
		return t.TypeSpec.Name.Name
	}

	return ""
}

// TypeName the type name of the typeSpec.
func (t *TypeSpecDef) TypeName() string {
	if ignoreNameOverride(t.TypeSpec.Name.Name) {
		return t.TypeSpec.Name.Name[1:]
	}

	var names []string
	if t.NotUnique {
		pkgPath := strings.Map(func(r rune) rune {
			if r == '\\' || r == '/' || r == '.' {
				return '_'
			}
			return r
		}, t.PkgPath)
		names = append(names, pkgPath)
	} else if t.File != nil {
		names = append(names, t.File.Name.Name)
	}
	if parentFun, ok := (t.ParentSpec).(*ast.FuncDecl); ok && parentFun != nil {
		names = append(names, parentFun.Name.Name)
	}
	names = append(names, t.TypeSpec.Name.Name)
	return fullTypeName(names...)
}

// FullPath return the full path of the typeSpec.
func (t *TypeSpecDef) FullPath() string {
	return t.PkgPath + "." + t.Name()
}

func (t *TypeSpecDef) Alias() string {
	return nameOverride(t.TypeSpec.Comment)
}

func (t *TypeSpecDef) SetSchemaName() {
	if alias := t.Alias(); alias != "" {
		t.SchemaName = alias
		return
	}

	t.SchemaName = t.TypeName()
}

// AstFileInfo information of an ast.File.
type AstFileInfo struct {
	//FileSet the FileSet object which is used to parse this go source file
	FileSet *token.FileSet

	// File ast.File
	File *ast.File

	// Path the path of the ast.File
	Path string

	// PackagePath package import path of the ast.File
	PackagePath string

	// ParseFlag determine what to parse
	ParseFlag ParseFlag
}

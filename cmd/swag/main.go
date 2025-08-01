package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/swaggo/swag"
	"github.com/swaggo/swag/format"
	"github.com/swaggo/swag/gen"
)

const (
	searchDirFlag            = "dir"
	excludeFlag              = "exclude"
	generalInfoFlag          = "generalInfo"
	pipeFlag                 = "pipe"
	propertyStrategyFlag     = "propertyStrategy"
	outputFlag               = "output"
	outputTypesFlag          = "outputTypes"
	parseVendorFlag          = "parseVendor"
	parseDependencyFlag      = "parseDependency"
	useStructNameFlag        = "useStructName"
	parseDependencyLevelFlag = "parseDependencyLevel"
	markdownFilesFlag        = "markdownFiles"
	codeExampleFilesFlag     = "codeExampleFiles"
	parseInternalFlag        = "parseInternal"
	generatedTimeFlag        = "generatedTime"
	requiredByDefaultFlag    = "requiredByDefault"
	parseDepthFlag           = "parseDepth"
	instanceNameFlag         = "instanceName"
	overridesFileFlag        = "overridesFile"
	parseGoListFlag          = "parseGoList"
	quietFlag                = "quiet"
	tagsFlag                 = "tags"
	parseExtensionFlag       = "parseExtension"
	templateDelimsFlag       = "templateDelims"
	packageName              = "packageName"
	collectionFormatFlag     = "collectionFormat"
	packagePrefixFlag        = "packagePrefix"
	stateFlag                = "state"
	parseFuncBodyFlag        = "parseFuncBody"
)

var initFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:    quietFlag,
		Aliases: []string{"q"},
		Usage:   "Make the logger quiet.",
	},
	&cli.StringFlag{
		Name:    generalInfoFlag,
		Aliases: []string{"g"},
		Value:   "main.go",
		Usage:   "Go file path in which 'swagger general API Info' is written",
	},
	&cli.StringFlag{
		Name:    searchDirFlag,
		Aliases: []string{"d"},
		Value:   "./",
		Usage:   "Directories you want to parse,comma separated and general-info file must be in the first one",
	},
	&cli.StringFlag{
		Name:  excludeFlag,
		Usage: "Exclude directories and files when searching, comma separated",
	},
	&cli.StringFlag{
		Name:    propertyStrategyFlag,
		Aliases: []string{"p"},
		Value:   swag.CamelCase,
		Usage:   "Property Naming Strategy like " + swag.SnakeCase + "," + swag.CamelCase + "," + swag.PascalCase,
	},
	&cli.StringFlag{
		Name:    outputFlag,
		Aliases: []string{"o"},
		Value:   "./docs",
		Usage:   "Output directory for all the generated files(swagger.json, swagger.yaml and docs.go)",
	},
	&cli.StringFlag{
		Name:    outputTypesFlag,
		Aliases: []string{"ot"},
		Value:   "go,json,yaml",
		Usage:   "Output types of generated files (docs.go, swagger.json, swagger.yaml) like go,json,yaml",
	},
	&cli.BoolFlag{
		Name:  parseVendorFlag,
		Usage: "Parse go files in 'vendor' folder, disabled by default",
	},
	&cli.IntFlag{
		Name:    parseDependencyLevelFlag,
		Aliases: []string{"pdl"},
		Usage:   "Parse go files inside dependency folder, 0 disabled, 1 only parse models, 2 only parse operations, 3 parse all",
	},
	&cli.BoolFlag{
		Name:    parseDependencyFlag,
		Aliases: []string{"pd"},
		Usage:   "Parse go files inside dependency folder, disabled by default",
	},
	&cli.BoolFlag{
		Name:    useStructNameFlag,
		Aliases: []string{"st"},
		Usage:   "Dont use those ugly full-path names when using dependency flag",
	},
	&cli.StringFlag{
		Name:    markdownFilesFlag,
		Aliases: []string{"md"},
		Value:   "",
		Usage:   "Parse folder containing markdown files to use as description, disabled by default",
	},
	&cli.StringFlag{
		Name:    codeExampleFilesFlag,
		Aliases: []string{"cef"},
		Value:   "",
		Usage:   "Parse folder containing code example files to use for the x-codeSamples extension, disabled by default",
	},
	&cli.BoolFlag{
		Name:  parseInternalFlag,
		Usage: "Parse go files in internal packages, disabled by default",
	},
	&cli.BoolFlag{
		Name:  generatedTimeFlag,
		Usage: "Generate timestamp at the top of docs.go, disabled by default",
	},
	&cli.IntFlag{
		Name:  parseDepthFlag,
		Value: 100,
		Usage: "Dependency parse depth",
	},
	&cli.BoolFlag{
		Name:  requiredByDefaultFlag,
		Usage: "Set validation required for all fields by default",
	},
	&cli.StringFlag{
		Name:  instanceNameFlag,
		Value: "",
		Usage: "This parameter can be used to name different swagger document instances. It is optional.",
	},
	&cli.StringFlag{
		Name:  overridesFileFlag,
		Value: gen.DefaultOverridesFile,
		Usage: "File to read global type overrides from.",
	},
	&cli.BoolFlag{
		Name:  parseGoListFlag,
		Value: true,
		Usage: "Parse dependency via 'go list'",
	},
	&cli.StringFlag{
		Name:  parseExtensionFlag,
		Value: "",
		Usage: "Parse only those operations that match given extension",
	},
	&cli.StringFlag{
		Name:    tagsFlag,
		Aliases: []string{"t"},
		Value:   "",
		Usage:   "A comma-separated list of tags to filter the APIs for which the documentation is generated.Special case if the tag is prefixed with the '!' character then the APIs with that tag will be excluded",
	},
	&cli.StringFlag{
		Name:    templateDelimsFlag,
		Aliases: []string{"td"},
		Value:   "",
		Usage:   "Provide custom delimiters for Go template generation. The format is leftDelim,rightDelim. For example: \"[[,]]\"",
	},
	&cli.StringFlag{
		Name:  packageName,
		Value: "",
		Usage: "A package name of docs.go, using output directory name by default (check `--output` option)",
	},
	&cli.StringFlag{
		Name:    collectionFormatFlag,
		Aliases: []string{"cf"},
		Value:   "csv",
		Usage:   "Set default collection format",
	},
	&cli.StringFlag{
		Name:  packagePrefixFlag,
		Value: "",
		Usage: "Parse only packages whose import path match the given prefix, comma separated",
	},
	&cli.StringFlag{
		Name:  stateFlag,
		Value: "",
		Usage: "Set host state for swagger.json",
	},
	&cli.BoolFlag{
		Name:  parseFuncBodyFlag,
		Usage: "Parse API info within body of functions in go files, disabled by default",
	},
}

func initAction(ctx *cli.Context) error {
	strategy := ctx.String(propertyStrategyFlag)

	switch strategy {
	case swag.CamelCase, swag.SnakeCase, swag.PascalCase:
	default:
		return fmt.Errorf("not supported %s propertyStrategy", strategy)
	}

	leftDelim, rightDelim := "{{", "}}"

	if ctx.IsSet(templateDelimsFlag) {
		delims := strings.Split(ctx.String(templateDelimsFlag), ",")
		if len(delims) != 2 {
			return fmt.Errorf(
				"exactly two template delimiters must be provided, comma separated",
			)
		} else if delims[0] == delims[1] {
			return fmt.Errorf("template delimiters must be different")
		}
		leftDelim, rightDelim = strings.TrimSpace(
			delims[0],
		), strings.TrimSpace(
			delims[1],
		)
	}

	outputTypes := strings.Split(ctx.String(outputTypesFlag), ",")
	if len(outputTypes) == 0 {
		return fmt.Errorf("no output types specified")
	}
	logger := log.New(os.Stdout, "", log.LstdFlags)
	if ctx.Bool(quietFlag) {
		logger = log.New(io.Discard, "", log.LstdFlags)
	}

	collectionFormat := swag.TransToValidCollectionFormat(
		ctx.String(collectionFormatFlag),
	)
	if collectionFormat == "" {
		return fmt.Errorf(
			"not supported %s collectionFormat",
			ctx.String(collectionFormat),
		)
	}

	var pdv = ctx.Int(parseDependencyLevelFlag)
	if pdv == 0 {
		if ctx.Bool(parseDependencyFlag) {
			pdv = 1
		}
	}
	return gen.New().Build(&gen.Config{
		SearchDir:           ctx.String(searchDirFlag),
		Excludes:            ctx.String(excludeFlag),
		ParseExtension:      ctx.String(parseExtensionFlag),
		MainAPIFile:         ctx.String(generalInfoFlag),
		PropNamingStrategy:  strategy,
		OutputDir:           ctx.String(outputFlag),
		OutputTypes:         outputTypes,
		ParseVendor:         ctx.Bool(parseVendorFlag),
		ParseDependency:     pdv,
		MarkdownFilesDir:    ctx.String(markdownFilesFlag),
		ParseInternal:       ctx.Bool(parseInternalFlag),
		UseStructNames:      ctx.Bool(useStructNameFlag),
		GeneratedTime:       ctx.Bool(generatedTimeFlag),
		RequiredByDefault:   ctx.Bool(requiredByDefaultFlag),
		CodeExampleFilesDir: ctx.String(codeExampleFilesFlag),
		ParseDepth:          ctx.Int(parseDepthFlag),
		InstanceName:        ctx.String(instanceNameFlag),
		OverridesFile:       ctx.String(overridesFileFlag),
		ParseGoList:         ctx.Bool(parseGoListFlag),
		Tags:                ctx.String(tagsFlag),
		LeftTemplateDelim:   leftDelim,
		RightTemplateDelim:  rightDelim,
		PackageName:         ctx.String(packageName),
		Debugger:            logger,
		CollectionFormat:    collectionFormat,
		PackagePrefix:       ctx.String(packagePrefixFlag),
		State:               ctx.String(stateFlag),
		ParseFuncBody:       ctx.Bool(parseFuncBodyFlag),
	})
}

func main() {
	app := cli.NewApp()
	app.Version = swag.Version
	app.Usage = "Automatically generate RESTful API documentation with Swagger 2.0 for Go."
	app.Commands = []*cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Create docs.go",
			Action:  initAction,
			Flags:   initFlags,
		},
		{
			Name:    "fmt",
			Aliases: []string{"f"},
			Usage:   "format swag comments",
			Action: func(c *cli.Context) error {

				if c.Bool(pipeFlag) {
					return format.New().Run(os.Stdin, os.Stdout)
				}

				searchDir := c.String(searchDirFlag)
				excludeDir := c.String(excludeFlag)
				mainFile := c.String(generalInfoFlag)

				return format.New().Build(&format.Config{
					SearchDir: searchDir,
					Excludes:  excludeDir,
					MainFile:  mainFile,
				})
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    searchDirFlag,
					Aliases: []string{"d"},
					Value:   "./",
					Usage:   "Directories you want to parse,comma separated and general-info file must be in the first one",
				},
				&cli.StringFlag{
					Name:  excludeFlag,
					Usage: "Exclude directories and files when searching, comma separated",
				},
				&cli.StringFlag{
					Name:    generalInfoFlag,
					Aliases: []string{"g"},
					Value:   "main.go",
					Usage:   "Go file path in which 'swagger general API Info' is written",
				},
				&cli.BoolFlag{
					Name:    "pipe",
					Aliases: []string{"p"},
					Value:   false,
					Usage:   "Read from stdin, write to stdout.",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

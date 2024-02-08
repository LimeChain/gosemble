package benchmarking

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
)

// template struct represents a template file that can be modified using go/parser and go/ast.
// the template file must be a compileable go file and defines a function and optionally constants for ref time, reads and writes.
// the template file must also start with a comment, which will be replaced with summary.
// see variables extrinsicTemplate and overheadTemplate defined below.
type template struct {
	filePath, fnName, refTimeVar, readsVar, writesVar string
}

var (
	_, f, _, _        = runtime.Caller(0) // gets the current file path
	extrinsicTemplate = template{filepath.Join(filepath.Dir(f), "weight_file_extrinsic_template.go"), "extrinsicWeightFn", "refTime", "reads", "writes"}
	overheadTemplate  = template{filepath.Join(filepath.Dir(f), "weight_file_overhead_template.go"), "overheadWeightFn", "refTime", "", ""}
)

func generateWeightFile(template template, outputPath, summary string, refTime, reads, writes uint64) error {
	// parse template file
	fset := token.NewFileSet()
	templateNode, err := parser.ParseFile(fset, template.filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing file: %v", err)
	}

	// find variable declarations and modify values
	ast.Inspect(templateNode, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for i, name := range valueSpec.Names {
					switch name.Name {
					case template.refTimeVar:
						valueSpec.Values[i] = &ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", refTime)}
					case template.readsVar:
						valueSpec.Values[i] = &ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", reads)}
					case template.writesVar:
						valueSpec.Values[i] = &ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", writes)}
					}
				}
			}
		}
		return true
	})

	// generate function name from output path
	functionName := strcase.ToLowerCamel(strings.TrimSuffix(filepath.Base(outputPath), ".go"))

	// find function declaration and modify the function name
	ast.Inspect(templateNode, func(n ast.Node) bool {
		if fnDecl, ok := n.(*ast.FuncDecl); ok {
			if fnDecl.Name.Name == template.fnName {
				fnDecl.Name.Name = functionName
			}
		}
		return true
	})

	// append info
	hostName, _ := os.Hostname()
	cpuInfo := runtime.GOARCH
	infoComment := fmt.Sprintf("// DATE: %s, STEPS: %d, REPEAT: %d, DBCACHE: %d, HEAPPAGES: %d, HOSTNAME: %s, CPU: %s, GC: %s, TINYGO VERSION: %s, TARGET: %s", time.Now(), *steps, *repeat, *dbCache, *heapPages, hostName, cpuInfo, *gc, *tinyGoVersion, *target)
	summaryComment := fmt.Sprintf("// %s", summary)
	generatedFileComment := "// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE"
	comment := fmt.Sprintf("%s\n%s\n\n// Summary:\n%s", generatedFileComment, infoComment, summaryComment)
	templateNode.Comments = []*ast.CommentGroup{{List: []*ast.Comment{{Text: comment}}}}

	// modify package name
	paths := strings.Split(filepath.Dir(outputPath), "/")
	packageName := paths[len(paths)-1]
	templateNode.Name.Name = packageName

	// create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	// write the modified AST to the output file
	if err := format.Node(outputFile, fset, templateNode); err != nil {
		return fmt.Errorf("error writing modified AST to file: %v", err)
	}

	return nil
}

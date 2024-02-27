package benchmarking

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"path/filepath"
	"runtime"
	"strings"
)

// weightFileTemplate struct represents a weightFileTemplate file that can be modified using go/parser and go/ast.
// the weightFileTemplate file must be a compileable go file and defines a function and optionally constants for ref time, reads and writes.
// the weightFileTemplate file must also start with a comment, which will be replaced with summary.
// see variables extrinsicTemplate and overheadTemplate defined below.
type weightFileTemplate struct {
	fSet                   *token.FileSet
	fileNode               *ast.File
	infoComment            *ast.Comment
	refTime, reads, writes *ast.BasicLit
	weightFn               *ast.FuncDecl
}

func InitOverheadWeightTemplate() (*weightFileTemplate, error) {
	return newWeightFileTemplate("weight_file_overhead_template.go", "overheadWeightFn", "refTime", "", "")
}

func newWeightFileTemplate(templateFile, weightFn, refTime, reads, writes string) (*weightFileTemplate, error) {
	template := &weightFileTemplate{}
	template.fSet = token.NewFileSet()

	// get current working directory
	_, cwd, _, ok := runtime.Caller(0)
	if !ok {
		return template, fmt.Errorf("error getting cwd from runtime")
	}

	// parse template file
	templatePath := filepath.Join(filepath.Dir(cwd), templateFile)
	fileNode, err := parser.ParseFile(template.fSet, templatePath, nil, parser.ParseComments)
	if err != nil {
		return template, fmt.Errorf("error parsing file: %v", err)
	}
	template.fileNode = fileNode

	// find info comment
	if len(template.fileNode.Comments) == 0 {
		return template, fmt.Errorf("error getting info comment from file template")
	}
	template.infoComment = &ast.Comment{}
	template.fileNode.Comments = []*ast.CommentGroup{{List: []*ast.Comment{template.infoComment}}}
	// find weightFn declaration
	ast.Inspect(template.fileNode, func(n ast.Node) bool {
		if fnDecl, ok := n.(*ast.FuncDecl); ok {
			if fnDecl.Name.Name == weightFn {
				template.weightFn = fnDecl
			}
		}
		return true
	})
	if template.weightFn == nil {
		return template, fmt.Errorf("error getting weightFn from file template")
	}

	// find variable declarations and modify values
	ast.Inspect(template.fileNode, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for i, name := range valueSpec.Names {
						switch name.Name {
						case refTime:
							template.refTime = &ast.BasicLit{Kind: token.INT, Value: ""}
							valueSpec.Values[i] = template.refTime
						case reads:
							template.reads = &ast.BasicLit{Kind: token.INT, Value: ""}
							valueSpec.Values[i] = template.reads
						case writes:
							template.writes = &ast.BasicLit{Kind: token.INT, Value: ""}
							valueSpec.Values[i] = template.writes
						}
					}
				}
			}
		}
		return true
	})
	if template.refTime == nil && refTime != "" {
		return template, fmt.Errorf("error getting 'refTime' variable from file template")
	}
	if template.reads == nil && reads != "" {
		return template, fmt.Errorf("error getting 'reads' variable from file template")
	}
	if template.writes == nil && writes != "" {
		return template, fmt.Errorf("error getting 'writes' variable from file template")
	}
	return template, nil
}

func (w *weightFileTemplate) SetInfoComment(info string) error {
	if len(strings.TrimSpace(info)) == 0 {
		return fmt.Errorf("info comment name cannot be empty")
	}

	infoCommentLines := strings.Split(info, "\n")
	for i, line := range infoCommentLines {
		if !strings.HasPrefix(line, "//") {
			infoCommentLines[i] = fmt.Sprintf("// %s", line)
		}
	}
	w.infoComment.Text = strings.Join(infoCommentLines, "\n")
	return nil
}

func (w *weightFileTemplate) SetWeightFnName(weightFn string) error {
	weightFn = strings.TrimSpace(weightFn)
	if len(weightFn) == 0 {
		return fmt.Errorf("weightFn name cannot be empty")
	}
	w.weightFn.Name.Name = weightFn
	return nil
}

func (w *weightFileTemplate) SetWeightValues(refTime, reads, writes uint64) error {
	if w.refTime != nil {
		w.refTime.Value = fmt.Sprintf("%d", refTime)
	}
	if w.reads != nil {
		w.reads.Value = fmt.Sprintf("%d", reads)
	}
	if w.writes != nil {
		w.writes.Value = fmt.Sprintf("%d", writes)
	}
	return nil
}

func (w *weightFileTemplate) SetPackageName(packageName string) error {
	packageName = strings.TrimSpace(packageName)
	if len(packageName) == 0 {
		return fmt.Errorf("package name cannot be empty")
	}
	w.fileNode.Name.Name = packageName
	return nil
}

func (w *weightFileTemplate) WriteGeneratedFile(output io.Writer) error {
	if err := format.Node(output, w.fSet, w.fileNode); err != nil {
		return fmt.Errorf("error writing modified AST to output: %v", err)
	}

	return nil
}

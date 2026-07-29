package only

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
)

// OnlyCall represents a tape.Only() call found in source.
type OnlyCall struct {
	Parent string // enclosing TestXxx function name
	Name   string // the name string argument
}

// FindOnlyCalls scans all .go files in dir for tape.Only() calls.
func FindOnlyCalls(dir string) ([]OnlyCall, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var all []OnlyCall
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		src, err := os.ReadFile(dir + "/" + e.Name())
		if err != nil {
			return nil, err
		}
		calls, err := FindOnlyCallsInSource(string(src))
		if err != nil {
			return nil, err
		}
		all = append(all, calls...)
	}
	return all, nil
}

// FindOnlyCallsInSource parses Go source and returns all tape.Only() calls.
func FindOnlyCallsInSource(src string) ([]OnlyCall, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		return nil, err
	}

	var calls []OnlyCall
	var currentFunc string

	ast.Inspect(f, func(n ast.Node) bool {
		// track enclosing TestXxx function
		if fn, ok := n.(*ast.FuncDecl); ok {
			if fn.Name != nil && strings.HasPrefix(fn.Name.Name, "Test") {
				currentFunc = fn.Name.Name
			}
			return true
		}

		// find tape.Only() or Only() calls
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if !isOnlyCall(call) {
			return true
		}

		// second argument is the name string literal
		if len(call.Args) >= 2 {
			if lit, ok := call.Args[1].(*ast.BasicLit); ok {
				name := strings.Trim(lit.Value, `"`)
				calls = append(calls, OnlyCall{
					Parent: currentFunc,
					Name:   name,
				})
			}
		}
		return true
	})

	return calls, nil
}

// isOnlyCall returns true if call is tape.Only() or Only().
func isOnlyCall(call *ast.CallExpr) bool {
	switch fn := call.Fun.(type) {
	case *ast.SelectorExpr:
		return fn.Sel.Name == "Only"
	case *ast.Ident:
		return fn.Name == "Only"
	}
	return false
}

// BuildRunPattern converts OnlyCalls to a go test -run pattern.
// Spaces become underscores (Go sub-test escaping).
// Multiple calls joined with | (OR).
func BuildRunPattern(calls []OnlyCall) string {
	if len(calls) == 0 {
		return ""
	}
	// Sort by Parent+Name for deterministic output
	sorted := make([]OnlyCall, len(calls))
	copy(sorted, calls)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Parent != sorted[j].Parent {
			return sorted[i].Parent < sorted[j].Parent
		}
		return sorted[i].Name < sorted[j].Name
	})
	patterns := make([]string, len(sorted))
	for i, c := range sorted {
		name := strings.ReplaceAll(c.Name, " ", "_")
		patterns[i] = c.Parent + "/" + name
	}
	return strings.Join(patterns, "|")
}

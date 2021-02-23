package deadcode

import "github.com/ctriv/php/ast"

// DeadFunctions returns a list of dead functions
func DeadFunctions(fs *ast.FileSet, entryPoints []string) []ast.Node {
	knownFunctions := AllTheFunctions(fs)

	for _, filename := range entryPoints {
		f, ok := fs.Files[filename]
		if !ok {
			continue
		}

		EliminateCalls(f.Nodes, knownFunctions)
	}

	nodes := []ast.Node{}
	for _, f := range knownFunctions {
		nodes = append(nodes, f)
	}

	return nodes
}

// EliminateCalls eliminates all dead calls
func EliminateCalls(nodes []ast.Node, knownFunctions map[string]ast.Node) {
	for _, node := range nodes {
		switch node := node.(type) {
		case ast.FunctionCallExpr:
			if static := ast.Static(node.FunctionName); static != nil {
				delete(knownFunctions, static.Value)
			}
		case *ast.FunctionCallExpr:
			if static := ast.Static(node.FunctionName); static != nil {
				delete(knownFunctions, static.Value)
			}
		}
		EliminateCalls(node.Children(), knownFunctions)
	}

}

// AllTheFunctions returns a list of all functions
func AllTheFunctions(fs *ast.FileSet) map[string]ast.Node {
	namedFunctions := map[string]ast.Node{}
	for f, n := range fs.GlobalNamespace.Functions {
		namedFunctions[f] = n
	}

	for _, class := range fs.GlobalNamespace.ClassesAndInterfaces {
		if class, ok := class.(*ast.Class); ok {
			for _, f := range class.Methods {
				namedFunctions[f.Name] = f.FunctionStmt
			}
		}
	}

	for _, ns := range fs.Namespaces {
		for f, n := range ns.Functions {
			namedFunctions[f] = n
		}
		for _, class := range ns.ClassesAndInterfaces {
			if class, ok := class.(*ast.Class); ok {
				for _, f := range class.Methods {
					namedFunctions[f.Name] = f.FunctionStmt
				}
			}
		}
	}
	return namedFunctions
}

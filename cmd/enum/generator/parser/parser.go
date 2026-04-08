package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"slices"
	"strings"
)

// Results represents the result of parsing a Go source file for enum types. It contains the package name,
// a list of parsed enums, and a set of found declarations to avoid generating code for existing methods or functions.
type Results struct {
	Package           string
	Enums             []*Enum
	FoundDeclarations map[string]struct{}
}

// HasDeclaration checks if a declaration with the given name exists in the FoundDeclarations set. This is used
// to determine if certain methods (like MarshalText, UnmarshalText, or Parse[Enum]) already exist for an enum type,
// so that the generator can skip generating those methods if they are already defined by the user.
func (p *Results) HasDeclaration(name string) bool {
	_, ok := p.FoundDeclarations[name]
	return ok
}

// Parse reads the specified Go source file, parses it for enum types defined in the file, and returns a
// Results struct containing the package name, a list of parsed enums, and a set of found declarations. It
// uses the Go AST to inspect the file and extract relevant information about enum types, their values, and
// any existing methods or functions.
func Parse(filename string, typeNames []string) (*Results, error) {
	// Parse the Go source file and create an AST (Abstract Syntax Tree) representation of it.
	f, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	res := Results{
		Package:           f.Name.Name,
		Enums:             []*Enum{},
		FoundDeclarations: map[string]struct{}{},
	}

	var (
		enums     = map[string]*Enum{}
		iotaState = IotaState{}
	)

	// Inspect the AST and find all structs.
	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.File, *ast.ValueSpec, *ast.TypeSpec, *ast.GenDecl:
			// We are interested in the children of these nodes,
			// so we return true to continue walking the tree
			return true
		case *ast.FuncDecl:
			// Check if it's a method declaration and if it has a receiver
			if node.Recv != nil && len(node.Recv.List) > 0 {
				if ident, ok := node.Recv.List[0].Type.(*ast.Ident); ok {
					// Receiver types are added as Type.MethodName to the FoundDeclarations set
					res.FoundDeclarations[ident.Name+"."+node.Name.Name] = struct{}{}
					return false
				}
			}

			// For regular function declarations, we add the function name to the FoundDeclarations set
			res.FoundDeclarations[node.Name.Name] = struct{}{}

		case *ast.Ident:
			// Make sure it's a Type Identifier
			if ident := node; node.Obj != nil {
				switch spec := ident.Obj.Decl.(type) {
				case *ast.TypeSpec:
					// Check if the type is one of the specified enum types
					if !slices.Contains(typeNames, ident.Name) {
						return false
					}

					if enums[ident.Name] == nil {
						enums[ident.Name] = newEnum(ident.Name)
					}

					// Check if the type is based on int or string and set the EnumType accordingly
					if typeIdent, ok := spec.Type.(*ast.Ident); ok &&
						(typeIdent.Name == "int" || typeIdent.Name == "string") {
						enums[ident.Name].Type = EnumType(typeIdent.Name)
					}
				case *ast.ValueSpec:
					// Add all declared names to the FoundDeclarations set
					for _, name := range spec.Names {
						res.FoundDeclarations[name.Name] = struct{}{}
					}

					// We only care about value specifications that are constants (i.e. part of an enum declaration)
					if spec.Names[0].Obj.Kind != ast.Con {
						return false
					}

					var typeName string
					if typeIdent, ok := spec.Type.(*ast.Ident); ok {
						// If the value specification has an explicit type, use it
						typeName = typeIdent.Name
					} else {
						// Otherwise, use the type from the iota state
						typeName = iotaState.Type
					}

					// If the type name is not empty and is not in the list of specified enum types, skip it
					if typeName != "" && !slices.Contains(typeNames, typeName) {
						return false
					}

					// If the enum type is not already in the enums map
					// e.g. if the values are declared before the type),
					// create a new Enum and add it to the map
					if enums[typeName] == nil {
						enums[typeName] = newEnum(typeName)
					}

					e := enums[typeName]

					// Process each name and value in the value specification. If a value is not provided,
					// use the next iota value.
					for i := range spec.Names {
						name := spec.Names[i].Name
						if i >= len(spec.Values) {
							e.AddValue(name, iotaState.NextValue())
							continue
						}

						switch value := spec.Values[i].(type) {
						case *ast.Ident:
							if value.Name == "iota" { // example: const ColorRed = iota
								iotaState.Reset(typeName)
								e.AddValue(name, iotaState.NextValue())
							} else { // example: const ColorRed = RedAlias
								for _, e := range enums {
									if e.AddAlias(name, value.Name) {
										break
									}
								}
							}
						case *ast.BasicLit: // example: const ColorRed = "red"
							if value.Kind == token.INT || value.Kind == token.STRING {
								e.AddValue(name, value.Value)
							}
						}
					}
				}
			}
		}

		// No need to look deeper into the identifier, so
		// stop processing this branch
		return false
	})

	// Remove any enums that don't have valid type definition or no values
	for name, e := range enums {
		if e.Type == EnumTypeInvalid || len(e.Values) == 0 {
			delete(enums, name)
		}
	}

	// Sort the enums by name and add them to the results
	res.Enums = slices.SortedFunc(maps.Values(enums), func(e1 *Enum, e2 *Enum) int {
		return strings.Compare(e1.Name, e2.Name)
	})

	return &res, nil
}

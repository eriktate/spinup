package gen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

// ServiceDef represents a gRPC service defintion.
type ServiceDef struct {
	PackageName string
	ImportPath  string
	Interface   string
	Implementor string
	Endpoints   []EndpointDef
}

// EndpointDef represents a gRPC endpoint defintion.
type EndpointDef struct {
	Name string
	In   string
	Out  string
}

func GenerateServiceDef(path string) ServiceDef {
	var fset token.FileSet
	f, err := parser.ParseFile(&fset, path, nil, parser.AllErrors)
	if err != nil {
		log.Println("Failed to parse file")
	}

	var serviceDef ServiceDef
	for _, decl := range f.Decls {
		if d, ok := decl.(*ast.GenDecl); ok {
			if d.Tok != token.TYPE {
				continue
			}

			for _, spec := range d.Specs {
				if spec, ok := spec.(*ast.TypeSpec); ok {
					if !strings.Contains(spec.Name.Name, "Server") {
						continue
					}

					serviceDef.Interface = spec.Name.Name
					if iface, ok := spec.Type.(*ast.InterfaceType); ok {
						for _, method := range iface.Methods.List {
							if len(method.Names) == 0 {
								continue
							}

							serviceDef.Endpoints = append(serviceDef.Endpoints, genEndpoint(method))
						}
					}
				}
			}
		}
	}

	return serviceDef
}

func genEndpoint(method *ast.Field) EndpointDef {
	var endpoint EndpointDef
	if len(method.Names) == 0 {
		return endpoint
	}

	endpoint.Name = method.Names[0].Name

	if fun, ok := method.Type.(*ast.FuncType); ok {
		// Need to get the Input struct.
		for _, param := range fun.Params.List {
			if t, ok := param.Type.(*ast.StarExpr); ok {
				// TODO: Not sure this is the cleanest conversion.
				endpoint.In = fmt.Sprintf("%s", t.X)
			}
		}

		// Need to get the Output struct.
		for _, result := range fun.Results.List {
			if t, ok := result.Type.(*ast.StarExpr); ok {
				endpoint.Out = fmt.Sprintf("%s", t.X)
			}
		}
	}

	return endpoint
}

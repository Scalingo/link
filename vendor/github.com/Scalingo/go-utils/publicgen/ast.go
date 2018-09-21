package publicgen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
	"io/ioutil"
)

type field struct {
	Name    string
	Type    string
	JSONTag string
	PkgPath string
	Pointer bool
	Slice   bool
}

func newField(field field) *ast.Field {
	ident := &ast.Ident{
		Name: field.Type,
	}

	fieldAst := &ast.Field{
		Names: []*ast.Ident{
			{
				Name: field.Name,
				Obj: &ast.Object{
					Kind: ast.Var,
					Name: field.Name,
					Decl: ident,
				},
			},
		},
	}

	fieldAst.Type = ident

	if field.Pointer {
		fieldAst.Type = &ast.StarExpr{
			Star: 1,
			X:    ident,
		}
	}

	if field.Slice {
		oldType := fieldAst.Type
		fieldAst.Type = &ast.ArrayType{
			Elt:    oldType,
			Lbrack: 1,
		}
	}

	if field.JSONTag != "" {
		fieldAst.Tag = &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("`json:\"%s\"`", field.JSONTag),
		}
	}

	return fieldAst
}

func newStructType(fields []field) (*ast.StructType, []string) {
	fieldList := make([]*ast.Field, 0)
	imports := make([]string, 0)

	for _, field := range fields {
		fieldList = append(fieldList, newField(field))
		if field.PkgPath != "" {
			imports = append(imports, field.PkgPath)
		}
	}

	return &ast.StructType{
		Fields: &ast.FieldList{
			List: fieldList,
		},
	}, imports
}

func newImports(imports []string) *ast.GenDecl {
	importDecls := make([]ast.Spec, 0)

	for _, imp := range imports {
		importDecls = append(importDecls,
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%s\"", imp),
				},
			})
	}

	return &ast.GenDecl{
		Rparen: 1,
		Tok:    token.IMPORT,
		Specs:  importDecls,
		Lparen: 1,
	}
}

func newAST(packageName string, structs map[string][]field) *ast.File {
	declsList := make([]ast.Decl, 0)
	imports := make([]string, 0)

	for name, fields := range structs {
		fields, imps := newStructType(fields)
		decl := &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: &ast.Ident{
						Name: name,
						Obj: &ast.Object{
							Kind: ast.Typ,
							Name: name,
							Decl: fields,
						},
					},
					Type: fields,
				},
			},
		}

		declsList = append(declsList, decl)
		imports = append(imports, imps...)
	}

	file := &ast.File{
		Package: 1,
		Name: &ast.Ident{
			Name: packageName,
		},
		Decls: declsList,
	}

	if len(imports) > 0 {
		file.Decls = append([]ast.Decl{newImports(imports)}, file.Decls...)
	}
	return file
}

func writeAst(ast *ast.File, filename string) error {
	printConfig := &printer.Config{Mode: printer.TabIndent, Tabwidth: 4}

	var buf bytes.Buffer
	err := printConfig.Fprint(&buf, token.NewFileSet(), ast)
	if err != nil {
		return err
	}
	out := buf.Bytes()

	out, err = format.Source(out)
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	err = ioutil.WriteFile(filename, out, 0644)
	return err
}

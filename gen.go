package main

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"log"
)

var (
	pkgInfo *build.Package
	pkgName string
)

type structType struct {
	name string
	node *ast.StructType
}

func init() {
	var err error

	pkgInfo, err = build.ImportDir(".", 0)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fset := token.NewFileSet()

	// main logic here
	for _, file := range pkgInfo.GoFiles {
		var contents interface{}

		node, err := parser.ParseFile(fset, file, contents, parser.ParseComments)

		pkgName = node.Name.Name

		if err != nil {
			log.Fatal(err)
		}

		structs := collectStructs(node)
		components, err := collectStructComponents(structs)

		if err != nil {
			log.Fatal(err)
		}

		genCode(file, components)
	}

}

func collectStructs(node ast.Node) []*structType {
	structs := make([]*structType, 0)

	collectStructs := func(n ast.Node) bool {
		var t ast.Expr
		var structName string

		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Type == nil {
				return true
			}

			structName = x.Name.Name
			t = x.Type
		case *ast.CompositeLit:
			t = x.Type
		case *ast.ValueSpec:
			structName = x.Names[0].Name
			t = x.Type
		case *ast.Field:
			if len(x.Names) != 0 {
				structName = x.Names[0].Name
			}
			t = x.Type
		}

		t = deref(t)

		x, ok := t.(*ast.StructType)
		if !ok {
			return true
		}

		structs = append(structs, &structType{
			name: structName,
			node: x,
		})
		return true
	}

	ast.Inspect(node, collectStructs)
	return structs
}

func deref(x ast.Expr) ast.Expr {
	switch t := x.(type) {
	case *ast.StarExpr:
		return deref(t.X)
	case *ast.ArrayType:
		return deref(t.Elt)
	}
	return x
}

type structComponets struct {
	Name       string
	FieldNames []string
	FieldTypes []string
	FieldTags  []string
}

func collectStructComponents(structs []*structType) (map[string]structComponets, error) {
	result := make(map[string]structComponets, len(structs))

	for _, s := range structs {
		sc := structComponets{
			Name: s.name,
		}

		fieldNames := make([]string, 0)
		fieldTypes := make([]string, 0)
		fieldTags := make([]string, 0)

		for _, f := range s.node.Fields.List {
			fieldName := ""

			if len(f.Names) != 0 {
				for _, field := range f.Names {
					fieldName = field.Name

				}
			}

			if f.Names == nil {
				ident, ok := f.Type.(*ast.Ident)
				if !ok {
					continue
				}
				fieldName = ident.Name
			}

			if fieldName == "" {
				continue
			}

			if f.Tag == nil {
				f.Tag = &ast.BasicLit{}
			}

			fieldNames = append(fieldNames, fieldName)

			// add Types
			typeExpr := f.Type
			typeString := types.ExprString(typeExpr)
			fieldTypes = append(fieldTypes, typeString)

			// add Tags
			tagExpr := f.Tag
			tagString := types.ExprString(tagExpr)
			fieldTags = append(fieldTags, tagString)

		}
		sc.FieldNames = fieldNames
		sc.FieldTypes = fieldTypes
		sc.FieldTags = fieldTags
		result[sc.Name] = sc
	}

	return result, nil
}

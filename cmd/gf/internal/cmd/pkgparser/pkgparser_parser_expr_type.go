package pkgparser

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

func (p *Parser) parseSyntaxTypeExpr(pkg *packages.Package, syntax *ast.File, expr ast.Expr) *ParsedType {
	if expr == nil {
		return nil
	}
	var parsedType *ParsedType
	switch t := expr.(type) {
	case *ast.Ident:
		parsedType = p.parseSyntaxTypeIdent(pkg, syntax, t)

	case *ast.StarExpr:
		parsedType = p.parseSyntaxTypeExpr(pkg, syntax, t.X)

	case *ast.SelectorExpr:
		parsedType = p.parseSyntaxTypeSelectorExpr(pkg, syntax, t)

	case *ast.StructType:
		parsedType = p.parseSyntaxTypeStructType(pkg, syntax, t)

	case *ast.ArrayType:
		parsedType = p.parseSyntaxTypeArrayType(pkg, syntax, t)

	case *ast.MapType:
		parsedType = p.parseSyntaxTypeMapType(pkg, syntax, t)

	case *ast.InterfaceType:
		parsedType = p.parseSyntaxTypeInterfaceType(pkg, syntax, t)

	case *ast.FuncType:
		parsedType = p.parseSyntaxTypeFuncType(pkg, syntax, t)

	default:
		panic("parseSyntaxTypeExpr failed")
	}
	return parsedType
}

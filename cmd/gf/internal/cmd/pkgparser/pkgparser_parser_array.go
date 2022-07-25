package pkgparser

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

func (p *Parser) parseSyntaxTypeArrayType(pkg *packages.Package, syntax *ast.File, expr *ast.ArrayType) *ParsedType {
	var (
		parsedType = &ParsedType{
			Name:     "",
			Kind:     KindArray,
			PkgPath:  pkg.ID,
			ElemType: &ParsedType{},
		}
		elemType = p.parseSyntaxTypeExpr(pkg, syntax, expr.Elt)
	)
	elemTypeId := p.makeTypeId(elemType.PkgPath, elemType.Name)
	if p.types[elemTypeId] != nil {
		parsedType.ElemType = &ParsedType{
			Name:  elemType.Name,
			Kind:  elemType.Kind,
			Refer: elemTypeId,
		}
	} else {
		*parsedType.ElemType = *elemType
	}
	return parsedType
}

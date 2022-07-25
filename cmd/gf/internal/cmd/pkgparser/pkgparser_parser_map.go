package pkgparser

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

func (p *Parser) parseSyntaxTypeMapType(pkg *packages.Package, syntax *ast.File, expr *ast.MapType) *ParsedType {
	var (
		parsedType = &ParsedType{
			Name:     "",
			Kind:     KindMap,
			PkgPath:  pkg.ID,
			Elems:    nil,
			KeyType:  &ParsedType{},
			ElemType: &ParsedType{},
		}
		keyType  = p.parseSyntaxTypeExpr(pkg, syntax, expr.Key)
		elemType = p.parseSyntaxTypeExpr(pkg, syntax, expr.Value)
	)
	// KeyType.
	keyTypeId := p.makeTypeId(keyType.PkgPath, keyType.Name)
	if p.types[keyTypeId] != nil {
		parsedType.KeyType = &ParsedType{
			Name:  keyType.Name,
			Kind:  keyType.Kind,
			Refer: keyTypeId,
		}
	} else {
		*parsedType.KeyType = *keyType
	}
	// ElemType.
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

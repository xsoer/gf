package pkgparser

import (
	"go/ast"

	"github.com/gogf/gf/v2/errors/gerror"
	"golang.org/x/tools/go/packages"
)

func (p *Parser) parseSyntaxTypeIdent(pkg *packages.Package, syntax *ast.File, ident *ast.Ident) (parsedType *ParsedType) {
	typeId := p.makeTypeId(pkg.ID, ident.Name)
	if _, ok := basicKindMap[ident.Name]; ok {
		typeId = ident.Name
	}
	if p.isStandardPackage(pkg) {
		return &ParsedType{
			Name:  ident.Name,
			Kind:  KindStd,
			Refer: typeId,
		}
	}

	if cachedType := p.types[typeId]; cachedType != nil && cachedType.Name != "" {
		tempType := &ParsedType{}
		*tempType = *cachedType
		return tempType
	}
	p.types[typeId] = &ParsedType{
		PkgPath: pkg.ID,
	}
	defer func() {
		p.saveParsedType(typeId, parsedType)
	}()

	if ident.Obj != nil {
		parsedType = p.parseSyntaxTypeTypeSpec(pkg, syntax, ident.Obj.Decl.(*ast.TypeSpec))
	} else {
		if _, ok := basicKindMap[ident.Name]; ok {
			parsedType = &ParsedType{
				Name:  ident.Name,
				Kind:  KindStd,
				Refer: typeId,
			}
		} else {
			parsedType = p.parsePackageType(pkg, ident.Name)
		}
	}
	if parsedType == nil {
		panic(gerror.Newf(`invalid kind "%s"`, ident.Name))
	}
	return parsedType
}

package pkgparser

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

func (p *Parser) parseSyntaxTypeInterfaceType(pkg *packages.Package, syntax *ast.File, expr *ast.InterfaceType) *ParsedType {
	parsedType := &ParsedType{
		Name:    "",
		Kind:    KindAny,
		PkgPath: pkg.ID,
	}
	return parsedType
}

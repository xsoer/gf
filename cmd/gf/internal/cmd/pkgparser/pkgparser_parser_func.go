package pkgparser

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

func (p *Parser) parseSyntaxTypeFuncType(pkg *packages.Package, syntax *ast.File, expr *ast.FuncType) *ParsedType {
	parsedType := &ParsedType{
		Name:    "",
		Kind:    KindFunc,
		PkgPath: pkg.ID,
	}
	return parsedType
}

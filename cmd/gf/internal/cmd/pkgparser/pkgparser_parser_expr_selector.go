package pkgparser

import (
	"go/ast"

	"github.com/gogf/gf/v2/text/gstr"
	"golang.org/x/tools/go/packages"
)

func (p *Parser) parseSyntaxTypeSelectorExpr(pkg *packages.Package, syntax *ast.File, expr *ast.SelectorExpr) *ParsedType {
	var parsedType *ParsedType
	for _, pkgImportSpec := range syntax.Imports {
		var (
			exprXName = expr.X.(*ast.Ident).Name
			pkgName   = gstr.Trim(pkgImportSpec.Path.Value, `"`)
			inputPkg  = p.getImportByPath(pkg, pkgImportSpec.Path.Value)
		)
		if gstr.Contains(pkgName, "/") {
			pkgName = gstr.SubStrFromREx(pkgName, "/")
		}
		// Search by alias name.
		if pkgImportSpec.Name != nil && pkgImportSpec.Name.Name == exprXName {
			parsedType = p.parsePackageType(inputPkg, expr.Sel.Name)
			break
		}
		// Search by package name.
		if pkgName == exprXName {
			parsedType = p.parsePackageType(inputPkg, expr.Sel.Name)
			break
		}
	}
	if parsedType == nil || parsedType.Name == "" {
		panic("parseSyntaxTypeSelectorExpr error")
	}
	return parsedType
}

func (p *Parser) getImportByPath(pkg *packages.Package, path string) *packages.Package {
	for _, pkgImport := range pkg.Imports {
		if pkgImport.ID == gstr.Trim(path, `"`) {
			return pkgImport
		}
	}
	return nil
}

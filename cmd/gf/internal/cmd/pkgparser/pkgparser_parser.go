package pkgparser

import (
	"fmt"
	"go/ast"
	"reflect"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"golang.org/x/tools/go/packages"
)

func (p *Parser) parsePackage(pkg *packages.Package) *ParsedResult {
	for _, syntax := range pkg.Syntax {
		for _, obj := range syntax.Scope.Objects {
			switch t := obj.Decl.(type) {
			case *ast.TypeSpec:
				p.parseSyntaxTypeTypeSpec(pkg, syntax, t)
			case *ast.ValueSpec:
				p.parseSyntaxTypeValueSpec(pkg, syntax, t)
			case *ast.FuncDecl:
				p.parseSyntaxTypeFuncDecl(pkg, syntax, t)
			default:
				panic(gerror.Newf(`parse package failed, unsupported type "%s"`, reflect.TypeOf(obj.Decl)))
			}
		}
	}
	return &ParsedResult{
		Types: p.types,
	}
}

func (p *Parser) parsePackageType(pkg *packages.Package, typeName string) (parsedType *ParsedType) {
	typeId := p.makeTypeId(pkg.ID, typeName)
	if p.isStandardPackage(pkg) {
		return &ParsedType{
			Name:  typeName,
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

	if len(pkg.Syntax) == 0 {
		return nil
	}
	for _, syntax := range pkg.Syntax {
		for _, obj := range syntax.Scope.Objects {
			if obj.Name == typeName {
				parsedType = p.parseSyntaxTypeTypeSpec(pkg, syntax, obj.Decl.(*ast.TypeSpec))
				if parsedType == nil || parsedType.Name == "" {
					panic(gerror.Newf(`type "%s" not found in package "%s"`, typeName, pkg.ID))
				}
				return parsedType
			}
		}
	}
	panic(gerror.Newf(`type "%s" not found in package "%s"`, typeName, pkg.ID))
	return nil
}

func (p *Parser) parseSyntaxTypeTypeSpec(pkg *packages.Package, syntax *ast.File, decl *ast.TypeSpec) (parsedType *ParsedType) {
	if decl.Name != nil {
		typeId := p.makeTypeId(pkg.ID, decl.Name.Name)
		if p.isStandardPackage(pkg) {
			return &ParsedType{
				Name:  decl.Name.Name,
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
		if _, ok := p.parsingStructs[typeId]; ok {
			return &ParsedType{
				Name:    decl.Name.Name,
				Kind:    KindRefer,
				PkgPath: pkg.ID,
				Refer:   typeId,
			}
		}
		p.parsingStructs[typeId] = struct{}{}
		defer func() {
			delete(p.parsingStructs, typeId)
			p.saveParsedType(typeId, parsedType)
		}()
	}

	parsedType = p.parseSyntaxTypeExpr(pkg, syntax, decl.Type)
	if decl.Name != nil {
		parsedType.Name = decl.Name.Name
		if parsedType.PkgPath == "" {
			parsedType.PkgPath = pkg.ID
		}
	}
	return parsedType
}

func (p *Parser) parseSyntaxTypeValueSpec(pkg *packages.Package, syntax *ast.File, decl *ast.ValueSpec) *ParsedType {
	parsedType := p.parseSyntaxTypeExpr(pkg, syntax, decl.Type)
	return parsedType
}

func (p *Parser) parseSyntaxTypeFuncDecl(pkg *packages.Package, syntax *ast.File, decl *ast.FuncDecl) *ParsedType {
	parsedType := p.parseSyntaxTypeExpr(pkg, syntax, decl.Type)
	return parsedType
}

func (p *Parser) makeTypeId(pkgPath, typeName string) string {
	return gstr.TrimLeft(fmt.Sprintf(`%s.%s`, pkgPath, typeName), ".")
}

func (p *Parser) saveParsedType(typeId string, parsedType *ParsedType) {
	if parsedType != nil {
		pkgPath := p.types[typeId].PkgPath
		*p.types[typeId] = *parsedType
		parsedTypeId := p.makeTypeId(parsedType.PkgPath, parsedType.Name)
		if parsedTypeId != typeId {
			p.types[parsedTypeId] = &ParsedType{}
			*p.types[parsedTypeId] = *parsedType
			// Change the `typeId` to reference to `parsedTypeId`
			p.types[typeId] = &ParsedType{
				Name:  parsedType.Name,
				Kind:  parsedType.Kind,
				Refer: parsedTypeId,
			}
			parsedType.PkgPath = pkgPath
		}
	} else {
		delete(p.types, typeId)
	}
}

func (p *Parser) isStandardPackage(pkg *packages.Package) bool {
	_, ok := standardPackages[pkg.ID]
	return ok
}

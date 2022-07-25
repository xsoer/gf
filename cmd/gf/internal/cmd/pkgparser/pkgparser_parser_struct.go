package pkgparser

import (
	"go/ast"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"golang.org/x/tools/go/packages"
)

func (p *Parser) parseSyntaxTypeStructType(pkg *packages.Package, syntax *ast.File, expr *ast.StructType) *ParsedType {
	parsedType := &ParsedType{
		Name:    "",
		Kind:    KindStruct,
		PkgPath: pkg.ID,
		Elems:   make([]*ParsedType, 0),
	}
	var hasEmbeddedStruct bool
	for _, field := range expr.Fields.List {
		tempParsedType := p.parseSyntaxTypeExpr(pkg, syntax, field.Type)
		// Embedded struct.
		if len(field.Names) == 0 {
			hasEmbeddedStruct = true
			parsedType.Elems = append(parsedType.Elems, tempParsedType.Elems...)
			continue
		}
		// Value copy, not reference.
		fieldParsedType := &ParsedType{}
		*fieldParsedType = *tempParsedType
		fieldTypeId := p.makeTypeId(fieldParsedType.PkgPath, fieldParsedType.Name)
		if p.types[fieldTypeId] != nil {
			fieldParsedType = &ParsedType{
				Name:  tempParsedType.Name,
				Kind:  tempParsedType.Kind,
				Refer: fieldTypeId,
			}
		}
		fieldParsedType.PkgPath = ""
		if len(field.Names) > 0 {
			fieldParsedType.Name = field.Names[0].Name
		}
		// Ignore unexported fields.
		if fieldParsedType.Name == "" {
			panic("fieldParsedType.Name is empty")
		}
		if !gstr.IsExported(fieldParsedType.Name) {
			continue
		}
		if field.Tag != nil {
			fieldParsedType.Tag = field.Tag.Value
		}
		if field.Doc != nil {
			for _, comment := range field.Doc.List {
				fieldParsedType.Comment += comment.Text
			}
		}
		if field.Comment != nil {
			for _, comment := range field.Comment.List {
				fieldParsedType.Comment += comment.Text
			}
		}
		if fieldParsedType.Comment != "" {
			fieldParsedType.Comment = gstr.ReplaceByArray(fieldParsedType.Comment, g.SliceStr{
				"// ", "\n",
				"//", "\n",
			})
			fieldParsedType.Comment = gstr.Trim(fieldParsedType.Comment)
		}
		parsedType.Elems = append(parsedType.Elems, fieldParsedType)
	}
	// Filter repeated elements.
	if hasEmbeddedStruct {
		for i := 0; i < len(parsedType.Elems); {
			var filtered bool
			for j := i + 1; j < len(parsedType.Elems); j++ {
				if parsedType.Elems[i].Name == parsedType.Elems[j].Name {
					filtered = true
					parsedType.Elems = append(parsedType.Elems[:i], parsedType.Elems[i+1:]...)
				}
			}
			if !filtered {
				i++
			}
		}
	}
	return parsedType
}

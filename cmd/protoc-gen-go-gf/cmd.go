package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	contextPackage      = protogen.GoImportPath("context")
	goframePackage      = protogen.GoImportPath("github.com/gogf/gf/v2/frame/g")
	methodMessageTpl, _ = template.New("methodMessageStruct").Parse(methodMessageStruct)
	goModPath           = getGoModImportName()
)

func process(genFile *protogen.Plugin, file *protogen.File) {
	gen := genFile.NewGeneratedFile(file.GeneratedFilenamePrefix+".gf.go", file.GoImportPath)
	processCopyrightAndVersion(gen, file)
	processContent(gen, file)
}

func processCopyrightAndVersion(gen *protogen.GeneratedFile, file *protogen.File) {
	gen.P()
	gen.P("package ", file.GoPackageName)
	gen.P()
	gen.P("var _ = ", contextPackage.Ident("Background"), "()")
	gen.P("var _ = ", goframePackage.Ident("Meta"), "{}")
	gen.P()
}

func processContent(gen *protogen.GeneratedFile, file *protogen.File) {
	for _, svc := range file.Services {
		for _, method := range svc.Methods {
			if method == nil || method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
				continue
			}
			processMethod(gen, method, svc)
		}
	}
}

func processMethod(g *protogen.GeneratedFile, method *protogen.Method, svc *protogen.Service) map[string]interface{} {
	processMessage(g, method)
	return map[string]interface{}{
		"method_name": string(method.Desc.Name()),
		"in_name":     string(method.Input.Desc.Name()),
		"out_name":    string(method.Output.Desc.Name()),
	}
}

func processMessage(gen *protogen.GeneratedFile, method *protogen.Method) {
	processMessageFunc := func(message *protogen.Message, needGenGMeta bool) {
		methodMessageTplBuffer := bytes.NewBuffer(nil)
		err := methodMessageTpl.Execute(methodMessageTplBuffer, map[string]interface{}{
			"message_name": string(message.Desc.Name()),
			"fields":       processMessageFields(message, gen),
		})
		if err != nil {
			info("gf-gen-go-http: Execute template error: %s\n", err.Error())
			panic(err.Error())
		}
		gen.P(methodMessageTplBuffer.String())
		gen.P()
	}
	processMessageFunc(method.Input, true)
	processMessageFunc(method.Output, false)
}

func processMessageFields(message *protogen.Message, gen *protogen.GeneratedFile) []string {
	var (
		tagContent       string
		fieldDefinitions = make([]string, 0)
	)
	tagContent = processComment(message.Comments)
	if tagContent != "" {
		fieldDefinitions = append(fieldDefinitions, fmt.Sprintf(
			"g.Meta %s", tagContent,
		))
	}
	for _, field := range message.Fields {
		fieldDefinition := fmt.Sprintf("%s %s", field.GoName, processFieldType(field, gen))
		tagContent = processComment(field.Comments)
		if tagContent != "" {
			fieldDefinition += " " + tagContent
		}
		fieldDefinitions = append(fieldDefinitions, fieldDefinition)
	}
	return fieldDefinitions
}

func processFieldType(field *protogen.Field, gen *protogen.GeneratedFile) string {
	if field.Desc.IsWeak() {
		return "struct{}"
	}
	goType := field.Desc.Kind().String()
	if field.Desc.Kind() == protoreflect.MessageKind {
		goType = field.Message.GoIdent.GoName
		if field.GoIdent.GoImportPath != field.Message.GoIdent.GoImportPath {
			goType = gen.QualifiedGoIdent(getFullImportPath(field.Message.GoIdent.GoImportPath).Ident(field.Message.GoIdent.GoName))
		}
		goType = "*" + goType
	} else if field.Desc.Kind() == protoreflect.EnumKind {
		goType = gen.QualifiedGoIdent(field.Enum.GoIdent)
	} else if field.Desc.HasPresence() && field.Desc.Kind() != protoreflect.MessageKind && field.Desc.Kind() != protoreflect.BytesKind {
		goType = "*" + goType
	}
	if field.Desc.IsList() {
		return "[]" + goType
	}
	if field.Desc.IsMap() {
		keyType := processFieldType(field.Message.Fields[0], gen)
		valType := processFieldType(field.Message.Fields[1], gen)
		return fmt.Sprintf("map[%v]%v", keyType, valType)
	}
	return goType
}

func getFullImportPath(path protogen.GoImportPath) protogen.GoImportPath {
	if fullImportPath == nil || !*fullImportPath || filePath == nil || *filePath == "*" {
		return path
	}
	pathStr := strings.Trim(path.String(), "\"")
	filePathStr := *filePath
	filePathStr = strings.TrimLeft(filePathStr, ".")
	filePathStr = strings.TrimRight(filePathStr, "/")
	return protogen.GoImportPath(goModPath + filePathStr + pathStr)
}

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
	for _, message := range file.Messages {
		var buffer = bytes.NewBuffer(nil)
		processMessage(gen, buffer, message)
		gen.P(buffer.String())
	}
}

func processMessage(gen *protogen.GeneratedFile, buffer *bytes.Buffer, message *protogen.Message) {
	methodMessageTplBuffer := bytes.NewBuffer(nil)
	err := methodMessageTpl.Execute(methodMessageTplBuffer, map[string]interface{}{
		"message_name": string(message.Desc.Name()),
		"fields":       processMessageFields(gen, message),
	})
	if err != nil {
		info("gf-gen-go-http: Execute template error: %s\n", err.Error())
		panic(err.Error())
	}
	buffer.WriteString(methodMessageTplBuffer.String())
}

func processMessageFields(gen *protogen.GeneratedFile, message *protogen.Message) []string {
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
		fieldDefinition := fmt.Sprintf("%s %s", field.GoName, processFieldType(gen, field))
		tagContent = processComment(field.Comments)
		if tagContent != "" {
			fieldDefinition += " " + tagContent
		}
		fieldDefinitions = append(fieldDefinitions, fieldDefinition)
	}
	return fieldDefinitions
}

func processFieldType(gen *protogen.GeneratedFile, field *protogen.Field) string {
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
		keyType := processFieldType(gen, field.Message.Fields[0])
		valType := processFieldType(gen, field.Message.Fields[1])
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

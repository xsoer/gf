package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	tagForComment = "dc"
)

// processComment handles the comments from proto message/field,
// it returns the field tag and comment content for struct/field.
func processComment(comments protogen.CommentSet) string {
	var (
		err          error
		match        []string
		tagMap       = gmap.NewListMap()
		leading      = gstr.Trim(comments.Leading.String())
		trailing     = gstr.Trim(comments.Trailing.String())
		commentLines = append(gstr.SplitAndTrim(leading, "\n"), trailing)
	)
	for _, line := range commentLines {
		line = gstr.Trim(line, "/")
		match, err = gregex.MatchString(`^(\w+): (.+)$`, line)
		if err != nil {
			panic(err)
		}
		var (
			key   = tagForComment
			value string
		)
		if len(match) == 3 {
			key = gstr.Trim(match[1])
			value = gstr.Trim(match[2])
		} else {
			value = line
		}
		tagMap.Set(key, tagMap.GetVar(key).String()+value)
	}
	var tagContent string
	tagMap.Iterator(func(key, value interface{}) bool {
		if g.IsEmpty(value) {
			return true
		}
		if tagContent != "" {
			tagContent += " "
		}
		tagContent += fmt.Sprintf(`%s:"%s"`, key, value)
		return true
	})
	if tagContent == "" {
		return ""
	}
	return fmt.Sprintf("`%s`", tagContent)
}

func getGoModImportName() string {
	path, _ := os.Getwd()
	path += "/go.mod"
	modBytes, err := os.ReadFile(path)
	if err != nil {
		panic("get go mod module fail: " + err.Error())
	}
	line := strings.Split(string(modBytes), "\n")
	mod := strings.TrimPrefix(line[0], "module ")
	return mod
}

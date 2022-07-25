package cmd

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/pkgparser"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

var (
	Parse = cParse{}
)

type cParse struct {
	g.Meta `name:"parse" usage:"{cParseUsage}" brief:"{cParseBrief}" eg:"{cParseEg}"`
}

type cParseInput struct {
	g.Meta `name:"parse"`
	Src    string `name:"SRC" arg:"true" brief:"{cParseSrcBrief}"`
	Dst    string `name:"DST" arg:"true" brief:"{cParseDstBrief}"`
}
type cParseOutput struct{}

func (c cParse) Index(ctx context.Context, in cParseInput) (out *cParseOutput, err error) {
	var (
		path      = `/Users/john/Workspace/Go/GOPATH/src/git.code.oa.com/Khaos/eros/app/khaos-oss/api`
		recursive = false
		allDirs   = []string{path}
	)
	if recursive {
		var subDirs []string
		subDirs, err = gfile.ScanDirFunc(path, "*", recursive, func(path string) string {
			if gfile.IsDir(path) {
				return path
			}
			return ""
		})
		if err != nil {
			return nil, err
		}
		allDirs = append(allDirs, subDirs...)
	}
	var (
		tempParsedResult *pkgparser.ParsedResult
		parsedResult     = pkgparser.ParsedResult{
			Types: make(map[string]*pkgparser.ParsedType),
		}
	)
	for _, dirPath := range allDirs {
		mlog.Printf(`parsing: %s`, dirPath)
		tempParsedResult, err = pkgparser.Parse(dirPath)
		if err != nil {
			return nil, err
		}
		for k, v := range tempParsedResult.Types {
			parsedResult.Types[k] = v
		}
	}
	gfile.PutContents(`/Users/john/Temp/t.json`, gjson.MustEncodeString(parsedResult))
	return
}

package pkgparser

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"golang.org/x/tools/go/packages"
)

const pkgLoadMode = packages.NeedImports | packages.NeedSyntax | packages.NeedDeps

const (
	KindArray  = "array"
	KindStruct = "struct"
	KindMap    = "map"
	KindAny    = "any"
	KindFunc   = "func"
	KindStd    = "std"
	KindRefer  = "refer"
)

var (
	basicKindMap = map[string]struct{}{
		"array":      {},
		"struct":     {},
		"map":        {},
		"any":        {},
		"func":       {},
		"bool":       {},
		"int":        {},
		"int8":       {},
		"int16":      {},
		"int32":      {},
		"int64":      {},
		"uint":       {},
		"uint8":      {},
		"uint16":     {},
		"uint32":     {},
		"uint64":     {},
		"string":     {},
		"byte":       {},
		"float32":    {},
		"float64":    {},
		"complex64":  {},
		"complex128": {},
		"uintptr":    {},
	}
	standardPackages = make(map[string]struct{})
)

type ParsedType struct {
	Name     string // Attribute name or type name.
	Kind     string
	PkgPath  string        `json:",omitempty"`
	Elems    []*ParsedType `json:",omitempty"` // Struct field elements.
	KeyType  *ParsedType   `json:",omitempty"` // Map       key   type.
	ElemType *ParsedType   `json:",omitempty"` // Map/Array value type.
	Tag      string        `json:",omitempty"`
	Comment  string        `json:",omitempty"`
	Refer    string        `json:",omitempty"`
}

type ParsedResult struct {
	Types map[string]*ParsedType // map[PkgPath.Type]*ParsedType
}

type Parser struct {
	types          map[string]*ParsedType // map[PkgPath.Type]*ParsedType
	parsingStructs map[string]struct{}    // To avoid recursively infinite parsing for struct attributes.
}

func init() {
	stdPackages, err := packages.Load(nil, "std")
	if err != nil {
		panic(err)
	}
	for _, p := range stdPackages {
		standardPackages[p.ID] = struct{}{}
	}
}

func Parse(path string) (*ParsedResult, error) {
	cfg := &packages.Config{
		Dir:   path,
		Mode:  pkgLoadMode,
		Tests: false,
	}
	loadedPackages, err := packages.Load(cfg)
	if err != nil {
		return nil, gerror.Wrapf(err, `load package failed for "%s"`, path)
	}
	p := newParser()
	return p.parsePackage(loadedPackages[0]), nil
}

func newParser() *Parser {
	return &Parser{
		types:          make(map[string]*ParsedType),
		parsingStructs: make(map[string]struct{}),
	}
}

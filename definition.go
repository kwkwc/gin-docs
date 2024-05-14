package gin_docs

const (
	PROJECT_NAME    = "Gin-Docs"
	PROJECT_VERSION = Version
)

type KVMap map[string]string
type KVMapSlice []KVMap

func (ks KVMapSlice) Len() int           { return len(ks) }
func (ks KVMapSlice) Less(i, j int) bool { return ks[i]["name"] < ks[j]["name"] }
func (ks KVMapSlice) Swap(i, j int)      { ks[i], ks[j] = ks[j], ks[i] }

type RouterMap map[string][]KVMap
type DataMap map[string]RouterMap

var rootPath string

var templateMap = KVMap{
	"index":              "",
	"css_template_cdn":   "",
	"css_template_local": "",
	"js_template_cdn":    "",
	"js_template_local":  "",
}

var docMap = make(map[string]KVMap)

var pkgMap = make(map[string][]string)

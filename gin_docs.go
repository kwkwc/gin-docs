package gin_docs

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ApiDoc struct {
	Ge   *gin.Engine
	Conf *Config
}

func (d ApiDoc) init() (err error) {
	rootPath = d.getRootPath()
	if err := d.readTemplate(rootPath); err != nil {
		return err
	}

	d.getDocData()

	return
}

func (d ApiDoc) OnlineHtml() (err error) {
	if err := d.init(); err != nil {
		return err
	}

	if !d.Conf.Enable {
		return
	}

	dataMap := d.getApiData()

	d.Ge.Static(d.Conf.UrlPrefix+"/static", filepath.Join(rootPath, "static"))

	d.Ge.GET(d.Conf.UrlPrefix+"/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, d.renderHtml())
	})

	d.Ge.GET(d.Conf.UrlPrefix+"/data",
		verifyPassword(d.Conf.PasswordSha2),
		func(c *gin.Context) {
			urlPrefix := d.Conf.UrlPrefix
			referer := c.Request.Header.Get("referer")
			if referer == "" {
				referer = "http://127.0.0.1"
			}
			host := strings.Split(referer, urlPrefix)[0]

			c.JSON(http.StatusOK, gin.H{
				"PROJECT_NAME":    PROJECT_NAME,
				"PROJECT_VERSION": PROJECT_VERSION,
				"host":            host,
				"title":           d.Conf.Title,
				"version":         d.Conf.Version,
				"description":     d.Conf.Description,
				"noDocText":       d.Conf.NoDocText,
				"data":            dataMap,
			})
		})

	return
}

func (d ApiDoc) OfflineHtml(out string, force bool) (err error) {
	if out == "" {
		out = "htmldoc"
	}

	if err := d.init(); err != nil {
		return err
	}

	htmlStr := d.renderHtml()

	dataMap := d.getApiData()
	data := gin.H{
		"PROJECT_NAME":    PROJECT_NAME,
		"PROJECT_VERSION": PROJECT_VERSION,
		"host":            "http://127.0.0.1",
		"title":           d.Conf.Title,
		"version":         d.Conf.Version,
		"description":     d.Conf.Description,
		"noDocText":       d.Conf.NoDocText,
		"data":            dataMap,
	}

	dest := filepath.Join(".", out)
	if ok, _ := pathExists(dest); ok {
		if !force {
			return fmt.Errorf("target `%s` exists, set `force=true` to override.", dest)
		}
		if err := os.RemoveAll(dest); err != nil {
			return err
		}
	}
	if err := os.Mkdir(dest, os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(
		filepath.Join(dest, "index.html"), []byte(htmlStr), 0644,
	); err != nil {
		return err
	}
	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := os.WriteFile(
		filepath.Join(dest, "data"), dataByte, 0644,
	); err != nil {
		return err
	}

	if err := copyFolder(
		filepath.Join(rootPath, "static"), filepath.Join(dest, "static"),
	); err != nil {
		return err
	}

	return
}

func (d ApiDoc) OfflineMarkdown(out string, force bool) (err error) {
	if out == "" {
		out = "doc.md"
	}

	if err := d.init(); err != nil {
		return err
	}

	dataMap := d.getApiData()

	dest := filepath.Join(".", out)
	if ok, _ := pathExists(dest); ok {
		if !force {
			return fmt.Errorf("target `%s` exists, set `force=true` to override.", dest)
		}
	}

	md := ""
	for fullName := range dataMap {
		md += "# " + fullName + "\n\n"
		for _, item := range dataMap[fullName]["children"] {
			md += "## " + item["name"]
			if item["name_extra"] != "" {
				md += "(" + item["name_extra"] + ")"
			}
			md += "\n\n"
			md = d.handleMd(md, item)
			md += item["doc_md"] + "\n\n\n"
		}
		md += "\n\n"
	}

	if err := os.WriteFile(dest, []byte(md), 0644); err != nil {
		return err
	}

	return
}

func (d ApiDoc) handleMd(md string, item KVMap) string {
	md += "### url" + "\n"
	urls := strings.Split(item["url"], " ")
	if len(urls) == 1 {
		urls = []string{strings.Split(urls[0], "\t")[0]}
	}
	for i := range len(urls) {
		md += "- " +
			strings.Replace(
				strings.Replace(
					strings.Replace(
						urls[i], "\t", " ", -1,
					), "<", "&lt;", -1,
				), ">", "&gt;", -1,
			) +
			"\n\n"
	}
	if item["api_type"] == "api" {
		md += "### method" + "\n"
		md += "- " + item["method"] + "\n\n"
	}
	if item["doc"] == d.Conf.NoDocText && item["doc_md"] != "" {
		//
	} else {
		md += "### doc" + "\n"
		md += "```doc\n" + item["doc"] + "\n```\n\n"
	}
	return md
}

func rootPathFunc() {}
func (d ApiDoc) getRootPath() string {
	funcValue := reflect.ValueOf(rootPathFunc)
	fn := runtime.FuncForPC(funcValue.Pointer())
	filePath, _ := fn.FileLine(0)
	rp := filepath.Dir(filePath)

	return rp
}

func (d ApiDoc) readTemplate(rp string) error {
	templatesPath := filepath.Join(rp, "templates")
	for k := range templateMap {
		tByte, err := os.ReadFile(
			filepath.Join(templatesPath, k+".html"),
		)
		if err != nil {
			return err
		}
		templateMap[k] = string(tByte)
	}

	return nil
}

func (d ApiDoc) renderHtml() string {
	htmlStr := templateMap["index"]
	if d.Conf.Cdn {
		cssTemplate := templateMap["css_template_cdn"]
		jsTemplate := templateMap["js_template_cdn"]

		if d.Conf.CdnCssTemplate != "" {
			cssTemplate = d.Conf.CdnCssTemplate
		}
		if d.Conf.CdnJsTemplate != "" {
			jsTemplate = d.Conf.CdnJsTemplate
		}

		return strings.Replace(
			strings.Replace(
				htmlStr, "<!-- ___CSS_TEMPLATE___ -->", cssTemplate, -1,
			), "<!-- ___JS_TEMPLATE___ -->", jsTemplate, -1,
		)
	} else {
		return strings.Replace(
			strings.Replace(
				htmlStr, "<!-- ___CSS_TEMPLATE___ -->", templateMap["css_template_local"], -1,
			), "<!-- ___JS_TEMPLATE___ -->", templateMap["js_template_local"], -1,
		)
	}
}

func (d ApiDoc) getDocData() {
	for _, r := range d.Ge.Routes() {
		funcValue := reflect.ValueOf(r.HandlerFunc)
		if funcValue.Kind() != reflect.Func {
			continue
		}

		fn := runtime.FuncForPC(funcValue.Pointer())
		filePath, _ := fn.FileLine(0)

		if _, ok := docMap[filePath]; ok {
			continue
		}
		docMap[filePath] = make(KVMap)

		node, err := parser.ParseFile(token.NewFileSet(), filePath, nil, parser.ParseComments)
		if err != nil {
			slog.Error(fmt.Sprintf("%s err: %s\n", PROJECT_NAME, err))
			continue
		}
		ast.Inspect(node, func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok {
				fnName := fn.Name.Name
				docMap[filePath][fnName] = fn.Doc.Text()
			}
			return true
		})
	}
}

func (d ApiDoc) getApiData() DataMap {
	dataMap := make(DataMap)
	for _, r := range d.Ge.Routes() {
		pkgName, funcName := d.splitHandler(r.Handler)

		if slices.Contains(d.Conf.Exclude, pkgName) {
			continue
		}

		if dataMap[pkgName] == nil {
			dataMap[pkgName] = make(RouterMap)
		}
		if _, ok := dataMap[pkgName]["children"]; !ok {
			dataMap[pkgName]["children"] = []KVMap{}
		}

		if !slices.Contains(d.Conf.MethodsList, r.Method) {
			continue
		}

		url := fmt.Sprintf("%s\t[%s]", r.Path, r.Method)

		apiData := KVMap{
			"name":     funcName,
			"url":      url,
			"method":   r.Method,
			"router":   pkgName,
			"api_type": "api",
		}

		d.addApiData(dataMap, apiData, r.HandlerFunc)
	}

	for k := range dataMap {
		if len(dataMap[k]["children"]) == 0 {
			delete(dataMap, k)
		} else {
			sort.Sort(KVMapSlice(dataMap[k]["children"]))
		}
	}

	return dataMap
}

func (d ApiDoc) splitHandler(handler string) (string, string) {
	handlerS := strings.Split(filepath.Base(handler), ".")
	pkgName := handlerS[0]
	funcName := handlerS[len(handlerS)-1]

	if pkgMap[pkgName] == nil {
		pkgMap[pkgName] = []string{}
	}

	dirPath := filepath.Dir(handler)
	if !slices.Contains(pkgMap[pkgName], dirPath) {
		pkgMap[pkgName] = append(pkgMap[pkgName], dirPath)
	}

	index := slices.Index(pkgMap[pkgName], dirPath)
	if index > 0 {
		pkgName = pkgName + "-" + strconv.Itoa(index+1)
	}

	return pkgName, funcName
}

func (d ApiDoc) addApiData(dataMap DataMap, apiData KVMap, hFunc gin.HandlerFunc) {
	router := apiData["router"]

	resultList := []KVMap{}
	for _, v := range dataMap[router]["children"] {
		if v["name"] == apiData["name"] {
			resultList = append(resultList, v)
		}
	}
	if len(resultList) > 0 {
		for _, v := range [][]string{{"url", apiData["url"]}, {"method", apiData["method"]}} {
			sList := strings.Split(
				strings.Join([]string{resultList[0][v[0]], v[1]}, " "), " ",
			)
			slices.Sort(sList)
			sList = slices.CompactFunc(sList, strings.EqualFold)
			resultList[0][v[0]] = strings.Join(sList, " ")
		}
		return
	}

	doc := d.getApiDoc(hFunc, apiData["name"])
	apiData["name_extra"], apiData["doc"], apiData["doc_md"] = d.splitDoc(doc)

	dataMap[router]["children"] = append(dataMap[router]["children"], apiData)
}

func (d ApiDoc) getApiDoc(hFunc gin.HandlerFunc, hFuncName string) string {
	funcValue := reflect.ValueOf(hFunc)
	filePath, _ := runtime.FuncForPC(funcValue.Pointer()).FileLine(0)
	funcDoc := docMap[filePath][hFuncName]
	funcDoc = strings.Replace(funcDoc, "\t", strings.Repeat(" ", 4), -1)

	return funcDoc
}

func (d ApiDoc) cleanStr(str string) string {
	return strings.TrimSpace(
		strings.TrimSuffix(
			strings.TrimSpace(
				strings.TrimSuffix(
					strings.TrimSpace(str), "\n\n",
				),
			), "\n",
		),
	)
}

func (d ApiDoc) getFirstLineOfDoc(docSrc string) string {
	return d.cleanStr(
		strings.Split(
			strings.Split(docSrc, "\n\n")[0], "\n",
		)[0],
	)
}

func (d ApiDoc) splitDoc(docSrc string) (nameExtra, doc, docMd string) {
	docSrcS := strings.Split(docSrc, "@@@")
	doc = docSrcS[0]

	if doc != "" {
		nameExtra = d.getFirstLineOfDoc(doc)
	} else {
		nameExtra = ""
	}

	doc = strings.TrimRight(
		strings.TrimSuffix(
			strings.TrimRight(
				strings.TrimSuffix(
					strings.TrimRight(
						strings.Replace(doc, nameExtra, "", 1), " ",
					), "\n\n",
				), " ",
			), "\n",
		), " ",
	)

	if len(docSrcS) >= 2 {
		docMd = strings.TrimSpace(docSrcS[1])
	} else if doc != "" && d.Conf.AllMd {
		docMd = strings.TrimSpace(doc)
		doc = ""
	} else {
		docMd = ""
	}

	if docMd != "" {
		splitN := strings.Split(docSrcS[0], "\n")
		spaceCount := strings.Count(splitN[len(splitN)-1], " ")
		docMdS := strings.Split(docMd, "\n"+strings.Repeat(" ", spaceCount))
		docMd = strings.Join(docMdS, "\n")
	}

	if doc == "" {
		doc = d.Conf.NoDocText
	}

	return nameExtra, doc, docMd
}

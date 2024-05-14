package gin_docs

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/add_data", AddData)
	r.PATCH("/add_data", AddData)
	r.POST("/post_data", AddData)
	r.PUT("/post_data", AddData)
	r.DELETE("/delete_data", DeleteData)
	r.PUT("/change_data", ChangeData)

	return r
}

func setupOnlineHtml(r *gin.Engine) error {
	c := &Config{}
	c = c.Default()

	c.MethodsList = []string{"GET", "POST", "DELETE"}

	c.Title = "Test App"
	apiDoc := ApiDoc{Ge: r, Conf: c}
	err := apiDoc.OnlineHtml()

	return err
}

func setupOnlineHtmlDisable(r *gin.Engine) error {
	c := &Config{}
	c = c.Default()

	c.Enable = false

	apiDoc := ApiDoc{Ge: r, Conf: c}
	err := apiDoc.OnlineHtml()

	return err
}

func setupOnlineHtmlCdn(r *gin.Engine) error {
	c := &Config{}
	c = c.Default()

	c.Cdn = true
	c.CdnCssTemplate = "test_css"
	c.CdnJsTemplate = "test_js"

	apiDoc := ApiDoc{Ge: r, Conf: c}
	err := apiDoc.OnlineHtml()

	return err
}

func setupOnlineHtmlUnauthorized(r *gin.Engine) error {
	c := &Config{}
	c = c.Default()

	c.PasswordSha2 = "8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918"

	apiDoc := ApiDoc{Ge: r, Conf: c}
	err := apiDoc.OnlineHtml()

	return err
}

func setupOfflineHtml(r *gin.Engine) error {
	c := &Config{}
	c = c.Default()
	c.Title = "Test App"
	apiDoc := ApiDoc{Ge: r, Conf: c}
	err := apiDoc.OfflineHtml("", false)

	return err
}

func setupOfflineHtmlExists(r *gin.Engine) error {
	c := &Config{}
	c = c.Default()
	c.Title = "Test App"
	apiDoc := ApiDoc{Ge: r, Conf: c}
	err := apiDoc.OfflineHtml("htmldoc_exists", false)

	return err
}

func setupOfflineHtmlForce(r *gin.Engine) error {
	c := &Config{}
	c = c.Default()
	c.Title = "Test App"
	apiDoc := ApiDoc{Ge: r, Conf: c}
	err := apiDoc.OfflineHtml("htmldoc_exists2", true)

	return err
}

func setupOfflineMarkdown(r *gin.Engine) error {
	c := &Config{}
	c = c.Default()
	c.Title = "Test App"
	apiDoc := ApiDoc{Ge: r, Conf: c}
	err := apiDoc.OfflineMarkdown("", false)

	return err
}

func setupOfflineMarkdownExists(r *gin.Engine) error {
	c := &Config{}
	c = c.Default()
	c.Title = "Test App"
	apiDoc := ApiDoc{Ge: r, Conf: c}
	err := apiDoc.OfflineMarkdown("doc_exists.md", false)

	return err
}

func setupOfflineMarkdownForce(r *gin.Engine) error {
	c := &Config{}
	c = c.Default()
	c.Title = "Test App"
	apiDoc := ApiDoc{Ge: r, Conf: c}
	err := apiDoc.OfflineMarkdown("doc_exists2.md", true)

	return err
}

/*
Submission of data

Extra notes:

	{
		"data":{
			"xx": "xxx"
		}
	}
*/
func AddData(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

/*
@@@
### markdown
@@@
*/
func DeleteData(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func ChangeData(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func TestOnlineHtml(t *testing.T) {
	r := setupRouter()
	err := setupOnlineHtml(r)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/docs/api/", nil)
	assert.NoError(t, err)

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestOnlineHtmlData(t *testing.T) {
	r := setupRouter()
	err := setupOnlineHtml(r)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/docs/api/data", nil)
	assert.NoError(t, err)

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestOnlineHtmlDisable(t *testing.T) {
	r := setupRouter()
	err := setupOnlineHtmlDisable(r)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/docs/api/", nil)
	assert.NoError(t, err)

	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
}

func TestOnlineHtmlCdn(t *testing.T) {
	r := setupRouter()
	err := setupOnlineHtmlCdn(r)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/docs/api/", nil)
	assert.NoError(t, err)

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestOnlineHtmlUnauthorized(t *testing.T) {
	r := setupRouter()
	err := setupOnlineHtmlUnauthorized(r)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/docs/api/data", nil)
	assert.NoError(t, err)

	r.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestOfflineHtml(t *testing.T) {
	r := setupRouter()
	err := setupOfflineHtml(r)
	assert.NoError(t, err)

	ok, _ := pathExists(filepath.Join("htmldoc", "index.html"))
	assert.Equal(t, true, ok)

	err = os.RemoveAll("htmldoc")
	assert.NoError(t, err)
}

func TestOfflineHtmlShouldErrorWhenExists(t *testing.T) {
	err := os.Mkdir("htmldoc_exists", os.ModePerm)
	assert.NoError(t, err)

	r := setupRouter()
	err = setupOfflineHtmlExists(r)
	assert.EqualError(t, err, "target `htmldoc_exists` exists, set `force=true` to override.")

	err = os.RemoveAll("htmldoc_exists")
	assert.NoError(t, err)
}

func TestOfflineHtmlShouldOverrideWhenForce(t *testing.T) {
	err := os.Mkdir("htmldoc_exists2", os.ModePerm)
	assert.NoError(t, err)

	r := setupRouter()
	err = setupOfflineHtmlForce(r)
	assert.NoError(t, err)

	ok, _ := pathExists(filepath.Join("htmldoc_exists2", "index.html"))
	assert.Equal(t, true, ok)

	err = os.RemoveAll("htmldoc_exists2")
	assert.NoError(t, err)
}

func TestOfflineMarkdown(t *testing.T) {
	r := setupRouter()
	err := setupOfflineMarkdown(r)
	assert.NoError(t, err)

	ok, _ := pathExists(filepath.Join(".", "doc.md"))
	assert.Equal(t, true, ok)

	err = os.RemoveAll("doc.md")
	assert.NoError(t, err)
}

func TestOfflineMarkdownShouldErrorWhenExists(t *testing.T) {
	err := os.WriteFile("doc_exists.md", []byte(""), 0644)
	assert.NoError(t, err)

	r := setupRouter()
	err = setupOfflineMarkdownExists(r)
	assert.EqualError(t, err, "target `doc_exists.md` exists, set `force=true` to override.")

	err = os.RemoveAll("doc_exists.md")
	assert.NoError(t, err)
}

func TestOfflineMarkdownShouldOverrideWhenForce(t *testing.T) {
	err := os.WriteFile("doc_exists2.md", []byte(""), 0644)
	assert.NoError(t, err)

	r := setupRouter()
	err = setupOfflineMarkdownForce(r)
	assert.NoError(t, err)

	ok, _ := pathExists(filepath.Join(".", "doc_exists2.md"))
	assert.Equal(t, true, ok)

	err = os.RemoveAll("doc_exists2.md")
	assert.NoError(t, err)
}

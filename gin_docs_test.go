package gin_docs

import (
	"net/http"
	"net/http/httptest"
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

func setupApiDoc(r *gin.Engine) error {
	c := &Config{}
	c = c.Default()
	c.MethodsList = []string{"GET", "POST", "DELETE"}
	c.Title = "Test App"
	apiDoc := ApiDoc{Ge: r, Conf: c}
	err := apiDoc.Init()

	return err
}

func setupApiDocDisable(r *gin.Engine) error {
	c := &Config{}
	c = c.Default()
	c.Enable = false
	apiDoc := ApiDoc{Ge: r, Conf: c}
	err := apiDoc.Init()

	return err
}

func setupApiDocCdn(r *gin.Engine) error {
	c := &Config{}
	c = c.Default()
	c.Cdn = true
	c.CdnCssTemplate = "test_css"
	c.CdnJsTemplate = "test_js"
	apiDoc := ApiDoc{Ge: r, Conf: c}
	err := apiDoc.Init()

	return err
}

func setupApiDocUnauthorized(r *gin.Engine) error {
	c := &Config{}
	c = c.Default()
	c.PasswordSha2 = "8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918"
	apiDoc := ApiDoc{Ge: r, Conf: c}
	err := apiDoc.Init()

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

func TestDocsApi(t *testing.T) {
	router := setupRouter()
	err := setupApiDoc(router)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/docs/api/", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestDocsApiData(t *testing.T) {
	router := setupRouter()
	err := setupApiDoc(router)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/docs/api/data", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestDocsApiDisable(t *testing.T) {
	router := setupRouter()
	err := setupApiDocDisable(router)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/docs/api/", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
}

func TestDocsApiCdn(t *testing.T) {
	router := setupRouter()
	err := setupApiDocCdn(router)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/docs/api/", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestDocsApiUnauthorized(t *testing.T) {
	router := setupRouter()
	err := setupApiDocUnauthorized(router)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/docs/api/data", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

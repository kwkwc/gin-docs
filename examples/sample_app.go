// go run examples/sample_app.go

package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	gd "github.com/kwkwc/gin-docs"
)

func main() {
	r := gin.Default()
	r.POST("/api/todo", AddTodo)
	r.GET("/api/todo", GetTodo)

	c := &gd.Config{}
	apiDoc := gd.ApiDoc{Ge: r, Conf: c.Default()}
	err := apiDoc.Init()
	if err != nil {
		slog.Error(fmt.Sprintf("Gin-Docs init err: %s", err))
		os.Exit(1)
	}

	err = r.Run()
	if err != nil {
		slog.Error(fmt.Sprintf("Start service err: %s", err))
		os.Exit(1)
	}
}

/*
Add todo

### args
|  args | required | location | type   |  help    |
|-------|----------|----------|--------|----------|
| name  |  true    |  json    | string | todo name |
| type  |  true    |  json    | string | todo type |

### request
```json
{"name": "xx", "type": "code"}
```

### response
```json
{"code": xxxx, "msg": "xxx", "data": null}
```
*/
func AddTodo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"todo": "post todo",
	})
}

/*
Get todo

### description
> Get todo

### args
|  args | required | location |  type  |  help    |
|-------|----------|----------|--------|----------|
|  name |  true    |  query   | string | todo name |
|  type |  false   |  query   | string | todo type |

### request
```
http://127.0.0.1:8080/api/todo?name=xxx&type=code
```

### response
```json
{"code": xxxx, "msg": "xxx", "data": null}
```
*/
func GetTodo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"todo": "get todo",
	})
}

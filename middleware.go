package gin_docs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func verifyPassword(passwordSha2 string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPasswordSha2 := c.Request.Header.Get("Auth-Password-SHA2")
		if passwordSha2 != "" && passwordSha2 != authPasswordSha2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
		}
	}
}

package middlewares

import "github.com/gin-gonic/gin"

func HeadersHandler(c *gin.Context) {
	c.Header("X-Powered-By", "gitpages")
	c.Next()
}

package middlewares

import "github.com/gin-gonic/gin"

func ErrorHandler(c *gin.Context) {
	c.Next()

	err := c.Errors.Last()
	if err != nil {
		_, _ = c.Writer.Write([]byte(err.Err.Error()))
	}
}

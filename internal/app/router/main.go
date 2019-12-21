package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oott123/gitpages/internal/app/repoman"
)

func Main(c *gin.Context) {
	repo := repoman.MatchHost(c.Request.Host)
	if repo == nil {
		_ = c.AbortWithError(404, fmt.Errorf("cannot find repo for host %s", c.Request.Host))
		return
	}

	repo.TreeRLock()
	defer repo.TreeRUnlock()

	handler := repo.HttpHandler()
	handler.ServeHTTP(c.Writer, c.Request)
}

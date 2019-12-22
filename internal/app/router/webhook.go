package router

import (
	"crypto/subtle"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oott123/gitpages/internal/app/repoman"
)

func WebHook(c *gin.Context) {
	repo := repoman.MatchHost(c.Request.Host)
	if repo == nil {
		_ = c.AbortWithError(404, fmt.Errorf("cannot find repo for host %s", c.Request.Host))
		return
	}

	secret := c.Param("secret")
	if secret == "" {
		_ = c.AbortWithError(403, fmt.Errorf("secret is needed for update"))
		return
	}

	if repo.WebHookSecret() == "" {
		_ = c.AbortWithError(400, fmt.Errorf("secret is not set in config file"))
		return
	}

	if !timingSafeCompareString(secret, repo.WebHookSecret()) {
		_ = c.AbortWithError(403, fmt.Errorf("secret is invalid"))
		return
	}

	go func() {
		err := repo.Update()
		if err != nil {
			log.Errorf("error while updating repo: %s", err)
		}
	}()

	c.String(202, "repo is scheduled to update")
}

func timingSafeCompareString(a, b string) bool {
	lena := len(a)
	lenb := len(b)

	if subtle.ConstantTimeEq(int32(lena), int32(lenb)) != 1 {
		return false
	}

	bytea := ([]byte)(a)
	byteb := ([]byte)(b)

	return subtle.ConstantTimeCompare(bytea, byteb) == 1
}

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oott123/gitpages/internal/app/repoman"
	"github.com/oott123/gitpages/internal/app/router"
	"github.com/oott123/gitpages/pkg/config"
	"github.com/oott123/gitpages/pkg/logger"
)

func main() {
	r := gin.Default()

	r.Use(router.Main)

	cfg := config.Get()
	log := logger.New()

	err := repoman.ReloadRepos()
	if err != nil {
		panic(err)
	}

	go log.Infof("trying to handle request on %s", cfg.Endpoint)
	err = r.Run(cfg.Endpoint)
	if err != nil {
		panic(err)
	}
}
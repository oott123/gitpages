package main

import (
	"github.com/bep/debounce"
	"github.com/gin-gonic/gin"
	"github.com/oott123/gitpages/internal/app/repoman"
	"github.com/oott123/gitpages/internal/app/router"
	"github.com/oott123/gitpages/internal/middlewares"
	"github.com/oott123/gitpages/pkg/config"
	"github.com/oott123/gitpages/pkg/logger"
	"time"
)

func main() {
	r := gin.Default()
	r.Use(middlewares.ErrorHandler, middlewares.HeadersHandler)

	r.GET("/_gitpages/update/:secret", router.WebHook)
	r.POST("/_gitpages/update/:secret", router.WebHook)
	r.PUT("/_gitpages/update/:secret", router.WebHook)

	r.Use(router.Main)

	cfg := config.Get()
	log := logger.New()

	err := repoman.ReloadRepos()
	if err != nil {
		panic(err)
	}

	reloadConfig := func() {
		log.Infof("config changed, reloading...")
		err := repoman.ReloadRepos()
		if err != nil {
			log.Errorf("failed to reload config: %s", err)
		} else {
			log.Infof("reloaded config")
		}
	}
	debounced := debounce.New(100 * time.Millisecond)

	config.Watch(func(c *config.Config) {
		debounced(reloadConfig)
	})

	go log.Infof("trying to handle request on %s", cfg.Endpoint)
	err = r.Run(cfg.Endpoint)
	if err != nil {
		panic(err)
	}
}

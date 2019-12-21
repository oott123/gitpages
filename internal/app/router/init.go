package router

import (
	"github.com/oott123/gitpages/pkg/logger"
	"go.uber.org/zap"
	"mime"
)

var log *zap.SugaredLogger

func init() {
	log = logger.New()
	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		panic(err)
	}
}

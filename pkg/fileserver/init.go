package fileserver

import (
	"github.com/oott123/gitpages/pkg/logger"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func init() {
	log = logger.New()
}

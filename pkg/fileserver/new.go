package fileserver

import (
	"github.com/oott123/gitpages/pkg/logger"
	"go.uber.org/zap"
)

type FileServer struct {
	serverConfig ServerConfig
	accessConfig AccessConfig
	log          *zap.SugaredLogger
}

type ServerConfig struct {
	Root         string
	AllowSymlink bool
}

func New(serverConfig *ServerConfig, accessConfig *AccessConfig) (*FileServer, error) {
	if serverConfig == nil {
		serverConfig = &ServerConfig{}
	}

	server := &FileServer{
		serverConfig: *serverConfig,
		accessConfig: *accessConfig,
		log:          logger.New(),
	}

	return server, nil
}

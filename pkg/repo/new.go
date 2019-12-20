package repo

import (
	"fmt"
	"github.com/oott123/gitpages/pkg/config"
	"github.com/oott123/gitpages/pkg/logger"
	"go.uber.org/zap"
	"gopkg.in/src-d/go-git.v4"
	"net/http"
	"sync"
)

type Repo struct {
	srv      *config.Server
	baseDir  string
	bareLock *sync.RWMutex
	treeLock *sync.RWMutex
	git      *git.Repository
	log      *zap.SugaredLogger
	httpHandler *http.Handler
}

func New(cfg *config.Server, baseDir string) (*Repo, error) {
	if cfg == nil {
		return nil, fmt.Errorf("cannot create repo for nil config")
	}

	remote := cfg.Remote
	if remote == "" {
		return nil, fmt.Errorf("cannot create repo %s for empty origin", cfg.Host)
	}


	repo := Repo{srv: cfg, bareLock: &sync.RWMutex{}, treeLock: &sync.RWMutex{}, log: logger.New(), baseDir: baseDir}

	httpHandler := http.FileServer(http.Dir(repo.ServeDir()))
	repo.httpHandler = &httpHandler
	
	return &repo, nil
}

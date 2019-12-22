package repo

import (
	"fmt"
	"github.com/oott123/gitpages/pkg/config"
	"github.com/oott123/gitpages/pkg/logger"
	"go.uber.org/zap"
	"gopkg.in/src-d/go-git.v4"
	"net/http"
	"sync"
	"time"
)

type Repo struct {
	srv                *config.Server
	baseDir            string
	bareLock           *sync.RWMutex
	treeLock           *sync.RWMutex
	git                *git.Repository
	log                *zap.SugaredLogger
	httpHandler        http.Handler
	updateTicker       *time.Ticker
	updateTickerCancel *chan bool
}

func New(cfg *config.Server, baseDir string) (*Repo, error) {
	if cfg == nil {
		return nil, fmt.Errorf("cannot create repo for nil config")
	}

	remote := cfg.Remote
	if remote == "" {
		return nil, fmt.Errorf("cannot create repo %s for empty origin", cfg.Host)
	}

	if cfg.Branch == "" {
		cfg.Branch = "master"
	}

	repo := Repo{srv: cfg, bareLock: &sync.RWMutex{}, treeLock: &sync.RWMutex{}, log: logger.New(), baseDir: baseDir}

	httpHandler := http.FileServer(http.Dir(repo.ServeDir()))
	repo.httpHandler = httpHandler

	if cfg.UpdateInterval > 1*time.Second {
		repo.updateTicker = time.NewTicker(cfg.UpdateInterval)
		ch := make(chan bool)
		repo.updateTickerCancel = &ch
		go repo.tick()
	}

	return &repo, nil
}

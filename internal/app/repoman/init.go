package repoman

import (
	"github.com/oott123/gitpages/pkg/logger"
	"github.com/oott123/gitpages/pkg/repo"
	"go.uber.org/zap"
	"sync"
)

var repos []*repo.Repo
var repoLock *sync.RWMutex
var log *zap.SugaredLogger

func init() {
	repoLock = &sync.RWMutex{}
	repos = make([]*repo.Repo, 0)
	log = logger.New()
}

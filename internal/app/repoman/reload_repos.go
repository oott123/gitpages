package repoman

import (
	"fmt"
	"github.com/oott123/gitpages/pkg/config"
	"github.com/oott123/gitpages/pkg/fileserver"
	"github.com/oott123/gitpages/pkg/repo"
)

func ReloadRepos() error {
	repoLock.Lock()
	defer repoLock.Unlock()

	cfg := config.Get()
	newRepos := make([]*repo.Repo, len(cfg.Servers))

	for i, c := range cfg.Servers {
		srv := c
		r, err := repo.New(&srv, cfg.StorageDir)
		if err != nil {
			return fmt.Errorf("reload create repo error: %w", err)
		}
		newRepos[i] = r
		err = r.CloneOrOpen()
		if err != nil {
			return fmt.Errorf("reload init repo error: %w", err)
		}
		err = r.Update()
		if err != nil {
			return fmt.Errorf("reload update repo error: %w", err)
		}

		accessConfig := ParseAccessConfig(r.ServeDir())
		serverConfig := fileserver.ServerConfig{
			Root:         r.ServeDir(),
			AllowSymlink: c.AllowSymlink,
		}

		fsrv, err := fileserver.New(&serverConfig, accessConfig)
		if err != nil {
			return fmt.Errorf("reload create fileserver error: %w", err)
		}

		r.SetHttpHandler(fsrv)
	}

	for i, r := range repos {
		r.Close()
		repos[i] = nil
	}

	repos = newRepos

	return nil
}

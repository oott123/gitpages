package repo

import (
	"fmt"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"os"
)

func (r *Repo) Open() error {
	r.bareLock.Lock()
	defer r.bareLock.Unlock()

	bareDir := r.bareDir()
	treeDir := r.treeDir()

	if s, err := os.Stat(bareDir); os.IsNotExist(err) || (s != nil && !s.IsDir()) {
		return fmt.Errorf("directory is not exists or it's not dir: %s", bareDir)
	}
	if s, err := os.Stat(treeDir); os.IsNotExist(err) || (s != nil && !s.IsDir()) {
		return fmt.Errorf("directory is not exists or it's not dir: %s", treeDir)
	}

	bare := filesystem.NewStorage(osfs.New(bareDir), cache.NewObjectLRUDefault())
	tree := osfs.New(treeDir)

	r.log.Debugf("opening git repo %s at %s", r.srv.Remote, bareDir)
	gr, err := git.Open(bare, tree)
	if err != nil {
		return fmt.Errorf("failed to open git repo: %w", err)
	}

	r.git = gr

	return nil
}

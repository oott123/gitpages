package repo

import (
	"fmt"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"os"
)

func (r *Repo) Clone() error {
	r.bareLock.Lock()
	defer r.bareLock.Unlock()

	bareDir := r.bareDir()
	treeDir := r.treeDir()

	if s, err := os.Stat(bareDir); !os.IsNotExist(err) && (s != nil && !s.IsDir()) {
		return fmt.Errorf("directory is exists and it's not dir: %s", bareDir)
	}
	if s, err := os.Stat(treeDir); !os.IsNotExist(err) && (s != nil && !s.IsDir()) {
		return fmt.Errorf("directory is exists and it's not dir: %s", treeDir)
	}

	bare := filesystem.NewStorage(osfs.New(bareDir), cache.NewObjectLRUDefault())
	tree := osfs.New(treeDir)

	r.log.Debugf("cloning git repo %s at %s", r.srv.Remote, bareDir)

	gr, err := git.Clone(bare, tree, &git.CloneOptions{
		URL:        r.srv.Remote,
		Auth:       nil,
		RemoteName: "origin",
		NoCheckout: true,
		// Don't use shallow clone; we can't update shallow clones https://github.com/src-d/go-git/issues/1143
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          nil,
		Tags:              git.NoTags,
	})
	if err != nil {
		return fmt.Errorf("git clone error: %w", err)
	}
	r.git = gr

	return nil
}

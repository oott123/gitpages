package repo

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (r *Repo) Update() error {
	r.bareLock.Lock()
	defer r.bareLock.Unlock()
	
	if r.git == nil {
		return fmt.Errorf("repo %s has not been loaded", r.srv.Remote)
	}

	r.log.Debugf("fetching git repo %s", r.srv.Remote)
	err := r.git.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Depth:      1,
		Progress:   nil,
		Tags:       git.NoTags,
		Force:      true,
	})

	r.log.Debugf("fetched git repo %s %s", r.srv.Remote, err)
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("%s fetch failed: %w", r.srv.Remote, err)
	}

	ref, err := r.git.Reference(plumbing.NewBranchReferenceName(r.srv.Branch), true)
	if err != nil {
		return fmt.Errorf("%s branch %s resolve failed: %w", r.srv.Remote, r.srv.Branch, err)
	}

	r.treeLock.Lock()
	defer r.treeLock.Unlock()

	worktree, err := r.git.Worktree()
	if err != nil {
		return fmt.Errorf("%s worktree error: %w", r.srv.Remote, err)
	}

	r.log.Debugf("checking out git repo %s", r.srv.Remote)
	err = worktree.Checkout(&git.CheckoutOptions{
		Hash: ref.Hash(),
		Create: false,
		Force:  false,
		Keep:   false,
	})
	if err != nil {
		return fmt.Errorf("%s checkout error: %w", r.srv.Remote, err)
	}

	return nil
}
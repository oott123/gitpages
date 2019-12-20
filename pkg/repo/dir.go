package repo

import (
	"github.com/oott123/gitpages/pkg/hash"
	"path"
)

func (r *Repo) hash() string {
	return hash.SHA1(r.srv.Remote)
}

func (r *Repo) bareDir() string {
	bareDir := path.Join(r.baseDir, r.hash(), "repo")
	return bareDir
}

func (r *Repo) treeDir() string {
	treeDir := path.Join(r.baseDir, r.hash(), "tree")
	return treeDir
}

func (r *Repo) ServeDir() string {
	return path.Join(r.treeDir(), r.srv.Dir)
}
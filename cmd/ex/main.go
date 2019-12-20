package main

import (
	"github.com/oott123/gitpages/pkg/config"
	"github.com/oott123/gitpages/pkg/repo"
)

func main() {
	cfg := config.Get()
	r, err := repo.New(&cfg.Servers[0], cfg.StorageDir)
	if err != nil {
		panic(err)
	}
	err = r.CloneOrOpen()
	if err != nil {
		panic(err)
	}
	err = r.Update()
	if err != nil {
		panic(err)
	}
}

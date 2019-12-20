package repoman

import "github.com/oott123/gitpages/pkg/repo"

func MatchHost(host string) *repo.Repo {
	repoLock.RLock()
	defer repoLock.RUnlock()

	for _, r := range repos {
		if r.MatchHost(host) {
			return r
		}
	}

	return nil
}

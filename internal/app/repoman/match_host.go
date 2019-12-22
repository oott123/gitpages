package repoman

import "github.com/oott123/gitpages/pkg/repo"

func MatchHost(host string) *repo.Repo {
	repoLock.RLock()
	defer repoLock.RUnlock()

	for _, r := range repos {
		if r.MatchHost(host) {
			log.Debugf("host %s matched repo %s", host, r.String())
			return r
		}
	}

	return nil
}

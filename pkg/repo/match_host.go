package repo

func (r *Repo) MatchHost(host string) bool {
	return r.srv.MatchHost(host)
}

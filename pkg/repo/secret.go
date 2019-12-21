package repo

func (r *Repo) WebHookSecret() string {
	return r.srv.WebHookSecret
}

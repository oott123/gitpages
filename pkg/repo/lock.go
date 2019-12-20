package repo

func (r *Repo) TreeRLock() {
	r.treeLock.RLock()
}

func (r *Repo) TreeRUnlock() {
	r.treeLock.RUnlock()
}

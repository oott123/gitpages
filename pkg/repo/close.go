package repo

func (r *Repo) Close() {
	r.treeLock.Lock()
	r.bareLock.Lock()
	defer r.treeLock.Unlock()
	defer r.bareLock.Unlock()

	if r.git != nil {
		r.git = nil
	}
	if r.updateTickerCancel != nil {
		*r.updateTickerCancel <- true
		r.updateTickerCancel = nil
	}
}

package repo

func (r *Repo) tick() {
	for {
		select {
		case <-r.updateTicker.C:
			if r.git != nil {
				r.log.Infof("updating repo %s", r.String())
				err := r.Update()
				if err != nil {
					r.log.Errorf("update repo %s failed: %s", r.String(), err)
				} else {
					r.log.Infof("updated repo %s", r.String())
				}
			}
		case <-*r.updateTickerCancel:
			return
		}
	}
}

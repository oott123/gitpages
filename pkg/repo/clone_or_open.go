package repo

import "os"

func (r *Repo) CloneOrOpen() error {
	r.bareLock.RLock()

	var shouldClone bool
	if _, err := os.Stat(r.bareDir()); os.IsNotExist(err) {
		shouldClone = true
	}

	r.bareLock.RUnlock()

	if shouldClone {
		return r.Clone()
	} else {
		return r.Open()
	}
}

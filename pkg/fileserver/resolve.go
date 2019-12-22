package fileserver

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (f *FileServer) Resolve(requestPath string) string {
	if containsDotDot(requestPath) {
		return ""
	}

	root, err := filepath.Abs(f.serverConfig.Root)
	if err != nil {
		f.log.Warnf("root dir %s abs failed: %s", root, err)
		return ""
	}

	filePath := path.Join(root, requestPath)
	filePath = filepath.Clean(filePath)

	var pathOk bool
	absPath, err := filepath.Abs(filePath)
	if err == nil {
		if f.serverConfig.AllowSymlink {
			// allow symlink, skip check
			pathOk = true
		} else {
			if _, err := os.Stat(absPath); os.IsNotExist(err) {
				// not exists, skip symlink check
				pathOk = true
			} else {
				// check symlink
				resolvedPath, err := filepath.EvalSymlinks(absPath)
				if err == nil {
					relPath, err := filepath.Rel(root, resolvedPath)
					if err == nil {
						if !strings.Contains(relPath, filepath.FromSlash("../")) {
							pathOk = true
						} else {
							f.log.Warnf("requested file %s is not inside repo root", resolvedPath)
						}
					} else {
						f.log.Warnf("failed to resolve file: %s", err)
					}
				} else {
					f.log.Warnf("failed to eval symlink for file: %s", err)
				}
			}
		}
	} else {
		f.log.Warnf("failed to find file %s: %s", filePath, err)
	}

	if pathOk {
		return absPath
	}
	return ""
}

func containsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}
	for _, ent := range strings.FieldsFunc(v, isSlashRune) {
		if ent == ".." {
			return true
		}
	}
	return false
}

func isSlashRune(r rune) bool { return r == '/' || r == '\\' }

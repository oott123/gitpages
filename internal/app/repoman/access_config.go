package repoman

import (
	"github.com/oott123/gitpages/pkg/fileserver"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	"path"
)

func ParseAccessConfig(root string) *fileserver.AccessConfig {
	config := fileserver.AccessConfig{}
	filename := path.Join(root, ".gitpagesfile")

	if s, err := os.Stat(filename); err == nil && s != nil && !s.IsDir() {
		buf, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Errorf("failed to read %s: %s", filename, err)
		}
		err = toml.Unmarshal(buf, &config)
		if err != nil {
			log.Errorf("failed to parse config file %s: %s", err)
		}
	} else {
		log.Debugf("gitpagesfile %s not exists; use default", filename)
	}

	return &config
}

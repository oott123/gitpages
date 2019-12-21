package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/minio/minio/pkg/wildcard"
	"github.com/oott123/gitpages/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
	"sync"
)

var latest *Config
var lock *sync.Mutex
var log *zap.SugaredLogger

func init() {
	lock = &sync.Mutex{}
	log = logger.New()
}

func Get() *Config {
	lock.Lock()
	defer lock.Unlock()

	if latest == nil {
		cfg, err := Read()
		if err != nil {
			log.Infof("config file load failed, using default instead: %s", err)
			d := Default()
			latest = &d
			return latest
		}

		latest = cfg
	}

	return latest
}

func Read() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to read in config: %w", err)
	}

	result := Default()
	err = viper.Unmarshal(&result)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return &result, nil
}

func Watch(callback func(c *Config)) {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		cfg, err := Read()
		if err != nil {
			log.Errorf("config read error while reloading: %w", err)
			return
		}
		lock.Lock()
		defer lock.Unlock()
		latest = cfg

		callback(cfg)
	})
}

func FindServer(host string) *Server {
	cfg := Get()

	for _, s := range cfg.Servers {
		matched := wildcard.MatchSimple(strings.ToLower(s.Host), strings.ToLower(host))
		if matched {
			return &s
		}
	}

	log.Warnf("cannot find server for host %s", host)
	return nil
}

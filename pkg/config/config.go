package config

type Config struct {
	// Endpoint is where GitPages will listened on
	Endpoint   string
	StorageDir string
	Servers    []Server
}

func Default() Config {
	return Config{
		Endpoint:   ":2289",
		StorageDir: "data",
		Servers: []Server{
			{
				Host:          "*",
				Remote:        "https://github.com/oott123/gitpages-example.git",
				WebHookSecret: "gitpages",
				Branch:        "master",
				Dir:           "/",
			},
		},
	}
}

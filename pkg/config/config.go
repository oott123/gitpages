package config

type Config struct {
	// Endpoint is where GitPages will listened on
	Endpoint string
	StorageDir string
	Servers []Server `toml:"server"`
}

func Default() Config {
	return Config{
		Endpoint: ":2289",
		StorageDir: "data",
		Servers: []Server{
			{
				Host: "*",
				Remote: "https://github.com/Yelp/yelp.github.io.git",
				WebHookSecret: "example",
				Branch: "master",
				Dir: "/",
			},
		},
	}
}

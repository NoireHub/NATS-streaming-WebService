package webservice

type Config struct {
	BindAddr      string `toml:"bind_addr"`
	LogLevel      string `toml:"log_level"`
	DatabaseURL   string `toml:"database_url"`
	StanClusterID string `toml:"stan_clusterID"`
	StanSubject   string `toml:"stan_subj"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}

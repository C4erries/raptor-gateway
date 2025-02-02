package server

type Config struct {
	Bindaddr       string `toml:"bind_addr"`
	UserServiceUrl string `toml:"user-service_url"`
	LogLevel       string `toml:"log_level"`
}

func NewConfig() *Config {
	return &Config{
		Bindaddr:       ":8080",
		UserServiceUrl: "localhost:5001",
		LogLevel:       "debug",
	}
}

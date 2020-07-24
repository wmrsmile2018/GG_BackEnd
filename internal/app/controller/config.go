package controller

//toml - библиотека

type Config struct {
	BindAddr    string `toml:"bind_addr"` //адрес, на котором запускается сервер
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
	SessionsKey string `toml:"sessions_key"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8000",
		LogLevel: "debug",
	}
}

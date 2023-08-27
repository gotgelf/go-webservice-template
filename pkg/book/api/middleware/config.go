package middleware

import (
	"io"
	"os"
)

type Config struct {
	Output io.Writer
}

var ConfigDefault = Config{
	Output: os.Stdout,
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.Output == nil {
		cfg.Output = ConfigDefault.Output
	}

	return cfg
}

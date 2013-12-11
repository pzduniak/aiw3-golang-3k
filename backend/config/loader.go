package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

func Load() Config {
	var cfg Config

	if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
		log.Panicln(err)
	}

	return cfg
}

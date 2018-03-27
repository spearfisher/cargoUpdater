package utils

import (
	"github.com/BurntSushi/toml"
)

type config struct {
	Login    string
	Password string
	Start    int
	Stop     int
	Period   int
}

// AppConfig - application configuration
var AppConfig config

func init() {
	_, err := toml.DecodeFile("./config.toml", &AppConfig)
	if err != nil {
		panic(err)
	}
}

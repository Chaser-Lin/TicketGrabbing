package config

import (
	"Project/MyProject/cache"
	"Project/MyProject/db"
	"Project/MyProject/event"
)

var Conf *Config

type Config struct {
	Mysql             db.Config    `yaml:"Mysql"`
	Redis             cache.Config `yaml:"Redis"`
	Kafka             event.Config `yaml:"Kafka"`
	ServerPort        string       `yaml:"ServerPort"`
	AdminName         string       `yaml:"AdminName"`
	AdminPassword     string       `yaml:"AdminPassword"`
	TokenSymmetricKey string       `yaml:"TokenSymmetricKey"`
}

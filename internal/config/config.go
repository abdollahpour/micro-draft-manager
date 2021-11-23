package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	WriteTimeout int32  `envconfig:"write_timeout" required:"true" default:"15"`
	ReadTimeout  int32  `envconfig:"read_timeout" required:"true" default:"15"`
	Port         int32  `envconfig:"port" required:"true" default:"8080"`
	Host         string `envconfig:"host" required:"true" default:"0.0.0.0"`
	MongoUri     string `envconfig:"mongo_uri" required:"true" default:"mongodb://localhost:27017/drafts"`
}

func NewEnvConfiguration() Configuration {
	var conf Configuration
	envconfig.Process("", &conf)
	return conf
}

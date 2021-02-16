package habitat

import (
	"os"

	log "github.com/sirupsen/logrus"
)

//Config is the goddamned config struct
type Config struct { 
	Env		string
}

var config *Config 

//GetConfig Loads and returns the configuration appropriate for the current environment
func GetConfig() *Config { 

	if config != nil { 
		return config
	}

	log.Debug("Loading config")

	env := "development"
	if val, ok := os.LookupEnv("HABITAT_ENV"); ok { 
		env = val 
	}

	config = &Config{
		Env: env,
	}

	return config
}
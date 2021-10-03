package habconf

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/komkom/toml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ConfigFile --
const ConfigFile = "habitat.toml"

// HabitatDir -- 
const HabitatDir = ".habitat"


type RenderHTTPErrors struct { 

	// MethodNotAllowed will cause Habitat to render HTTP 405 where applicable - i.e. for unimplemented HTTP methods where an implemented method is configured on the same path.
	// If set to `false`, Habitat will render HTTP 404 instead. 
	MethodNotAllowed 	bool
	
	// UnprocessableEntity will cause Habitat to render HTTP 422 for requests that are syntactically correct but malformed due to missing or mistyped parameters. 
	// If set to `false`, Habitat will render HTTP 400 instead. 
	// HTTP 422 was added as a part of WebDAV in RFC 4918 (2007) for this situation, and I like it, but later RFC 7231 (2014) clarified that HTTP 400 is acceptable for this situation too.  See this post for more: https://softwareengineering.stackexchange.com/a/342896
	UnprocessableEntity	bool
}

var Errors = struct { 
	RenderHTTPErrors RenderHTTPErrors

	FallbackToOtherTypes bool 
	FallbackToHabitatTemplate bool 

} { 
	
	RenderHTTPErrors: RenderHTTPErrors { 
		MethodNotAllowed: 		true,
		UnprocessableEntity:	true,
	},

	FallbackToOtherTypes: true,
	FallbackToHabitatTemplate: true,

}


var MaxFormSize = int64(1<<63 - 1)





//Config is the goddamned config struct
type Config struct { 
	Env				string

	RootDir 		string
	ProjectDirs 	[]string


	binSass			string 
	binWebpack 		string 

	toml			tomlRoot
}

type tomlRoot struct { 

}


var config *Config 

// GetConfig presumes that the config has already been loaded with LoadConfig() and returns it
func GetConfig() * Config { 
	if config == nil { 
		panic("Config not loaded")
	}
	return config 
}

//GetConfig loads and returns the configuration appropriate for the current environment, or returns the memoized config if it has already been loaded
func LoadConfig() (*Config, error) { 

	if config != nil { 
		return config, nil 
	}


	// Environment
	env := "development"
	if val, ok := os.LookupEnv("HABITAT_ENV"); ok { 
		env = val 
	}


	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Debug().Str("env", env).Msg("Loading config")


	// Root
	cwd, err := os.Getwd()
	if err != nil { 
		panic(err) 
	}

	projectRoot := findRoot(cwd) 
	log.Debug().Msgf("Found project root: %s", projectRoot)
	os.Chdir( projectRoot ) 



	// Read config file
	tomlVal := tomlRoot{} 

	cfp := filepath.Join(projectRoot, ConfigFile)
	tomlFile, err := os.Open(cfp)
	if err != nil { 
		log.Err(err).Msgf("Could not open config file '%s'", ConfigFile)
		return nil, err
	}
	defer tomlFile.Close()

	tomlDecoder := json.NewDecoder( toml.New( tomlFile ) )
	err = tomlDecoder.Decode( &tomlVal )	
	if err != nil { 
		log.Err(err).Msgf("Could not read config file '%s'", ConfigFile)
		return nil, err
	}



	config = &Config{
		Env: 			env,

		RootDir: 		projectRoot,
		ProjectDirs: 	[]string { "src/" },

		toml: 			tomlVal,
	}


	if err := config.loadNode(projectRoot); err != nil { 
		return nil, err 
	}

	return config, nil 
}




// Project root
func findRoot( dir string ) string { 
	fmt.Println(dir)

	if _, err := os.Stat( path.Join(dir, ConfigFile ) ); err == nil { 
		return dir
	}


	if dir == "/" { 
		log.Fatal().Msgf("Could not find `%s`", ConfigFile)
		panic(0);
	}

	return findRoot( filepath.Dir( dir ) ) // Parent directory
}






// GetProjectRootDir --
func (c Config) GetProjectRootDir() string { 
	return c.RootDir
}


func (c Config) getProjectDir(path ...string) string { 
	res := filepath.Join( append( []string{ c.RootDir } , path[:]... )... )
	os.MkdirAll(res, os.FileMode(int(0777)))
	return res 
}


// GetDir --
func (c Config) GetDir() string { 
	return filepath.Join( c.RootDir, HabitatDir )
}
// GetDirCache --
func (c Config) GetDirCache() string { 
	return c.getProjectDir( HabitatDir, "tmp", "cache" )
}
// GetDirOutput --
func (c Config) GetDirOutput() string { 
	return c.getProjectDir( HabitatDir, "tmp", "output")
}




// OpenLogFile opens the given path within the log directory.
func (c Config) OpenLogFile(path string, flag int) (*os.File, string, error) { 
	path = filepath.Join(c.getProjectDir( HabitatDir, "logs" ), path) 
	file, error := os.OpenFile(path, flag, os.FileMode(int(0777)))

	if error != nil { 
		log.Err(error).Msgf("Error opening file '%s'", path)
	}

	return file, path, error
}

func (c Config) OpenLogFileTruncate(path string) (*os.File, string, error) { 
	return c.OpenLogFile(path, os.O_RDWR | os.O_CREATE | os.O_TRUNC)
}



// GetNodeSass Path to Sass executable
func (c Config) GetNodeSass() string { 
	return c.binSass
}
// GetNodeWebpack Path to webpack executable
func (c Config) GetNodeWebpack() string { 
	return c.binWebpack
}
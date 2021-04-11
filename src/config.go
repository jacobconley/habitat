package habitat

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"

	"github.com/komkom/toml"
	log "github.com/sirupsen/logrus"
)

// ConfigFile --
const ConfigFile = "habitat.toml"

// HabitatDir -- 
const HabitatDir = ".habitat"


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

	Database 		tomlDB // `json:"database"`

}


var config *Config 

//GetConfig Loads and returns the configuration appropriate for the current environment
func GetConfig() (*Config, error) { 
	log.SetLevel(log.DebugLevel)

	if config != nil { 
		return config, nil 
	}
	log.Debug("Loading config")


	// Environment
	env := "development"
	if val, ok := os.LookupEnv("HABITAT_ENV"); ok { 
		env = val 
	}


	// Root
	cwd, err := os.Getwd()
	if err != nil { 
		panic(err) 
	}

	projectRoot := findRoot(cwd) 
	log.Debugf("Found project root: %s", projectRoot)
	os.Chdir( projectRoot ) 



	// Read config file
	tomlVal := tomlRoot{} 

	cfp := filepath.Join(projectRoot, ConfigFile)
	tomlFile, err := os.Open(cfp)
	if err != nil { 
		log.Errorf("Could not open config file '%s'", ConfigFile)
		log.Error(err)
		return nil, err
	}
	defer tomlFile.Close()

	tomlDecoder := json.NewDecoder( toml.New( tomlFile ) )
	err = tomlDecoder.Decode( &tomlVal )	
	if err != nil { 
		log.Errorf("Could not read config file '%s'", ConfigFile)
		log.Error(err)
		return nil, err
	}
	


	config = &Config{
		Env: 			env,

		RootDir: 		projectRoot,
		ProjectDirs: 	[]string { "src/" },

		toml: 			tomlVal,
	}

	if err := config.processDB(); err != nil { 
		log.Error("Could not process DB config")
		return nil, err 
	}


	if err := config.loadNode(projectRoot); err != nil { 
		return nil, err 
	}

	return config, nil 
}




// Project root
func findRoot( dir string ) string { 

	if _, err := os.Stat( path.Join(dir, ConfigFile ) ); err == nil { 
		return dir
	}


	if dir == "/" { 
		log.Errorf("Could not find `%s`", ConfigFile)
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


func (c Config) GetDirDB() string { 
	return c.getProjectDir( "db" )
}
func (c Config) GetDirMigrations() string { 
	return c.getProjectDir( "db", "migrations" )
}

func (c Config) GetDirModels() string { 
	return c.getProjectDir( "db", "models" )
}




// OpenLogFile opens the given path within the log directory.
func (c Config) OpenLogFile(path string, flag int) (*os.File, string, error) { 
	path = filepath.Join(c.getProjectDir( HabitatDir, "logs" ), path) 
	file, error := os.OpenFile(path, flag, os.FileMode(int(0777)))

	if error != nil { 
		log.Errorf("Error opening file '%s'", path)
		log.Error(error)
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
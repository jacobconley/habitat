package auxiliary

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/jacobconley/habitat/habconf"

	log "github.com/sirupsen/logrus"
)

// CSS --
type CSS struct {

	SourceDirs 		[]string 
	SourceFiles 	[]string

	TargetFile 			string
}

// NewCSSFromConfig Creates the default loader implied by the given habconf.Config
func NewCSSFromConfig(config * habconf.Config) CSS { 
	return CSS { 
		SourceDirs: 	config.ProjectDirs,
		SourceFiles: 	[]string{},

		TargetFile: 	".habitat/out/css",
	}
}

func GetCssLoaders(config * habconf.Config) []CSS { 
	return []CSS { 
		NewCSSFromConfig(config),
	}
}

func BuildCSS() (error, []CSS) { 
	config, err := habconf.LoadConfig()
	if err != nil { 
		return err, nil
	}

	loader := NewCSSFromConfig(config)
	return loader.Build(), []CSS{ loader }
}


// PathMatches CSS or Sass files 
func (css CSS) PathMatches(path string) bool { 
	re, _ := regexp.Compile("\\.(c|sa|sc)ss$")
	return re.MatchString(path)
}




// Build finds all source files and executes `sass` on them with the confingured target file 
func (css CSS) Build() error { 

	// Find files

	files := css.SourceFiles

	for _, dir := range css.SourceDirs { 

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error { 

			if css.PathMatches(path) { 
				files = append(files, path)
			}
			return nil 

		} )

		if err != nil {
			log.Error("Error looking for source files", err) 
			return err 
		}
	}


	config, err := habconf.LoadConfig()
	if err != nil { 
		return err 
	}

	log.Debugf("[SASS] Building to %s", css.TargetFile)

	


	// Do the buildin' 

	if err := os.MkdirAll(filepath.Dir(css.TargetFile), os.FileMode(int(0777))); err != nil { 
		log.Error("[SASS] Could not create output directory; ", err) 
		return err
	}


	// Output (target) file
	targetFile, tferr := os.OpenFile(css.TargetFile, os.O_CREATE | os.O_RDWR | os.O_TRUNC, os.FileMode(int(0666)))
	if tferr != nil { 
		log.Error("[SASS] Could not open target file ^ :", tferr)
		return tferr
	}

	defer targetFile.Close()


	// Log file
	logFile, logFilepath, logErr := config.OpenLogFileTruncate("sass.log")
	if logErr != nil { 
		return logErr
	}
	log.Debugf("[SASS] Logging to '%s'", logFilepath)

	defer logFile.Close()



	// Each input file
	for _, fpath := range files { 

		log.Debugf("[SASS] Building '%s'", fpath)

		// The CSS parser we use in tests gets confused about comments
		// targetFile.WriteString( fmt.Sprintf("\n/* --- %s --- */\n\n", fpath) )

		cmd := exec.Command( config.GetNodeSass(), fpath)
		stdout, outErr := cmd.StdoutPipe()
		stderr, errErr := cmd.StderrPipe()


		if outErr != nil || errErr != nil { 
			log.Error("-> Error initializing pipes")
			log.Error("   stdOut: ", outErr)
			log.Error("   stdErr: ", errErr) 
			
			if outErr != nil { 
				return outErr
			} else { 
				return errErr 
			}
		}

		cmd.Start()		

		_, outErr = io.Copy(targetFile, stdout)
		_, errErr = io.Copy(logFile, stderr)

		if outErr != nil || errErr != nil { 
			log.Error("-> Error piping output")
			log.Error("   stdOut: ", outErr)
			log.Error("   stdErr: ", errErr) 
			
			if outErr != nil { 
				return outErr
			} else { 
				return errErr 
			}
		}
	}


	return nil
}

// Watch -- 
func (css CSS) Watch() error { 
	return nil 
}
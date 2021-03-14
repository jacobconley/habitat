package loaders

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	habitat "habitat/src"

	log "github.com/sirupsen/logrus"
)

// CSS --
type CSS struct {

	SourceDirs 		[]string 
	SourceFiles 	[]string

	TargetFile 			string
}

// NewCSSFromConfig Creates the default loader implied by the given habitat.Config
func NewCSSFromConfig(config * habitat.Config) CSS { 
	return CSS { 
		SourceDirs: 	config.ProjectDirs,
		SourceFiles: 	[]string{},

		TargetFile: 	".habitat/out/css",
	}
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
				log.Debugf("[SASS] (-> %s) %s", css.TargetFile, path) 
				files = append(files, path)
			}
			return nil 

		} )

		if err != nil {
			log.Error("Error looking for source files", err) 
			return err 
		}
	}


	config, err := habitat.GetConfig()
	if err != nil { 
		return err 
	}

	log.Debugf("[SASS] Building to %s", css.TargetFile)

	


	// Do the buildin' 

	// Output (target) file
	targetFile, tferr := os.OpenFile(css.TargetFile, os.O_CREATE | os.O_RDWR | os.O_TRUNC, os.FileMode(int(0666)))
	if tferr != nil { 
		log.Error("[SASS] Could not open target file ^ :", tferr)
		return tferr
	}

	defer targetFile.Close()
	targetWriter := bufio.NewWriter(targetFile)


	// Log file
	logFile, logFilepath, logErr := config.OpenLogFileTruncate("sass.log")
	if logErr != nil { 
		log.Error("[SASS] Could not open log file", logErr)
		return logErr
	}
	log.Debugf("[SASS] Logging to '%s'", logFilepath)

	defer logFile.Close()
	logWriter := bufio.NewWriter(logFile)




	var retval error 


	// Each input file
	for _, fpath := range files { 

		cmd := exec.Command( config.GetNodeSass(), fpath)
		stdout, outErr := cmd.StdoutPipe()
		stderr, errErr := cmd.StderrPipe()


		if outErr != nil || errErr != nil { 
			log.Error("[SASS] Error initializing pipes")
			log.Error("[SASS] stdOut: ", outErr)
			log.Error("[SASS] stdErr: ", errErr) 
			
			if outErr != nil { 
				return outErr
			} else { 
				return errErr 
			}
		}


		outReader := bufio.NewReader(stdout)
		errReader := bufio.NewReader(stderr) 
		outBuf := make([]byte, 256)
		errBuf := make([]byte, 256)

		readOut := true
		readErr := true 

		cmd.Start()		

		for { 
			if readOut && retval == nil { 
				outN, outErr := outReader.Read(outBuf)

				if outN == 0 { 
					log.Debug("[SASS] -> STDOUT Done")
					readOut = false 
				} else { 

					if outErr != nil { 
						log.Error("[SASS] Error piping output to target file: ", outErr)
						retval = outErr
					} else { 
						_, tferr = targetWriter.Write(outBuf[0:outN])

						if tferr != nil { 
							log.Error("[SASS] Error writing to target file: ", tferr)
							retval = tferr
						}
					}

				}
			}


			if readErr && retval == nil { 
				errN, errErr := errReader.Read(errBuf)

				if errN == 0 { 
					log.Debug("[SASS] -> STDERR Done")
					readErr = false 
				} else { 

					if errErr != nil { 
						log.Error("[SASS] Error piping stderr to log file: ", errErr)
						retval = errErr
					} else { 
						_, logErr = logWriter.Write(errBuf[0:errN])

						if logErr != nil { 
							log.Error("[SASS] Error writing to log file: ", logErr)
							retval = logErr
						}
					}

				}
			}


			if (readOut || readErr) == false { 
				break
			}
		}
	}


	return retval
}

// Watch -- 
func (css CSS) Watch() error { 
	return nil 
}
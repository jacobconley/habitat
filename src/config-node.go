package habitat

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Node executables - except we presume yarn :)
func (c Config) loadNode(dir string) error { 
	log.Debug("Looking for Node executables (via Yarn)...")

	cmd := exec.Command("yarn", "bin", "--silent")
	stdout, err := cmd.Output()

	if err != nil { 
		if err.(* exec.ExitError).ExitCode() == 127 { 
			return errors.New("Could not find Yarn! Is it installed?")
		}

		log.Error("Could not run Yarn", err)
		return err 
	}

	yarnBin := strings.TrimSpace(string(stdout))
	_, err = os.Stat(yarnBin)
	if err != nil { 
		return errors.New("Could not verify the Yarn installation - try running `yarn bin` and checking for warnings")
	}


	var binErr error = nil 
	getBin := func(module string, into * string) { 
		if binErr != nil { 
			return
		}

		cmd := exec.Command("yarn", "bin", "--silent", module)
		stdout, err := cmd.Output()
		res := strings.TrimSpace((string(stdout)))

		if err != nil { 
			log.Errorf("Could not find Yarn executable `%s`", module) 
			log.Error(err) 
			binErr = err
			return 
		}

		if _, err = os.Stat(res); err != nil { 
			log.Errorf("Could not find Yarn executable `%s` at expected path `%s`", module, res)
			binErr = err 
			return 
		}

		log.Debugf("-> Found %s: %s", module, res) 
		*into = res
	}

	getBin("sass", 		&config.binSass)
	getBin("webpack", 	&config.binWebpack)
	return binErr
}




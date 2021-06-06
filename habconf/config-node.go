package habconf

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

// Node executables - except we presume yarn :)
func (c Config) loadNode(dir string) error { 
	log.Debug().Msg("Looking for Node executables (via Yarn)...")

	cmd := exec.Command("yarn", "bin", "--silent")
	stdout, err := cmd.Output()

	if err != nil { 
		if err.(* exec.ExitError).ExitCode() == 127 { 
			return errors.New("Could not find Yarn! Is it installed?")
		}

		log.Err(err).Msg("Could not run Yarn")
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
			log.Err(err).Msgf("Could not find Yarn executable `%s`", module) 
			binErr = err
			return 
		}

		if _, err = os.Stat(res); err != nil { 
			log.Err(err).Msgf("Could not find Yarn executable `%s` at expected path `%s`", module, res)
			binErr = err 
			return 
		}

		log.Debug().Msgf("-> Found %s: %s", module, res) 
		*into = res
	}

	getBin("sass", 		&config.binSass)
	getBin("webpack", 	&config.binWebpack)
	return binErr
}




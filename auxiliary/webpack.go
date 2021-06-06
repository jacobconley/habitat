package auxiliary

import (
	"os"
	"os/exec"

	"github.com/jacobconley/habitat/habconf"

	"github.com/rs/zerolog/log"
)

func BuildWebpack() error {
	config, err := habconf.LoadConfig()
	if err != nil { 
		return err
	}

	logFile, logFilepath, logErr := config.OpenLogFileTruncate("webpack.log")
	if logErr != nil { 
		return logErr 
	}
	log.Debug().Msgf("[WBPK] Logging to %s", logFilepath)

	defer logFile.Close()


	os.Chdir( config.RootDir )

	cmd := exec.Command( config.GetNodeWebpack() )
	output, err := cmd.CombinedOutput()

	if err != nil { 
		log.Err(err).Msg("[WBPK] Error executing build")
		return err 
	}

	_, err = logFile.Write(output)
	if err != nil { 
		log.Warn().Err(err).Msg("[WBPK] Could not write to log file")
	}

	return nil
}
package peripherals

import (
	habitat "habitat/src"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func BuildWebpack() error {
	config, err := habitat.GetConfig()
	if err != nil { 
		return err
	}

	logFile, logFilepath, logErr := config.OpenLogFileTruncate("webpack.log")
	if logErr != nil { 
		return logErr 
	}
	log.Debugf("[WBPK] Logging to %s", logFilepath)

	defer logFile.Close()


	os.Chdir( config.RootDir )

	cmd := exec.Command( config.GetNodeWebpack() )
	output, err := cmd.CombinedOutput()

	if err != nil { 
		log.Error("[WBPK] Error executing build: ", err)
		return err 
	}

	_, err = logFile.Write(output)
	if err != nil { 
		log.Warn("[WBPK] Could not write to log file: ", err)
	}

	return nil
}
package logging

import (
	"fmt"
	"time"

	"transfDoc/conf"
)

// getLogFilePath get the log file save path
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", conf.GetConfig().RuntimeRootPath, conf.GetConfig().LogSavePath)
}

// getLogFileName get the save name of the log file
func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		conf.GetConfig().LogSaveName,
		time.Now().Format(conf.GetConfig().TimeFormat),
		conf.GetConfig().LogFileExt,
	)
}

package logging

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"io"
	"log"
	"os"
	"time"
	"transfDoc/conf"
)

var logger zerolog.Logger

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		"log-",
		time.Now().Format("2006-01-02"),
		"log",
	)
}

func Setup() {
	conf := conf.GetConfig()
	var lever zerolog.Level
	switch conf.LogLevel {
	case "debug":
		lever = zerolog.DebugLevel
	case "info":
		lever = zerolog.InfoLevel
	case "error":
		lever = zerolog.ErrorLevel
	case "warn":
		lever = zerolog.WarnLevel
	case "trace":
		lever = zerolog.TraceLevel
	case "fatal":
		lever = zerolog.FatalLevel
	default:
		lever = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(lever)
	logger = zerolog.New(os.Stdout).With().CallerWithSkipFrameCount(3).Timestamp().Logger()
	if conf.LogSavePath != "" { //为空串输出到控制台
		var err error
		fileName := getLogFileName()
		var f *os.File
		f, err = MustOpen(fileName, conf.LogSavePath)
		if err != nil {
			log.Fatalf("logging.Setup err: %v", err)
		}
		logger = zerolog.New(io.MultiWriter(f, os.Stdout)).
			With().Str("app", conf.AppName).CallerWithSkipFrameCount(3).Timestamp().Logger()
	}
}

func Debug(v ...interface{}) {
	Debugf(v...)
}
func Debugf(v ...interface{}) {
	requestId, format, msg := handleLogfParameter(v...)
	logger.Debug().Str("requestId", requestId).Msgf(format, msg...)
}

func Errorf(v ...interface{}) {
	requestId, format, msg := handleLogfParameter(v...)
	logger.Info().Str("requestId", requestId).Msgf(format, msg...)
}
func Error(v ...interface{}) {
	Errorf(v...)
}
func Fatalf(v ...interface{}) {
	requestId, format, msg := handleLogfParameter(v...)
	logger.Fatal().Str("requestId", requestId).Msgf(format, msg...)
}

func Info(v ...interface{}) {
	Infof(v...)
}
func Infof(v ...interface{}) {
	requestId, format, msg := handleLogfParameter(v...)
	logger.Info().Str("requestId", requestId).Msgf(format, msg...)
}
func handleLogfParameter(args ...interface{}) (string, string, []interface{}) {

	if len(args) < 1 {
		return "", "", args
	}
	switch args[0].(type) {
	case *gin.Context:
		if len(args) < 2 {
			return args[0].(*gin.Context).Request.Header.Get("requestId"), "", args[1:]
		}
		switch args[1].(type) {
		case string:
			return args[0].(*gin.Context).Request.Header.Get("requestId"), args[1].(string), args[2:]
		default:
			return args[0].(*gin.Context).Request.Header.Get("requestId"), "", args[1:]
		}
	case string:
		return "", args[0].(string), args[1:]
	default:
		return "", "", args
	}
}

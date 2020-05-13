package mylog

import (
	"EX_okexquant/config"
	"EX_okexquant/data"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

var (
	Logger     zerolog.Logger
	DataLogger zerolog.Logger
)

func SgNow() time.Time {
	return time.Now().In(data.Location)
}

func ConfigLoggers() {
	var perm = os.ModePerm
	dir, err := filepath.Abs(filepath.Dir(config.Config.Server.BaseLogPath))
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(dir, perm)
	if err != nil {
		panic(err)
	}
	//zerolog.CallerSkipFrameCount = 3
	zerolog.TimestampFunc = SgNow

	log.Logger = zerolog.New(&lumberjack.Logger{
		Filename:   config.Config.Server.BaseLogPath + "/" + config.Config.BaseLog.FileName,
		MaxSize:    config.Config.BaseLog.FileMaxSize,
		MaxBackups: config.Config.BaseLog.FileMaxBackups,
		MaxAge:     config.Config.BaseLog.FileMaxAge,
		Compress:   true,
	}).With().Caller().Timestamp().Logger()
	fmt.Println("Default logger init succeed.")

	Logger = zerolog.New(&lumberjack.Logger{
		Filename:   config.Config.Server.BaseLogPath + "/" + config.Config.BaseLog.FileName,
		MaxSize:    config.Config.BaseLog.FileMaxSize,
		MaxBackups: config.Config.BaseLog.FileMaxBackups,
		MaxAge:     config.Config.BaseLog.FileMaxAge,
		Compress:   true,
	}).With().Caller().Timestamp().Logger()
	fmt.Println("Logger init succeed.")

	DataLogger = zerolog.New(&lumberjack.Logger{
		Filename:   config.Config.Server.BaseLogPath + "/" + config.Config.DataLog.FileName,
		MaxSize:    config.Config.DataLog.FileMaxSize,
		MaxBackups: config.Config.DataLog.FileMaxBackups,
		MaxAge:     config.Config.DataLog.FileMaxAge,
		Compress:   true,
	}).With().Caller().Timestamp().Logger()
	fmt.Println("Logger init succeed.")
}

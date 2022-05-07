package config

import (
	"errors"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"time"
)

func init() {
	err := os.Mkdir(LogPath, 0777)
	if err != nil {
		if !errors.Is(err, fs.ErrExist) {
			log.Errorf("创建日志目录失败：%v", err)
		}
	}
	err = os.Mkdir(ChunkPath, 0777)
	if err != nil {
		if !errors.Is(err, fs.ErrExist) {
			log.Errorf("创建chunk目录失败：%v", err)
		}
	}
	log.AddHook(newLogHook())
}

func newLogHook() log.Hook {
	logPath := LogPath
	infoWriter, err := rotatelogs.New(
		logPath+"/Info/"+"Info"+".%Y%m%d%H%M",
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 保留1周内的日志
		rotatelogs.WithRotationTime(24*time.Hour), // 分割日期一天
	)
	if err != nil {
		log.Infof("failed to log to file, err:%v", err)
	}
	warnWriter, err := rotatelogs.New(
		logPath+"/Warn/"+"Warn"+".%Y%m%d%H%M",
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 保留1周内的日志
		rotatelogs.WithRotationTime(24*time.Hour), // 分割日期一天
	)
	if err != nil {
		log.Infof("failed to log to file, err:%v", err)
	}
	errWriter, err := rotatelogs.New(
		logPath+"/Err/"+"Err"+".%Y%m%d%H%M",
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 保留1周内的日志
		rotatelogs.WithRotationTime(24*time.Hour), // 分割日期一天
	)
	if err != nil {
		log.Infof("failed to log to file, err:%v", err)
	}
	log.SetLevel(log.InfoLevel)
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		log.InfoLevel:  infoWriter,
		log.WarnLevel:  warnWriter,
		log.ErrorLevel: errWriter,
	}, &log.TextFormatter{DisableColors: true})
	return lfsHook
}

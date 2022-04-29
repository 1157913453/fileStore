package config

import (
	"filestore/util"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var logpath = "/log/fileStore"

func init() {
	exits, err := util.PathExists(logpath)
	if err != nil {
		log.Errorf("查询日志目录出错：%v", err)
	}
	if !exits {
		err = os.Mkdir(logpath, 0777)
		if err != nil {
			log.Errorf("创建日志目录失败：%v", err)
		}
	}
	log.AddHook(newLogHook())

}

func newLogHook() log.Hook {
	logPath := "/log/fileStore"
	infoWriter, err := rotatelogs.New(
		logPath+"/Info/"+"Info"+".%Y%m%d%H%M",
		rotatelogs.WithMaxAge(7*24*time.Hour),      // 保留1周内的日志
		rotatelogs.WithRotationTime(3*time.Second), // 3秒切换
	)
	if err != nil {
		log.Infof("failed to log to file, err:%v", err)
	}
	warnWriter, err := rotatelogs.New(
		logPath+"/Warn/"+"Warn"+".%Y%m%d%H%M",
		rotatelogs.WithMaxAge(7*24*time.Hour),      // 保留1周内的日志
		rotatelogs.WithRotationTime(3*time.Second), // 3秒切换
	)
	if err != nil {
		log.Infof("failed to log to file, err:%v", err)
	}
	errWriter, err := rotatelogs.New(
		logPath+"/Err/"+"Err"+".%Y%m%d%H%M",
		rotatelogs.WithMaxAge(7*24*time.Hour),      // 保留1周内的日志
		rotatelogs.WithRotationTime(3*time.Second), // 3秒切换
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

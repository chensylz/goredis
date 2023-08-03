package logger

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type Logrus struct {
	log *logrus.Logger
}

func NewLogrus() *Logrus {
	log := logrus.New()
	logFilePath := "logs/log"
	// 设置日志格式，可以根据需要进行调整
	log.SetFormatter(&logrus.TextFormatter{})

	// 使用 rotatelogs 创建轮转日志文件
	rotateLog, err := rotatelogs.New(
		logFilePath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(logFilePath),
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 最大保留时间
		rotatelogs.WithRotationTime(24*time.Hour), // 轮转周期
	)
	if err != nil {
		log.Panicf("config local file system logger error. %v", err)
		return nil
	}
	log.SetOutput(rotateLog)

	return &Logrus{log: log}
}

func (l *Logrus) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *Logrus) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *Logrus) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"

	"github.com/chensylz/goredis/internal/global/constants"
)

type Log interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

var rLog Log

type Logrus struct {
	log *logrus.Entry
}

func SetupLog() {
	rLog = NewLogrus()
}

func Info(args ...interface{}) {
	rLog.Info(args...)
}

func Error(args ...interface{}) {
	rLog.Error(args...)
}

func Debug(args ...interface{}) {
	rLog.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	rLog.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	rLog.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	rLog.Errorf(format, args...)
}

func NewLogrus() *Logrus {
	logureLog := logrus.New()
	logFilePath := "logs/log"
	// 设置日志格式，可以根据需要进行调整
	logureLog.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: constants.TimestampFormat,
		FullTimestamp:   true,
		ForceColors:     true,
	})

	// 创建一个 WriteSyncer，用于同时输出到文件和终端
	fileWriter, err := newFileWriter(logFilePath)
	if err != nil {
		log.Panicf("config local file system logger error. %v", err)
		return nil
	}
	multiWriter := io.MultiWriter(os.Stdout, fileWriter)
	logureLog.SetOutput(multiWriter)
	logureLog.AddHook(&CallerHook{})

	return &Logrus{log: logureLog.WithFields(logrus.Fields{})}
}

func newFileWriter(logFilePath string) (io.Writer, error) {
	writer, err := rotatelogs.New(
		logFilePath+".%Y%m%d%H%M",
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		return nil, err
	}
	return writer, nil
}

type CallerHook struct{}

func (hook *CallerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *CallerHook) Fire(entry *logrus.Entry) error {
	// 获取调用者信息
	_, file, line, _ := runtime.Caller(9) // 调整参数以适应你的调用层数
	parts := strings.Split(file, "/")
	callerInfo := fmt.Sprintf("%s:%d", parts[len(parts)-1], line)

	// 将调用者信息添加到日志字段中
	entry.Data["caller"] = callerInfo
	return nil
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

func (l *Logrus) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *Logrus) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *Logrus) Error(args ...interface{}) {
	l.log.Error(args...)
}

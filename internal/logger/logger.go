package logger

import (
	"io"
	"log"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"

	"github.com/chensylz/goredis/internal/global/constants"
)

type Log interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

var rLog Log

type Logrus struct {
	log *logrus.Logger
}

func SetupLog() {
	rLog = NewLogrus()
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

	return &Logrus{log: logureLog}
}

func newFileWriter(logFilePath string) (io.Writer, error) {
	writer, err := rotatelogs.New(
		logFilePath+".%Y%m%d%H%M",
		//rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		return nil, err
	}
	return writer, nil
}

type limitedWriter struct {
	file      *os.File
	limitSize int
	written   int
}

func (w *limitedWriter) Write(p []byte) (n int, err error) {
	if w.written+len(p) > w.limitSize {
		// 达到文件大小限制时，重新创建文件
		w.file.Close()
		newFile, err := os.OpenFile(w.file.Name(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return 0, err
		}
		w.file = newFile
		w.written = 0
	}

	n, err = w.file.Write(p)
	if err == nil {
		w.written += n
	}
	return
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

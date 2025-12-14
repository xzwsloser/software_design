package utils

import (
	"bytes"
	"fmt"
	"os"
	"io"
	"path"
	"github.com/sirupsen/logrus"
)

// 自定义 logrus 配置

const (
	RED    = 31
	YELLOW = 33 
	BLUE   = 36
	GRAY   = 37
)

type LogFormatter struct {

}

var (
	logger *logrus.Logger
)

func (t *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 日志信息颜色设置
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = GRAY
	case logrus.WarnLevel:
		levelColor = YELLOW
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = RED
	default:
		levelColor = BLUE
	}

	// 字节缓冲区设置
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// 自定义日志格式
	timestamp := entry.Time.Format("2006-01-02 15:04:06")
	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", 
					path.Base(entry.Caller.File), 
					entry.Caller.Line)
		// 自定义输出格式
		fmt.Fprintf(b,"[%s] \033[%dm[%s]\033[0m %s %s %s \n",timestamp,levelColor,entry.Level,fileVal,funcVal,entry.Message)
    } else {
        fmt.Fprintf(b,"[%s] \033[%dm[%s]\033[0m %s\n",timestamp,levelColor,entry.Level,entry.Message)
    }
	return b.Bytes(), nil
}

func InitLogger() {
	logger = logrus.New()
	file, _  := os.OpenFile("backend.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	logger.SetOutput(io.MultiWriter(file, os.Stdout))
	logger.SetReportCaller(true)
	logger.SetFormatter(&LogFormatter{})
	logger.SetLevel(logrus.DebugLevel)
}

func GetLogger() *logrus.Logger {
	return logger
}

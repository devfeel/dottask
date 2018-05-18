package task

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

const (
	defaultDateFormatForFileName = "2006_01_02"
	defaultDateLayout            = "2006-01-02"
	defaultFullTimeLayout        = "2006-01-02 15:04:05.999999"
)

type (
	Logger interface {
		Error(v ...interface{})
		Warn(v ...interface{})
		Info(v ...interface{})
		Debug(v ...interface{})
	}

	fileLogger struct {
		filePath string
	}
	fmtLogger struct{}
)

/************************* FmtLogger *******************************/

func NewFmtLogger() *fmtLogger {
	return &fmtLogger{}
}

func (logger *fmtLogger) Debug(v ...interface{}) {
	logger.writeLog(fmt.Sprint(v...), "debug")
}
func (logger *fmtLogger) Info(v ...interface{}) {
	logger.writeLog(fmt.Sprint(v...), "info")
}
func (logger *fmtLogger) Warn(v ...interface{}) {
	logger.writeLog(fmt.Sprint(v...), "warn")
}
func (logger *fmtLogger) Error(v ...interface{}) {
	logger.writeLog(fmt.Sprint(v...), "error")
}

func (logger *fmtLogger) writeLog(log string, level string) {
	logStr := time.Now().Format(defaultFullTimeLayout) + " " + log
	logStr += "\r\n"
	fmt.Println(logStr)
}

/************************* FileLogger *******************************/

func NewFileLogger(filePath string) *fileLogger {
	if filePath == "" || !fileExists(filePath) {
		filePath = getCurrentDirectory()
	}
	logger := &fileLogger{filePath: filePath}
	return logger
}

func (logger *fileLogger) Debug(v ...interface{}) {
	logger.writeLog(fmt.Sprint(v...), "debug")
}
func (logger *fileLogger) Info(v ...interface{}) {
	logger.writeLog(fmt.Sprint(v...), "info")
}
func (logger *fileLogger) Warn(v ...interface{}) {
	logger.writeLog(fmt.Sprint(v...), "warn")
}
func (logger *fileLogger) Error(v ...interface{}) {
	logger.writeLog(fmt.Sprint(v...), "error")
}

func (logger *fileLogger) writeLog(log string, level string) {
	filePath := logger.filePath + "dottask_" + level
	filePath += "_" + time.Now().Format(defaultDateFormatForFileName) + ".log"
	logStr := time.Now().Format(defaultFullTimeLayout) + " " + log
	logStr += "\r\n"
	var mode os.FileMode
	flag := syscall.O_RDWR | syscall.O_APPEND | syscall.O_CREAT
	mode = 0666
	file, err := os.OpenFile(filePath, flag, mode)
	defer file.Close()
	if err != nil {
		fmt.Println(filePath, err)
		return
	}
	file.WriteString(logStr)
}

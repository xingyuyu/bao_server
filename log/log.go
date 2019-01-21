package log

import "os"

//notice log
var noticeFd *os.File
var wfFd *os.File

type LogConfig struct {
	filePath string
	fileName string
}

type LogContext struct {
	uA       string
	clientIP string
	uri      string
}

func InitLog() error {

}

func CloseLog() error {

}

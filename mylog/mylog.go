package mylog

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os"
)

//notice log
var noticeFd *os.File
var wfFd *os.File

var noticeFile = "bao_server.log"
var wfFile = "bao_server.log.wf"

type LogContext struct {
	UA       string
	clientIP string
	uri      string
}

type BaoLog struct {
	logContext LogContext
	noticeMap  map[string]string
	warningMap map[string]string
	loger      *log.Logger
	buf        bytes.Buffer
	wfLogger   *log.Logger
	wfBuf      bytes.Buffer
}

func InitLog() error {
	var err1 error
	var err2 error

	noticeFd, err1 = os.OpenFile(noticeFile, os.O_RDWR|os.O_CREATE, 0755)
	if err1 != nil {
		return errors.New("notice file open fail")
	}
	wfFd, err2 = os.OpenFile(wfFile, os.O_RDWR|os.O_CREATE, 0755)
	if err2 != nil {
		return errors.New("wf file open fail")
	}
	return nil
}
func New() *BaoLog {
	baoLog := new(BaoLog)
	baoLog.loger = log.New(&baoLog.buf, "NOTICE", log.Lshortfile)
	baoLog.wfLogger = log.New(&baoLog.wfBuf, "WARNING", log.Lshortfile)
	return baoLog
}

func (BaoLog *BaoLog) SetContext(r *http.Request) {
	BaoLog.logContext.clientIP = r.Header.Get("clent_ip")
	BaoLog.logContext.UA = r.Header.Get("User-Agent")
	BaoLog.logContext.uri = r.URL.String()
	BaoLog.AddNotice("clent_ip", BaoLog.logContext.clientIP)
	BaoLog.AddNotice("UA", BaoLog.logContext.UA)
	BaoLog.AddNotice("uri", BaoLog.logContext.uri)
}

func CloseLog() {
	if noticeFd != nil {
		noticeFd.Close()
	}
	if wfFd != nil {
		wfFd.Close()
	}
}

func (baoLog *BaoLog) Flush() {
	noticeFd.Write(baoLog.buf.Bytes())
	wfFd.WriteString(baoLog.buf.String())
}

func (baoLog *BaoLog) AddNotice(key string, value string) {
	baoLog.loger.Print("[%s:%s]", key, value)
}

func (baoLog *BaoLog) AddWarning(key string, value string) {
	baoLog.wfLogger.Print("[%s:%s]", key, value)
}

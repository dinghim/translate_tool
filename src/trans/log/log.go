package log

import (
	_log "log"
	"os"
	"sync"
)

const (
	LOG_FILE = 1 << iota
	LOG_PRINT
)

const (
	LOG_INFO  string = "[INFO]  "
	LOG_ERROR string = "[ERROR] "
)

type log struct {
	fhandle  *os.File
	logFile  *_log.Logger
	logPrint *_log.Logger
}

var errlog error

const (
	const_log_file string = "log.txt"
)

var instance *log
var once sync.Once

func getinstance() *log {
	once.Do(func() {
		instance = &log{}
		instance.fhandle, errlog = os.Create(const_log_file)
		if errlog != nil {
			panic(errlog)
		}
		instance.logFile = _log.New(instance.fhandle, "", _log.LstdFlags)
		instance.logPrint = _log.New(os.Stdout, "", _log.LstdFlags)
	})
	return instance
}

func (l *log) Close() {
	if l.fhandle != nil {
		l.fhandle.Close()
	}
}

func WriteLog(flag int, level string, v ...interface{}) {
	l := getinstance()
	if l.logFile != nil && flag&LOG_FILE != 0 {
		l.logFile.SetPrefix(level)
		l.logFile.Println(v...)
	}
	if l.logPrint != nil && flag&LOG_PRINT != 0 {
		l.logPrint.SetPrefix(level)
		l.logPrint.Println(v...)
	}
}

func CloseLog() {
	getinstance().Close()
}

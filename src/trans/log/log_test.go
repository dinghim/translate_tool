package log_test

import (
	"testing"
	"trans/log"
)

func Test_example(t *testing.T) {
	defer log.CloseLog()
	log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, "this is a test information")
}

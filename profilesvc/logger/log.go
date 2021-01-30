package logger

import (
	"github.com/go-kit/kit/log"
	"os"
	"time"
)

var Log log.Logger

func InitLog() {
	Log = log.NewLogfmtLogger(os.Stderr)
	Log = log.With(Log, "ts", time.Now().Format("2006-01-02 15:04:05"))
	Log = log.With(Log, "caller", log.DefaultCaller)
}

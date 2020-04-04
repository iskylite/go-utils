package main

import (
	"os"

	"github.com/iskylite/go-utils/ilog"
)

// Logger 全局日志解析器
var Logger *ilog.Logger

func init() {
	Logger = ilog.Default()
	logdir, _ := os.Getwd()
	Logger.SetLevel(ilog.INFO)
	Logger.SetLogDir(logdir)
	Logger.SetLogFile("example.log")
	Logger.SetRotateSize(1024 * 1024)
	Logger.SetStdoutPrint(false)
}

func main() {
	for i := 0; i < 201; i++ {
		if i < 100 {
			Logger.Debug("example: Test ilog")
		}
		if i < 200 {
			Logger.Info("example: Test ilog")
		}
		if i == 200 {
			Logger.Fatal("example: Test ilog")
		}
	}
}

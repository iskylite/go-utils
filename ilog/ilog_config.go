// Package ilog 自定义日志输出
// ilog_config.go 日志配置
package ilog

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/iskylite/go-utils/itools"
)

// 规范日志等级
const (
	DEBUG int = iota
	INFO
	WARN
	ERROR
	FATAL
	PANIC
)

// Config 日志配置结构体
type Config struct {
	Level       int       // 日志最低输出等级
	Stdout      io.Writer // 默认日志输出方向
	LogDir      string    // 日志输出到文件中时指定的目录
	LogFile     string    // 日志输出到文件中时指定的文件名
	Prefix      string    // 日志输出字段前缀
	StdoutPrint bool      // 是否执行标准输出
	// Rotate       bool      // 是否分段
	RotateSize   int64 // 日志文件分段依据 文件大小
	RotateBackup bool  // 日志文件分段保存后，是否保留过去的备份
	Mill         bool  // 时间输出时是否精确到毫秒
}

// DefaultConfig 返回默认的配置结构体
func DefaultConfig() Config {
	c := Config{
		Level:       DEBUG,
		Stdout:      os.Stdout,
		StdoutPrint: true,
		// Rotate:       false,
		Mill:         false,
		RotateBackup: true,
	}
	return c
}

// SetLevel 设定当前日志输出等级
func (l *Logger) SetLevel(level int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.Level = level
}

// SetStdout 设定默认输出方向
func (l *Logger) SetStdout(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.Stdout = w
}

// SetStdoutPrint 是否执行默认输出
func (l *Logger) SetStdoutPrint(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.StdoutPrint = enabled
}

// SetLogDir 设定日志输出目录
func (l *Logger) SetLogDir(dir string) error {
	// 传参目录为空
	if dir == "" {
		return errors.New("LogDir is empty")
	}
	// 目录不存在
	if !itools.Exist(dir) {
		if err := itools.Mkdir(dir); err != nil {
			return err
		}
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.LogDir = strings.TrimRight(dir, itools.Sep)
	return nil
}

// SetLogFile 设定日志输出文件名字
func (l *Logger) SetLogFile(fname string) error {
	if fname == "" {
		return errors.New("LogFile is empty")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.LogFile = fname
	return nil
}

// SetPrefix 设定日志输出前缀 用来做标识
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.Prefix = prefix
}

// SetRotateSize 设定日志文件分段的依据  文件大小
func (l *Logger) SetRotateSize(size int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.RotateSize = size
}

// SetRotateBackup 日志分段后,是否保留过去的日志
func (l *Logger) SetRotateBackup(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.RotateBackup = enabled
}

// SetMill 设定日志输出的日期时间是否精确到毫秒
func (l *Logger) SetMill(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.Mill = enabled
}

// Package ilog 自定义日志
// 2020-04-01 iskylite
// 将标准包log.go进行修改，并参考了goframe（https://github.com/gogf/gf）内的日志库glog的实现，
// 通过精简不必要的一些配置完成的
// 仅仅字面上支持debug、info、warn、error、fatal和panic
// 且支持等级过滤
// 默认仅仅支持os.Stdout作为日志输出方向，通过配置可还以支持文件输出
// 文件切割暂不支持
// 日志输出格式：（无颜色支持）
// [2020-04-01T22:24:30] INFO ilog.go:6 => 具体信息
// 且暂时不支持json输出
package ilog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/iskylite/go-utils/itools"
)

// Logger 日志管理结构体
type Logger struct {
	mu     sync.Mutex
	buf    []byte
	config Config
}

// Default 工厂函数返回默认值 日志等级DEBUG 输出方向标准输出os.Stdout
func Default() *Logger {
	return &Logger{
		config: DefaultConfig(),
	}
}

// New 使用新的配置生成Logger
func New(config Config) *Logger {
	fmt.Println("new start...")
	return &Logger{
		config: config,
	}
}

// SetConfig 更改配置
func (l *Logger) SetConfig(c Config) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config = c
}

// GetFilepath 获取当前日志输出完整路径
func (l *Logger) GetFilepath() string {
	return filepath.Join(l.config.LogDir, l.config.LogFile)
}

// Itoa Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// getLevel 输出日志输出级别的字符串
func getLevel(lv int) string {
	level := map[int]string{
		0: "DEBUG",
		1: "INFO",
		2: "WARN",
		3: "ERROR",
		4: "FATAL",
		5: "PANIC",
	}

	return level[lv]
}

// formatHeader 按照我自己的方式格式化
func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int, level int) {
	// 前缀
	if l.config.Prefix != "" {
		*buf = append(*buf, l.config.Prefix...)
	}
	*buf = append(*buf, '[')
	// 日期
	t = t.UTC()
	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '-')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '-')
	itoa(buf, day, 2)
	*buf = append(*buf, 'T')

	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)
	if l.config.Mill {
		*buf = append(*buf, '.')
		itoa(buf, t.Nanosecond()/1e3, 6)
	}
	*buf = append(*buf, "] "...)

	*buf = append(*buf, fmt.Sprintf("%s ", getLevel(level))...)

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}

	file = short
	*buf = append(*buf, file...)
	*buf = append(*buf, ':')
	itoa(buf, line, -1)
	*buf = append(*buf, " => "...)
}

// Output 输出日志到stdout和文件中
func (l *Logger) Output(calldepth, level int, s string) error {
	fmt.Println("output start...")
	// 判断当前输出的日志等级是否满足Logger中定义的Level
	if level < l.config.Level {
		return nil
	}

	now := time.Now() // get this early.
	var file string
	var line int
	var ok bool
	_, file, line, ok = runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, file, line, level)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}

	// 输出
	if l.config.StdoutPrint {
		l.WriteToStdout()
	}
	if l.config.LogDir != "" && l.config.LogFile != "" {
		l.WriteToFile(now)
	}
	return nil
}

// WriteToStdout 向标准输出中写入日志
func (l *Logger) WriteToStdout() {
	fmt.Println("writetostdout start...")
	_, err := l.config.Stdout.Write(l.buf)
	if err != nil {
		panic(err)
	}
}

// WriteToFile 向文件中写入日志信息
func (l *Logger) WriteToFile(now time.Time) {
	fmt.Println("writetofile start...")
	logfile := l.GetFilepath()
	file := itools.OpenFile(logfile)
	defer file.Close()

	if l.config.RotateSize > 0 {
		state, err := file.Stat()
		if err != nil {
			panic(err)
		}

		if state.Size() >= l.config.RotateSize {
			file.Close()
			l.RotateFileBySize(now)
			file = itools.OpenFile(logfile)
			defer file.Close()
		}

	}
	_, err := file.Write(l.buf)
	if err != nil {
		panic(err)
	}
}

// RotateFileBySize 根据文件大小来进行日志分段
func (l *Logger) RotateFileBySize(now time.Time) {
	fmt.Println("rotatebysize start...")
	if !l.config.RotateBackup {
		fmt.Println("remove start...")
		err := os.RemoveAll(l.GetFilepath())
		if err != nil {
			panic(err)
		}
	} else {
		var logfilepath string
		for {
			time := itools.FmtTimeWithMillOut(now)
			if strings.HasSuffix(l.GetFilepath(), "log") {
				logfilepath = fmt.Sprintf("%s%s.log", strings.TrimSuffix(l.GetFilepath(), "log"), time)
			} else {
				logfilepath = fmt.Sprintf("%s%s", l.GetFilepath(), time)
			}
			// logfilepath = fmt.Sprintf("%s%s", l.GetFilepath(), time)
			if !itools.Exist(logfilepath) {
				break
			}
		}
		fmt.Println("rename start...")
		err := os.Rename(l.GetFilepath(), logfilepath)
		if err != nil {
			panic(err)
		}
	}
}

// Debug 输出debug级别的日志
func (l *Logger) Debug(v ...interface{}) {
	fmt.Println("debug start...")
	l.Output(2, 0, fmt.Sprint(v...))
}

// Debugf 输出debug级别的日志
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Output(2, 0, fmt.Sprintf(format, v...))
}

// Panic 输出Panic级别的日志
func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Output(2, 5, s)
	panic(s)
}

// Panicf 输出Panic级别的日志
func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output(2, 5, s)
	panic(s)
}

// Info 输出Info级别的日志
func (l *Logger) Info(v ...interface{}) {
	l.Output(2, 1, fmt.Sprint(v...))
}

// Infof 输出Info级别的日志
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Output(2, 1, fmt.Sprintf(format, v...))
}

// Warn 输出Warn级别的日志
func (l *Logger) Warn(v ...interface{}) {
	l.Output(2, 2, fmt.Sprint(v...))
}

// Warnf 输出Warn级别的日志
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Output(2, 2, fmt.Sprintf(format, v...))
}

// Error 输出Error级别的日志
func (l *Logger) Error(v ...interface{}) {
	l.Output(2, 3, fmt.Sprint(v...))
}

// Errorf 输出Error级别的日志
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Output(2, 3, fmt.Sprintf(format, v...))
}

// Fatal 输出Fatal级别的日志
func (l *Logger) Fatal(v ...interface{}) {
	l.Output(2, 4, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf 输出Fatalf级别的日志
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(2, 4, fmt.Sprintf(format, v...))
	os.Exit(1)
}

package itools

import (
	"os"
	"path/filepath"
)

const (
	// Sep 当前操作系统分割符
	Sep = string(filepath.Separator)
)

// Exist 判断路径是否存在
func Exist(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

// Mkdir 创建目录
func Mkdir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// Create 创建文件
func Create(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if !Exist(dir) {
		if err := Mkdir(dir); err != nil {
			return nil, err
		}
	}
	return os.Create(path)
}

// OpenFile 包装标准库 os.OpenFile
func OpenFile(fname string) *os.File {
	file, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return file
}

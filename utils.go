package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func exists(path string) bool {
	_, err := os.Stat(path) // os.Stat 获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func isFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// 判断所给路径是否为文件
func getFileMode(path string) os.FileMode {
	s, err := os.Stat(path)
	if err != nil {
		return 0766
	}
	return s.Mode()
}

func isLink(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.Mode()&os.ModeSymlink != 0
}

func printError(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func copyFile(from, to string, mode os.FileMode) (err error) {
	src, err := os.Open(from)
	if err != nil {
		return err
	}
	defer func() {
		_ = src.Close()
	}()

	dest, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE, mode)
	if err != nil {
		return err
	}
	defer func() {
		_ = dest.Close()
	}()

	_, err = io.Copy(dest, src)
	return err
}

func copyDir(from string, to string) error {
	return filepath.Walk(from, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(from, p)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return os.MkdirAll(path.Join(to, rel), info.Mode())
		} else {
			return copyFile(p, path.Join(to, rel), info.Mode())
		}
	})
}

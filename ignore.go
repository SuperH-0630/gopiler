package main

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type IgnoreFileType int8

const (
	Dir IgnoreFileType = iota
	File
	Prefix
	Suffix
	Substring
	Equal
)

type IgnoreFile struct {
	ft   IgnoreFileType
	path string
}

var ignoreFile []IgnoreFile

func init() {
	ignoreFile = make([]IgnoreFile, 100)
	ignoreFile = append(ignoreFile, IgnoreFile{ft: Equal, path: "__pycache__"})
}

func getIgnoreFile() (err error) {
	gopilerignore := path.Join(projectPath, ".gopilerignore")
	if !isFile(gopilerignore) {
		return
	}

	file, err := os.Open(gopilerignore)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		f := strings.Trim(fileScanner.Text(), " ")
		if len(f) == 0 {
			continue
		}

		if strings.HasPrefix(f, "=") {
			ignoreFile = append(ignoreFile, IgnoreFile{ft: Equal, path: f[1:]})
		} else if strings.HasPrefix(f, "*") && strings.HasPrefix(f, "*") {
			ignoreFile = append(ignoreFile, IgnoreFile{ft: Substring, path: f[1 : len(f)-1]})
		} else if strings.HasPrefix(f, "*") {
			ignoreFile = append(ignoreFile, IgnoreFile{ft: Suffix, path: f[1:]})
		} else if strings.HasSuffix(f, "*") {
			ignoreFile = append(ignoreFile, IgnoreFile{ft: Prefix, path: f[:len(f)-1]})
		} else {
			if filepath.IsAbs(f) && exists(f) {
				if isDir(f) {
					ignoreFile = append(ignoreFile, IgnoreFile{ft: Dir, path: f})
				} else {
					ignoreFile = append(ignoreFile, IgnoreFile{ft: File, path: f})
				}
			} else if exists(path.Join(projectPath, f)) {
				f, err = filepath.Abs(path.Join(projectPath, f))
				if err != nil {
					continue
				}

				if isDir(f) {
					ignoreFile = append(ignoreFile, IgnoreFile{ft: Dir, path: f})
				} else {
					ignoreFile = append(ignoreFile, IgnoreFile{ft: File, path: f})
				}
			}
		}
	}

	return nil
}

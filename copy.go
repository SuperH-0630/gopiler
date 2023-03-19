package main

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func copyAndCompilerFile() (err error) {
	err = os.MkdirAll(path.Join(sharePath, "src"), os.ModePerm)
	if err != nil {
		return err
	}

	return filepath.Walk(projectPath, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if isHidden(p) {
			if info.IsDir() {
				ignoreFile = append(ignoreFile, IgnoreFile{ft: Dir, path: p})
			}
			return nil // 忽略此文件
		}

		for _, f := range ignoreFile {
			if len(f.path) == 0 {
				continue
			}

			if f.ft == Equal {
				if Base(p) == f.path {
					if info.IsDir() {
						ignoreFile = append(ignoreFile, IgnoreFile{ft: Dir, path: p})
					}

					return nil // 忽略此文件
				}
			} else if f.ft == Substring {
				if strings.Contains(Base(p), f.path) {
					if info.IsDir() {
						ignoreFile = append(ignoreFile, IgnoreFile{ft: Dir, path: p})
					}

					return nil // 忽略此文件
				}
			} else if f.ft == Prefix {
				if strings.HasPrefix(Base(p), f.path) {
					if info.IsDir() {
						ignoreFile = append(ignoreFile, IgnoreFile{ft: Dir, path: p})
					}

					return nil // 忽略此文件
				}
			} else if f.ft == Suffix {
				if strings.HasSuffix(Base(p), f.path) {
					if info.IsDir() {
						ignoreFile = append(ignoreFile, IgnoreFile{ft: Dir, path: p})
					}

					return nil // 忽略此文件
				}
			} else if f.ft == Dir {
				rel, err := filepath.Rel(f.path, p)
				if err != nil {
					return err
				}

				if !strings.HasPrefix(rel, "..") {
					// .. 开头表示path不是f的子目录或不是f
					return nil // 忽略此文件
				}
			} else {
				if p == f.path {
					return nil // 忽略此文件
				}
			}
		}

		if info.IsDir() {
			return nil // 接下来的操作跳过目录
		}

		rel, err := filepath.Rel(projectPath, p)
		if err != nil {
			return err
		}

		to, err := filepath.Abs(path.Join(sharePath, "src", rel)) // 目标文件夹
		if err != nil {
			return err
		}

		if strings.HasSuffix(Base(p), ".py") {
			return compilerFile(p, to+"c") // 编译成pyc
		} else {
			return copyFile(p, to, info.Mode())
		}
	})
}

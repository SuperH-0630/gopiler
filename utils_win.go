//go:build windows

package main

import (
	"strings"
	"syscall"
)

func getFileAttributes(filePath string) (uint32, error) {
	filePtr, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return 0, err
	}
	fileInfo, err := syscall.GetFileAttributes(filePtr)
	if err != nil {
		return 0, err
	}
	return fileInfo, nil
}

func isHidden(p string) bool {
	if strings.HasPrefix(Base(p), ".") {
		return true
	}

	fileInfo, err := getFileAttributes(p)
	if err != nil {
		return true
	}

	return fileInfo&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}

/* Base
 * 功能类似于path.Base，但可以处理\分隔符
 */
func Base(p string) string {
	if p == "" {
		return "."
	}
	// Strip trailing slashes.
	for len(p) > 0 && (p[len(p)-1] == '/' || p[len(p)-1] == '\\') {
		p = p[0 : len(p)-1]
	}
	// Find the last element
	if i := lastSlash(p); i >= 0 {
		p = p[i+1:]
	}
	// If empty now, it had only slashes.
	if p == "" {
		return "/"
	}
	return p
}

func lastSlash(s string) int {
	i := len(s) - 1
	for i >= 0 && (s[i] != '/' && s[i] != '\\') {
		i--
	}
	return i
}

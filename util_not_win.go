//go:build !windows

package main

import (
	"path"
	"strings"
)

func isHidden(p string) bool {
	return strings.HasPrefix(path.Base(p), ".")
}

func Base(p string) string {
	return path.Base(p)
}

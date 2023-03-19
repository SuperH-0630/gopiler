package main

import (
	"os"
	"os/exec"
)

func runEditBin(args, exe string) bool {
	cmd := exec.Command(editbin, args, exe)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}

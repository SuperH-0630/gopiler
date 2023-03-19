package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func getGoVersion(goPath string) (x, y, z int) {
	cmd := exec.Command(goPath, "version")
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, 0
	}

	n, err := fmt.Sscanf(string(out), "go version go%d.%d.%d", &x, &y, &z)
	if err != nil || n != 3 {
		return 0, 0, 0
	}

	return x, y, z
}

func createStartCode() (string, error) {
	codeFormat := `package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

const ProjectName = "%s"

func main() {
	executableFile, err := os.Executable()
	if err != nil {
		os.Exit(1)
	}

	binPath := executableFile[:lastSlash(executableFile)]
	sharePath, err := filepath.Abs(path.Join(binPath[:lastSlash(binPath)], "share", ProjectName))
	if err != nil {
		os.Exit(1)
	}

	venvPath := path.Join(sharePath, "venv")
	if !isDir(venvPath) {
		os.Exit(1)
	}
	
	var pythonPath string
	if runtime.GOOS == "windows" {
		pythonPath = path.Join(venvPath, "Scripts", "python.exe")
	} else {
		pythonPath = path.Join(venvPath, "bin", "python")
	}
	if !isFile(pythonPath) {
		os.Exit(1)
	}

	srcPath := path.Join(sharePath, "src")
	if !isDir(srcPath) {
		os.Exit(1)
	}

	args := []string{%s}
	if len(os.Args) > 1 {
		for _, a := range os.Args[1:] {
			args = append(args, a)
		}
	}

	for i, a := range args { // 装饰
		if strings.HasPrefix(a, "<gopiler:src>") {
			args[i] = path.Join(srcPath, a[len("<gopiler:src>"):])
		} else if strings.HasPrefix(a, "<gopiler:venv>") {
			args[i] = path.Join(venvPath, a[len("<gopiler:src>"):])
		} else if strings.HasSuffix(a, ".py") && isFile(a+"c") { // 转换为pyc文件
			args[i] += "c"
		}
	}

	pwd, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	cmd := exec.Command(pythonPath, args...)
	if %s {
		cmd.Stdin = nil
		cmd.Stdout = os.NewFile(uintptr(syscall.Stdin), os.DevNull)
		cmd.Stderr = os.NewFile(uintptr(syscall.Stderr), os.DevNull)

	} else {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
	}
	cmd.Dir = pwd

	err = cmd.Run()
	if err != nil {
		if errExit, ok := err.(*exec.ExitError); ok {
			if status, ok := errExit.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		os.Exit(1)
	}

	os.Exit(0)
}

func lastSlash(s string) int {
	i := len(s) - 1
	if runtime.GOOS == "windows" {
		for i >= 0 && (s[i] != '/' && s[i] != '\\') {
			i--
		}
	} else {
		for i >= 0 && (s[i] != '/') {
			i--
		}
	}
	return i
}

func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func isFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

`
	args := make([]string, 0, 10)
	for _, i := range flag.Args() {
		args = append(args, fmt.Sprintf("\"%s\"", strings.Replace(i, "\"", "\\\"", -1)))
	}

	if len(args) == 0 {
		args = []string{"<gopiler:src>main.pyc"}
	}

	argsString := strings.Join(args, ", ")

	var code string
	if xwindows {
		code = fmt.Sprintf(codeFormat, projectName, argsString, "true")
	} else {
		code = fmt.Sprintf(codeFormat, projectName, argsString, "false")
	}

	codePath := path.Join(tmpPath, fmt.Sprintf("%s_gopiler.go", projectName))
	err := os.WriteFile(codePath, []byte(code), 0666)
	if err != nil {
		return "", err
	}

	return codePath, nil
}

func buildGoCode(from, to string, xwindows bool) error {
	var cmd *exec.Cmd
	if xwindows {
		cmd = exec.Command(goPath, "build", "-o", to, `-ldflags=-H windowsgui`, from)
	} else {
		cmd = exec.Command(goPath, "build", "-o", to, from)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%e", err)
	}

	return nil
}

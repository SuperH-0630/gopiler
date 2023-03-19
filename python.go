package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func getPythonVersion(pythonPath string) (x, y, z int) {
	cmd := exec.Command(pythonPath, "--version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, 0
	}

	n, err := fmt.Sscanf(string(out), "Python %d.%d.%d", &x, &y, &z)
	if err != nil || n != 3 {
		return 0, 0, 0
	}

	return x, y, z
}

func testPythonPip() (bool, error) {
	return testPythonModel("pip")
}

func testPythonVenv() (bool, error) {
	return testPythonModel("venv")
}

func testPythonModel(model string) (bool, error) {
	codeFormat := `# -*- coding: utf-8 -*-
try:
    import %s
except:
    print(1)
else:
    print(0)
`
	testVenv := path.Join(tmpPath, fmt.Sprintf("test_%s.py", model))
	err := os.WriteFile(testVenv, []byte(fmt.Sprintf(codeFormat, model)), 0666)
	if err != nil {
		return false, err
	}

	cmd := exec.Command(pythonPath, testVenv)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("%e", err)
	}

	if strings.HasPrefix(string(out), "1") {
		return false, nil
	}

	return true, nil
}

func installPythonVenv() error {
	return installPythonModel("venv")
}

func installPythonModel(model string) error {
	cmd := exec.Command(pythonPath, "-m", "pip", "install", model)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%e", err)
	}

	return nil
}

func createVenv() (python string, err error) {
	venv := path.Join(sharePath, "venv")

	if runtime.GOOS == "windows" {
		python = path.Join(venv, "Scripts", "python.exe")
	} else {
		python = path.Join(venv, "bin", "python")
	}

	if exists(venv) {
		if isFile(python) {
			return python, nil
		}

		err = os.RemoveAll(venv)
		if err != nil {
			return "", err
		}
	}

	cmd := exec.Command(pythonPath, "-m", "venv", venv)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%e", err)
	}

	if !isFile(python) {
		return "", fmt.Errorf("fail to find python")
	}

	return python, nil
}

func checkVenvFile() (err error) {
	type File struct {
		path string
		mode os.FileMode
	}
	link := make([]File, 0, 100)

	err = filepath.Walk(path.Join(sharePath, "venv"), func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		for _, l := range link {
			if strings.HasPrefix(p, l.path) {
				return nil // 跳过这个文件
			}
		}

		if isLink(p) {
			link = append(link, File{path: p, mode: info.Mode()})
		}

		return nil
	})
	if err != nil {
		return err
	}

	for _, l := range link {
		var err error

		if len(l.path) == 0 {
			continue
		}

		target, err := os.Readlink(l.path)
		if err != nil {
			return err
		}

		err = os.RemoveAll(l.path)
		if err != nil {
			return err
		}

		if isDir(target) {
			return copyDir(target, l.path)
		} else {
			return copyFile(target, l.path, l.mode)
		}
	}

	return nil
}

func createCompiler() error {
	code := `# -*- coding: utf-8 -*-
import sys
import py_compile
py_compile.compile(sys.argv[1], sys.argv[2])
`

	compiler := path.Join(tmpPath, "compiler.py")
	err := os.WriteFile(compiler, []byte(code), 0666)
	if err != nil {
		return err
	}

	return nil
}

func compilerFile(from, to string) error {
	compiler := path.Join(tmpPath, "compiler.py")

	cmd := exec.Command(vpythonPath, compiler, from, to)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%e", err)
	}

	return nil
}

func installPythonRequirements() error {
	requirements := path.Join(projectPath, "requirements.txt")
	if !isFile(requirements) {
		return nil
	}

	cmd := exec.Command(vpythonPath, "-m", "pip", "install", "-i", pip, "-r", requirements)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%e", err)
	}

	return nil
}

/*
	Gopiler 工作原理

1. 找到Python可执行程序
2. 验证virtualenv模块是否存在，不存在则安装
3. 在工作目录/share/{project-name}/venv创建虚拟环境
5. 读取.gopilerignore文件，获取忽略文件
4. 复制文件到工作目录/share/{project-name}/src，.py编译为pyz，.开头目录和隐藏目录不复制
6. 创建可执行文件编译源码到/tmp/{project-name}.go目录，windows包括资源文件rc
5. 创建可执行程序到/bin/{project-name} （windows支持自定义图标）
*/
package main

import (
	"flag"
	"fmt"
	"github.com/flytam/filenamify"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

const VERSION string = "0.0.1"

var (
	help        bool
	pythonPath  string
	vpythonPath string
	goPath      string
	workPath    string
	projectPath string
	projectName string
	xwindows    bool
	editbin     string
	pip         string
)

var (
	tmpPath   string
	sharePath string
	binPath   string
)

func init() {
	flag.BoolVar(&help, "help", false, "Print for the help.")
	flag.StringVar(&pythonPath, "python", "python", "The location of the python executable.")
	flag.StringVar(&goPath, "go", "go", "The location of the golang executable.")
	flag.StringVar(&workPath, "work", "", "The work path.")
	flag.StringVar(&projectPath, "project", "", "The project path.")
	flag.StringVar(&projectName, "pname", "", "The project name.")
	flag.BoolVar(&xwindows, "xwindows", false, "As an x-windows project.")
	flag.StringVar(&editbin, "editbin", "editbin.exe", "The editbin.exe path.")
	flag.StringVar(&pip, "pip", "https://pypi.python.org/simple", "The pip source.")
	flag.Usage = usage
}

func main() {
	var err error
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if len(projectPath) == 0 {
		printError("Please set the project path by --project\n")
	}

	if !isDir(projectPath) {
		printError("Bad project (only a file)\n")
	}

	if len(projectName) == 0 {
		printError("Please set the project name by --pname\n")
	}

	projectName, err = filenamify.Filenamify(projectName, filenamify.Options{
		Replacement: "?",
	}) // 转换为合法文件名

	if len(workPath) != 0 {
		if !filepath.IsAbs(workPath) { // 相对路径转换为绝对路径
			workPath, err = filepath.Abs(workPath)
			if err != nil {
				printError("Get work path fail: %e", err)
			}
		}
	} else {
		workPath, err = os.Getwd()
		if err != nil {
			printError("Get work path fail: %e", err)
		}
		workPath = path.Join(workPath, projectName)
	}

	if !exists(workPath) {
		err = os.MkdirAll(workPath, os.ModePerm)
		if err != nil {
			printError("Create work path fail: %e", err)
		}
	} else if isFile(workPath) {
		printError("The work path is a file")
	}

	tmpPath = path.Join(workPath, "tmp")
	sharePath = path.Join(workPath, "share", projectName)
	binPath = path.Join(workPath, "bin")

	err = os.MkdirAll(tmpPath, os.ModePerm)
	if err != nil {
		printError("Create work tmp path fail: %e", err)
	}

	err = os.MkdirAll(sharePath, os.ModePerm)
	if err != nil {
		printError("Create work share path fail: %e", err)
	}

	err = os.MkdirAll(binPath, os.ModePerm)
	if err != nil {
		printError("Create work bin path fail: %e", err)
	}

	fmt.Printf("Work Path: %s\n", workPath)
	fmt.Printf("Work Temp Path: %s\n", tmpPath)
	fmt.Printf("Work Share Path: %s\n", sharePath)
	fmt.Printf("Work Bin Path: %s\n", binPath)

	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(pythonPath, ".exe") {
			pythonPath = pythonPath + ".exe"
		}

		if !strings.HasSuffix(goPath, ".exe") {
			goPath = goPath + ".exe"
		}
	}

	var x, y, z int

	x, y, z = getPythonVersion(pythonPath)
	if x == 0 {
		printError("python not found")
	}

	fmt.Printf("Python Path: %s\n", pythonPath)
	fmt.Printf("Python Version: %d.%d.%d\n", x, y, z)

	x, y, z = getGoVersion(goPath)
	if x == 0 {
		printError("go not found")
	}

	fmt.Printf("Golang Path: %s\n", goPath)
	fmt.Printf("Golang Version: %d.%d.%d\n", x, y, z)

	hasPip, err := testPythonPip()
	if err != nil {
		printError("Test python (pip) error: %e\n", err)
	}

	hasVenv, err := testPythonVenv()
	if err != nil {
		printError("Test python (venv) error: %e\n", err)
	}

	if !hasVenv {
		if !hasPip {
			printError("Please install pip first\n")
		}

		err = installPythonVenv()
		if err != nil {
			printError("Install venv in python fail: %e\n", err)
		}

		hasVenv, err = testPythonVenv()
		if !hasVenv {
			printError("Check after install venv in python fail: %e\n", err)
		}
	}

	vpythonPath, err = createVenv()
	if err != nil {
		printError("Create venv fail: %e\n", err)
	}

	fmt.Printf("Venv python: %s\n", vpythonPath)

	err = checkVenvFile()
	if err != nil {
		printError("Check venv file error: %e\n", err)
	}

	err = installPythonRequirements()
	if err != nil {
		printError("Install venv requirements pack error: %e\n", err)
	}

	err = getIgnoreFile()
	if err != nil {
		printError("Get ignore file error: %e\n", err)
	}

	err = createCompiler()
	if err != nil {
		printError("Create compiler fail: %e\n", err)
	}

	err = copyAndCompilerFile()
	if err != nil {
		printError("Copy and compile file fail: %e\n", err)
	}

	goCode, err := createStartCode()
	if err != nil {
		printError("Create start code file fail: %e\n", err)
	}

	fmt.Printf("Start Code: %s\n", goCode)

	target := path.Join(binPath, projectName)
	if runtime.GOOS == "windows" {
		target += ".exe"
	}

	err = buildGoCode(goCode, target, xwindows)
	if err != nil {
		printError("Build go code error: %e\n", err)
	}

	fmt.Printf("Start Target: %s\n", target)

	if runtime.GOOS == "windows" && xwindows {
		success := runEditBin("/subsystem:windows", vpythonPath)
		if !success {
			fmt.Printf("Edit %s as xwindows project fail\n", vpythonPath)
		} else {
			fmt.Printf("Edit %s as xwindows project success\n", vpythonPath)
		}
	}

	fmt.Printf("Finish!\n")
}

func usage() {
	var err error
	_, err = fmt.Fprintf(os.Stderr, "Gopiler %s\n", VERSION)
	if err != nil {
		panic(err)
	}

	flag.PrintDefaults()
}

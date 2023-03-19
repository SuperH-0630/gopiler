# Gopiler
## 介绍
Gopiler是一款Python打包工具，基于golang。其原理是：生成python字节码和虚拟环境，并创建一个启动他们的可执行程序。

## 参数介绍
1. `--help` 查看帮助 【bool参数】
2. `--python=xxx` python可执行程序路径
3. `--go=xxx` golang编译器路径
4. `--work=xxx` 工作目录，内容会输出在该目录下项目名称的文件夹
5. `--project=xxx` 工程文件路径
6. `--pname=xxx` 项目名称
7. `--xwindows` 是否启用GUI模式（默认无命令行）【bool参数】
8. `--editbin` 在windows下，启用xwindows模式最好指定editbin.exe的位置，位于visual studio文件夹内
9. `--pip` pip源
10. `--ico` 可执行程序图标，windows下必须设置，同时项目文件必须配备versioninfo.json

其余参数会作为python的启动参数，传给启动器。
这些参数支持：
1. `.py`转换为`.pyc`，前提是该文件存在
2. `<gopiler:src>`开头表示相对路径，相对于源码文件夹
3. `<gopiler:venv>`开头表示相对路径，相对于虚拟环境文件夹

## 项目文件
项目文件中的`.py`会全部编译为`.pyc`文件，拷贝到`src`目录。
其余文件会原封不动的拷贝。
除了：
1. `.`开头的文件和文件夹
2. windows下的隐藏文件和文件夹
3. `.gopilerignore`忽略的文件和文件夹

## .gopilerignore文件
1. 可列出绝对路径文件和文件夹
2. 相对于工程目录的文件和文件夹
3. =开头表示文件名匹配
4. *开头表示前缀匹配，例如`*ab`会忽略文件`cab`
5. *结尾表示后缀匹配，例如`ab*`将会忽略文件`abb`
6. *开头和结尾表示字串匹配，例如`*ab*`将会忽略文件`ab`，`cab`，`abb`

## versioninfo.json
必须设置在项目文件夹下。
一个可行的版本是：
```json
{
    "FixedFileInfo": {
        "FileVersion": {
            "Major": 1,
            "Minor": 0,
            "Patch": 0,
            "Build": 0
        },
        "ProductVersion": {
            "Major": 1,
            "Minor": 0,
            "Patch": 0,
            "Build": 0
        },
        "FileFlagsMask": "3f",
        "FileFlags ": "00",
        "FileOS": "040004",
        "FileType": "01",
        "FileSubType": "00"
    },
    "StringFileInfo": {
        "Comments": "",
        "CompanyName": "",
        "FileDescription": "",
        "FileVersion": "",
        "InternalName": "",
        "LegalCopyright": "",
        "LegalTrademarks": "",
        "OriginalFilename": "",
        "PrivateBuild": "",
        "ProductName": "",
        "ProductVersion": "v1.0.0.0",
        "SpecialBuild": ""
    },
    "VarFileInfo": {
        "Translation": {
            "LangID": "0409",
            "CharsetID": "04B0"
        }
    }
}
```
具体的含义可参考：[learn-microsoft](https://learn.microsoft.com/zh-cn/windows/win32/menurc/versioninfo-resource?redirectedfrom=MSDN)
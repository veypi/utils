// Package utils 是一套自己构建的
//
// O功能
//
// - 随机构建字符串和Byte序列
//
// -
package utils

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	Version = "v0.2.2"
)

func CallPath(s int) string {
	_, f, l, _ := runtime.Caller(s + 1)
	return f + ":" + strconv.Itoa(l)
}

func PathJoin(paths ...string) string {
	return filepath.Join(paths...)
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//检测文件夹路径时候存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Abs(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Abs(filepath.Join(usr.HomeDir, path[1:]))
}

func MkFile(dest string) (*os.File, error) {
	if temp, err := Abs(dest); err == nil {
		dest = temp
	}
	//分割path目录
	destSplitPathDirs := strings.Split(dest, string(filepath.Separator))
	//检测时候存在目录
	destSplitPath := ""
	for _, dir := range destSplitPathDirs[:len(destSplitPathDirs)-1] {
		destSplitPath = destSplitPath + dir + string(filepath.Separator)
		b, _ := PathExists(destSplitPath)
		if !b {
			//创建目录
			_ = os.Mkdir(destSplitPath, 0755)
		}
	}
	// 覆写模式
	return os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
}

// CopyFile 拷贝源文件到指定位置
// 若目标文件夹不存在 则新建文件夹
//
func CopyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()
	dstFile, err := MkFile(dest)
	if err != nil {
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

// GetRunnerPath 获取运行程序所在的绝对路径
func GetRunnerPath() string {
	if path, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		return path
	}
	return os.Args[0]
}

// IsWindows 判断运行平台是否为Windows
func IsWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

// ChMod 修改文件权限 仅对类Unix 系统有效。
// name 文件路径， mode 修改的模式
func ChMod(name string, mode os.FileMode) {
	if !IsWindows() {
		os.Chmod(name, mode)
	}
}
func Exec(acts ...string) (string, error) {
	if len(acts) == 0 {
		return "", nil
	}

	//First argv must be executable,second must be argv,no space in it
	cmd := exec.Command(acts[0], acts[1:]...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

/*
Home 获取当前用户目录。

Windows系统:

1. 找到 HOMEDRIVE 和 HOMEPATH 变量，拼接为用户路径。

2. 找到 USERPROFILE 变量。

3. 都不存在时，抛出 error

类 Unix 系统:

1. 找到 HOME 变量

2. 利用 命令行 "sh -c eval echo ~$USER" 获取

3. 抛出异常
*/
func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// 不同操作系统使用不同方式获取
	if IsWindows() {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

// InList 判断列表是否含有某元素
func InList(str string, list []string) bool {
	for _, temp := range list {
		if str == temp {
			return true
		}
	}
	return false
}

// RemoveRep 通过map主键唯一的特性过滤重复元素
func RemoveRep(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

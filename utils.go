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
	Version = "v0.4.0"
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

func PathIsDir(p string) bool {
	s, err := os.Stat(p)
	if err != nil {
		return false
	}
	return s.IsDir()
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

//生成目录并拷贝文件
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

//Get the absolute path to the running directory
func GetRunnerPath() string {
	if path, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		return path
	}
	return os.Args[0]
}

//Determine whether the current system is a Windows system?
func IsWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

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

func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
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

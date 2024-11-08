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
	"unicode"
)

const (
	Version = "v0.4.2"
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

// 检测文件夹路径时候存在
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

// 生成目录并拷贝文件
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

// Get the absolute path to the running directory
func GetRunnerPath() string {
	if path, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		return path
	}
	return os.Args[0]
}

// Determine whether the current system is a Windows system?
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

func ToTitle(str string) string {
	if str == "" {
		return str
	}
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

// ToLowerFirst 将字符串的第一个字母转换为小写
func ToLowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	firstRune := unicode.ToLower(rune(s[0]))
	return string(firstRune) + s[1:]
}

func SnakeToPrivateCamel(input string) string {
	parts := strings.Split(input, "_")
	for i := 1; i < len(parts); i++ {
		parts[i] = ToTitle(parts[i])
	}
	parts[0] = strings.ToLower(parts[0])
	return strings.Join(parts, "")
}
func SnakeToCamel(input string) string {
	parts := strings.Split(input, "_")
	for i := 0; i < len(parts); i++ {
		parts[i] = ToTitle(parts[i])
	}
	return strings.Join(parts, "")
}

// CamelToSnake 将驼峰命名法转换为下划线命名法
// 例如：CamelToSnake("CamelToSnake") => "camel_to_snake"
// special case: CaseID => case_id
func CamelToSnake(input string) string {
	var result []rune
	input = strings.ReplaceAll(input, "ID", "Id")
	for i, r := range input {
		if unicode.IsUpper(r) {
			// 在大写字母前面添加下划线，除非它是第一个字母
			if i > 0 {
				result = append(result, '_')
			}
			// 转换为小写字母
			r = unicode.ToLower(r)
		}
		result = append(result, r)
	}
	return string(result)
}

package utils

import (
	"fmt"
	"testing"
)

func TestGetRunnerPath(t *testing.T) {
	t.Log(GetRunnerPath())
}

func ExampleGetRunnerPath() {

	// 获取当前运行路径
	fmt.Println(GetRunnerPath())

	// Output:
	// C:\Users\yss93\AppData\Local\Temp
}

// 在windows 下的运行情况
func ExampleHome() {
	fmt.Println(Home())

	// Output:
	// C:\Users\yss93 <nil>
}

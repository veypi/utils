//
// system.go
// Copyright (C) 2025 veypi <i@veypi.com>
// 2025-04-07 16:56
// Distributed under terms of the MIT license.
//

package utils

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/veypi/utils/logv"
)

// OpenURL 使用系统默认浏览器打开指定的 URL
func OpenBrowser(url string) error {
	var cmd *exec.Cmd

	// 根据操作系统选择合适的命令
	logv.Warn().Msg(runtime.GOOS)
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// 执行命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open URL: %w", err)
	}

	fmt.Printf("Successfully opened URL: %s\n", url)
	return nil
}

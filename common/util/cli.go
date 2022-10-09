package util

// 控制台程序常用工具

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// 颜色枚举
const (
	ColorBlack = iota
	ColorBlue
	ColorGreen
	ColorCyan
	ColorRed
	ColorPurple
	ColorYellow
	ColorLightGray
	ColorGray
	ColorLightBlue
	ColorLightGreen
	ColorLightCyan
	ColorLightRed
	ColorLightPurple
	ColorLightYellow
	ColorWhite
	ColorDefault = ColorWhite
)

// ===================================================

// ColorPrintf 在控制台打印特定颜色的字体
func ColorPrintf(color int, format string, arg ...interface{}) {
	fmt.Printf(format, arg...)
}

func ColorPrintln(color int, any interface{}) {
	fmt.Println(any)
}

// ClearConsole 清空控制台
func ClearConsole() error {
	env := runtime.GOOS
	var cmd *exec.Cmd
	switch env {
	case "linux":
		cmd = exec.Command("clear") // Linux example, its tested
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls") // windows example, its tested
	default:
		return fmt.Errorf("unexpire env: %v", env)
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
	return nil
}

// ScanInput 从控制台获取输入,(非整行,遇到空格会中断)
func ScanInput() string {
	var input string
	fmt.Scanf("%s", &input)
	return strings.TrimSpace(input)
}

// ScanStdLine 从控制台获取整行输入,带提示
func ScanStdLine() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please input a line: # ")
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

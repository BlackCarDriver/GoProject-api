package util

// 控制台程序常用工具

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

var (
	kernel32 *syscall.LazyDLL  = syscall.NewLazyDLL(`kernel32.dll`)
	proc     *syscall.LazyProc = kernel32.NewProc(`SetConsoleTextAttribute`)
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
	proc.Call(uintptr(syscall.Stdout), uintptr(color))
	fmt.Printf(format, arg...)
	proc.Call(uintptr(syscall.Stdout), uintptr(ColorDefault))
}

func ColorPrintln(color int, any interface{}) {
	proc.Call(uintptr(syscall.Stdout), uintptr(color))
	fmt.Println(any)
	proc.Call(uintptr(syscall.Stdout), uintptr(ColorDefault))
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

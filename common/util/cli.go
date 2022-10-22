package util

// 控制台程序常用工具

import (
	"bufio"
	"fmt"
	"golang.org/x/text/width"
	"os"
	"os/exec"
	"regexp"
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

// ============ ColumnPrinter =============

// ColumnPrinter 用于在控制台打印表格时对其各列
type ColumnPrinter struct {
	space       int
	table       [][]string
	colToMaxLen map[int]int // 列id->最大长度
}

// Write 输入一行,可用任意数量的空格隔开每个元素
func (c *ColumnPrinter) Write(str string) {
	reg := regexp.MustCompile("\\s+")
	cols := reg.Split(strings.TrimSpace(str), -1)
	if c.colToMaxLen == nil {
		c.colToMaxLen = map[int]int{}
	}
	// 更新各列的最大长度
	for i, item := range cols {
		t := getStringWidth(item)
		if c.colToMaxLen[i] < t {
			c.colToMaxLen[i] = t
		}
	}
	c.table = append(c.table, cols)
}

// Print 打印表格
func (c *ColumnPrinter) Print() {
	for _, line := range c.table {
		tmp := ""
		for j, col := range line {
			colWidth := c.colToMaxLen[j] + c.space
			tmp += fmt.Sprintf("%s%s", col, strings.Repeat(" ", colWidth-getStringWidth(col)))
		}
		//tmp := strings.Join(line, ",")
		fmt.Println(tmp)
	}
	c.table = [][]string{}
	c.colToMaxLen = map[int]int{}
}

func (c *ColumnPrinter) Debug() {
	for i, r := range c.table {
		fmt.Printf("table_%d: ", i)
		for _, item := range r {
			fmt.Printf("%s(%d), ", item, getStringWidth(item))
		}
		fmt.Println()
	}
	fmt.Printf("lenMap: %+v \n", c.colToMaxLen)
}

// NewColumnPrinter 创建一个实例, space控制每列元素的间距
func NewColumnPrinter(space int) (p ColumnPrinter) {
	if space <= 0 {
		space = 1
	}
	return ColumnPrinter{space: space}
}

// 获取字符串的宽度 (目的是解决中文和英文字符宽度不一导致显示对不齐)
func getStringWidth(str string) int {
	w := 0
	//reg := regexp.MustCompile("[^\\x00-\\xFF]+") // 包含非ascii字符
	//if !reg.MatchString(str) {
	//	return len(str)
	//}
	for _, c := range str {
		prot := width.LookupRune(c)
		if prot.Kind() == width.EastAsianWide {
			w = w + 2
		} else {
			w = w + 1
		}
	}
	return w
}

// ============  =============

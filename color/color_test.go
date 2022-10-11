package color

import (
	"fmt"
	"testing"
)

// 打印所有颜色
func TestDefault(t *testing.T) {
	Black("Black")
	HiBlack("HiBlack")
	White("White")
	HiWhite("HiWhite")
	Yellow("Yellow")
	HiYellow("HiYellow")
	Red("Red")
	HiRed("HiRed")
	Green("Green")
	HiGreen("HiGreen")
	Blue("Blue")
	HiBlue("HiBlue")
	Cyan("Cyan")
	HiCyan("HiCyan")
	Magenta("Magenta")
	HiMagenta("HiMagenta")
}

// 颜色字符串组合
func TestString(t *testing.T) {
	y := YellowString("===%s===", "YellowString")
	g := GreenString("===%s===", "GreenString")
	r := RedString("===%s===", "RedString")
	fmt.Printf("%s %s %s \n", y, g, r)
}

// 创建color对象
func TestMix(t *testing.T) {
	c := New(FgCyan).Add(Underline)
	c.Println("New(FgCyan).Add(Underline)")

	d := New(FgCyan, Bold)
	d.Printf("New(FgCyan, Bold)\n")

	red := New(FgRed)
	boldRed := red.Add(Bold)
	boldRed.Println("New(FgRed).Add(Bold)")
}

// 创建打印函数
func TestFunction(t *testing.T) {
	myFun := New(FgRed).Add(Bold, Underline).PrintfFunc() // PrintlnFunc
	myFun("PrintfFunc()")
	myFun("PrintfFunc()")
	myFun("PrintfFunc()")
	myFun("PrintfFunc()")
}

// 取消颜色
func TestUnSet(t *testing.T) {
	c := New(FgCyan)
	c.Println("Prints cyan text")

	c.DisableColor()
	c.Println("This is printed without any color")

	c.EnableColor()
	c.Println("This prints again cyan...")
}

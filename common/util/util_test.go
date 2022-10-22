package util

import "testing"

func TestColumnPrinter(t *testing.T) {
	pt := NewColumnPrinter(5)
	pt.Write("中证500 5.96 0.00%     2.23%    -4.19%")
	pt.Write("中概互联 0.91  -0.55% 7.24% 0.99%")
	pt.Write("招商银行    29.58 -1.63%    12.47% -1.80%")
	pt.Write("白酒   1.02 -1.36% 7.96% -1.57%")
	pt.Print()
}

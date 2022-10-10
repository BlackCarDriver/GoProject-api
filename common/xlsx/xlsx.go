package xlsx

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/tealeg/xlsx"
	"net/http"
)

// ExportXlsxWithHttp 传入二维数组，生成xlsx文件并通过http响应体返回文件
func ExportXlsxWithHttp(w http.ResponseWriter, form [][]string, fileName string) (err error) {
	binary, err := ParseFormToXlsx(form)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprint(len(binary)))
	w.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	_, err = w.Write(binary)
	return err
}

// ParseFormToXlsx 将二维数组表示的表格转换成二进制xlsx文件
func ParseFormToXlsx(before [][]string) (after []byte, err error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("main")
	if err != nil {
		return
	}
	for _, row := range before {
		sheet.AddRow().WriteSlice(&row, len(row))
	}
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)
	err = file.Write(bw)
	if err != nil {
		return
	}
	after = b.Bytes()
	return
}

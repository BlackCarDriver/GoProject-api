package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// http://push2his.eastmoney.com/api/qt/stock/trends2/get?fields1=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13&
// fields2=f51,f52,f53,f54,f55,f56,f57,f58&
// ut=fa5fd1943c7b386f172d6893dbfba10b&
// iscr=0&
// ndays=1&
// secid=1.000001&
// cb=jQuery3510005388373829342763_1665403466103&
// _=1665403466104

type kLineParams struct {
	Fields1   string `json:"fields1"`
	Fields2   string `json:"fields2"`
	UT        string `json:"ut"`
	ISCR      int    `json:"iscr"`
	NDays     int    `json:"ndays"`
	SECID     string `json:"secid"`
	CB        string `json:"-"`
	Timestamp int64  `json:"_"`
}

type Params2 struct {
	Fields    string `json:"fields"`
	CB        string `json:"cb"`
	Fltt      int    `json:"fltt"`
	SECIDS    string `json:"secids"` // 选中的股票ID, 参考: "1.000001,0.399001"
	UT        string `json:"ut"`
	Timestamp int64  `json:"_"`
}

// https://push2.eastmoney.com/api/qt/ulist.np/get?
// cb=jQuery112305582875234802821_1665462176107&
// fltt=2&
// secids=1.000001,0.399001&
// fields=f1,f2,f3,f4,f6,f12,f13,f104,f105,f106&
// ut=b2884a393a59ad64002292a3e90d46a5&
// _=1665462176108

type UtilRespLayer1 struct {
	Data UtilRespLayer2 `json:"data"`
}

type UtilRespLayer2 struct {
	Total int            `json:"total"`
	Diff  []UtilRespData `json:"diff"`
}

type UtilRespData struct {
	F2   float64 `json:"f2"`   // 现价
	F3   float64 `json:"f3"`   // 增幅
	F4   float64 `json:"f4"`   // 增值
	F6   float64 `json:"f6"`   // 总市值
	F12  string  `json:"f12"`  // 股票代码
	F104 int     `json:"f104"` // 涨_数量
	F105 int     `json:"f105"` // 跌_数量
	F106 int     `json:"f106"` // 平_数量
}

func GetSummaryInfo() {
	timestamp := time.Now().Unix()
	params := Params2{
		Fltt:      2,
		SECIDS:    "1.000001,0.399001",
		UT:        "b2884a393a59ad64002292a3e90d46a5",
		Fields:    "f1,f2,f3,f4,f6,f12,f13,f104,f105,f106",
		Timestamp: timestamp,
	}
	resp, err := GetRequireWithParams("https://push2.eastmoney.com/api/qt/ulist.np/get", params)
	if err != nil {
		fmt.Printf("get data fail: err=%v \n", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("unexpect status: %d \n", resp.StatusCode)
		return
	}
	data, _ := ioutil.ReadAll(resp.Body)
	var ret UtilRespLayer1
	err = json.Unmarshal(data, &ret)
	if err != nil {
		fmt.Printf("unmarshal fail: err=%v data=%s \n", err, string(data))
		return
	}
	fmt.Printf("response=%+v", ret.Data)
}

func TestGetRequireWithParams(t *testing.T) {
	GetSummaryInfo()
}

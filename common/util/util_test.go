package util

import (
	"testing"
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

// https://push2.eastmoney.com/api/qt/ulist.np/get?
// cb=jQuery112305582875234802821_1665462176107&
// fltt=2&
// secids=1.000001,0.399001&
// fields=f1,f2,f3,f4,f6,f12,f13,f104,f105,f106&
// ut=b2884a393a59ad64002292a3e90d46a5&
// _=1665462176108
func TestGetRequireWithParams(t *testing.T) {

}

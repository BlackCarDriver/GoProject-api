package util

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

func MyPrint() {
	fmt.Println("successs!!!")
}

/*HttpValuesMustUnmarshalJson 将Get请求或Post请求中传输的参数赋值到结构体里面的字段中,以下为注意事项：
- 目标结构体中的字段类型目前仅支持：int,int32,int64,float32,float64,bool,string, 包含其他类型将返回错误;
- 目标结构体中的所有字段必须导出(大写开头),否则将返回错误;
- ptrToTargetb类型必须为指向目标结构体的指针, 目标结构体中所有字段的'json标签'必须在表单中找得到相应字段;
- 只有结构体中全部字段都在表单中找到并且转换成功才返回nil;
- 返回error时代表参数错误或未完成全部字段的赋值，不代表目标结构体未发生变化;
- post请求解析失败时检查：前端使用post请求时加上'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8';
*/
func HttpValuesMustUnmarshalJson(req *http.Request, ptrToTarget interface{}) (err error) {
	defer func() {
		msg, ok := recover().(interface{})
		if ok {
			var errptr = &err
			*errptr = fmt.Errorf("catch panic: err=%v", msg)
		}
		// log.Info("convert result: error=%v  target=%v", err, ptrToTarget)
	}()

	if req == nil {
		return fmt.Errorf("params request is null")
	}
	req.ParseForm()
	if len(req.Form) == 0 {
		return fmt.Errorf("request Form is empty")
	}
	// 确认是指针
	rType := reflect.TypeOf(ptrToTarget)
	rValue := reflect.ValueOf(ptrToTarget)
	if rType.Kind() != reflect.Ptr {
		return fmt.Errorf("target not a pointer: type=%v", rType.Kind())
	}
	// 确认指向结构体
	rType = rType.Elem()
	rValue = rValue.Elem()
	if rType.Kind() != reflect.Struct {
		return fmt.Errorf("target is not a pointer to struct: kind=%v", rType.Kind())
	}
	// 确认能被修改
	if !rValue.CanSet() {
		return fmt.Errorf("target can't be changed")
	}

	// 遍历结构体中的字段并从表单中获取相应值
	for i := 0; i < rValue.NumField(); i++ {
		tmpVal := rValue.Field(i)
		vname := rType.Field(i).Name

		if !tmpVal.CanSet() {
			return fmt.Errorf("field can't be set, index=%d name=%s", i, vname)
		}

		jsTag, found := rType.Field(i).Tag.Lookup("json")
		if !found {
			return fmt.Errorf("json tag not found, index=%d name=%v tag=%v", i, vname, rType.Field(i).Tag)
		}

		// 从请求表单中获取相应值
		rawStr := req.Form.Get(jsTag)
		if rawStr == "" && tmpVal.Kind() != reflect.String {
			return fmt.Errorf("field not found in query form: index=%d name=%s jsTag=%s form=%v", i, vname, jsTag, req.Form)
		}

		// 检查字段类型并赋值
		switch tmpVal.Kind() {
		case reflect.String:
			tmpVal.SetString(rawStr)
		case reflect.Int, reflect.Int32, reflect.Int64:
			tmpInt, err := strconv.ParseInt(rawStr, 10, 64)
			if err != nil {
				return fmt.Errorf("parse Int fail: index=%d name=%s rawStr=%s err=%v", i, vname, rawStr, err)
			}
			tmpVal.SetInt(tmpInt)

		case reflect.Float32, reflect.Float64:
			tmpFloat, err := strconv.ParseFloat(rawStr, 64)
			if err != nil {
				return fmt.Errorf("parse Float fail: index=%d name=%s rawStr=%s err=%v", i, vname, rawStr, err)
			}
			tmpVal.SetFloat(tmpFloat)

		case reflect.Bool:
			tmpBool, err := strconv.ParseBool(rawStr)
			if err != nil {
				return fmt.Errorf("parse Bool fail: index=%d name=%s rawStr=%s err=%v", i, vname, rawStr, err)
			}
			tmpVal.SetBool(tmpBool)

		default:
			return fmt.Errorf("unsupport kind of field, index=%d name=%v kind=%v", i, vname, tmpVal.Kind())
		}
	}
	return nil
}

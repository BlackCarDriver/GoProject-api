package binder

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

// SmartBinder binding.Binding接口的一种实现,可用于将常见形式的请求参解析到结构体
type SmartBinder struct{}

func (b SmartBinder) Name() string {
	return "SmartBinder"
}

// Bind 将Get请求或Post请求中传输的参数赋值到结构体里面的字段中, 解析方法与 util.HttpValuesMustUnmarshalJson 一致
func (b SmartBinder) Bind(req *http.Request, ptrToStruct interface{}) (err error) {
	defer func() {
		msg, ok := recover().(interface{})
		if ok {
			var errptr = &err
			*errptr = fmt.Errorf("catch panic: err=%v", msg)
		}
	}()

	if req == nil {
		return fmt.Errorf("params request is null")
	}
	req.ParseForm()
	if len(req.Form) == 0 {
		return fmt.Errorf("request Form is empty")
	}
	// 确认是指针
	rType := reflect.TypeOf(ptrToStruct)
	rValue := reflect.ValueOf(ptrToStruct)
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

// GetSmartBinderInstances 获取SmartBinder实例
func GetSmartBinderInstances() SmartBinder {
	return SmartBinder{}
}

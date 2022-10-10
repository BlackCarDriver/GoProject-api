package util

// http请求相关工具

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

// GetIpAndPortFromRequire 从请求中获取IP和端口
func GetIpAndPortFromRequire(r *http.Request) (remoteAddr, port string) {
	remoteAddr = r.RemoteAddr
	port = "?"
	XForwardedFor := "X-Forwarded-For"
	XRealIP := "X-Real-IP"
	if ip := r.Header.Get(XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = r.Header.Get(XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, port, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return
}

// WriteJson json格式响应
func WriteJson(w http.ResponseWriter, data interface{}) error {
	if bytes, err := json.Marshal(data); err != nil {
		return err
	} else if _, err = w.Write(bytes); err != nil {
		return err
	}
	return nil
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
			var tmpInt int64
			tmpInt, err = strconv.ParseInt(rawStr, 10, 64)
			if err != nil {
				return fmt.Errorf("parse Int fail: index=%d name=%s rawStr=%s err=%v", i, vname, rawStr, err)
			}
			tmpVal.SetInt(tmpInt)

		case reflect.Float32, reflect.Float64:
			var tmpFloat float64
			tmpFloat, err = strconv.ParseFloat(rawStr, 64)
			if err != nil {
				return fmt.Errorf("parse Float fail: index=%d name=%s rawStr=%s err=%v", i, vname, rawStr, err)
			}
			tmpVal.SetFloat(tmpFloat)

		case reflect.Bool:
			var tmpBool bool
			tmpBool, err = strconv.ParseBool(rawStr)
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

// GetRequireWithParams 发起Get请求,请求参数的字段和值从params中获取
func GetRequireWithParams(rawURL string, params interface{}) (resp *http.Response, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	if !u.IsAbs() {
		err = fmt.Errorf("url must use absolute path")
		return
	}
	res, err := parseObjToURLParams(params)
	if err != nil {
		err = fmt.Errorf("unexpect params: err=%v", err)
		return
	}
	u.RawQuery = res.Encode()
	fixedURL := u.String()

	resp, err = http.Get(fixedURL)
	return
}

// 结构体转url传参
func parseObjToURLParams(obj interface{}) (res url.Values, err error) {
	res = url.Values{}
	// 确认是结构体
	rType := reflect.TypeOf(obj)
	rValue := reflect.ValueOf(obj)
	if rType.Kind() != reflect.Struct {
		err = fmt.Errorf("target type must be struct: type=%v", rType.Kind())
		return
	}

	// 遍历结构体中的字段并从表单中获取相应值
	for i := 0; i < rValue.NumField(); i++ {
		jsTag, found := rType.Field(i).Tag.Lookup("json")
		if !found || jsTag == "-" {
			continue
		}
		// 判断是否类型是否支持
		isSupportKind := false
		supportKind := []reflect.Kind{reflect.String, reflect.Bool, reflect.Int64, reflect.Int32, reflect.Int}
		for _, k := range supportKind {
			if k == rType.Field(i).Type.Kind() {
				isSupportKind = true
				break
			}
		}
		if !isSupportKind {
			continue
		}
		res.Add(jsTag, fmt.Sprint(rValue.Field(i)))
	}
	return
}

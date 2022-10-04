package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// GetResponseBodyByURL 请求指定的url并将响应主体返回
func GetResponseBodyByURL(url string) (data []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// DownloadSourceByURL 将一个url的响应主体保存到本地文件
func DownloadSourceByURL(url, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Response code is %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodySize := len(body)
	if bodySize == 0 {
		return fmt.Errorf("body size = 0")
	}

	imgPath := fmt.Sprintf("./%s", fileName)
	out, err := os.Create(imgPath)
	if err != nil {
		return fmt.Errorf("create file fail: %v", err)
	}

	defer out.Close()

	_, err = io.Copy(out, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("save file fail: %v", err)
	}

	return nil
}

// CheckAndFixUrlToAbs 检验一个url, 且将相对地址转换为绝对地址
func CheckAndFixUrlToAbs(baseUrl string, targetUrl *string) error {
	if baseUrl == "" {
		return fmt.Errorf("baseUrl is empty")
	}
	if targetUrl == nil {
		return fmt.Errorf("targetUrl is nil")
	}
	target, err := url.Parse(*targetUrl)
	if err != nil {
		return err
	}
	//若本身为绝对路径则无需转换
	if target.IsAbs() {
		return nil
	}
	base, err := url.Parse(baseUrl)
	if err != nil {
		return fmt.Errorf("BaseUrl not right. err=%v", err)
	}
	if !base.IsAbs() {
		return fmt.Errorf("BaseUrl is not absolute url")
	}
	*targetUrl = fmt.Sprintf("%s://%s%s", base.Scheme, base.Host, target.Path)
	return nil
}

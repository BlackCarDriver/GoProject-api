package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// GetFileListByPath 根据文件夹路径,返回下面的文件列表
func GetFileListByPath(dir string) (fileList []string) {
	file, err := os.Open(dir)
	if err != nil {
		return nil
	}
	defer file.Close()
	fi, err := file.Readdir(0)
	if err != nil {
		return nil
	}
	fileList = make([]string, 0)
	for _, info := range fi {
		fileList = append(fileList, info.Name())
	}
	return fileList
}

// ParseFileToString 读取指定路径的文件,返回其文本类容
func ParseFileToString(path string) (text string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("Open %s fall: %v", path, err)
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	bytes, err := ioutil.ReadAll(buf)
	if err != nil {
		return "", fmt.Errorf("ioutil.ReadAll fall : %v", err)
	}
	return string(bytes), nil
}

// ClearFile 清空指定路径的文件
func ClearFile(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.WriteString(""); err != nil {
		return err
	}
	return nil
}

// CheckFileIsExist 检查指定路径的文件是否已存在的文件
func CheckFileIsExist(path string) (isExist bool, err error) {
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return true, err
	}
	return true, nil
}

// CheckDirIsExist 检查指定的路径是否已存在的文件夹
func CheckDirIsExist(dir string) (isExist bool, err error) {
	info, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
		return false, nil
	}
	return info.IsDir(), nil
}

// SaveContentToFile 保存文本到指定文件
func SaveContentToFile(content, path string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file '%s' is alread exist", path)
	} else if os.IsNotExist(err) == false {
		return err
	}
	err := ioutil.WriteFile(path, []byte(content), 0644)
	return err
}

// SaveDataToFile 将指定数据保存到文件
func SaveDataToFile(data []byte, path string) (err error) {
	if _, err = os.Stat(path); err == nil || !os.IsNotExist(err) {
		err = fmt.Errorf("file may already exist")
		return err
	}
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, bytes.NewReader(data))
	if err != nil {
		return err
	}
	return nil
}

// UnmarshalJsonFromFile 将指定文件的文本内容以json格式解析到结构体中
func UnmarshalJsonFromFile(path string, v interface{}) (err error) {
	content, err := ParseFileToString(path)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(content), &v)
	return
}

// MarshalJsonToFile 将内容以json格式存储到指定文件中
func MarshalJsonToFile(path string, v interface{}) (err error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return
	}
	err = SaveContentToFile(string(bs), path)
	return
}
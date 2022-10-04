package util

// 系统相关工具

import (
	"bytes"
	"os/exec"
)

// LinuxExec 执行一指定命令并返回输出结果
func LinuxExec(name string, arg ...string) (string, error) {
	var out, stderr bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	var err error
	if err = cmd.Start(); err != nil {
		return "", err
	} else if err = cmd.Wait(); err != nil {
		return "", err
	}
	return out.String(), nil
}

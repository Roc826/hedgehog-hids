package common

import (
	"os/exec"
	"runtime"
	"strings"
	"github.com/axgle/mahonia"
)

var (
	// Config 配置信息
	Config ClientConfig
	// LocalIP 本机活跃IP
	LocalIP string
	// ServerInfo 主机相关信息
	ServerInfo ComputerInfo
	// ServerIPList 服务端列表
	ServerIPList []string
)

func InSlice(str string, slice []string) bool {
	for _, a := range slice {
		if str == a {
			return true
		}
	}
	return false

}



func CmdExec(cmd string) (string, error) {
	var c *exec.Cmd
	var data string
	var err error
	system := runtime.GOOS
	if system == "windows" {
		argArray := strings.Split("/c" + cmd, " ")
		c = exec.Command("cmd",argArray...)
	}else{
		c = exec.Command("/bin/sh","-c", cmd)
	}
	out,err := c.CombinedOutput()
	if err != nil{
		return data,err
	}
	data = string(out)
	if system == "windows"{
		dec :=mahonia.NewDecoder("gbk")
		data = dec.ConvertString(data)
	}
	return data,err
}

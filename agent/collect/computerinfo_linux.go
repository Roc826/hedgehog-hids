// +build linux

package collect

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"hedgehog-hids-agent/common"
)

// GetComInfo 获取计算机信息

func GetComInfo() (info common.ComputerInfo) {
	var issue string
	info.IP = common.LocalIP
	info.Hostname, _ = os.Hostname()
	out ,err:= common.CmdExec("uname -r")
	dat, err := ioutil.ReadFile("/etc/redhat-release")
	if err != nil {
		dat, _ = ioutil.ReadFile("/etc/issue")
		reg:= regexp.MustCompile("^[a-zA-Z0-9_ ]*")
		result := reg.FindStringSubmatch(string(dat))
		if len(result)!=0{
			issue = result[0]
		}else {
			issue = strings.SplitN(string(dat), "\\n", 2)[0]
		}
		out2,_:= common.CmdExec("uname -m")
		info.System = strings.Replace( issue + out + " "+ out2,"\n","",-1)
	} else {
		info.System = string(dat) + " " + out
	}
	return info
}

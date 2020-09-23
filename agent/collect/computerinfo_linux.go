// +build linux

package collect

import (
	"io/ioutil"
	"os"
	"strings"
	"hedgehog-hids-agent/common"
)

// GetComInfo 获取计算机信息

//这里需要检查
func GetComInfo() (info common.ComputerInfo) {
	info.IP = common.LocalIP
	info.Hostname, _ = os.Hostname()
	out ,err:= common.CmdExec("uname -r")
	dat, err := ioutil.ReadFile("/etc/redhat-release")
	if err != nil {
		dat, _ = ioutil.ReadFile("/etc/issue")
		issue := strings.SplitN(string(dat), "\n", 2)[0]
		out2,_:= common.CmdExec("uname -m")
		info.System = issue + " " + out + out2
	} else {
		info.System = string(dat) + " " + out
	}
	//discern(&info)
	return info
}

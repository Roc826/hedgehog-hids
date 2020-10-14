package collect

import (
	"io/ioutil"
	"regexp"
	"strings"
)

func GetCrontab() (resultData []map[string]string) {
	//系统计划任务
	dat, err := ioutil.ReadFile("/etc/crontab")
	if err != nil {
		return resultData
	}
	cronList := strings.Split(string(dat), "\n")
	for _, info := range cronList {
		if strings.HasPrefix(info, "#") || strings.Count(info, " ") < 6 {
			continue
		}
		reg := regexp.MustCompile("(([0-9/*]+\\s*){5})\\s+(\\w*)\\s+(.*)")
		submatchs := reg.FindStringSubmatch(info)
		m := map[string]string{"command": submatchs[4], "user": submatchs[3], "rule": submatchs[1]}
		resultData = append(resultData, m)
	}

	//用户计划任务
	dir, err := ioutil.ReadDir("/var/spool/cron/")
	if err != nil {
		return resultData
	}
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		dat, err = ioutil.ReadFile("/var/spool/cron/" + f.Name())
		if err != nil {
			continue
		}
		cronList = strings.Split(string(dat), "\n")
		for _, info := range cronList {
			if strings.HasPrefix(info, "#") || strings.Count(info, " ") < 5 {
				continue
			}
			reg := regexp.MustCompile("(([0-9/*]+\\s*){5})\\s+(\\w*)\\s+(.*)")
			submatchs := reg.FindStringSubmatch(info)
			m := map[string]string{"command": submatchs[4], "user": submatchs[3], "rule": submatchs[1]}
			resultData = append(resultData, m)
		}
	}
	return resultData
}
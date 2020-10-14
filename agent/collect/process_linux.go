package collect

import (
	"fmt"
	"hedgehog-hids-agent/log"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"hedgehog-hids-agent/common"
)

// Process is an implementation of Process that contains Unix-specific
type Process struct {
	pid    int
	ppid   int
	state  rune
	pgrp   int
	sid    int
	ruid   int
	euid   int
	suid   int
	fsuid  int
	status map[string]string

	binary string
}

func (p *Process) Pid() int {
	return p.pid
}

func (p *Process) Ppid() int {
	return p.ppid
}

func (p *Process) Binary() string {
	return p.binary
}

func (p *Process) Pgrp() int {
	return p.pgrp
}

func (p *Process) Ruid() int {
	return p.ruid
}

func Processes() ([]Process, error) {
	d, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer d.Close()

	results := make([]Process, 0, 50)
	for {
		names, err := d.Readdirnames(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, name := range names {
			// We only care if the name starts with a numeric
			if name[0] < '0' || name[0] > '9' {
				continue
			}

			// From this point forward, any errors we just ignore, because
			// it might simply be that the process doesn't exist anymore.
			pid, err := strconv.ParseInt(name, 10, 0)
			if err != nil {
				continue
			}

			p, err := newProcess(int(pid))
			if err != nil {
				continue
			}

			results = append(results, *p)
		}
	}

	return results, nil
}

func GetUserProcess(user string) ([]Process, error) {
	results := make([]Process, 0, 10)
	cmd := fmt.Sprintf("ps -u %s |awk '{print $1}'|grep -v -i \"pid\"", user)
	out ,err:= common.CmdExec(cmd)
	if err != nil{
		log.Error("exec fail:",err.Error())
		return results,err
	}
	pids := strings.Split(out, "\n")
	for _, pid := range pids {
		pid, err := strconv.Atoi(pid)
		if err != nil {
			continue
		}

		p, err := newProcess(pid)
		if err != nil {
			continue
		}

		results = append(results, *p)
	}
	return results, nil
}

func newProcess(pid int) (*Process, error) {
	p := &Process{pid: pid}
	return p, p.Refresh()
}

// Refresh reloads all the data associated with this process.
func (p *Process) Refresh() error {
	status := p.GetStatus()
	p.status = status
	statPath := fmt.Sprintf("/proc/%d/stat", p.pid)
	dataBytes, err := ioutil.ReadFile(statPath)
	if err != nil {
		return err
	}

	// First, parse out the image name
	data := string(dataBytes)
	binStart := strings.IndexRune(data, '(') + 1
	binEnd := binStart
	for {
		End := strings.IndexRune(data[binEnd:], ')')
		if End != -1 {
			binEnd += End
			if strings.IndexRune(data[binEnd+1:], ')') != -1 {
				binEnd += 1
				continue
			}
		}
		break
	}
	p.binary = data[binStart:binEnd]

	// Move past the image name and start parsing the rest
	data = data[binEnd+2:]
	_, err = fmt.Sscanf(data,
		"%c %d %d %d",
		&p.state,
		&p.ppid,
		&p.pgrp,
		&p.sid)

	_, err = fmt.Sscanf(status["Uid"], "%d %d %d %d",
		&p.ruid,
		&p.euid,
		&p.suid,
		&p.fsuid)

	return err
}

func (p *Process) GetStatu(key string) string {
	if value, ok := p.status[key]; ok {
		return value
	}
	return ""
}

func (p *Process) GetParentPs() (*Process, error) {
	pp := &Process{pid: p.ppid}
	return pp, pp.Refresh()
}

func (p *Process) getProcessChain() []Process {
	var processChainTemp []Process
	var processChain []Process
	var err error
	p_temp := p
	processChainTemp = append(processChain, *p_temp)
	for {
		ppid := p_temp.Ppid()
		if ppid != 0 {
			p_temp, err = p_temp.GetParentPs()
			if err != nil {
				log.Error("获取父进程失败")
			}
			processChainTemp = append(processChainTemp, *p_temp)
			continue
		}
		break
	}
	processChain = reverse(processChainTemp)
	return processChain

}

func GetPsFromPipeId(uid string, pipeid int) ([]Process, error) {
	results := make([]Process, 0)
	cmd := fmt.Sprintf("lsof -u %s 2>/dev/null |grep \"%d pipe\"|awk '{print $2}' ", uid, pipeid)
	out,err:= common.CmdExec(cmd)
	if err != nil{
		log.Error("exec fail:",err.Error())
		return results,err
	}
	pids := strings.Split(out, "\n")
	for _, pid := range pids {
		pid, err := strconv.Atoi(pid)
		if err != nil {
			continue
		}
		p, err := newProcess(pid)
		if err != nil {
			continue
		}
		results = append(results, *p)
	}

	return results, nil
}

func (p *Process) GetCmdline() string {
	cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", p.Pid())
	cmdline, err := ioutil.ReadFile(cmdlinePath)
	if err != nil {
		log.Error("读取进程%dCmdline失败\n", p.Pid())
	}
	return string(cmdline)

}

func (p *Process) GetStatus() map[string]string {
	res := make(map[string]string)
	statusReg := "([a-zA-z_]+?):\\s*(.*)"
	re := regexp.MustCompile(statusReg)
	statusPath := fmt.Sprintf("/proc/%d/status", p.Pid())
	contents, err := ioutil.ReadFile(statusPath)
	if err != nil {
		log.Error("读取文件%s失败\n", statusPath)
	}
	matches := re.FindAllStringSubmatch(string(contents), -1)
	for _, m := range matches {
		res[m[1]] = m[2]
	}
	return res
}

func (p *Process) ShowProcessChain() {
	PSC := p.getProcessChain()
	for _, ps := range PSC {
		fmt.Printf("Pid:%-10d\tName:%-30sCmdline:%s\n", ps.Pid(), ps.Binary(), ps.GetCmdline())
	}
}

func (p *Process) GetAllFd() []int {
	var fd []int
	// 读取当前目录中的所有文件和子目录
	files, err := ioutil.ReadDir(fmt.Sprintf("/proc/%d/fd", p.Pid()))
	if err != nil {
		panic(err)
	}

	// 获取文件，并输出它们的名字
	for _, file := range files {
		name, _ := strconv.Atoi(file.Name())
		fd = append(fd, name)
	}
	return fd
}

func reverse(s []Process) []Process {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

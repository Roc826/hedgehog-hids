package common

type ClientConfig struct {
	Cycle  int    // 信息传输频率，单位：分钟
	UDP    bool   // 是否记录UDP请求
	LAN    bool   // 是否本地网络请求
	Mode   string // 模式，考虑中
	Filter struct {
		File    []string // 文件hash、文件名
		IP      []string // IP地址
		Process []string // 进程名、参数
	}                    // 直接过滤不回传的规则
	MonitorPath []string // 监控目录列表
	Lasttime    string   // 最后一条登录日志时间
}

// ComputerInfo 计算机信息结构
type ComputerInfo struct {
	IP       string   // IP地址
	System   string   // 操作系统
	Hostname string   // 计算机名
	Type     string   // 服务器类型
	Path     []string // WEB目录
}

type DataInfo struct {
	Ip     string
	Type   string
	System string
	Data   []map[string]string
}

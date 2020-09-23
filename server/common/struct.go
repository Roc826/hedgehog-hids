package common

import "time"

type ComputerInfo struct {
	IP       string
	System   string
	Hostname string
	Type     string
	Path     []string

	Uptime time.Time
}
// ClientConfig 客户端配置信息结构
type ClientConfig struct {
	Cycle       int      `bson:"cycle"` // 信息传输频率，单位：分钟
	UDP         bool     `bson:"udp"`   // 是否记录UDP请求
	LAN         bool     `bson:"lan"`   // 是否本地网络请求
	Mode        string   `bson:"mode"`  // 模式，考虑中
	Filter      filter   // 直接过滤不回传的数据
	MonitorPath []string `bson:"monitorPath"` // 监控目录列表
	Lasttime    string   // 最后一条登录日志时间
}

type filter struct {
	File    []string `bson:"file"`    // 文件hash、文件名
	IP      []string `bson:"ip"`      // IP地址
	Process []string `bson:"process"` // 进程名、参数
}
type blackList struct {

}
type whiteList struct {
}
type intelligence struct {

}
type notice struct {

}

type serverConfig struct {
	Learn        bool         `bson:"learn"`        // 是否为观察模式
	OfflineCheck bool         `bson:"offlinecheck"` // 开启离线主机检测和通知
	BlackList    blackList    // 黑名单
	WhiteList    whiteList    // 白名单
	Private      string       `bson:"privatekey"` // 加密秘钥
	Cert         string       `bson:"cert"`       // TLS加密证书
	Intelligence intelligence // 威胁情报
	Notice       notice       // 通知
}
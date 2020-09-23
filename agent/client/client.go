package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
	"hedgehog-hids-agent/collect"
	"hedgehog-hids-agent/common"
	"hedgehog-hids-agent/log"
	"hedgehog-hids-agent/monitor"
	"net"
	"strings"
	"sync"
)

var err error

type Agent struct {
	ServerNetLoc string          // 服务端地址 IP:PORT
	Client       client.XClient  // RPC 客户端
	ServerList   []string        // 存活服务端集群列表
	PutData      common.DataInfo // 要传输的数据
	Reply        int             // RPC Server 响应结果
	Mutex        *sync.Mutex     //安全操作锁
	IsDebug      bool            //是否开启debug模式
	ctx          context.Context
}

func (a *Agent) init() {
	a.ServerList, err = a.getServerList()
	if err != nil {
		a.log(fmt.Sprint("get ServerList fail:", err))
		panic(1)
	}
	a.ctx = context.WithValue(context.Background(), share.ReqMetaDataKey, make(map[string]string))
	a.log(fmt.Sprint("Available server node:", a.ServerList))
	if len(a.ServerList) == 0 {
		a.log("No server node available")
		panic(1)
	}
	a.newClient()
	if common.LocalIP == "" {
		a.log("Can not get local address")
		panic(1)
	}
	a.Mutex = new(sync.Mutex)
	err := a.Client.Call(a.ctx, "GetInfo", &common.ServerInfo, &common.Config)
	if err != nil {
		a.log("RPC Client Call Error:", err.Error())
		panic(1)
	}
	a.log("Common Client Config:", common.Config)

}

func (a *Agent) monitor() {
	resultChan := make(chan map[string]string, 16)
	monitor.FileMonitor(resultChan)
	//go
	//这一段是RPC用的后期在开发
	//go func(result chan map[string]string) {
	//	var resultdata []map[string]string
	//	var data map[string]string
	//	for {
	//		data = <-result
	//		data["time"] = fmt.Sprintf(fmt.Sprintf("%d", time.Now().Unix()))
	//		source := data["source"]
	//		delete(data,"source")
	//	}
	//}(resultChan)
}

func (a *Agent) getServerList() ([]string, error) {
	var serlist []string
	//先作为测试,写死代码
	serlist = append(serlist, "localhost:8972")

	return serlist, nil
}

//通过链接一台服务器来获取自己的ip
func (a Agent) setLocalIP(ip string) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		a.log("Net.Dial:", ip)
		a.log("Error:", err)
		print(err)
		panic(1)
	}
	defer conn.Close()
	common.LocalIP = strings.Split(conn.LocalAddr().String(), ":")[0]
}

func (a *Agent) newClient() {
	var servers []*client.KVPair
	for _, server := range a.ServerList {
		common.ServerIPList = append(common.ServerIPList, strings.Split(server, ":")[0])
		s := client.KVPair{Key: server}
		servers = append(servers, &s)
		if common.LocalIP == "" {
			a.setLocalIP(server)
			common.ServerInfo = collect.GetComInfo()
			a.log("Host Information:", common.ServerInfo)
		}
	}
	conf := &tls.Config{
		InsecureSkipVerify: true, //不开启证书验证
	}
	option := client.DefaultOption
	option.TLSConfig = conf
	serverd := client.NewMultipleServersDiscovery(servers)
	a.Client = client.NewXClient("Watcher", FAILMODE, client.RandomSelect, serverd, option)
	a.Client.Auth(AUTH_TOKEN)
}
func (a *Agent) Run() {
	//初始化数据
	a.init()
	//开启监控
	a.monitor()
}

func (a *Agent) log(info ...interface{}) {
	if a.IsDebug {
		log.Debug(info...)
	}
}

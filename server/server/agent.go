package server

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"hedgehog-hids-server/common"
	"hedgehog-hids-server/log"
)

const authToken string = "es03aq9tdu9gxyk6wbyvi0nbzqgeorup"

type Watcher int

func (w *Watcher) PutInfo(ctx context.Context, info *common.ComputerInfo, result *common.ClientConfig) error{
	return nil
}

func (w *Watcher) GetInfo(ctx context.Context, info *common.ComputerInfo, result *common.ClientConfig) error{
	print("getinfo")
	return nil
}

func auth(ctx context.Context, req *protocol.Message, token string) error {
	if token == authToken {
		return nil
	}
	return errors.New("invalid token")
}



func serverInit() {

}
func Run(addr string) {
	serverInit()
	startServer(addr)
}


func startServer(addr string){
	cert,err := tls.LoadX509KeyPair("cert.pem","key.pem")
	if err != nil{
		log.Error("cert error!",err)
		return
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	s := server.NewServer(server.WithTLSConfig(config))
	//s.AuthFunc = auth
	s.RegisterName("Watcher",new(Watcher),"")
	log.Info("RPC Server started")
	err = s.Serve("tcp",addr)
	if err != nil{
		log.Error(err.Error())
	}

}
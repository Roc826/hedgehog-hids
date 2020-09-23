package main

import (
	"flag"
	"hedgehog-hids-server/server"
)
func init(){

}

func main(){
	var addr = flag.String("addr", "localhost:8972", "server address")
	flag.Parse()
	server.Run(*addr)
}

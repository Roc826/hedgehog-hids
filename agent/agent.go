package main

import (
	"fmt"
	"hedgehog-hids-agent/client"
	"hedgehog-hids-agent/log"
	"os"
)

var err error
func main(){
	if len(os.Args) <=1 {
		fmt.Println("Usage:agent ServerIp [debug]")
		fmt.Println("Example:agent 8.8.8.8 debug")
		return
	}

	var agent client.Agent
	if len(os.Args) == 3 && os.Args[2] == "debug"{
		log.Debug("DEBUG MODE")
		agent.IsDebug = true
	}
	agent.Run()
	select{
	}
}

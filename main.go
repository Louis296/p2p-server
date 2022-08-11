package main

import (
	"github.com/louis296/p2p-server/config"
	"github.com/louis296/p2p-server/pkg/server"
	"os"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		//log
		os.Exit(1)
	}

	wsServer := server.NewP2PServer(WebSocketHandler)

	defaultConfig := server.DefaultConfig()

	defaultConfig.Host = conf.Server.Host
	defaultConfig.Port = conf.Server.Port
	defaultConfig.HTMLRoot = conf.Server.HTMLRoot

	wsServer.Bind(defaultConfig)
}

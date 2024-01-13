package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"proximity-chat/config"
	"proximity-chat/controller"
	"proximity-chat/discover"
	"proximity-chat/message"
	"proximity-chat/service"

	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	networkConfig := config.DefaultNetworkConfig()
	config.SetConfig(networkConfig)
	minPort := config.GetConfig().GetMinPort()
	maxPort := config.GetConfig().GetMaxPort()
	discover.SetServiceCIDR(config.GetConfig().GetCIDR())
	if minPort > maxPort {
		minPort = maxPort
	}
	var l net.Listener
	var err error
	var addr string
	for port := minPort; port <= maxPort; port++ {
		discover.SetServicePort(port)
		addr = discover.GetAddr()
		l, err = net.Listen("tcp", addr)
		if err != nil {
			fmt.Println("failed to listen: ", err)
			continue
		} else {
			break
		}
	}
	if l == nil {
		log.Fatal("unable to start server")
	}
	fmt.Println("listening on " + addr)
	go discover.BeginDiscoverService()
	w := &controller.OutputController{}
	discover.SetWriter(w)
	s := grpc.NewServer()
	go w.Start()
	message.RegisterChatServer(s, &service.Message{Writer: w})
	go func() {
		if err := s.Serve(l); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	var msg string
	for {
		fmt.Scanln(&msg)
		controller.Publish(msg)
	}
}

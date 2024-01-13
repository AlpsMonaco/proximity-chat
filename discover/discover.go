package discover

import (
	"context"
	"fmt"
	"proximity-chat/config"
	"proximity-chat/controller"
	"proximity-chat/message"
	"proximity-chat/service"
	"strconv"
	"time"

	"github.com/korylprince/ipnetgen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serviceCIDR string = ""
var serviceIP string = ""
var servicePort = 0
var writer service.MessageWriter = nil

func SetServiceCIDR(cidr string) {
	serviceCIDR = cidr
}

func SetServicePort(port int) {
	servicePort = port
}

func SetServiceIP(ip string) {
	serviceIP = ip
}

func GetServiceIP() string {
	return serviceIP
}

func GetServicePort() int {
	return servicePort
}

func GetServiceCIDR() string {
	return serviceCIDR
}

func GetAddr() string {
	return GetServiceIP() + ":" + strconv.Itoa(GetServicePort())
}

func SetWriter(wr service.MessageWriter) {
	writer = wr
}

func BeginDiscoverService() {
	minPort := config.GetConfig().GetMinPort()
	maxPort := config.GetConfig().GetMaxPort()
	if minPort > maxPort {
		minPort = maxPort
	}
	for {
		time.Sleep(time.Second)
		gen, err := ipnetgen.New(config.GetConfig().GetCIDR())
		if err != nil {
			panic(err)
		}
		for ip := gen.Next(); ip != nil; ip = gen.Next() {
			for i := minPort; i <= maxPort; i++ {
				addr := fmt.Sprintf("%s:%d", ip.String(), i)
				if addr == GetAddr() {
					continue
				}
				if controller.IsChatNodeExist(addr) {
					continue
				}
				conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					fmt.Printf("did not connect: %v\n", err)
					continue
				}
				client := message.NewChatClient(conn)
				cli, err := client.NewNode(context.Background())
				if err != nil {
					continue
				}
				err = cli.Send(&message.NodeRequest{Msg: GetAddr()})
				if err != nil {
					writer.Write(fmt.Sprint(err))
					continue
				}
				resp, err := cli.Recv()
				if err != nil {
					cli.CloseSend()
					writer.Write(fmt.Sprint(err))
					continue
				}
				if resp.GetMsg() != "ok" {
					cli.CloseSend()
					continue
				}
				if !controller.AddChatNode(&service.ClientChatNode{C: cli}, addr) {
					cli.CloseSend()
					continue
				}
				writer.Write("discover " + addr)
				go func() {
					for {
						msg, err := cli.Recv()
						if err != nil {
							writer.Write(fmt.Sprint(err))
							controller.RemoveNode(addr)
							return
						}
						writer.Write(msg.GetMsg())
					}
				}()
			}
		}
	}
}

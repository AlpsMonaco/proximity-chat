package main

import (
	"context"
	"fmt"
	"log"
	"proximity-chat/message"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGrpcClientSend(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:7789", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := message.NewChatClient(conn)
	cli, err := client.NewNode(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		for {
			msg, err := cli.Recv()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(msg)
		}
	}()
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
		cli.Send(&message.NodeRequest{Msg: "Hello"})
	}
}

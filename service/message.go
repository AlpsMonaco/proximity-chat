package service

import (
	"fmt"
	"proximity-chat/controller"
	"proximity-chat/message"
)

type MessageWriter interface {
	Write(string)
}

type Message struct {
	Writer MessageWriter
	message.UnimplementedChatServer
}

type ServerChatNode struct{ s message.Chat_NewNodeServer }

func (s *ServerChatNode) SendChatMsg(msg string) error {
	return s.s.Send(&message.NodeReply{Msg: msg})
}

func (s *ServerChatNode) RecvChatMsg() (string, error) {
	msg, err := s.s.Recv()
	if err != nil {
		return "", err
	}
	return msg.GetMsg(), nil
}

type ClientChatNode struct{ C message.Chat_NewNodeClient }

func (c *ClientChatNode) SendChatMsg(msg string) error {
	return c.C.Send(&message.NodeRequest{Msg: msg})
}

func (c *ClientChatNode) RecvChatMsg() (string, error) {
	msg, err := c.C.Recv()
	if err != nil {
		return "", err
	}
	return msg.GetMsg(), nil
}

func (m *Message) NewNode(ss message.Chat_NewNodeServer) error {
	head, err := ss.Recv()
	if err != nil {
		m.Writer.Write(fmt.Sprint(err))
		return err
	}
	addr := head.GetMsg()
	if controller.IsChatNodeExist(addr) {
		return nil
	}
	if !controller.AddChatNode(&ServerChatNode{s: ss}, addr) {
		return nil
	}
	err = ss.Send(&message.NodeReply{Msg: "ok"})
	if err != nil {
		return err
	}
	m.Writer.Write("new node " + addr + " has joined")
	for {
		msg, err := ss.Recv()
		if err != nil {
			controller.RemoveNode(addr)
			fmt.Println(err)
			return err
		}
		m.Writer.Write(msg.GetMsg())
	}
}

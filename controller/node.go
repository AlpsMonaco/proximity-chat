package controller

import (
	"sync"
)

type ChatNode interface {
	SendChatMsg(string) error
	RecvChatMsg() (string, error)
}

var nodeMap map[string]ChatNode = make(map[string]ChatNode)
var nodeMapLock sync.RWMutex

func AddChatNode(node ChatNode, addr string) bool {
	nodeMapLock.Lock()
	defer nodeMapLock.Unlock()
	_, ok := nodeMap[addr]
	if !ok {
		nodeMap[addr] = node
		return true
	}
	return false
}

func RemoveNode(addr string) {
	nodeMapLock.Lock()
	defer nodeMapLock.Unlock()
	delete(nodeMap, addr)
}

func IsChatNodeExist(addr string) bool {
	nodeMapLock.RLock()
	defer nodeMapLock.RUnlock()
	_, ok := nodeMap[addr]
	return ok
}

func Publish(s string) {
	nodeMapLock.RLock()
	defer nodeMapLock.RUnlock()
	for _, n := range nodeMap {
		n.SendChatMsg(s)
	}
}

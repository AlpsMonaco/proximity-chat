package controller

import (
	"fmt"
	"sync/atomic"
	"time"
)

type OutputController struct {
	ch   chan string
	stop atomic.Bool
}

func (oc *OutputController) Start() {
	if oc.ch == nil {
		oc.ch = make(chan string, 100)
	}
	oc.stop.Store(false)
	go func() {
		for {
			if oc.stop.Load() {
				fmt.Println("stop")
				return
			}
			select {
			case data := <-oc.ch:
				fmt.Println(data)
			case <-time.After(time.Second):
				continue
			}
		}
	}()
}

func (oc *OutputController) Write(s string) {
	oc.ch <- s
}

func (oc *OutputController) Stop() {
	oc.stop.Store(true)
}

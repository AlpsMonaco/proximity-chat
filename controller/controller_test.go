package controller

import (
	"testing"
	"time"
)

func TestController(t *testing.T) {
	var oc OutputController
	oc.Start()
	for {
		time.Sleep(time.Second)
		oc.Write("Hello")
	}
}

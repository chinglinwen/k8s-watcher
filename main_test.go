package main

import (
	"fmt"
	"testing"
)

func TestNodeList(t *testing.T) {
	err := nodeList()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSetReceiver(t *testing.T) {
	c := &sendconfig{
		touser:  *receiver,
		toparty: *receiverParty,
	}
	SetReceiver("hello")(c)
	if c.touser != "hello" || c.toparty != "" {
		t.Error("set err")
		fmt.Println("c ", c)
		return
	}

	SetReceiver("3")(c)
	if c.touser != "" || c.toparty != "3" {
		t.Error("set err")
		fmt.Println("c ", c)
		return
	}

}

package main

import "testing"

func TestNodeList(t *testing.T) {
	err := nodeList()
	if err != nil {
		t.Error(err)
		return
	}
}

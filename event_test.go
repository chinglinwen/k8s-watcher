package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"

	coreevent "github.com/ericchiang/k8s/apis/events/v1beta1"
)

func init() {
	*receiver = "wenzhenglin"
	*receiverParty = ""
}
func TestFormatAndSend(t *testing.T) {
	b := `{"metadata":{"name":"172.31.90.51.1590a4350b8276d0","generateName":"","namespace":"default","selfLink":"/apis/events.k8s.io/v1beta1/namespaces/default/events/172.31.90.51.1590a4350b8276d0","uid":"0c4361f1-5641-11e9-8fd4-1e5e900bfc2b","resourceVersion":"68210680","generation":0,"creationTimestamp":"2019-04-04T02:48:19+08:00","clusterName":""},"eventTime":{},"reportingController":"","reportingInstance":"","action":"","reason":"ImageGCFailed","regarding":{"kind":"Node","namespace":"","name":"172.31.90.51","uid":"172.31.90.51","apiVersion":"","resourceVersion":"","fieldPath":""},"note":"(combined from similar events): wanted to free 2300926361 bytes, but freed 0 bytes space with errors in image deletion: [rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete f4cb5e83f0a4 (cannot be forced) - image is being used by running container 8b894819fa43, rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete 9b1ea3f29465 (cannot be forced) - image is being used by running container 3f9d6f96e14e, rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete ea6441073322 (cannot be forced) - image is being used by running container 599c217fb553, rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete 9ca888fe33b2 (cannot be forced) - image is being used by running container 64338392316c, rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete 06b8f3008f78 (cannot be forced) - image has dependent child images, rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete 99e59f495ffa (cannot be forced) - image is being used by running container e0dd362fe894, rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete f32a589c03a7 (cannot be forced) - image has dependent child images, rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete 305ee5b8952c (cannot be forced) - image is being used by running container 2c9f47bb72fe, rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete 3712315a6b39 (must be forced) - image is being used by stopped container aa4ff392bfa5, rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete f5e2026b197e (cannot be forced) - image is being used by running container af8cfa576ba8, rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete 3bea3bff0190 (cannot be forced) - image is being used by running container 762ed0bebd58, rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete 6a67380bb061 (cannot be forced) - image is being used by running container 8eabc70bbcf3, rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete a3e95f74984e (cannot be forced) - image is being used by running container ac8aeca453e9, rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete de5d42976498 (cannot be forced) - image is being used by running container 8aeb2fdd9cb3, rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete 9ca6fd371ca6 (cannot be forced) - image is being used by running container de5e15f2cd96, rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete fce289e99eb9 (must be forced) - image is being used by stopped container afdb45fd0648]","type":"Warning","deprecatedSource":{"component":"kubelet","host":"172.31.90.51"},"deprecatedFirstTimestamp":"2019-03-30T13:06:59+08:00","deprecatedLastTimestamp":"2019-04-04T10:03:31+08:00","deprecatedCount":256}`

	e := &coreevent.Event{}
	err := json.Unmarshal([]byte(b), e)
	if err != nil {
		t.Error("unmarshal err", err)
		return
	}
	message := formatevent(e)
	fmt.Println(message)

	reply, err := send(message)
	if err != nil {
		log.Printf("send err: %v\n", err)
	}

	log.Printf("send reply: %v\n", strings.Split(reply, ",")[0])
}

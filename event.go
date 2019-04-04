package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	coreevent "github.com/ericchiang/k8s/apis/events/v1beta1"
)

/*

(*v1beta1.Event)(0xc0002741b0)(metadata:<name:"adm-old-online-c58969bc6-nnknb.1591da15b31f3c4e" generateName:"" namespace:"xindaiquan" selfLink:"/apis/events.k8s.io/v1beta1/namespaces/xindaiquan/events/adm-old-online-c58969bc6-nnknb.1591da15b31f3c4e" uid:"eee587a4-55c2-11e9-8fd4-1e5e900bfc2b" resourceVersion:"67969050" generation:0 creationTimestamp:<seconds:1554263133 nanos:0 > clusterName:"" > eventTime:<> reportingController:"" reportingInstance:"" action:"" reason:"FailedMount" regarding:<kind:"Pod" namespace:"xindaiquan" name:"adm-old-online-c58969bc6-nnknb" uid:"a569e310-55c2-11e9-bf66-cef85e680407" apiVersion:"v1" resourceVersion:"67936340" fieldPath:"" > note:"Unable to mount volumes for pod \"adm-old-online-c58969bc6-nnknb_xindaiquan(a569e310-55c2-11e9-bf66-cef85e680407)\": timeout expired waiting for volumes to attach/mount for pod \"xindaiquan\"/\"adm-old-online-c58969bc6-nnknb\". list of unattached/unmounted volumes=[adm-public adm-bank adm-xindaiyuan adm-common]" type:"Warning" deprecatedSource:<component:"kubelet" host:"172.31.82.85" > deprecatedFirstTimestamp:<seconds:1554263133 nanos:0 > deprecatedLastTimestamp:<seconds:1554272651 nanos:0 > deprecatedCount:71 )

*/
func watchevent() {

start:
	var e coreevent.Event
	watcher, err := client.Watch(context.Background(), "", &e)
	if err != nil {
		log.Println("watch err:", err)
		log.Printf("interval: %v, reconnect: %v\n", interval, reconnect)
		time.Sleep(time.Duration(interval) * time.Second)
		interval = interval * 2
		reconnect += 1
		goto start
	}
	defer watcher.Close()

	log.Println("watcher setup ok, started listening")

	// m := make(map[string]string)
	for {
		e := new(coreevent.Event)
		eventType, err := watcher.Next(e)
		if err != nil {
			// watcher encountered and error, exit or create a new watcher
			log.Println("watcher next err: ", err)
			log.Println("try create new water")
			goto start
		}
		_ = eventType
		// fmt.Println("eventType: ", eventType)

		// ignore normal action
		if e.GetType() == "Normal" {
			log.Println("ignore normal event")
			continue
		}
		// spew.Dump("e", e)

		message := formatevent(e)
		log.Printf("%v", message)

		reply, err := checkandsend(message)
		if err != nil {
			log.Printf("send err: %v\n", err)
		}
		log.Printf("send reply: %v\n", strings.Split(reply, ",")[0])
	}
}

func formatevent(e *coreevent.Event) string {
	// a, _ := json.Marshal(e)
	// fmt.Printf("json: %v", string(a))
	t := `时间: %v
类别: %v
名字: %v/%v
-----
来源: %v
原因: %v
内容: %v`
	now := time.Now().Format("2006-1-2 15:04:05")
	// remove useless suffix
	a := strings.Split(e.Metadata.GetName(), ".")
	name := strings.Join(a[:len(a)-1], ".")

	msg := e.GetNote()
	if len(msg) > 300 {
		msg = msg[:300] + "... (omited)"
	}
	kind := e.GetRegarding().GetKind()
	return fmt.Sprintf(t, now, e.GetType(), e.Metadata.GetNamespace(), name, kind, e.GetReason(), msg)
}

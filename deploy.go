package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	deployv1beta1 "github.com/ericchiang/k8s/apis/extensions/v1beta1"
)

func watchdeploy() {

start:
	var e deployv1beta1.Deployment
	watcher, err := client.Watch(context.Background(), "", &e)
	if err != nil {
		log.Println("watch deploy err:", err)
		log.Printf("interval: %v, deploy reconnect: %v\n", interval, reconnect)
		time.Sleep(time.Duration(interval) * time.Second)
		interval = interval * 2
		reconnect += 1
		goto start
	}
	defer watcher.Close()

	log.Println("deploy watcher setup ok, started listening")

	// m := make(map[string]string)
	for {
		e := new(deployv1beta1.Deployment)
		deployType, err := watcher.Next(e)
		if err != nil {
			// watcher encountered and error, exit or create a new watcher
			log.Println("watcher next err: ", err)
			log.Println("try create new water")
			goto start
		}
		// _ = deployType
		// fmt.Println("deployType: ", deployType)

		// ignore normal action
		an := e.GetMetadata().GetAnnotations()
		person := an[*annotationName]
		if person == "" {
			continue
		}

		// spew.Dump(e.GetStatus())

		// a, _ := json.Marshal(e)
		// fmt.Printf("json: %v", string(a))

		// not ready yet
		if deployType == "MODIFIED" {
			if e.GetStatus().GetUnavailableReplicas() != 0 || e.GetStatus().GetReplicas() == 0 {
				continue
			}
			// deployType="ready"
		}

		message := formatdeploy(e, deployType)
		fmt.Printf("%v\n", message)

		reply, err := checkandsend(message, SetReceiver(person))
		if err != nil {
			log.Printf("send err: %v\n", err)
		}
		log.Printf("send reply: %v\n", strings.Split(reply, ",")[0])
	}
}

func formatdeploy(e *deployv1beta1.Deployment, deployType string) string {
	// a, _ := json.Marshal(e)
	// fmt.Printf("json: %v", string(a))
	t := `时间: %v
名字: %v
部署状态: %v%v`
	now := time.Now().Format("2006-1-2 15:04:05")
	name := e.Metadata.GetNamespace() + "/" + e.Metadata.GetName()

	ready := e.GetStatus().GetReadyReplicas()
	want := e.GetStatus().GetReplicas()
	var msg string

	if deployType == "MODIFIED" {
		if ready == want {
			msg = "\n-----\n内容: ready"
		} else {
			msg = fmt.Sprintf("\n-----\n内容: replicas want %v, ready %v", want, ready)
		}
	}

	return fmt.Sprintf(t, now, name, strings.ToLower(deployType), msg)
}

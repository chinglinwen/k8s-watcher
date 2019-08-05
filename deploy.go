package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	deployv1beta1 "github.com/ericchiang/k8s/apis/extensions/v1beta1"
)

func watchDeploy() {

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

		// ignore old event
		ts := e.GetMetadata().GetCreationTimestamp()
		t := time.Unix(ts.GetSeconds(), int64(ts.GetNanos()))
		now := time.Now()
		timeRange := 1 * time.Minute
		if t.Add(timeRange).Before(now) {
			log.Printf("ignore old event than %v, created: %v, now: %v\n\n",
				timeRange,
				t.Format(layout),
				now.Format(layout))
			continue
		}

		message := formatdeploy(e, deployType)
		fmt.Printf("%v\n", message)

		// send to ops
		reply, err := checkandsend(message)
		if err != nil {
			log.Printf("send to ops err: %v\n", err)
		}
		log.Printf("send to ops reply: %v\n", strings.Split(reply, ",")[0])

		// if set, send to specific receiver ( service owner )
		an := e.GetMetadata().GetAnnotations()
		person := an[*annotationName]
		if person == "" {
			continue
		}
		reply, err = checkandsend(message, SetReceiver(person))
		if err != nil {
			log.Printf("send to %v err: %v\n", person, err)
		}
		log.Printf("send to %v reply: %v\n", person, strings.Split(reply, ",")[0])
	}
}

func formatdeploy(e *deployv1beta1.Deployment, deployType string) string {
	// a, _ := json.Marshal(e)
	// fmt.Printf("json: %v", string(a))
	t := `名字: %v
部署状态: %v%v`

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

	return fmt.Sprintf(t, name, strings.ToLower(deployType), msg)
}

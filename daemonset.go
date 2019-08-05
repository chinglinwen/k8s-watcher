package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	dsv1beta1 "github.com/ericchiang/k8s/apis/extensions/v1beta1"
)

func watchDaemon() {

start:
	var e dsv1beta1.DaemonSet
	watcher, err := client.Watch(context.Background(), "", &e)
	if err != nil {
		log.Println("watch daemonset err:", err)
		log.Printf("interval: %v, daemonset reconnect: %v\n", interval, reconnect)
		time.Sleep(time.Duration(interval) * time.Second)
		interval = interval * 2
		reconnect += 1
		goto start
	}
	defer watcher.Close()

	log.Println("deploy watcher setup ok, started listening")

	// m := make(map[string]string)
	for {
		e := new(dsv1beta1.DaemonSet)
		dsType, err := watcher.Next(e)
		if err != nil {
			// watcher encountered and error, exit or create a new watcher
			log.Println("watcher next err: ", err)
			log.Println("try create new water")
			goto start
		}
		// _ = deployType
		// fmt.Println("deployType: ", deployType)

		// spew.Dump(e.GetStatus())

		// a, _ := json.MarshalIndent(e, "", "  ")
		// fmt.Printf("json: %v", string(a))

		// not ready yet
		if dsType == "MODIFIED" {
			if e.GetStatus().GetNumberUnavailable() != 0 {
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

		message := formatds(e, dsType)
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

func formatds(e *dsv1beta1.DaemonSet, deployType string) string {
	// a, _ := json.Marshal(e)
	// fmt.Printf("json: %v", string(a))
	t := `名字: %v
Daemonset状态: %v%v`

	name := e.Metadata.GetNamespace() + "/" + e.Metadata.GetName()

	ready := e.GetStatus().GetNumberReady()
	want := e.GetStatus().GetDesiredNumberScheduled()
	var msg string

	if deployType == "MODIFIED" {
		if ready == want {
			msg = "\n-----\n内容: ready"
		} else {
			msg = fmt.Sprintf("\n-----\n内容: daemonset want %v, ready %v", want, ready)
		}
	}

	return fmt.Sprintf(t, name, strings.ToLower(deployType), msg)
}

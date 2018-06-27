package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	"github.com/ghodss/yaml"
)

var (
	wechatNotifyURL = flag.String("w", "http://localhost:8001", "wechat notify service url")
	receiver        = flag.String("r", "wenzhenglin", "default receiver")

	// init interval
	interval = 2

	// number of reconnect times
	reconnect = 0

	starttime = time.Now()
	startsend bool
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	//client, err := k8s.NewInClusterClient()
	client, err := loadClient(*kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	var nodes corev1.NodeList
	if err = client.List(context.Background(), "", &nodes); err != nil {
		log.Fatal(err)
	}
	for _, node := range nodes.Items {
		log.Printf("name=%q schedulable=%t\n", *node.Metadata.Name, !*node.Spec.Unschedulable)
	}

start:
	// Watch configmaps in the "kube-system" namespace
	var pod corev1.Pod
	watcher, err := client.Watch(context.Background(), "", &pod)
	if err != nil {
		log.Println("watch err:", err)
		log.Printf("interval: %v, reconnect: %v\n", interval, reconnect)
		time.Sleep(time.Duration(interval) * time.Second)
		interval = interval * interval
		reconnect += 1
		goto start
	}
	defer watcher.Close()

	log.Println("watcher setup, started listening")

	m := make(map[string]string)
	for {
		pod := new(corev1.Pod)
		eventType, err := watcher.Next(pod)
		if err != nil {
			// watcher encountered and error, exit or create a new watcher
			log.Println("watcher next err: ", err)
			log.Println("try create new water")
			goto start
		}

		podname := *(pod.Metadata.Name)
		preStatus := m[podname]

		if len(pod.Status.ContainerStatuses) == 0 {
			continue
		}
		status := pod.Status.ContainerStatuses[0]
		if *status.Ready {
			t, _ := status.State.Running.StartedAt.MarshalJSON()
			log.Println("ok", *status.Name, "running: ready is", *status.Ready, "start at:", string(t), "restart: ", *status.RestartCount)
			continue
		}
		state, _ := json.Marshal(status.State)
		message := fmt.Sprintln("err", *status.Name, string(state), *status.Ready, "restart: ", *status.RestartCount)
		log.Printf("%v", message)

		reply, err := checkandsend(message)
		if err != nil {
			log.Printf("send err: %v\n", err)
		}
		log.Printf("send reply: %v\n", strings.TrimSpace(reply))

		if *(pod.Status.Phase) == preStatus && eventType == "MODIFIED" {
			continue
		}

		message = fmt.Sprintln("event ", eventType, podname, *(pod.Status.Phase), *(pod.Status.Message), *(pod.Status.Reason))
		log.Printf("%v", message)
		reply, err = checkandsend(message)
		if err != nil {
			log.Printf("send err: %v\n", err)
		}
		log.Printf("send reply: %v\n", strings.TrimSpace(reply))

		m[podname] = *(pod.Status.Phase)
	}
}

func loadClient(kubeconfigPath string) (*k8s.Client, error) {
	data, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("read kubeconfig: %v", err)
	}

	// Unmarshal YAML into a Kubernetes config object.
	var config k8s.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unmarshal kubeconfig: %v", err)
	}
	return k8s.NewClient(&config)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

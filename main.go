package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	"github.com/pkg/errors"

	"github.com/ghodss/yaml"
)

var (
	wechatNotifyURL = flag.String("w", "http://localhost:8001", "wechat notify service url")
	receiver        = flag.String("r", "", "default receiver")
	receiverParty   = flag.String("party", "3", "default receiver party")
	expire          = flag.String("e", "10m", "default expire time duration")

	annotationName = flag.String("a", "wechat", "deploy annotation name")

	// init interval
	interval = 2

	// number of reconnect times
	reconnect = 0

	starttime = time.Now()
	startsend bool

	client *k8s.Client
)

func init() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	//client, err := k8s.NewInClusterClient()
	var err error
	client, err = loadClient(*kubeconfig)
	if err != nil {
		client, err = k8s.NewInClusterClient()
	}
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	if *receiver == "" && *receiverParty == "" {
		log.Println("args receiver and party is empty")
		return
	}
	go watchdeploy()
	nodeList()
	watchevent()
}

func nodeList() error {
	var nodes corev1.NodeList
	if err := client.List(context.Background(), "", &nodes); err != nil {
		return errors.Wrap(err, "client list")
	}
	for _, node := range nodes.Items {
		log.Printf("name=%q schedulable=%t\n", *node.Metadata.Name, !*node.Spec.Unschedulable)
	}
	return nil
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

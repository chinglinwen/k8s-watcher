package main

import (
	"context"
	"fmt"

	corev1 "github.com/ericchiang/k8s/apis/core/v1"
)

func Pod(name string) (pod corev1.Pod, err error) {
	err = client.Get(context.Background(), "", name, &pod)
	if err != nil {
		err = fmt.Errorf("get pod err %v", err)
		return
	}
	return
}

func PodListAll() (pods []*corev1.Pod, err error) {
	return PodList("")
}

func CheckPodExist(ns, pod string) (exist bool, err error) {
	pods, err := PodList(ns)
	if err != nil {
		err = fmt.Errorf("PodList err %v", err)
		return
	}
	exist = IsPodExist(pods, pod)
	return
}

func IsPodExist(pods []*corev1.Pod, name string) bool {
	if len(pods) == 0 {
		return false
	}
	for _, v := range pods {
		if v.GetMetadata().GetName() == name {
			return true
		}
	}
	return false
}

func PodList(ns string) (pods []*corev1.Pod, err error) {
	var slist corev1.PodList
	err = client.List(context.Background(), ns, &slist)
	if err != nil {
		err = fmt.Errorf("get secret err %v", err)
		return
	}
	pods = slist.GetItems()
	return
}

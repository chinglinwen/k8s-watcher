package main

import (
	"fmt"
	"testing"
)

func TestPodList(t *testing.T) {
	ss, err := PodList("xindaiquan")
	if err != nil {
		t.Error("PodList err", err)
		return
	}
	// b, _ := json.MarshalIndent(ss, "", "  ")
	// fmt.Println(string(b))
	for _, v := range ss {
		fmt.Println(v.GetMetadata().GetNamespace(), v.GetMetadata().GetName())
	}
}

func TestCheckPodExist(t *testing.T) {
	exist, err := CheckPodExist("xindaiquan", "adm-old-online-5749bcbd7b-cxzr6")
	if err != nil {
		t.Error("CheckPodExist err", err)
		return
	}
	// pod may change overtime, here we just print result
	fmt.Printf("exist: %v\n", exist)
}

func TestPodListAll(t *testing.T) {
	ss, err := PodListAll()
	if err != nil {
		t.Error("PodListAll err", err)
		return
	}
	// b, _ := json.MarshalIndent(ss, "", "  ")
	// fmt.Println(string(b))
	for _, v := range ss {
		fmt.Println(v.GetMetadata().GetNamespace(), v.GetMetadata().GetName())
	}
}

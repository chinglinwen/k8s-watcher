package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	coreevent "github.com/ericchiang/k8s/apis/events/v1beta1"
)

var connFluentd net.Conn

func init() {
	flag.Parse()
	var err error
	connFluentd, err = net.Dial("udp", *fluentd)
	if err != nil {
		log.Fatal(fmt.Sprintf("fluentd:%s ,err:%s", *fluentd, err))
	}
}

func consumerFluentd(e *coreevent.Event) (err error) {
	buf, err := json.MarshalIndent(e, "", " ")
	if err != nil {
		err = fmt.Errorf("marshal event err", err)
		return
	}
	log.Printf("event: %v\n", string(buf))

	ts := e.GetMetadata().GetCreationTimestamp()
	tsStr := time.Unix(ts.GetSeconds(), int64(ts.GetNanos())).Format("2006-01-02T15:04:05Z07:00")

	buf = bytes.Replace(buf, []byte(`"eventTime":`), []byte(fmt.Sprintf(`"@timestamp":"%s","eventTime":`, tsStr)), 1)
	_, err = connFluentd.Write(buf)

	return
}

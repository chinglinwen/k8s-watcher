package main

import (
	coreevent "github.com/ericchiang/k8s/apis/events/v1beta1"
	"net"
	"log"
	"encoding/json"
	"github.com/pkg/errors"
	"time"
	"fmt"
	"flag"
	"bytes"
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

func consumerFluentd(e *coreevent.Event) error {
	buf, err := json.Marshal(e)
	if err != nil {
		return errors.New(err.Error())
	}

	ts := e.GetMetadata().GetCreationTimestamp()
	tsStr := time.Unix(ts.GetSeconds(), int64(ts.GetNanos())).Format("2006-01-02T15:04:05Z07:00")

	buf = bytes.Replace(buf, []byte(`"eventTime":`), []byte(fmt.Sprintf(`"@timestamp":"%s","eventTime":`, tsStr)), 1)
	connFluentd.Write(buf)

	return nil
}

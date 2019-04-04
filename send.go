package main

import (
	"strings"
	"time"

	resty "gopkg.in/resty.v1"
)

func send(message string) (reply string, err error) {
	r := strings.NewReplacer("\"", " ", "{", "", "}", "")
	message = r.Replace(message)

	resp, e := resty.R().
		SetQueryParams(map[string]string{
			"user":    *receiver,
			"toparty": *receiverParty,
			"content": message,
			"expire":  *expire,
		}).
		Get(*wechatNotifyURL)

	if e != nil {
		err = e
		return
	}
	reply = string(resp.Body())
	return
}

func checkandsend(message string) (reply string, err error) {
	if !startsend {
		if !time.Now().After(starttime.Add(5 * time.Second)) {
			return "skip send at start time", nil
		}
	}
	return send(message)
}

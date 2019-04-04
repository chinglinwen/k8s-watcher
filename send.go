package main

import (
	"regexp"
	"strings"
	"time"

	resty "gopkg.in/resty.v1"
)

type sendconfig struct {
	touser  string
	toparty string
}

type sendoption func(*sendconfig)

// both touser and toparty
func SetReceiver(receiver string) sendoption {
	return func(c *sendconfig) {
		if regexp.MustCompile(`^[0-9]+$`).MatchString(receiver) {
			c.toparty = receiver
			c.touser = ""
			return
		}
		c.touser = receiver
		c.toparty = ""
	}
}

func send(message string, options ...sendoption) (reply string, err error) {
	c := &sendconfig{
		touser:  *receiver,
		toparty: *receiverParty,
	}
	for _, option := range options {
		option(c)
	}

	r := strings.NewReplacer("\"", " ", "{", "", "}", "")
	message = r.Replace(message)

	resp, e := resty.R().
		SetQueryParams(map[string]string{
			"user":    c.touser,
			"toparty": c.toparty,
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

func checkandsend(message string, options ...sendoption) (reply string, err error) {
	if !startsend {
		if !time.Now().After(starttime.Add(5 * time.Second)) {
			return "skip send at start time", nil
		}
	}
	return send(message, options...)
}

package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mediocregopher/radix/v3"
)

var channel = "0000000000001"

type indmessage struct {
	UN   string
	Msg  string
	To   string
	Time time.Time
}

type tryme struct {
	Stuff indmessage
}

//doing it like this so that when Redis is implemented in a bit I don't have to change a lot of code
func pushToQueue(topic string, message map[string]string) {
	var structmsg indmessage
	var rdb = getClient()
	structmsg.Time = time.Now()
	structmsg.Msg = message["message"]
	structmsg.UN = message["name"]
	structmsg.To = message["for"]

	tryme, err := json.Marshal(structmsg)
	if err != nil {
		println(err)
		panic(err)
	}
	hashme := asSha256(string(tryme))
	sendme := make(map[string]string)
	sendme[hashme] = string(tryme)
	final, err := json.Marshal(sendme)
	rdb.Do(radix.Cmd(nil, "publish", "chat"+topic+channel+structmsg.UN, string(final))) //, //string(tryme)))
	rdb.Do(radix.Cmd(nil, "lpush", "chat"+topic+channel+structmsg.UN, string(final)))
}

//Index handles messages going to the index
func Index(w http.ResponseWriter, r *http.Request) {
	Messaging(w, r)
	//return the web page here
}

//Messaging handles all messaging requests, pushing them to either fork, message, or
func Messaging(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	var msg map[string]string
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	switch r.URL.Path {
	case "/fork":
		pushToQueue("fork", msg)
	case "/message":
		pushToQueue("message", msg)
	case "/history":
		CheckMessages(w, msg)
	}

}

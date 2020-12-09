package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mediocregopher/radix/v3"
)

var channel = "0000000000001"

type indmessage struct {
	UN  string
	Msg string
	To  string
	// time    time.Time
}

type tryme struct {
	Stuff indmessage
}

var ctx = context.Background()

var q = map[string][]indmessage{
	"message": make([]indmessage, 1000),
	"fork":    make([]indmessage, 1000)}

func getClient() *radix.Pool {
	pool, err := radix.NewPool("tcp", "localhost:6379", 10)
	if err != nil {
		// handle error
	}
	return pool
}

//doing it like this so that when Redis is implemented in a bit I don't have to change a lot of code
func pushToQueue(topic string, message map[string]string) {
	var structmsg indmessage
	var rdb = getClient()
	structmsg.Msg = message["message"]
	structmsg.UN = message["name"]
	structmsg.To = message["for"]

	tryme, err := json.Marshal(structmsg)
	if err != nil {
		//stuff
	}
	rdb.Do(radix.Cmd(nil, "publish", topic+channel+structmsg.UN, string(tryme))) //, //string(tryme)))
	rdb.Do(radix.Cmd(nil, "lpush", topic+channel+structmsg.UN, string(tryme)))
}

//Index handles messages going to the index
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
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

	pushToQueue(strings.Split(r.URL.Path, "/")[2], msg)
}

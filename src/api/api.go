package api

import (
	"context"
	"encoding/json"
	"net/http"

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
	// err := rdb.Set(topic, message, 0).Err() //set key to ????, and set value to message
	// if err != nil {
	// 	panic(err)
	// }
	var structmsg indmessage
	var rdb = getClient()
	structmsg.Msg = message["message"]
	structmsg.UN = message["name"]
	structmsg.To = message["for"]

	tryme, err := json.Marshal(structmsg)
	if err != nil {
		//stuff
	}
	// println(rdb.NumAvailConns())
	rdb.Do(radix.Cmd(nil, "publish", channel+structmsg.UN, string(tryme))) //, //string(tryme)))
	rdb.Do(radix.Cmd(nil, "lpush", channel+structmsg.UN, string(tryme)))
	// print("CAUGHT THIS")
	// println(stuf)
	// println("End catch")
	// println()
}

//Index handles messages going to the index
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	//return the web page here

}

//Fork puts a message on the fork queue
func Fork(w http.ResponseWriter, r *http.Request) {
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
	pushToQueue("fork", msg)
}

//Message puts a message on the message queue
func Message(w http.ResponseWriter, r *http.Request) {
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
	// println(msg["name"])
	pushToQueue("message", msg)
}

//Queue returns the value of the queues
func Queue(w http.ResponseWriter, r *http.Request) {
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
	var rdb = getClient()
	rdb.Do(radix.Cmd(nil, "get", channel+msg["name"]))
}

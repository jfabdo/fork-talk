package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mediocregopher/radix/v3"
)

type indmessage struct {
	un  string
	msg string
	// time    time.Time
}

type tryme struct {
	Stuff []indmessage
}

var ctx = context.Background()

var q = map[string][]indmessage{
	"message": make([]indmessage, 1000),
	"fork":    make([]indmessage, 1000)}

func getClient() *radix.Pool {
	pool, err := radix.NewPool("tcp", "127.0.0.1:6379", 10)
	if err != nil {
		// handle error
	}
	return pool
}

var rdb = getClient()

//doing it like this so that when Redis is implemented in a bit I don't have to change a lot of code
func pushToQueue(topic string, message map[string]string) {
	// err := rdb.Set(topic, message, 0).Err() //set key to ????, and set value to message
	// if err != nil {
	// 	panic(err)
	// }
	var structmsg indmessage

	structmsg.msg = message["message"]
	structmsg.un = message["name"]

	q[topic] = append(q[topic], structmsg)
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
	// var things tryme
	// things.Stuff = q["un"]
	// r, err := json.Marshal(things)
	// if err != nil {
	// 	http.Error(w, "Method is not supported.", http.StatusNotFound)
	// 	return
	// }
	for i := 0; i < len(q); i++ {
		r, err := json.Marshal(q["message"])
		if err != nil {
			http.Error(w, "Fuck", http.StatusNotFound)
		}
		s, err := json.Marshal(q["message"])
		if err != nil {
			http.Error(w, "Fuck", http.StatusNotFound)
		}
		io.WriteString(w, string(r))
		io.WriteString(w, string(s))
	}

	// return q
}

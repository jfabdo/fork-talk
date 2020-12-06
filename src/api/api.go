package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

type indmessage struct {
	name    string
	message string
	time    time.Time
}

var ctx = context.Background()

var q = map[string][]indmessage{
	"message": make([]indmessage, 1000),
	"fork":    make([]indmessage, 1000)}

func getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

var rdb = getClient()

//doing it like this so that when Redis is implemented in a bit I don't have to change a lot of code
func pushToQueue(topic string, message indmessage) {
	// err := rdb.Set(topic, message, 0).Err() //set key to ????, and set value to message
	// if err != nil {
	// 	panic(err)
	// }
	q[topic] = append(q[topic], message)
}

//Index handles messages going to the index
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	//return the web page here

}

func fork(w http.ResponseWriter, r *http.Request) {

	var msg indmessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	pushToQueue("fork", msg)
}

func message(w http.ResponseWriter, r *http.Request) {
	var msg indmessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	pushToQueue("message", msg)
}

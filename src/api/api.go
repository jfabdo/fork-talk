package api

import (
	"context"
	"net/http"

	"github.com/go-redis/redis"
)

var ctx = context.Background()

func getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

//Index handles messages going to the index
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
}

func fork(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	rdb := getClient()

	err := rdb.Set("key", "value", 0).Err() //set key to ????, and set value to message
	if err != nil {
		panic(err)
	}
}

func message(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	rdb := getClient()

	err := rdb.Set("key", "value", 0).Err() //set key to ????, and set value to message
	if err != nil {
		panic(err)
	}
}

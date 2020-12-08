package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jfabdo/fork-talk/src/api"
)

func main() {
	http.HandleFunc("/", api.Index)
	http.HandleFunc("/fork", api.Fork)
	http.HandleFunc("/message", api.Message)
	http.HandleFunc("/queue", api.Queue)
	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

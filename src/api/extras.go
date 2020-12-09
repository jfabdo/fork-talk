package api

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mediocregopher/radix/v3"
)

var ctx = context.Background()

func asSha256(o interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func getClient() *radix.Pool {
	pool, err := radix.NewPool("tcp", "localhost:6379", 10)
	if err != nil {
		println(err)
		panic(err)
	} else {
		return pool
	}
}

// CheckMessages checks the incoming list of messages and makes sure that
func CheckMessages(w http.ResponseWriter, messages map[string]string) {
	var rdb = getClient()

	// var messagedict map[string]string
	var messagelist []string
	err := rdb.Do(radix.Cmd(&messagelist, "lrange", "chat"+"message"+channel+messages["name"], "0", "10"))
	if err != nil {
		println(err)
		panic(err)
	}
	for i := range messages {
		for j := range messagelist {
			if messages[i] == messagelist[j] {
				copy(messagelist[j:], messagelist[j+1:])
				messagelist = messagelist[:len(messagelist)-1]
			}
		}
	}

	rspns, err := json.Marshal(&messagelist)
	if err != nil {
		//do somethign
	}

	fmt.Fprintf(w, string(rspns))
}

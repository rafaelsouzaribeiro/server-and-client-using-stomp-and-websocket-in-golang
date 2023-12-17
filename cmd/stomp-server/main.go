package main

import (
	"github.com/go-stomp/stomp/server"
)

func main() {

	err := server.ListenAndServe("localhost:61613")

	if err != nil {
		panic(err)
	}

}

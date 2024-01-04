package main

import (
	"github.com/go-stomp/stomp/server"
)

func main() {

	err := server.ListenAndServe("springboot:8080")

	if err != nil {
		panic(err)
	}

}

package main

import (
	"fmt"
	"github.com/louis296/p2p-server/pkg/server"
	"net/http"
)

func WebSocketHandler(conn *server.WebSocketConn, request *http.Request) {
	fmt.Println("On Open ", request)
	conn.On("message", func(message []byte) {
		fmt.Println(string(message))
		err := conn.Send("reply: " + string(message))
		if err != nil {
			fmt.Println(err.Error())
		}
	})
	conn.On("close", func(code int, text string) {
		fmt.Println("connection close with ", code, text)
	})
}

package server

import (
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

type P2PServerConfig struct {
	Host          string
	Port          int
	CertFile      string
	KeyFile       string
	HTMLRoot      string
	WebSocketPath string
}

func DefaultConfig() P2PServerConfig {
	return P2PServerConfig{
		Host:          "0.0.0.0",
		Port:          8000,
		HTMLRoot:      "html",
		WebSocketPath: "/ws",
	}
}

type P2PServer struct {
	handleWebSocket func(we *WebSocketConn, request *http.Request)
	upgrader        websocket.Upgrader
}

func NewP2PServer(wsHandler func(ws *WebSocketConn, request *http.Request)) *P2PServer {
	server := &P2PServer{
		handleWebSocket: wsHandler,
	}
	server.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return server
}

func (server *P2PServer) handleWebSocketRequest(writer http.ResponseWriter, request *http.Request) {
	responseHeader := http.Header{}
	socket, err := server.upgrader.Upgrade(writer, request, responseHeader)
	if err != nil {
		//log
	}
	wsTransport := NewWebSocketConn(socket)
	server.handleWebSocket(wsTransport, request)
	wsTransport.ReadMessage()
}

func (server *P2PServer) Bind(conf P2PServerConfig) {
	http.HandleFunc(conf.WebSocketPath, server.handleWebSocketRequest)
	http.Handle("/", http.FileServer(http.Dir(conf.HTMLRoot)))
	//log

	//start as http
	err := http.ListenAndServe(conf.Host+":"+strconv.Itoa(conf.Port), nil)
	if err != nil {
		panic("server start error")
	}
	//start as https
	//panic(http.ListenAndServeTLS(conf.Host+":"+strconv.Itoa(conf.Port),conf.CertFile,conf.KeyFile,nil))
}

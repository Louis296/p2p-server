package server

import (
	"errors"
	"github.com/chuckpreslar/emission"
	"github.com/gorilla/websocket"
	"github.com/louis296/p2p-server/pkg/util"
	"net"
	"sync"
	"time"
)

const pingPeriod = 5 * time.Second

type WebSocketConn struct {
	emission.Emitter
	socket *websocket.Conn
	mutex  *sync.Mutex
	closed bool
}

func NewWebSocketConn(socket *websocket.Conn) *WebSocketConn {
	var conn WebSocketConn
	conn.Emitter = *emission.NewEmitter()
	conn.socket = socket
	conn.mutex = new(sync.Mutex)
	conn.closed = false
	conn.socket.SetCloseHandler(func(code int, text string) error {
		//log
		conn.Emit("close", code, text)
		conn.closed = true
		return nil
	})
	return &conn
}

func (conn *WebSocketConn) ReadMessage() {
	in := make(chan []byte)
	stop := make(chan struct{})
	pingTicker := time.NewTicker(pingPeriod)

	var c = conn.socket
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				//log
				if c, ok := err.(*websocket.CloseError); ok {
					conn.Emit("close", c.Code, c.Text)
				} else {
					if c, ok := err.(*net.OpError); ok {
						conn.Emit("close", 1008, c.Error())
					}
				}
				close(stop)
				break
			}
			in <- message
		}
	}()

	for {
		select {
		case _ = <-pingTicker.C:
			//log
			heartPackage := map[string]interface{}{
				"type": "heartPackage",
				"data": "",
			}
			if err := conn.Send(util.Marshal(heartPackage)); err != nil {
				//log
				pingTicker.Stop()
				return
			}
		case message := <-in:
			//log
			conn.Emit("message", message)
		case <-stop:
			return
		}
	}
}

func (conn *WebSocketConn) Send(message string) error {
	//log
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	if conn.closed {
		return errors.New("websocket: write closed")
	}
	return conn.socket.WriteMessage(websocket.TextMessage, []byte(message))
}

func (conn *WebSocketConn) Close() {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	if !conn.closed {
		//log
		conn.socket.Close()
		conn.closed = true
	} else {
		//log
	}
}

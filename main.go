package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	connectionInitMsg      = "connection_init"      // Client -> Server
	connectionTerminateMsg = "connection_terminate" // Client -> Server
	startMsg               = "start"                // Client -> Server
	stopMsg                = "stop"                 // Client -> Server
	connectionAckMsg       = "connection_ack"       // Server -> Client
	connectionErrorMsg     = "connection_error"     // Server -> Client
	dataMsg                = "data"                 // Server -> Client
	errorMsg               = "error"                // Server -> Client
	completeMsg            = "complete"             // Server -> Client
	//connectionKeepAliveMsg = "ka"                 // Server -> Client  TODO: keepalives
)

type operationMessage struct {
	Payload json.RawMessage `json:"payload,omitempty"`
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
}

func wsConnect(url string) *websocket.Conn {
	c, _, err := websocket.DefaultDialer.Dial(strings.Replace(url, "http://", "ws://", -1), nil)
	if err != nil {
		panic(err)
	}
	return c
}

func writeRaw(conn *websocket.Conn, msg string) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		panic(err)
	}
}

func readOp(conn *websocket.Conn) operationMessage {
	var msg operationMessage
	if err := conn.ReadJSON(&msg); err != nil {
		panic(err)
	}
	return msg
}

func main() {

	c := wsConnect("http://localhost:4466/")
	defer c.Close()

	c.WriteJSON(&operationMessage{Type: connectionInitMsg})
	log.Println(readOp(c).Type)
	c.WriteJSON(&operationMessage{
		Type: startMsg,
		ID:   "test_1",
		Payload: json.RawMessage(`{"query": "subscription {
			post {
			  node {
				id
				title
			  }
			}
		  }"}`),
	})

	msg := readOp(c)
	log.Println(msg.Type)
	log.Println(msg.ID)
	log.Println(string(msg.Payload))

}

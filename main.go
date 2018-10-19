package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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

type PostSubscriptionResponse struct {
	Data struct {
		Publish struct {
			ID          string    `json:"id"`
			CreatedAt   time.Time `json:"createdAt"`
			UpdatedAt   time.Time `json:"updatedAt"`
			IsPublished bool      `json:"isPublished"`
			Title       string    `json:"title"`
			Content     string    `json:"content"`
			Author      struct {
				Name string `json:"name"`
			} `json:"author"`
		} `json:"publish"`
	} `json:"data"`
}

func wsConnect(url string) *websocket.Conn {
	headers := make(http.Header)
	headers.Add("Sec-Websocket-Protocol", "graphql-ws")
	c, _, err := websocket.DefaultDialer.Dial(url, headers)

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

func BytesToString(data []byte) string {
	return string(data[:])
}

func main() {

	c := wsConnect("ws://localhost:4466/ws")
	defer c.Close()

	c.WriteJSON(&operationMessage{Type: connectionInitMsg})
	log.Println(readOp(c).Type)
	c.WriteJSON(&operationMessage{
		Type:    startMsg,
		ID:      "test_1",
		Payload: json.RawMessage(`{"query": "subscription { post { node { id title } } }"}`),
	})

	msg := readOp(c)
	log.Println(msg.Type)
	log.Println(msg.ID)
	//log.Println(string(json.RawMessage(msg.Payload)))
	postSubscriptionResponse := new(PostSubscriptionResponse)
	s := json.Unmarshal(msg.Payload, postSubscriptionResponse)
	log.Println(s)
	select {}

}

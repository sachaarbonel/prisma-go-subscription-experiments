package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Session struct {
	ws *websocket.Conn
	//errChan chan error
}

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
		Post struct {
			Node struct {
				ID    string `json:"id"`
				Title string `json:"title"`
			} `json:"node"`
		} `json:"post"`
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

func (s *Session) ReadOp() operationMessage {
	var msg operationMessage
	if err := s.ws.ReadJSON(&msg); err != nil {
		panic(err)
	}
	return msg
}

func (s *Session) Subscribe(query string) {
	// refactor from this
	s.ws.WriteJSON(&operationMessage{
		Type:    startMsg,
		ID:      "test_1",
		Payload: json.RawMessage(query),
	})

	// to this. In a function that return chan
	// example promise := ps.Sub(query)
	// select {
	// 	fmt.PrintLn(prosmise.(string))
	// }

	msg := s.ReadOp()
	log.Println(msg.Type)
	log.Println(msg.ID)
	rawPayload := json.RawMessage(msg.Payload)
	//log.Println(rawPayload)
	str := string(rawPayload[:])
	log.Println(str)
	//select {}

}

func main() {

	c := wsConnect("ws://localhost:4466")
	defer c.Close()

	session := &Session{
		ws: c,
	}
	session.ws.WriteJSON(&operationMessage{Type: connectionInitMsg})
	log.Println(session.ReadOp().Type)

	query := string(`{"query": "subscription { post { node { id title } } }"}`)
	session.Subscribe(query)
	//select {}

}

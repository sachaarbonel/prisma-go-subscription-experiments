package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Session struct {
	ws      *websocket.Conn
	errChan chan error
}

type futurePayload chan string

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

func (s *Session) ReadOp() (operationMessage, error) {
	var msg operationMessage
	// if err := s.ws.ReadJSON(&msg); err != nil {
	// 	panic(err)
	// }
	err := s.ws.ReadJSON(&msg)
	return msg, err
}

func (s *Session) Subscribe(query string) futurePayload {

	channel := make(futurePayload)

	s.ws.WriteJSON(&operationMessage{
		Type:    startMsg,
		ID:      "test_1", // Do I need to generate a random ID here
		Payload: json.RawMessage(query),
	})

	for {
		//defer close(channel)
		msg, err := s.ReadOp()
		if err != nil {
			s.errChan <- err
			//TODOD :switch
		}

		log.Println(msg.Type)
		log.Println(msg.ID)
		rawPayload := json.RawMessage(msg.Payload)
		strPayload := string(rawPayload[:])
		log.Println(strPayload)

		channel <- strPayload

	}

	return channel
}

func main() {

	c := wsConnect("ws://localhost:4466")
	defer c.Close()

	session := &Session{
		ws: c,
	}
	session.ws.WriteJSON(&operationMessage{Type: connectionInitMsg})
	msg, _ := session.ReadOp()
	log.Println(msg.Type)

	query := string(`{"query": "subscription { post { node { id title } } }"}`)
	session.Subscribe(query) // goroutine here ??

}

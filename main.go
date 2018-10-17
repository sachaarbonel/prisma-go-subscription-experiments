package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseStruct struct {
	Data struct {
		Users []struct {
			ID string `json:"id"`
		} `json:"users"`
	} `json:"data"`
}

type Payload struct {
	Query string `json:"query"`
}

func main() {

	data := Payload{
		Query: "query {users {id }}",
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:4466/", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Origin", "http://localhost:4466")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	//body, _ = ioutil.ReadAll(resp.Body)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()

	fmt.Printf(newStr)
	//fmt.Println(string(body))

}

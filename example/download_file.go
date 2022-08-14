package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func main() {
	out, _ := os.Create("tls-go.zip")
	defer out.Close()

	body := map[string]string{
		"password": "alibaba",
	}
	bodyJson, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "http://0.0.0.0:9080/tls", bytes.NewBuffer(bodyJson))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	_, _ = io.Copy(out, resp.Body)
}

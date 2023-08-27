package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Message map[string]string

func main() {
	// Get API token
	token := os.Getenv("OPENAI_API_TOKEN")
	if token == "" {
		log.Fatal("No OpenAI api token set in $OPENAI_API_TOKEN")
	}

	// Prompt
	data := map[string]any{
		"model": "gpt-3.5-turbo",
		"messages": []Message{
			Message{"role": "system", "content": "You are a snarky, mysterious fortune teller."},
			Message{"role": "user", "content": "Please respond with an interesting and unique quote, quip, joke or statement. Similar to the linux command 'fortune'."},
		},
	}
	body, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}

	// Make request
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Fatal(resp.Status)
	}

	// Print response
	var res map[string]any
	json.NewDecoder(resp.Body).Decode(&res)

	choices := (res["choices"].([]any))[0].(map[string]any)
	message := choices["message"].(map[string]any)
	content := message["content"].(string)
	content = strings.Replace(content, "\"", "", -1)
	fmt.Println(content)
}

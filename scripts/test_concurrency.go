package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type UpdatePostPayload struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

func updatePost(postId int, p UpdatePostPayload, wg *sync.WaitGroup) {
	defer wg.Done()

	// construc url for updated endpoint
	url := fmt.Sprintf("http://localhost:8080/v1/posts/%d", postId)

	// marshal the payload
	b, _ := json.Marshal(p)

	req, err := http.NewRequest("PATHCT", url, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Error updating resource: ", err)
		return
	}

	// set headers as needed for example
	req.Header.Set("Content-Type", "application/json")

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("update response status:", resp.Status)
}

func main() {
	var wg sync.WaitGroup

	// assuming the post id to update is 1
	postId := 1

	// simulate user A and user B updating the same post concurrently
	wg.Add(2)
	content := "NEW CONTENT FROM USER B"
	title := "NEW TITLE FROM USER A"

	go updatePost(postId, UpdatePostPayload{Title: &title}, &wg)
	go updatePost(postId, UpdatePostPayload{Content: &content}, &wg)
	wg.Wait()
}

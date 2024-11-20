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

func updatePost(postID int, p UpdatePostPayload, wg *sync.WaitGroup) {
	defer wg.Done()

	// Construct the URL for endpoint
	url := fmt.Sprintf("http://localhost:8080/v1/posts/%d", postID)

	// Create json payload
	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error on marshaling payload", err)
		return
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Error on creating request", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on sending request", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Update response status", resp.Status)
}

func main() {
	var wg sync.WaitGroup

	// Post ID to update
	postID := 1

	// Simulate User A and User B updating the same post concurently
	wg.Add(2)
	content := "Content from user B"
	title := "Title from user A"

	go updatePost(postID, UpdatePostPayload{Title: &title}, &wg)
	go updatePost(postID, UpdatePostPayload{Content: &content}, &wg)
	wg.Wait()
}

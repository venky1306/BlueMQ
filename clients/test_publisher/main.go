package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func main() {
	url := "http://localhost:5001/publish/test"
	stringData := ""
	size := 102400
	for i := 0; i < size; i++ {
		stringData += "V"
	}

	// duration := 10 * time.Second // Send requests for a total of 5 minutes

	// endTime := time.Now().Add(duration)

	for i := 0; i < 1000; i++ {
		sendRequest(url, stringData)
		if ((i + 1) % 100) == 0 {
			fmt.Println(time.Now().UnixNano() / 1000000)
		}
	}
}

func sendRequest(url string, data string) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()
}

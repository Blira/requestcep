package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type ApiReturn struct {
	url  string
	data string
}

func main() {
	start := time.Now()

	channel := make(chan ApiReturn)

	go requestApi(channel, "https://cdn.apicep.com/file/apicep/52050-355.json")
	go requestApi(channel, "https://viacep.com.br/ws/52050355/json/")

	select {
	case result := <-channel:
		fmt.Printf("Url: %v\nData: %v\n", result.url, result.data)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout reached!")
	}
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}

func requestApi(channel chan ApiReturn, url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("request failed with status %v", res.StatusCode)
	}
	jsonBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	channel <- ApiReturn{
		url:  url,
		data: string(jsonBody),
	}
	return nil
}

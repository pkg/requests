package main

import (
	"fmt"
	"log"

	"github.com/pkg/requests"
)

func main() {
	var client requests.Client
	resp, err := client.Get("https://httpbin.org/get")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	fmt.Println(resp.Request.Method, resp.Request.URL, resp.Status.Code)
}

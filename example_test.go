package requests_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/pkg/requests"
)

func ExampleClient_Get() {
	var c requests.Client

	resp, err := c.Get("https://www.example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Request.Method, resp.Request.URL, resp.Status.Code)
}

func ExampleClient_Post() {
	var c requests.Client

	body := strings.NewReader("Hello there!")
	resp, err := c.Post("https://www.example.com", body, requests.WithHeader("Content-Type", "application/x-www-form-urlencoded"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Request.Method, resp.Request.URL, resp.Status.Code)
}

func ExampleBody_JSON() {
	var c requests.Client

	resp, err := c.Get("https://frinkiac.com/api/search?q=burn+that+seat")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Close()

	if !resp.IsSuccess() {
		log.Fatalf("%s: expected 200, got %v", resp.Request.URL, resp.Status)
	}

	var results []struct {
		Id        int    `json:"Id"`
		Episode   string `json:"Episode"`
		Timestamp int    `json:"Timestamp"`
	}

	err = resp.JSON(&results)
	fmt.Printf("%#v\n%v", results, err)
}

var response = requests.Response{
	Headers: []requests.Header{
		{Key: "Server", Values: []string{"nginx/1.2.1"}},
		{Key: "Connection", Values: []string{"keep-alive"}},
		{Key: "Content-Type", Values: []string{"text/html; charset=UTF-8"}},
	},
}

func ExampleResponse_Header() {
	fmt.Println(response.Header("Server"))
	fmt.Println(response.Header("Content-Type"))

	// Output:
	// nginx/1.2.1
	// text/html; charset=UTF-8
}

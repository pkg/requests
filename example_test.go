package requests_test

import (
	"fmt"
	"log"

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

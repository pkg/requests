package requests_test

import (
	"fmt"

	"github.com/pkg/requests"
)

func ExampleClient_Get() {
	var c requests.Client

	resp, err := c.Get("https://www.example.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Request.Method, resp.Request.URL, resp.Status.Code)
}

package main

import (
	"fmt"
	"log"

	"github.com/pkg/requests"
)

func main() {
	var client requests.Client
	resp, err := client.Get("https://httpbin.org/get")
	check(err)

	m := make(map[string]interface{})
	err = resp.JSON(&m)
	check(err)

	fmt.Printf("%#v\n", m)
}

func check(err error) {
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

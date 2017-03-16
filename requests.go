package requests

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// Client is a HTTP Client.
type Client struct {
}

// Header is a HTTP header.
type Header struct {
	Key    string
	Values []string
}

// Get issues a GET to the specified URL.
func (c *Client) Get(url string, options ...func(*Request) error) (*Response, error) {
	req := Request{
		Method: "GET",
		URL:    url,
	}
	if err := c.applyOptions(&req, options...); err != nil {
		return nil, err
	}
	return c.do(&req)
}

func (c *Client) do(request *Request) (*Response, error) {
	req, err := http.NewRequest(request.Method, request.URL, request.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	r := Response{
		Request: request,
		Status: Status{
			Code:   resp.StatusCode,
			Reason: resp.Status[4:],
		},
		Headers: headers(resp.Header),
		Body: Body{
			ReadCloser: resp.Body,
		},
	}
	return &r, nil
}

func (c *Client) applyOptions(req *Request, options ...func(*Request) error) error {
	for _, opt := range options {
		if err := opt(req); err != nil {
			return err
		}
	}
	return nil
}

// Request is a HTTP request.
type Request struct {
	Method  string
	URL     string
	Headers []Header
	Body    io.Reader
}

// Status is a HTTP reponse status.
type Status struct {
	Code   int
	Reason string
}

// Response is a HTTP response.
type Response struct {
	*Request
	Status
	Headers []Header
	Body
}

type Body struct {
	io.ReadCloser

	json *json.Decoder
}

// JSON decodes the next JSON encoded object in the body to v.
func (b *Body) JSON(v interface{}) error {
	if b.json == nil {
		b.json = json.NewDecoder(b)
	}
	return b.json.Decode(v)
}

// return the body as a string, or bytes, or somethign

func headers(h map[string][]string) []Header {
	headers := make([]Header, 0, len(h))
	for k, v := range h {
		headers = append(headers, Header{
			Key:    k,
			Values: v,
		})
	}
	return headers
}

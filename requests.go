package requests

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

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

// Response is a HTTP response.
type Response struct {
	*Request
	Status
	Headers []Header
	Body
}

// Header returns the canonicalised version of a response header as a string
// If there is no key present in the response the empty string is returned.
// If multiple headers are present, they are canonicalised into as single string
// by joining them with a comma. See RFC 2616 § 4.2.
func (r *Response) Header(key string) string {
	var vals []string
	for _, h := range r.Headers {

		// TODO(dfc) § 4.2 states that not all header values can be combined, but equally those
		// that cannot be combined with a comma may not be present more than once in a
		// header block.
		if h.Key == key {
			vals = append(vals, h.Values...)
		}
	}
	return strings.Join(vals, ",")
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

// return the body as a string, or bytes, or something

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

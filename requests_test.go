package requests

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestBodyJSON(t *testing.T) {
	jsonbody := func(s string) io.ReadCloser {
		type rc struct {
			io.Reader
			io.Closer
		}
		return rc{
			Reader: strings.NewReader(s),
		}
	}

	type T struct {
		A string `json:"a"`
	}

	tests := []struct {
		Body
		want []T
	}{{
		Body: Body{
			ReadCloser: jsonbody(`{"a":"hello"}`),
		},
		want: []T{{A: "hello"}},
	}, {
		Body: Body{
			ReadCloser: jsonbody(`{"a":"first"}{"a":"second"}`),
		},
		want: []T{{A: "first"}, {A: "second"}},
	}}

	for _, tt := range tests {
		got := make([]T, len(tt.want))
		for i := range got {
			if err := tt.Body.JSON(&got[i]); err != nil {
				t.Fatal(err)
			}
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got: %v, want: %v", got, tt.want)
		}
	}
}

func TestResponseHeader(t *testing.T) {
	header := func(key, val string, vals ...string) Header {
		v := []string{val}
		return Header{
			Key:    key,
			Values: append(v, vals...),
		}
	}

	resp := &Response{
		Headers: []Header{
			header("foo", "bar"),
			header("quxx", "frob", "frob"),
		},
	}

	tests := []struct {
		*Response
		key  string
		want string
	}{{
		resp,
		"foo",
		"bar",
	}, {
		resp,
		"quxx",
		"frob,frob",
	}, {
		resp,
		"flimm",
		"",
	}}

	for i, tc := range tests {
		got := tc.Header(tc.key)
		if got != tc.want {
			t.Errorf("%d: Header(%q): got %q, want %v", i, tc.key, got, tc.want)
		}
	}
}

func TestRequestToHeaders(t *testing.T) {
	tests := []struct {
		Headers []Header
		want    map[string][]string
	}{{
		Headers: []Header{
			{Key: "foo", Values: []string{"bar"}},
			{Key: "cram", Values: []string{"witt", "jannet"}},
		},
		want: map[string][]string{
			"foo":  []string{"bar"},
			"cram": []string{"witt", "jannet"},
		},
	}}

	for i, tc := range tests {
		got := toHeaders(tc.Headers)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%d: %v.toHeaders(): got: %v, want: %v", i, tc.Headers, got, tc.want)
		}
	}
}

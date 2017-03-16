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

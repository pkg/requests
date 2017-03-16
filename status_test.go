package requests

import "testing"

func TestStatusString(t *testing.T) {
	tests := []struct {
		Status
		want string
	}{
		{Status{Code: 200, Reason: "OK"}, "200 OK"},
		{Status{Code: 418, Reason: "I'm a teapot"}, "418 I'm a teapot"},
	}

	for _, tt := range tests {
		got := tt.Status.String()
		if got != tt.want {
			t.Errorf("got: %q, want: %q", got, tt.want)
		}
	}
}

func TestStatusMethods(t *testing.T) {
	tests := []struct {
		Status
		informational, success, redirect, error, clienterr, servererr bool
	}{
		{Status{Code: INFO_CONTINUE}, true, false, false, false, false, false},
		{Status{Code: SUCCESS_OK}, false, true, false, false, false, false},
		{Status{Code: REDIRECTION_MULTIPLE_CHOICES}, false, false, true, false, false, false},
		{Status{Code: CLIENT_ERROR_BAD_REQUEST}, false, false, false, true, true, false},
		{Status{Code: SERVER_ERROR_INTERNAL}, false, false, false, true, false, true},
	}

	for _, tt := range tests {
		if info := tt.Status.IsInformational(); info != tt.informational {
			t.Errorf("Status(%q).Informational: expected %v, got %v", tt.Status, tt.informational, info)
		}
		if success := tt.Status.IsSuccess(); success != tt.success {
			t.Errorf("Status(%q).Success: expected %v, got %v", tt.Status, tt.success, success)
		}
		if redirect := tt.Status.IsRedirect(); redirect != tt.redirect {
			t.Errorf("Status(%q).Redirect: expected %v, got %v", tt.Status, tt.redirect, redirect)
		}
		if error := tt.Status.IsError(); error != tt.error {
			t.Errorf("Status(%q).IsError: expected %v, got %v", tt.Status, tt.error, error)
		}
		if error := tt.Status.IsClientError(); error != tt.clienterr {
			t.Errorf("Status(%q).IsError: expected %v, got %v", tt.Status, tt.clienterr, error)
		}
		if error := tt.Status.IsServerError(); error != tt.servererr {
			t.Errorf("Status(%q).IsError: expected %v, got %v", tt.Status, tt.servererr, error)
		}
	}
}

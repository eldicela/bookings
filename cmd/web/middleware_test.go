package main

import (
	"net/http"
	"testing"
)

func TestNoSurve(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// Do nothing
	default:
		// t.Error(fmt.Sprintf("Type is not http.Handler, but is %T", v))
		t.Errorf("Type is not http.Handler, but is %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// Do nothing
	default:
		// t.Error(fmt.Sprintf("Type is not http.Handler, but is %T", v))
		t.Errorf("Type is not http.Handler, but is %T", v)
	}
}

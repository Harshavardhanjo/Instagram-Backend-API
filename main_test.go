package main

import (
	"Strings"
	"go"
	"net/http/httptest"
	"testing"
)

func TestGetUserEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/?id=6161835c7f83f57c9e98b85a", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUserEndpoint)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"_id":"6161835c7f83f57c9e98b85a","Username":"Bob","Password":"���h�f^C~�F�}(b&S�i�#�\bMtM����","Email":"hi@gmail.com"}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

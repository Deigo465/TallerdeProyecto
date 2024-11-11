package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHome(t *testing.T) {
	// GIVEN
	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	rw := httptest.NewRecorder()
	// homeTemplatePath = "../../web/views/landing.html.tmpl"
	SetBasePath("../../web/views/")

	// WHEN
	Home(rw, req)

	// THEN
	resp := rw.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be 200, got %v", resp.StatusCode)
	}

	body := rw.Body.String()

	if body == "" {
		t.Errorf("Expected body to not be empty")
	}

	if !strings.Contains(body, "BlockeHR") {
		t.Errorf("Expected body to include 'BlockeHR', got %s", body)
	}
}

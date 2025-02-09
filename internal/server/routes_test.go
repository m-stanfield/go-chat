package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-chat-react/internal/database"
)

var port = 8080

func setupSuite(tb testing.TB) (*Server, func(tb testing.TB)) {
	server := database.NewInMemory()
	if server == nil {
		tb.Fatal("failed to create server")
	}
	return &Server{port: port, db: server}, func(tb testing.TB) {
		server.Close()
	}
}

// Function to perform login and retrieve the login cookie
func getLoginCookie(username, password string) (*http.Cookie, error) {
	loginReq := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: username,
		Password: password,
	}
	loginData, err := json.Marshal(loginReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal login request: %v", err)
	}

	// Send the login POST request
	resp, err := http.Post(
		"http://localhost:8080/api/login",
		"application/json",
		bytes.NewBuffer(loginData),
	)
	if err != nil {
		return nil, fmt.Errorf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check if login was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status: %s", resp.Status)
	}

	// Retrieve the login cookie from the response
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "login" { // assuming the cookie is named "login"
			return cookie, nil
		}
	}

	return nil, fmt.Errorf("login cookie not found in the response")
}

func TestHandler(t *testing.T) {
	s := &Server{}
	server := httptest.NewServer(http.HandlerFunc(s.HelloWorldHandler))
	defer server.Close()
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	expected := "{\"message\":\"Hello World\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestCreateNewServer(t *testing.T) {
	s, teardown := setupSuite(t)
	defer teardown(t)
	userid := database.Id(1)
	server := httptest.NewServer(http.HandlerFunc(s.createNewServer))
	defer server.Close()
	payload := map[string]string{"servername": "testserver"}
	jsonData, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/newserver", bytes.NewBuffer(jsonData))
	req = req.WithContext(context.WithValue(context.Background(), "userid", userid))
	req.Header.Set("Content-Type", "application/json")

	// Step 3: Create a ResponseRecorder to capture the response
	resp := httptest.NewRecorder()
	s.createNewServer(resp, req)
	// Assertions
	if resp.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Code)
	}
	expected_id := 3
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	result := struct {
		ServerId database.Id `json:"serverid"`
	}{}
	json.Unmarshal(body, &result)

	if expected_id != int(result.ServerId) {
		t.Errorf("expected response body to be %d; got %d", expected_id, result)
	}
	serverinfo, err := s.db.GetServer(database.Id(3))
	if err != nil {
		t.Fatalf("error getting server info. Err: %v", err)
	}
	if serverinfo.ServerName != "testserver" {
		t.Errorf("expected server name to be %v; got %v", "testserver", serverinfo.ServerName)
	}
}

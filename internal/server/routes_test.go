package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"go-chat-react/internal/database"
)

var port = 8080

type TestServer struct {
	server *httptest.Server
}

func setupTest(tb testing.TB) (*TestServer, func(tb testing.TB)) {
	server := database.NewInMemory()
	if server == nil {
		tb.Fatal("failed to create server")
	}

	s := Server{port: port, db: server}
	httpserver := httptest.NewServer(s.RegisterRoutes())
	return &TestServer{server: httpserver}, func(tb testing.TB) {
		server.Close()
	}
}

func (s *TestServer) buildRequest(
	method string,
	endpoint string,
	payload map[string]string,
) (*http.Request, error) {
	var jsonData []byte = nil
	var buffer io.Reader = nil
	if len(payload) > 0 {
		jsonData, _ = json.Marshal(payload)
		buffer = bytes.NewBuffer(jsonData)
	}
	req, err := http.NewRequest(method, s.server.URL+endpoint, buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (s *TestServer) sendRequest(
	method string,
	endpoint string,
	payload map[string]string,
) (*http.Response, error) {
	session, err := s.buildRequest(method, endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	return s.server.Client().Do(session)
}

func (s *TestServer) sendAuthRequest(
	method string,
	endpoint string,
	payload map[string]string,
	usernameOverride *string,
	passwordOverride *string,
) (*http.Response, error) {
	session, err := s.buildRequest(method, endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	if usernameOverride == nil {
		temp := "u1"
		usernameOverride = &temp
	}
	if passwordOverride == nil {
		temp := "1"
		passwordOverride = &temp
	}
	cookie, err := s.getLoginCookie(*usernameOverride, *passwordOverride)
	if err != nil {
		return nil, fmt.Errorf("failed to get login cookie: %v", err)
	}
	session.AddCookie(cookie)
	return s.server.Client().Do(session)
}

// Function to perform login and retrieve the login cookie
func (s *TestServer) getLoginCookie(username, password string) (*http.Cookie, error) {
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
	resp, err := s.server.Client().Post(
		s.server.URL+"/api/auth/login",
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
		if cookie.Name == "token" {
			return cookie, nil
		}
	}

	return nil, fmt.Errorf("login cookie not found in the response")
}

func TestCreateNewServer_Error_ShortName(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/servers"
	payload := map[string]string{"servername": "te"}
	resp, err := s.sendAuthRequest(http.MethodPost, endpoint, payload, nil, nil)
	if err != nil {
		t.Fatalf("error creating session. Err: %v", err)
	}
	// Assertions
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status StatusBadRequest; got %v", resp.Status)
	}
}

func TestCreateNewServer_Error_LongName(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/servers"

	payload := map[string]string{
		"servername": "wwwwwwwwwwwwwwwwwwwwwwwwwwwwwww",
	}
	resp, _ := s.sendAuthRequest(http.MethodPost, endpoint, payload, nil, nil)
	// Assertions
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status StatusBadRequest; got %v", resp.Status)
	}
}

func TestCreateNewServer_Valid(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/servers"
	payload := map[string]string{"servername": "testserver"}
	expected_server_id := database.Id(3)
	resp, _ := s.sendAuthRequest(http.MethodPost, endpoint, payload, nil, nil)
	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	result := struct {
		ServerId database.Id `json:"serverid"`
	}{}
	json.Unmarshal(body, &result)

	if expected_server_id != result.ServerId {
		t.Errorf("expected response body to be %d; got %d", expected_server_id, result)
	}
	getresp, err := s.sendRequest(
		http.MethodGet,
		"/api/servers/"+strconv.Itoa(int(expected_server_id)),
		nil,
	)
	if err != nil {
		t.Fatalf("error getting server info. Err: %v", err)
	}

	body, err = io.ReadAll(getresp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	getresult := struct {
		ServerId   database.Id `json:"serverid"`
		ServerName string      `json:"servername"`
	}{}
	json.Unmarshal(body, &getresult)
	if getresult.ServerName != "testserver" {
		t.Errorf("expected server name to be %v; got %v", "testserver", getresult.ServerName)
	}
}

func TestGetUser_Valid(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/users/1"
	resp, _ := s.sendRequest(http.MethodGet, endpoint, nil)
	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	result := struct {
		UserName string `json:"username"`
	}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("error unmarshalling response body. Err: %v", err)
	}
	if result.UserName != "u1" {
		t.Errorf("expected username to be %v; got %v", "u1", result.UserName)
	}
}

func TestGetUser_Invalid(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/users/100"
	resp, _ := s.sendRequest(http.MethodGet, endpoint, nil)
	// Assertions
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status NotFound; got %v", resp.Status)
	}
}

func TestCreateUser(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/users"
	payload := map[string]string{"username": "new_user", "password": "new_user_password"}
	resp, _ := s.sendRequest(http.MethodPost, endpoint, payload)
	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	create_result := struct {
		UserId database.Id `json:"userid"`
	}{}
	err = json.Unmarshal(body, &create_result)
	if err != nil {
		t.Fatalf("error unmarshalling response body. Err: %v", err)
	}
	if create_result.UserId != 4 {
		t.Errorf("expected userid to be %v; got %v", 4, create_result.UserId)
	}
	user_resp, err := s.sendRequest(http.MethodGet, "/api/users/4", nil)
	if err != nil {
		t.Fatalf("error getting user info. Err: %v", err)
	}
	user_result := struct {
		UserName string `json:"username"`
	}{}
	body, err = io.ReadAll(user_resp.Body)
	err = json.Unmarshal(body, &user_result)
	if err != nil {
		t.Fatalf("error unmarshalling response body. Err: %v", err)
	}

	if user_result.UserName != payload["username"] {
		t.Errorf("expected userid to be %v; got %v", user_result.UserName, payload["username"])
	}
}

func TestLogin_Valid(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/auth/login"
	payload := map[string]string{"username": "u1", "password": "1"}
	resp, _ := s.sendRequest(http.MethodPost, endpoint, payload)
	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	result := struct {
		UserId          database.Id `json:"userid"`
		TokenExpireTime string      `json:"token_expire_time"`
	}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("error unmarshalling response body. Err: %v", err)
	}
	if result.UserId != 1 {
		t.Errorf("expected userid to be %v; got %v", 1, result.UserId)
	}
}

func TestLogin_Invalid(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/auth/login"
	payload := map[string]string{"username": "u1", "password": "notcorrectpassword"}
	resp, _ := s.sendRequest(http.MethodPost, endpoint, payload)
	// Assertions
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status expected %v; got %v", http.StatusBadRequest, resp.Status)
	}
}

func TestGetServerChannels_Valid(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/users/1/servers"
	payload := map[string]string{"serverid": "1"}
	resp, err := s.sendAuthRequest(http.MethodGet, endpoint, payload, nil, nil)
	if err != nil {
		t.Fatalf("error creating session. Err: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	result := struct {
		Servers []database.Server `json:"servers"`
	}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("error unmarshalling response body. Err: %v", err)
	}
	if len(result.Servers) != 2 {
		t.Errorf("expected number of channels to be %v; got %v", 2, len(result.Servers))
	}
	if result.Servers[0].ServerId != 1 {
		t.Errorf("expected server id to be %v; got %v", 1, result.Servers[0].ServerId)
	}
	if result.Servers[1].ServerId != 2 {
		t.Errorf("expected server id to be %v; got %v", 2, result.Servers[1].ServerId)
	}
}

func TestGetServerChannels_Unauthed(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/users/1/servers"
	payload := map[string]string{"serverid": "1"}
	unauthed_user := "u2"
	password := "2"
	resp, err := s.sendAuthRequest(http.MethodGet, endpoint, payload, &unauthed_user, &password)
	if err != nil {
		t.Fatalf("error creating session. Err: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %v; got %v", http.StatusBadRequest, resp.Status)
	}
}

func TestGetServerInformation_Valid(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/servers/1"
	resp, err := s.sendRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		t.Fatalf("error creating session. Err: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	result := database.Server{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("error unmarshalling response body. Err: %v", err)
	}
	if result.ServerId != 1 {
		t.Errorf("expected server id to be %v; got %v", 1, result.ServerId)
	}
	if result.ServerName != "server1" {
		t.Errorf("expected server name to be %v; got %v", "server1", result.ServerName)
	}
	if result.OwnerId != 1 {
		t.Errorf("expected owner id to be %v; got %v", 1, result.OwnerId)
	}
}

func TestGetServerChannels(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/servers/1/channels"
	resp, err := s.sendAuthRequest(http.MethodGet, endpoint, nil, nil, nil)
	if err != nil {
		t.Fatalf("error creating session. Err: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	result := struct {
		Channels []database.Channel `json:"channels"`
	}{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("error unmarshalling response body. Err: %v", err)
	}
	if len(result.Channels) != 2 {
		t.Errorf("expected number of channels to be %v; got %v", 2, len(result.Channels))
	}
	if result.Channels[0].ServerId != 1 {
		t.Errorf("expected server id to be %v; got %v", 1, result.Channels[0].ServerId)
	}
	if result.Channels[1].ServerId != 1 {
		t.Errorf("expected server id to be %v; got %v", 1, result.Channels[1].ServerId)
	}
	if result.Channels[0].Timestamp.IsZero() {
		t.Errorf("expected timestamp to be non-zero")
	}
	if result.Channels[1].Timestamp.IsZero() {
		t.Errorf("expected timestamp to be non-zero")
	}
	if result.Channels[0].ChannelId != 1 {
		t.Errorf("expected channel id to be %v; got %v", 1, result.Channels[0].ChannelId)
	}
	if result.Channels[1].ChannelId != 2 {
		t.Errorf("expected channel id to be %v; got %v", 2, result.Channels[1].ChannelId)
	}
	if result.Channels[0].ChannelName != "a" {
		t.Errorf(
			"expected channel name to be %v; got %v",
			"b",
			result.Channels[0].ChannelName,
		)
	}
	if result.Channels[1].ChannelName != "b" {
		t.Errorf(
			"expected channel name to be %v; got %v",
			"b",
			result.Channels[1].ChannelName,
		)
	}
}

func TestGetServerMessages(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/servers/1/messages"
	resp, err := s.sendAuthRequest(http.MethodGet, endpoint, nil, nil, nil)
	if err != nil {
		t.Fatalf("error creating session. Err: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	result := struct {
		Messages []Message `json:"messages"`
	}{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("error unmarshalling response body. Err: %v", err)
	}
	for _, message := range result.Messages {
		if len(message.Message) <= 0 {
			t.Errorf("expected message to be non-empty")
		}
	}
	if len(result.Messages) != 4 {
		t.Errorf("expected number of messages to be %v; got %v", 4, len(result.Messages))
	}
}

func TestGetServerMembers_Valid(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/servers/1/members"
	resp, err := s.sendAuthRequest(http.MethodGet, endpoint, nil, nil, nil)
	if err != nil {
		t.Fatalf("error creating session. Err: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	result := struct {
		Users []database.User `json:"users"`
	}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("error unmarshalling response body. Err: %v", err)
	}
	if len(result.Users) != 2 {
		t.Errorf("expected number of channels to be %v; got %v", 2, len(result.Users))
	}
	if result.Users[0].UserId != 1 {
		t.Errorf("expected user id to be %v; got %v", 1, result.Users[0].UserId)
	}
	if result.Users[1].UserId != 3 {
		t.Errorf("expected user id to be %v; got %v", 2, result.Users[1].UserId)
	}
	if result.Users[0].UserName != "u1" {
		t.Errorf("expected user name to be %v; got %v", "u1", result.Users[0].UserName)
	}
	if result.Users[1].UserName != "u3" {
		t.Errorf("expected user name to be %v; got %v", "u2", result.Users[1].UserName)
	}
}

func TestUpdateServer_Valid(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/servers/1"
	new_server_name := "new_server_name"
	payload := map[string]string{"servername": new_server_name}
	resp, err := s.sendAuthRequest(http.MethodPatch, endpoint, payload, nil, nil)
	if err != nil {
		t.Fatalf("error creating session. Err: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	getresp, err := s.sendRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		t.Fatalf("error getting server info. Err: %v", err)
	}
	if getresp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", getresp.Status)
	}
	body, err := io.ReadAll(getresp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	result := struct {
		ServerId   database.Id `json:"serverid"`
		OwnerId    database.Id `json:"ownerid"`
		ServerName string      `json:"servername"`
	}{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("error unmarshalling response body. Err: %v", err)
	}
	if result.ServerName != new_server_name {
		t.Errorf("expected server name to be %v; got %v", "new_server_name", result.ServerName)
	}
}

func Test_UpdateServer_InvalidPermissions(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/servers/1"
	new_server_name := "new_server_name"
	payload := map[string]string{"servername": new_server_name}
	username := "u2"
	password := "2"
	resp, err := s.sendAuthRequest(http.MethodPatch, endpoint, payload, &username, &password)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("got unexpected status code Err: %v", err)
	}
}

func Test_DeleteServer_Valid(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/servers/1"
	resp, err := s.sendAuthRequest(http.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		t.Fatalf("error creating session. Err: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	getresp, err := s.sendRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		t.Fatalf("error getting server info. Err: %v", err)
	}
	if getresp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status OK; got %v", getresp.Status)
	}
}

func Test_UpdateChannel_Valid(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/channels/1"
	new_channel_name := "new_channel_name"
	payload := map[string]string{"channelname": new_channel_name}
	resp, err := s.sendAuthRequest(http.MethodPatch, endpoint, payload, nil, nil)
	if err != nil {
		t.Fatalf("error creating session. Err: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	getresp, err := s.sendAuthRequest(http.MethodGet, endpoint, nil, nil, nil)
	if err != nil {
		t.Fatalf("error getting server info. Err: %v", err)
	}
	if getresp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", getresp.Status)
	}
	body, err := io.ReadAll(getresp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	result := struct {
		ChannelId   database.Id `json:"channelid"`
		ServerId    database.Id `json:"serverid"`
		ChannelName string      `json:"channelname"`
		TimeStamp   string      `json:"timestamp"`
	}{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("error unmarshalling response body. Err: %v", err)
	}
	if result.ChannelName != new_channel_name {
		t.Errorf("expected server name to be %v; got %v", new_channel_name, result.ChannelName)
	}
}

func Test_AddChannelMember_Valid(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/channels/1/members"
	payload := map[string]string{"userid": "3"}
	resp, err := s.sendAuthRequest(http.MethodPost, endpoint, payload, nil, nil)
	if err != nil {
		t.Fatalf("error creating session. Err: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	getresp, err := s.sendAuthRequest(http.MethodGet, endpoint, nil, nil, nil)
	if err != nil {
		t.Fatalf("error getting server info. Err: %v", err)
	}
	if getresp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", getresp.Status)
	}
	body, err := io.ReadAll(getresp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	result := struct {
		Users []User `json:"users"`
	}{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("error unmarshalling response body. Err: %v", err)
	}
	for _, user := range result.Users {
		if user.UserID == 3 {
			return
		}
	}
	t.Errorf("expected user id to be %v; got %v", 3, result.Users)
}

func TestServer_GetChannelMembers(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/channels/1/members"
	resp, err := s.sendAuthRequest(http.MethodGet, endpoint, nil, nil, nil)
	if err != nil {
		t.Fatalf("error creating session. Err: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	result := struct {
		Users []User `json:"users"`
	}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("error unmarshalling response body. Err: %v", err)
	}
	if result.Users[0].UserID != database.Id(1) {
		t.Errorf("expected user id to be %v; got %v", 1, result.Users[0].UserID)
	}
	if result.Users[0].UserName != "u1" {
		t.Errorf("expected user name to be %v; got %v", "u1", result.Users[0].UserName)
	}
}

func TestServer_RemoveChannelMember_Valid(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/channels/1/members"
	payload := map[string]string{"userid": "2"}
	resp, err := s.sendAuthRequest(http.MethodDelete, endpoint, payload, nil, nil)
	if err != nil {
		t.Fatalf("error creating session. Err: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status OK; got %v", resp.Status)
	}
	getresp, err := s.sendAuthRequest(http.MethodGet, endpoint, nil, nil, nil)
	if err != nil {
		t.Fatalf("error getting server info. Err: %v", err)
	}
	if getresp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", getresp.Status)
	}
	body, err := io.ReadAll(getresp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	result := struct {
		Users []User `json:"users"`
	}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("error unmarshalling response body. Err: %v", err)
	}
	for _, user := range result.Users {
		if user.UserID == 2 {
			t.Errorf("expected user id to be %v; got %v", 2, user.UserID)
		}
	}
}

func TestServer_RemoveChannelMember_NotInChannel(t *testing.T) {
	s, teardown := setupTest(t)
	defer teardown(t)
	endpoint := "/api/channels/1/members"
	payload := map[string]string{"userid": "20"}
	resp, err := s.sendAuthRequest(http.MethodDelete, endpoint, payload, nil, nil)
	if err != nil {
		t.Errorf("error creating session. Err: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status BadRequest; got %v", resp.Status)
	}
}

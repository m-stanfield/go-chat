package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/coder/websocket"

	"go-chat-react/internal/database"
)

var (
	clients         = make([]chan Message, 0)
	incomingChannel = make(chan Message)
	startTime       = time.Now()
)

type Message struct {
	UserName  string      `json:"username"`
	UserId    database.Id `json:"userid"`
	MessageID database.Id `json:"messageid"`
	ChannelId database.Id `json:"channelid"`
	Message   string      `json:"message"`
	Date      string      `json:"date"`
}

type User struct {
	UserID   database.Id `json:"userid"`
	UserName string      `json:"username"`
}

type SubmittedMessage struct {
	UserID    string      `json:"userid"`
	ChannelId database.Id `json:"channelid"`
	Token     string      `json:"token"`
	Message   string      `json:"message"`
}

func (s *Server) RegisterRoutes() http.Handler {
	go s.handleIncomingMessages(incomingChannel)
	mux := http.NewServeMux()

	// Register routes

	mux.HandleFunc("/", s.redirectToReact)
	mux.HandleFunc("/websocket", s.WithAuthUser(s.websocketHandler))

	mux.HandleFunc("POST /api/auth/login", s.loginHandler)
	mux.HandleFunc("POST /api/auth/logout", s.WithAuthUser(s.LogoutHandler))

	mux.HandleFunc("POST /api/users", s.createUserHandler)
	mux.HandleFunc("GET /api/users/{userid}", s.GetUserHandler)
	mux.HandleFunc("PATCH /api/users/{userid}", s.WithAuthUser(s.UpdateUser))
	mux.HandleFunc("GET /api/users/{userid}/servers", s.WithAuthUser(s.GetServersOfUser))

	mux.HandleFunc("POST /api/servers", s.WithAuthUser(s.createNewServer))
	mux.HandleFunc("GET /api/servers/{serverid}", s.GetServerInformation)
	mux.HandleFunc("PATCH /api/servers/{serverid}", s.WithAuthUser(s.UpdateServer))
	mux.HandleFunc("DELETE /api/servers/{serverid}", s.WithAuthUser(s.DeleteServer))
	mux.HandleFunc("GET /api/servers/{serverid}/channels", s.WithAuthUser(s.GetServerChannels))
	mux.HandleFunc("POST /api/servers/{serverid}/channels", s.WithAuthUser(s.CreateChannel))
	mux.HandleFunc("GET /api/servers/{serverid}/members", s.WithAuthUser(s.GetServerMembersHandler))
	mux.HandleFunc("GET /api/servers/{serverid}/messages", s.WithAuthUser(s.GetServerMessages))

	mux.HandleFunc("GET /api/channels/{channelid}", s.WithAuthUser(s.GetChannel))
	mux.HandleFunc("PATCH /api/channels/{channelid}", s.WithAuthUser(s.UpdateChannel))
	mux.HandleFunc("DELETE /api/channels/{channelid}", s.WithAuthUser(s.DeleteChannel))
	mux.HandleFunc("POST /api/channels/{channelid}/members", s.WithAuthUser(s.AddChannelMember))
	mux.HandleFunc("GET /api/channels/{channelid}/members", s.WithAuthUser(s.GetChannelMembers))
	mux.HandleFunc(
		"DELETE /api/channels/{channelid}/members",
		s.WithAuthUser(s.RemoveChannelMember),
	)

	mux.HandleFunc("GET /api/channels/{channelid}/messages", s.GetChannelMessages)
	mux.HandleFunc("POST /api/channels/{channelid}/messages", s.CreateChannelMessage)
	mux.HandleFunc("GET /api/channels/{channelid}/messages/{messageid}", s.GetMessage)
	mux.HandleFunc(
		"PATCH /api/channels/{channelid}/messages/{messageid}",
		s.WithAuthUser(s.UpdateMessage),
	)
	mux.HandleFunc(
		"DELETE /api/channels/{channelid}/messages/{messageid}",
		s.WithAuthUser(s.DeleteMessage),
	)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(s.logEndpoint(mux))
}

// redirect so I only have to remember one port during development
func (s *Server) redirectToReact(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:5173", http.StatusTemporaryRedirect)
}

func (s *Server) UpdateServer(w http.ResponseWriter, r *http.Request) {
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serverid_str := r.PathValue("serverid")
	serverid, err := strconv.Atoi(serverid_str)
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	if serverid <= 0 {
		http.Error(w, "invalid request: invalid server id", http.StatusBadRequest)
		return
	}
	server_info, err := s.db.GetServer(database.Id(serverid))
	if err != nil {
		http.Error(w, "error: unable to locate server", http.StatusBadRequest)
		return
	}
	if server_info.OwnerId != userid {
		http.Error(w, "error: user not owner of server", http.StatusBadRequest)
		return
	}

	new_server_name := struct {
		ServerName string `json:"servername"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&new_server_name)
	if err != nil {
		http.Error(w, "error: unable to parse request", http.StatusBadRequest)
		return
	}

	err = s.db.UpdateServerName(database.Id(serverid), new_server_name.ServerName)
	if err != nil {
		http.Error(w, "error: unable to update server name", http.StatusBadRequest)
		return
	}
}

func (s *Server) DeleteServer(w http.ResponseWriter, r *http.Request) {
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serverid_str := r.PathValue("serverid")
	serverid, err := strconv.Atoi(serverid_str)
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	if serverid <= 0 {
		http.Error(w, "invalid request: invalid server id", http.StatusBadRequest)
		return
	}
	server_info, err := s.db.GetServer(database.Id(serverid))
	if err != nil {
		http.Error(w, "error: unable to locate server", http.StatusBadRequest)
		return
	}
	if server_info.OwnerId != userid {
		http.Error(w, "error: user not owner of server", http.StatusBadRequest)
		return
	}
	err = s.db.DeleteServer(database.Id(serverid))
	if err != nil {
		http.Error(w, "error: unable to locate server", http.StatusBadRequest)
		return
	}
}

func (s *Server) UpdateChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	channelid_str := r.PathValue("channelid")
	channelid, err := strconv.Atoi(channelid_str)
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	if channelid <= 0 {
		http.Error(w, "invalid request: invalid server id", http.StatusBadRequest)
		return
	}
	channel_info, err := s.db.GetChannel(database.Id(channelid))
	if err != nil {
		http.Error(w, "error: unable to locate channel", http.StatusBadRequest)
		return
	}
	server_info, err := s.db.GetServer(channel_info.ServerId)
	if server_info.OwnerId != userid {
		http.Error(w, "error: user not owner of channel", http.StatusBadRequest)
		return
	}

	new_channel_name := struct {
		UpdatedChannelName string `json:"channelname"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&new_channel_name)
	if err != nil {
		http.Error(w, "error: unable to parse request", http.StatusBadRequest)
		return
	}
	err = s.db.UpdateChannel(database.Id(channelid), new_channel_name.UpdatedChannelName)
	if err != nil {
		http.Error(w, "error: unable to update channel", http.StatusBadRequest)
		return

	}
}

func (s *Server) GetChannelMembers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	channelid_str := r.PathValue("channelid")
	channelid, err := strconv.Atoi(channelid_str)
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	if channelid <= 0 {
		http.Error(w, "invalid request: invalid server id", http.StatusBadRequest)
		return
	}
	inchannel, err := s.db.IsUserInChannel(userid, database.Id(channelid))
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), http.StatusBadRequest)
		return
	}
	if !inchannel {
		http.Error(w, "user not in channel", http.StatusBadRequest)
		return
	}
	users, err := s.db.GetUsersInChannel(database.Id(channelid))
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	// convert users from database.User to server.User
	newusers := make([]User, len(users))
	for i, user := range users {
		newusers[i] = User{
			UserID:   user.UserId,
			UserName: user.UserName,
		}
	}

	resp := map[string]interface{}{"users": newusers}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) AddChannelMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	post_data := struct {
		UserId string `json:"userid"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&post_data)
	if err != nil {
		fmt.Printf("error: unable to parse request %s", err)
		http.Error(w, fmt.Sprintf("error: unable to parse request %s", err), http.StatusBadRequest)
		return
	}
	newuserid_str, err := strconv.Atoi(post_data.UserId)
	if err != nil {
		http.Error(w, "invalid request: unable to parse user id", http.StatusBadRequest)
		return
	}
	if newuserid_str <= 0 {
		http.Error(w, "invalid request: invalid user id", http.StatusBadRequest)
		return
	}
	newuserid := database.Id(newuserid_str)

	channelid_str := r.PathValue("channelid")
	channelid, err := strconv.Atoi(channelid_str)
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	if channelid <= 0 {
		http.Error(w, "invalid request: invalid server id", http.StatusBadRequest)
		return
	}
	channel, err := s.db.GetChannel(database.Id(channelid))
	if err != nil {
		http.Error(w, "error: unable to locate server", http.StatusBadRequest)
		return
	}
	server_info, err := s.db.GetServer(channel.ServerId)
	if err != nil {
		http.Error(w, "error: unable to locate server", http.StatusBadRequest)
		return
	}
	if server_info.OwnerId != userid {
		http.Error(w, "error: user not owner of server", http.StatusBadRequest)
		return
	}
	inserver, err := s.db.IsUserInServer(database.Id(newuserid), server_info.ServerId)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), http.StatusBadRequest)
		return
	}
	if !inserver {
		http.Error(w, "user not in server", http.StatusBadRequest)
		return
	}
	err = s.db.AddUserToChannel(database.Id(newuserid), server_info.ServerId)
	if err != nil {
		http.Error(w, "error: unable to add user to channel", http.StatusBadRequest)
		return
	}
}

func (s *Server) RemoveChannelMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// get serverid for the channel and make sure the user is the owner
	channelid_str := r.PathValue("channelid")
	channelid, err := strconv.Atoi(channelid_str)
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	if channelid <= 0 {
		http.Error(w, "invalid request: invalid server id", http.StatusBadRequest)
		return
	}
	channel, err := s.db.GetChannel(database.Id(channelid))
	if err != nil {
		http.Error(w, "error: unable to locate server", http.StatusBadRequest)
		return
	}
	server_info, err := s.db.GetServer(channel.ServerId)
	if err != nil {
		http.Error(w, "error: unable to locate server", http.StatusBadRequest)
		return
	}
	if server_info.OwnerId != userid {

		http.Error(w, "error: user not owner of server", http.StatusBadRequest)
		return
	}
	post_data := struct {
		UserId string `json:"userid"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&post_data)
	if err != nil {
		fmt.Printf("error: unable to parse request %s", err)
		http.Error(w, fmt.Sprintf("error: unable to parse request %s", err), http.StatusBadRequest)
		return
	}
	newuserid_str, err := strconv.Atoi(post_data.UserId)
	if err != nil {
		http.Error(w, "invalid request: unable to parse user id", http.StatusBadRequest)
		return
	}
	if newuserid_str <= 0 {
		http.Error(w, "invalid request: invalid user id", http.StatusBadRequest)
		return
	}

	newuserid := database.Id(newuserid_str)
	err = s.db.RemoveUserFromChannel(database.Id(channelid), newuserid)
	if errors.Is(err, database.ErrRecordNotFound) {
		http.Error(w, "user not in channel", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "error: unable to remove user from channel", http.StatusBadRequest)
		return
	}
}

func (s *Server) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	message_id_str := r.PathValue("messageid")
	messageid, err := strconv.Atoi(message_id_str)
	if err != nil {
		http.Error(w, "error: unable to parse messageid", http.StatusBadRequest)
		return
	}
	if messageid < 0 {
		http.Error(w, "error: invalid messageid", http.StatusBadRequest)
		return
	}
	message, err := s.db.GetMessage(database.Id(messageid))
	if err != nil {
		http.Error(w, "error: unable to fetch message", http.StatusBadRequest)
		return
	}
	if message.UserId != userid {
		http.Error(w, "error: attempting to modify different user message", http.StatusBadRequest)
		return
	}
	in_channel, err := s.db.IsUserInChannel(userid, message.ChannelId)
	if err != nil {
		http.Error(w, "error: unable to verify channel permissions", http.StatusBadRequest)
		return
	}
	if !in_channel {
		http.Error(w, "error: user not in channel", http.StatusBadRequest)
		return
	}

	message_data := struct {
		Message string `json:"message"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&message_data)
	if err != nil {
		http.Error(w, "error: unable to parse message from body", http.StatusBadRequest)
		return
	}
	err = s.db.UpdateMessage(message.MessageId, message_data.Message)
	if err != nil {
		http.Error(w, "error: issue while updating message", http.StatusBadRequest)
		return
	}
}

func (s *Server) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) CreateChannel(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) CreateChannelMessage(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) GetChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	_, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	channelid_str := r.PathValue("channelid")
	channelid, err := strconv.Atoi(channelid_str)
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	if channelid <= 0 {
		http.Error(w, "invalid request: invalid server id", http.StatusBadRequest)
		return
	}
	channel_info, err := s.db.GetChannel(database.Id(channelid))
	if err != nil {
		http.Error(w, "error: unable to locate channel", http.StatusBadRequest)
		return
	}

	payload := struct {
		ChannelId   database.Id `json:"channelid"`
		ServerId    database.Id `json:"serverid"`
		ChannelName string      `json:"channelname"`
		Timestamp   time.Time   `json:"timestamp"`
	}{
		ChannelId:   channel_info.ChannelId,
		ServerId:    channel_info.ServerId,
		ChannelName: channel_info.ChannelName,
		Timestamp:   channel_info.Timestamp,
	}
	jsonResp, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) GetMessage(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) GetChannelMessages(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data (default for HTML form submission)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	loginData := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&loginData)
	username := loginData.Username
	password := loginData.Password
	userid, err := s.db.GetUserIDFromUserName(username)
	if err != nil {
		http.Error(w, "unable to locate username", http.StatusBadRequest)
		return
	}

	valid, err := s.db.ValidateUserLoginInfo(userid, password)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	token, _, err := s.db.UpdateUserSessionToken(userid)
	if err != nil {
		http.Error(w, "unable to update session token", http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{
		"userid": userid,
	}

	// Set a cookie (you can modify the cookie as needed)
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect the user to /chat
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(
			w,
			"internal error: unable to send encoded response",
			http.StatusInternalServerError,
		)
		return
	}
}

func (s *Server) createNewServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	newServerData := struct {
		ServerName string `json:"servername"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&newServerData)
	if err != nil {
		http.Error(w, "unable to parse request", http.StatusBadRequest)
		return
	}
	if len(newServerData.ServerName) > 30 {
		http.Error(w, "server name too long", http.StatusBadRequest)
		return
	}
	if len(newServerData.ServerName) < 3 {
		http.Error(w, "server name too short", http.StatusBadRequest)
		return
	}

	serverid, err := s.db.CreateServer(userid, newServerData.ServerName)
	if err != nil {
		http.Error(w, "unable to create server", http.StatusBadRequest)
		return
	}
	resp := map[string]interface{}{
		"serverid": serverid,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(
			w,
			"internal error: unable to send encoded response",
			http.StatusInternalServerError,
		)
		return
	}

	return
}

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data (default for HTML form submission)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	loginData := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&loginData)
	username := loginData.Username
	password := loginData.Password
	userid, err := s.db.CreateUser(username, password)
	if errors.Is(err, database.ErrRecordAlreadyExists) {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "unable to create user", http.StatusBadRequest)
		return
	}

	valid_user, err := s.db.ValidateUserLoginInfo(userid, password)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if !valid_user {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	token, _, err := s.db.UpdateUserSessionToken(userid)
	if err != nil {
		http.Error(w, "unable to update session token", http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{
		"userid": userid,
	}

	// Set a cookie (you can modify the cookie as needed)
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect the user to /chat
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(
			w,
			"internal error: unable to send encoded response",
			http.StatusInternalServerError,
		)
		return
	}
}

func (s *Server) GetServerChannels(w http.ResponseWriter, r *http.Request) {
	serverid_str := r.PathValue("serverid")
	serverid, err := strconv.Atoi(serverid_str)
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	if serverid <= 0 {
		http.Error(w, "invalid request: invalid server id", http.StatusBadRequest)
		return
	}

	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userinfo, err := s.db.GetUser(userid)
	isServerMember, err := s.db.IsUserInServer(userinfo.UserId, database.Id(serverid))
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	if !isServerMember {
		http.Error(w, "user not member of server", http.StatusNetworkAuthenticationRequired)
		return
	}

	channels, err := s.db.GetChannelsOfServer(database.Id(serverid))
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{"channels": channels}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) GetServersOfUser(w http.ResponseWriter, r *http.Request) {
	userid_str := r.PathValue("userid")
	userid_int, err := strconv.Atoi(userid_str)
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	if userid_int <= 0 {
		http.Error(w, "invalid request: invalid server id", http.StatusBadRequest)
		return
	}
	userid := database.Id(userid_int)

	autheduserid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if userid != autheduserid {
		http.Error(w, "invalid request: invalid user id", http.StatusBadRequest)
		return
	}

	servers, err := s.db.GetServersOfUser(database.Id(userid))
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{"servers": servers}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) GetServerInformation(w http.ResponseWriter, r *http.Request) {
	serverid_str := r.PathValue("serverid") // r.URL.Query().Get("serverid")
	serverid, err := strconv.Atoi(serverid_str)
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	server, err := s.db.GetServer(database.Id(serverid))
	if errors.Is(err, database.ErrRecordNotFound) {
		// TODO: figure out proper method for valid resopnse but no data
		http.Error(w, "server not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "invalid request: unable to find server", http.StatusNotFound)
		return
	}

	jsonstruct := struct {
		ServerId   database.Id `json:"serverid"`
		OwnerId    database.Id `json:"ownerid"`
		ServerName string      `json:"servername"`
	}{
		ServerId:   server.ServerId,
		OwnerId:    server.OwnerId,
		ServerName: server.ServerName,
	}
	jsonResp, err := json.Marshal(jsonstruct)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) GetServerMessages(w http.ResponseWriter, r *http.Request) {
	serverid_str := r.PathValue("serverid")
	serverid, err := strconv.Atoi(serverid_str)
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	if serverid <= 0 {
		http.Error(w, "invalid request: invalid server id", http.StatusBadRequest)
		return
	}
	count_str := r.URL.Query().Get("count")
	var count uint = 30
	if count_str != "" {
		tempcount, err := strconv.Atoi(count_str)
		if err != nil {
			http.Error(w, "invalid request: unable to parse count", http.StatusBadRequest)
			return
		}
		if tempcount > 0 {
			count = uint(tempcount)
		} else {
			http.Error(w, "invalid request: invalid count", http.StatusBadRequest)
			return
		}
	}

	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	isServerMember, err := s.db.IsUserInServer(userid, database.Id(serverid))
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	if !isServerMember {
		http.Error(w, "user not member of server", http.StatusNetworkAuthenticationRequired)
		return
	}

	channels, err := s.db.GetChannelsOfServer(database.Id(serverid))
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	var messages []Message
	for _, channel := range channels {
		db_messages, err := s.db.GetMessagesInChannel(channel.ChannelId, count)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		tempmsgs := make([]Message, len(db_messages))
		for i, dbmsg := range db_messages {
			userinfo, err := s.db.GetUser(dbmsg.UserId)
			if err != nil {
				http.Error(w, "database error", http.StatusInternalServerError)
				return
			}
			username := userinfo.UserName
			tempmsgs[i] = Message{
				UserId:    dbmsg.UserId,
				UserName:  username,
				MessageID: dbmsg.MessageId,
				ChannelId: dbmsg.ChannelId,
				Message:   dbmsg.Contents,
				Date:      dbmsg.Timestamp.Format(time.UnixDate),
			}
		}
		messages = append(messages, tempmsgs...)

	}
	resp := map[string]interface{}{"serverid": serverid, "messages": messages}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) GetServerMembersHandler(w http.ResponseWriter, r *http.Request) {
	serverid_str := r.PathValue("serverid") // r.URL.Query().Get("serverid")
	serverid, err := strconv.Atoi(serverid_str)
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	if serverid <= 0 {
		http.Error(w, "invalid request: invalid server id", http.StatusBadRequest)
		return
	}

	users, err := s.db.GetUsersOfServer(database.Id(serverid))
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{"users": users, "serverid": serverid}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userid_str := r.PathValue("userid") // r.URL.Query().Get("userid")
	userid, err := strconv.Atoi(userid_str)
	if err != nil {
		http.Error(w, "invalid request: unable to parse user id", http.StatusBadRequest)
		return
	}
	if userid <= 0 {
		http.Error(w, "invalid request: invalid userid", http.StatusBadRequest)
		return
	}

	user, err := s.db.GetUser(database.Id(userid))
	if errors.Is(err, database.ErrRecordNotFound) {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	resp := map[string]string{"username": user.UserName}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	passinfo, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userinfo, err := s.db.GetUser(passinfo)
	if err != nil {
		http.Error(w, "error fetching user", http.StatusInternalServerError)
		return
	}

	opts := websocket.AcceptOptions{InsecureSkipVerify: true}
	socket, err := websocket.Accept(w, r, &opts)
	if err != nil {
		http.Error(w, "Failed to open websocket", http.StatusInternalServerError)
		return
	}
	defer socket.Close(websocket.StatusGoingAway, "Server closing websocket")
	outgoingChannel := make(chan Message)
	go s.handleMessages(socket, outgoingChannel)
	clients = append(clients, outgoingChannel)

	fmt.Printf("starting websocket loop: %d ms",
		time.Since(startTime).Milliseconds(),
	)
	newerr := websocket.CloseError{}
	for {
		_, message, err := socket.Read(r.Context())
		if errors.As(err, &newerr) {
			if newerr.Code == websocket.StatusGoingAway {
				fmt.Println("socket closed by user with status go away")
			} else if newerr.Code == websocket.StatusNoStatusRcvd {
				fmt.Println("socket closed due to not getting a response")
			} else {
				fmt.Printf("websocketHandler error: %v", err)
			}
			break
		} else if err != nil {
			fmt.Printf("error getting message from websocket: %v", err)
			break
		}

		// todo add message parsing
		payload := struct {
			Message    string      `json:"message"`
			ChannnelId database.Id `json:"channel_id"`
		}{
			Message:    "",
			ChannnelId: 0,
		}
		err = json.Unmarshal(message, &payload)
		if err != nil {
			fmt.Printf("error getting message from websocket: %e", err)
			continue
		}
		if payload.ChannnelId <= 0 {
			fmt.Printf("websocketHandler: invalid channel id channe_id=%d", payload.ChannnelId)
			continue
		}
		if len(payload.Message) > 1000 {
			fmt.Printf("format error: length of message to large length=%d", len(payload.Message))
			continue
		}
		messageid, err := s.db.AddMessage(payload.ChannnelId, userinfo.UserId, payload.Message)
		if err != nil {
			fmt.Printf("error saving message: %e", err)
			continue
		}
		dbmsg, err := s.db.GetMessage(messageid)
		if err != nil {
			fmt.Printf("error saving message: %e", err)
			continue
		}

		msg := Message{
			UserName:  userinfo.UserName,
			UserId:    dbmsg.UserId,
			MessageID: messageid,
			ChannelId: dbmsg.ChannelId,
			Message:   dbmsg.Contents,
			Date:      dbmsg.Timestamp.Format(time.UnixDate),
		}
		incomingChannel <- msg
	}
}

func (s *Server) handleIncomingMessages(broadcast chan Message) {
	for {
		message := <-broadcast
		for _, ch := range clients {
			ch <- message
		}
	}
}

func (s *Server) handleMessages(client *websocket.Conn, broadcast chan Message) {
	for {
		msg := <-broadcast
		jsondata, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("handleMessage %v", err)
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 1.0*time.Second)
		defer cancel()
		client.Write(ctx, websocket.MessageText, jsondata)
	}
}

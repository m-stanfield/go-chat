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
	clients         = make(map[database.Id]chan ServerMessage, 0)
	incomingChannel = make(chan ServerMessage)
	startTime       = time.Now()
)

type ServerMessage struct {
	UserId    database.Id `json:"userid"`
	MessageID database.Id `json:"messageid"`
	ChannelId database.Id `json:"channelid"`
	ServerId  database.Id `json:"serverid"`
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

func (s *Server) RegisterRoutes(logserver bool) http.Handler {
	go s.handleIncomingMessages(incomingChannel)
	mux := http.NewServeMux()

	// Register routes

	mux.HandleFunc("/", s.redirectToReact)
	mux.HandleFunc("/websocket", s.WithAuthUser(s.websocketHandler))

	mux.HandleFunc("POST /api/auth/login", s.loginHandler)
	mux.HandleFunc("POST /api/auth/session", s.WithAuthUser(s.sessionHandler))
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

	handler := http.Handler(mux)
	if logserver {
		handler = s.logEndpoint(handler)
	}
	// Wrap the mux with CORS middleware
	return s.corsMiddleware(handler)
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

	serverid, err := parsePathFromID(r, "serverid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	server_info, err := s.db.GetServer(serverid)
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

	err = s.db.UpdateServerName(serverid, new_server_name.ServerName)
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

	serverid, err := parsePathFromID(r, "serverid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	server_info, err := s.db.GetServer(serverid)
	if err != nil {
		http.Error(w, "error: unable to locate server", http.StatusBadRequest)
		return
	}
	if server_info.OwnerId != userid {
		http.Error(w, "error: user not owner of server", http.StatusBadRequest)
		return
	}
	err = s.db.DeleteServer(serverid)
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

	channelid, err := parsePathFromID(r, "channelid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	channel_info, err := s.db.GetChannel(channelid)
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
	err = s.db.UpdateChannel(channelid, new_channel_name.UpdatedChannelName)
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
	channelid, err := parsePathFromID(r, "channelid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	inchannel, err := s.db.IsUserInChannel(userid, channelid)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), http.StatusBadRequest)
		return
	}
	if !inchannel {
		http.Error(w, "user not in channel", http.StatusBadRequest)
		return
	}
	users, err := s.db.GetUsersInChannel(channelid)
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

	channelid, err := parsePathFromID(r, "channelid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	channel, err := s.db.GetChannel(channelid)
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
	inserver, err := s.db.IsUserInServer(newuserid, server_info.ServerId)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), http.StatusBadRequest)
		return
	}
	if !inserver {
		http.Error(w, "user not in server", http.StatusBadRequest)
		return
	}
	err = s.db.AddUserToChannel(newuserid, server_info.ServerId)
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
	channelid, err := parsePathFromID(r, "channelid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	channel, err := s.db.GetChannel(channelid)
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
	messageid, err := parsePathFromID(r, "messageid")
	if err != nil {
		http.Error(w, "error: unable to parse messageid", http.StatusBadRequest)
		return
	}
	message, err := s.db.GetMessage(messageid)
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
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	messageid, err := parsePathFromID(r, "messageid")
	if err != nil {
		http.Error(w, "error: unable to parse messageid", http.StatusBadRequest)
		return
	}
	message, err := s.db.GetMessage(messageid)
	if err != nil {
		http.Error(w, "error: unable to fetch message", http.StatusBadRequest)
		return
	}
	if message.UserId != userid {

		http.Error(w, "error: attempting to modify different user message", http.StatusBadRequest)
		return
	}
	err = s.db.DeleteMessage(message.MessageId)
	if err != nil {
		http.Error(w, "error: issue while deleting message", http.StatusBadRequest)
		return
	}
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := s.db.GetUser(userid)
	if err != nil {
		http.Error(w, "error: unable to fetch user", http.StatusBadRequest)
		return
	}
	err = s.db.UpdateUserName(user.UserId, user.UserName)
	if err != nil {
		http.Error(w, "error: unable to update username", http.StatusBadRequest)
		return
	}
}

func (s *Server) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	channelid, err := parsePathFromID(r, "channelid")
	if err != nil {
		http.Error(w, "error: unable to parse channelid", http.StatusBadRequest)
		return
	}
	channel, err := s.db.GetChannel(channelid)
	if err != nil {
		http.Error(w, "error: unable to fetch channel", http.StatusBadRequest)
		return
	}
	server, err := s.db.GetServer(channel.ServerId)
	if err != nil {
		http.Error(w, "error: unable to fetch server", http.StatusBadRequest)
		return
	}
	if server.OwnerId != userid {
		http.Error(w, "error: user not owner of server", http.StatusBadRequest)
		return
	}
	err = s.db.DeleteChannel(channel.ChannelId)
	if errors.Is(err, database.ErrRecordNotFound) {
		http.Error(w, "error: unable to locate channel", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "error: unable to delete channel", http.StatusBadRequest)
		return
	}
}

func (s *Server) CreateChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	channel_data := struct {
		ChannelName string `json:"channelname"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&channel_data)
	if err != nil {
		http.Error(w, "error: unable to parse request", http.StatusBadRequest)
		return
	}
	channelid, err := s.db.AddChannel(userid, channel_data.ChannelName)
	if err != nil {
		http.Error(w, "error: unable to create channel", http.StatusBadRequest)
		return
	}
	resp := map[string]interface{}{
		"channelid": channelid,
	}
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

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.db.DeleteUserSessionToken(userid)
	if err != nil {
		http.Error(w, "error: unable to delete session token", http.StatusBadRequest)
		return
	}
	cookie := &http.Cookie{
		Name:     "token", // Replace with your actual cookie name
		Value:    "",
		Path:     "/", // Ensure this matches the cookie's original path
		HttpOnly: true,
		Secure:   false, // Set to true if your site uses HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-time.Hour), // Set the expiration time to the past
		MaxAge:   -1,                         // Set MaxAge to 0 or a negative value to delete the cookie immediately
	}

	// Set the expired cookie in the response header.
	http.SetCookie(w, cookie)
	return
}

func (s *Server) CreateChannelMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	channelid, err := parsePathFromID(r, "channelid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	inchannel, err := s.db.IsUserInChannel(userid, channelid)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), http.StatusBadRequest)
		return
	}
	if !inchannel {
		http.Error(w, "user not in channel", http.StatusBadRequest)
		return
	}
	message_data := struct {
		Message string `json:"message"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&message_data)
	if err != nil {
		http.Error(w, "error: unable to parse request", http.StatusBadRequest)
		return
	}
	messageid, err := s.db.AddMessage(userid, channelid, message_data.Message)
	if err != nil {
		http.Error(w, "error: unable to create message", http.StatusBadRequest)
		return
	}
	resp := map[string]interface{}{
		"messageid": messageid,
	}
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

func (s *Server) GetChannelMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	channelid, err := parsePathFromID(r, "channelid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}

	inchannel, err := s.db.IsUserInChannel(userid, channelid)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), http.StatusBadRequest)
		return
	}
	if !inchannel {
		http.Error(w, "user not in channel", http.StatusBadRequest)
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

	messages, err := s.db.GetMessagesInChannel(channelid, count)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{"messages": messages}
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

	channelid, err := parsePathFromID(r, "channelid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	channel_info, err := s.db.GetChannel(channelid)
	if errors.Is(err, database.ErrRecordNotFound) {
		http.Error(w, "error: unable to locate channel", http.StatusNotFound)
		return
	}
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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	messageid, err := parsePathFromID(r, "messageid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	dbmessage, err := s.db.GetMessage(messageid)
	if errors.Is(err, database.ErrRecordNotFound) {
		http.Error(w, "error: unable to locate message", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error: internal server error", http.StatusBadRequest)
		return
	}
	message := fromDBMessageToSeverMessage(dbmessage)
	jsonResp, err := json.Marshal(message)
	if err != nil {
		http.Error(
			w,
			"error: internal server error. Unable to process request",
			http.StatusBadRequest,
		)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) sessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, "Unable to get user from context", http.StatusUnauthorized)
		return
	}

	user, err := s.db.GetUser(userid)
	if err != nil {
		http.Error(w, "Unable to find user", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"userid":   userid,
		"username": user.UserName,
	}

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
	serverid, err := parsePathFromID(r, "serverid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}

	userid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userinfo, err := s.db.GetUser(userid)
	isServerMember, err := s.db.IsUserInServer(userinfo.UserId, serverid)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	if !isServerMember {
		http.Error(w, "user not member of server", http.StatusNetworkAuthenticationRequired)
		return
	}

	channels, err := s.db.GetChannelsOfServer(serverid)
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
	userid, err := parsePathFromID(r, "userid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	autheduserid, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if userid != autheduserid {
		http.Error(w, "invalid request: invalid user id", http.StatusBadRequest)
		return
	}

	servers, err := s.db.GetServersOfUser(userid)
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
	serverid, err := parsePathFromID(r, "serverid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	server, err := s.db.GetServer(serverid)
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
	serverid, err := parsePathFromID(r, "serverid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
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
	isServerMember, err := s.db.IsUserInServer(userid, serverid)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	if !isServerMember {
		http.Error(w, "user not member of server", http.StatusNetworkAuthenticationRequired)
		return
	}

	channels, err := s.db.GetChannelsOfServer(serverid)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	var messages []ServerMessage
	for _, channel := range channels {
		db_messages, err := s.db.GetMessagesInChannel(channel.ChannelId, count)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		tempmsgs := make([]ServerMessage, len(db_messages))
		for i, dbmsg := range db_messages {
			tempmsgs[i] = ServerMessage{
				UserId:    dbmsg.UserId,
				MessageID: dbmsg.MessageId,
				ServerId:  serverid,
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
	serverid, err := parsePathFromID(r, "serverid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}

	users, err := s.db.GetUsersOfServer(serverid)
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
	userid, err := parsePathFromID(r, "userid")
	if err != nil {
		http.Error(w, "invalid request: unable to parse user id", http.StatusBadRequest)
		return
	}

	user, err := s.db.GetUser(userid)
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
	if _, ok := clients[userinfo.UserId]; ok {
		close(clients[userinfo.UserId])
		delete(clients, userinfo.UserId)
		log.Printf("new connected for %d. deleted previous channel", userinfo.UserId)
	}
	clients[userinfo.UserId] = make(chan ServerMessage, 0)
	defer func() {
		if _, ok := clients[userinfo.UserId]; ok {
			close(clients[userinfo.UserId])
			delete(clients, userinfo.UserId)
		}
		log.Printf("deleted channel for user %d", userinfo.UserId)
	}()
	go s.handleMessages(userinfo.UserId, socket, clients[userinfo.UserId])

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

		msg := ServerMessage{
			UserId:    dbmsg.UserId,
			MessageID: messageid,
			ServerId:  dbmsg.ServerId,
			ChannelId: dbmsg.ChannelId,
			Message:   dbmsg.Contents,
			Date:      dbmsg.Timestamp.Format(time.UnixDate),
		}
		incomingChannel <- msg
	}

	log.Printf("websocketHandler: closing channel for user %d", userinfo.UserId)
}

func (s *Server) handleIncomingMessages(messages chan ServerMessage) {
	for {
		message := <-messages
		for _, ch := range clients {
			ch <- message
		}
	}
}

func handleMessageContext(user database.Id, client *websocket.Conn, msg ServerMessage) error {
	if user == msg.UserId {
		log.Printf("skipping message to self %d", user)
	}
	log.Printf("sending message to %d", user)
	jsondata, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1.0*time.Second)
	defer cancel()
	err = client.Write(ctx, websocket.MessageText, jsondata)
	log.Printf("message sent %d", user)
	return err
}

func (s *Server) handleMessages(
	user database.Id,
	client *websocket.Conn,
	messageChan chan ServerMessage,
) {
	for {
		msg := <-messageChan
		err := handleMessageContext(user, client, msg)
		if err != nil {
			log.Printf("handleMessage %v", err)
		}
	}
}

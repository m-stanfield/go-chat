package server

import (
	"context"
	"encoding/json"
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
	counter         = 0
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

	mux.HandleFunc("/health", s.healthHandler)

	mux.HandleFunc("/websocket", s.websocketHandler)

	mux.HandleFunc("POST /api/login", s.loginHandler)
	mux.HandleFunc("GET /api/user/{userid}", s.GetUserHandler)
	mux.HandleFunc("GET /api/server/{serverid}/members", s.GetServerMembersHandler)
	mux.HandleFunc("GET /api/user/{userid}/servers", s.GetMemberServers)

	mux.HandleFunc("GET /api/server/{serverid}", s.GetServerInformation)
	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().
			Set("Access-Control-Allow-Origin", "*")
			// Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().
			Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().
			Set("Access-Control-Allow-Credentials", "false")
			// Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) validSession(userinfo database.UserLoginInfo, usertoken string) bool {
	if time.Now().After(userinfo.TokenExpireTime) {
		return false
	}
	return userinfo.Token == usertoken
}

func (s *Server) comparePassword(userinfo database.UserLoginInfo, password string) bool {
	return (password + userinfo.Salt) == (userinfo.PasswordHash + userinfo.Salt)
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
	passwordInfo, err := s.db.GetUserLoginInfo(userid)
	if err != nil {
		http.Error(w, "unable to locate password", http.StatusBadRequest)
		return
	}

	if !s.comparePassword(passwordInfo, password) {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	token, expire_time, err := s.db.UpdateUserSessionToken(userid)
	if err != nil {
		http.Error(w, "unable to update session token", http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{
		"token":             token,
		"token_expire_time": expire_time,
		"userid":            userid,
	}

	// Set a cookie (you can modify the cookie as needed)

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

func (s *Server) GetMemberServers(w http.ResponseWriter, r *http.Request) {
	userid_str := r.PathValue("userid")
	userid, err := strconv.Atoi(userid_str)
	if err != nil {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	if userid <= 0 {
		http.Error(w, "invalid request: invalid server id", http.StatusBadRequest)
		return
	}

	servers, err := s.db.GetServersOfUser(database.Id(userid))
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{"servers": servers, "userid": userid}
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

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}
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

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	userid_str := r.URL.Query().Get("userid")

	userid, err := strconv.Atoi(userid_str)
	if err != nil || userid < 0 {
		http.Error(w, "invalid request: unable to parse server id", http.StatusBadRequest)
		return
	}
	passwordInfo, err := s.db.GetUserLoginInfo(database.Id(userid))
	if err != nil {
		http.Error(w, "unable to locate password", http.StatusBadRequest)
		return
	}

	if token != passwordInfo.Token || passwordInfo.TokenExpireTime.Before(time.Now()) {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}
	userinfo, err := s.db.GetUser(database.Id(userid))
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

	for {
		_, message, err := socket.Read(r.Context())
		if err != nil {
			fmt.Printf("error getting message from websocket: %e", err)
			break
		}

		// todo add message parsing
		payload := struct {
			Message string `json:"message"`
		}{
			Message: "",
		}
		err = json.Unmarshal(message, &payload)
		if err != nil {
			fmt.Printf("error getting message from websocket: %e", err)
			continue
		}
		msg := Message{
			UserName: userinfo.UserName,
			UserId:   userinfo.UserId,
			MessageID: database.Id(
				counter,
			), Message: payload.Message, Date: time.Now().Format(time.UnixDate),
		}
		counter = counter + 1
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
			fmt.Println(err)
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		client.Write(ctx, websocket.MessageText, jsondata)
		if err != nil {
			fmt.Println(err)
		}
	}
}

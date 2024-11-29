package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/m-stanfield/go-chat/internal/database"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]int)
var broadcast = make(chan Message)

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

type ChatServer struct {
	db     *database.Database
	tokens map[database.Id]string
}

func main() {
	filename := "./database.db"
	sqldb, err := sql.Open("sqlite3", filename)
	if err != nil {
		return
	}
	var db = database.NewDatabase(sqldb)
	defer db.Close()
	tokens := make(map[database.Id]string)
	server := ChatServer{db: db, tokens: tokens}
	user, err := db.GetUser(1)
	if err != nil {
		return
	}
	fmt.Println(user.UserName)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/home.html")
	})
	http.HandleFunc("/server", server.handleConnections)
	http.HandleFunc("/users", server.getUsers)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/chat.html")
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/login.html")
	})
	http.HandleFunc("POST /login", server.handleLogin)
	go server.handleMessages()

	fmt.Println("Server started on :8080")
	// Initialize CORS with options
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},                   // Allow React app on port 3000
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed HTTP methods
		AllowedHeaders:   []string{"Content-Type"},                            // Allowed headers
		AllowCredentials: true,                                                // If you need credentials (cookies, auth headers)
	})

	// Wrap the handler with CORS middleware
	handler := corsHandler.Handler(http.DefaultServeMux)
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}

func (s *ChatServer) comparePassword(userinfo database.UserLoginInfo, password string) bool {
	return (password + userinfo.Salt) == (userinfo.PasswordHash + userinfo.Salt)
}

func (s *ChatServer) validSession(userinfo database.UserLoginInfo, usertoken string) bool {
	if time.Now().After(userinfo.TokenExpireTime) {
		return false
	}
	return userinfo.Token == usertoken
}
func (s *ChatServer) getUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "error with getting json users", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	users, err := s.db.GetUsersOfServer(1)
	if err != nil {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var jsonusers []User
	for _, u := range users {
		user := User{UserID: u.UserId, UserName: u.UserName}
		jsonusers = append(jsonusers, user)
	}
	jsonstr, err := json.Marshal(jsonusers)
	if err != nil {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write(jsonstr)

}
func (s *ChatServer) homePage(w http.ResponseWriter, r *http.Request) {

	// Redirect the user to /chat
	http.Redirect(w, r, "/chat", http.StatusSeeOther) // 303 redirect
}
func (s *ChatServer) handleLogin(w http.ResponseWriter, r *http.Request) {
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

	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password: %s\n", password)
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

	// Set a cookie (you can modify the cookie as needed)
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  expire_time,
		HttpOnly: false,
		Secure:   false, // Use false for local development, true for production with HTTPS
	})

	// Set a cookie (you can modify the cookie as needed)
	useridCookie := strconv.Itoa(int(userid))
	http.SetCookie(w, &http.Cookie{
		Name:     "userid",
		Value:    useridCookie,
		Path:     "/",
		Expires:  expire_time,
		HttpOnly: false,
		Secure:   false, // Use false for local development, true for production with HTTPS
	})

	// Redirect the user to /chat
	http.Redirect(w, r, "/chat", http.StatusSeeOther) // 303 redirect
}

func (s *ChatServer) handleConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("making connections")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	client_number := len(clients) + 1
	clients[conn] = client_number

	for {
		var msg SubmittedMessage
		_, msgText, err := conn.ReadMessage()
		fmt.Println("received message: ", string(msgText))
		if err != nil {
			fmt.Println(err)
			delete(clients, conn)
			return
		}
		err = json.Unmarshal(msgText, &msg)
		if err != nil {
			fmt.Println("Unable to parse", string(msgText), err)
			return
		}
		channelid := msg.ChannelId
		if err != nil {
			fmt.Println("Unable to create new message", string(msgText))
			return
		}
		userid, err := strconv.ParseUint(msg.UserID, 10, 64)
		if err != nil {
			fmt.Println("Unable to create new message", string(msgText))
			return
		}

		messageid, err := s.db.AddMessage(uint(channelid), uint(userid), msg.Message)
		if err != nil {
			fmt.Println("Unable to create new message", string(msgText))
			return
		}
		dbmessage, err := s.db.GetMessage(messageid)
		if err != nil {
			fmt.Println("error", err)
			return
		}
		user, err := s.db.GetUser(dbmessage.UserId)
		if err != nil {
			fmt.Println("error", err)
			return
		}
		var message Message
		message.UserName = user.UserName
		message.UserId = user.UserId
		message.MessageID = dbmessage.MessageId
		message.ChannelId = dbmessage.ChannelId
		message.Message = dbmessage.Contents
		message.Date = dbmessage.Timestamp.Format(time.UnixDate)

		broadcast <- message
	}
}

func (s *ChatServer) handleMessages() {
	for {
		msg := <-broadcast
		fmt.Println("Sending Message", msg)
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

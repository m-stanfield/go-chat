package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"go-chat-react/internal/database"
)

type UserService interface {
	GetUserIDFromUserName(username string) (database.Id, error)
	UpdateUserSessionToken(userid database.Id) (string, time.Time, error)
	DeleteUserSessionToken(userid database.Id) error
	GetUserLoginInfoFromToken(token string) (database.UserLoginInfo, error)
	GetUserLoginInfo(userid database.Id) (database.UserLoginInfo, error)
	ValidateUserLoginInfo(userid database.Id, password string) (bool, error)

	GetUser(userid database.Id) (database.User, error)
	CreateUser(username string, password string) (database.Id, error)
	UpdateUserName(userid database.Id, username string) error
	GetRecentUsernames(userid database.Id, number uint) ([]database.UsernameLogEntry, error)
}

type ServerService interface {
	GetUsersOfServer(serverid database.Id) ([]database.User, error)
	GetServersOfUser(userid database.Id) ([]database.Server, error)
	GetServer(serverid database.Id) (database.Server, error)
	CreateServer(ownerid database.Id, servername string) (database.Id, error)
	DeleteServer(serverid database.Id) error
	UpdateServerName(serverid database.Id, servername string) error
	IsUserInServer(userid database.Id, serverid database.Id) (bool, error)
}

type ChannelService interface {
	AddChannel(serverid database.Id, channelname string) (database.Id, error)
	DeleteChannel(channelid database.Id) error
	GetChannel(channelid database.Id) (database.Channel, error)
	GetChannelsOfServer(serverid database.Id) ([]database.Channel, error)
	UpdateChannel(channelid database.Id, username string) error
	AddUserToChannel(channelid database.Id, userid database.Id) error
	RemoveUserFromChannel(channelid database.Id, userid database.Id) error
	GetUsersInChannel(channelid database.Id) ([]database.User, error)
	IsUserInChannel(userid database.Id, channelid database.Id) (bool, error)
}

type MessageService interface {
	GetMessage(messageid database.Id) (database.Message, error)
	GetMessagesInChannel(channelid database.Id, number uint) ([]database.Message, error)
	AddMessage(channelid database.Id, userid database.Id, message string) (database.Id, error)
	UpdateMessage(messageid database.Id, message string) error
	DeleteMessage(messageid database.Id) error
}

type LifecycleService interface {
	Close() error
}

type (
	AtomicService interface {
		Service() Service
		Commit() error
		Rollback() error
	}

	Service interface {
		UserService
		ServerService
		ChannelService
		MessageService
		LifecycleService
	}
)

var (
	dburl      = os.Getenv("BLUEPRINT_DB_URL")
	dbInstance *database.DBService
)

func executeSQLFile(db *sql.DB, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	_, err = db.Exec(string(data))
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	return nil
}

func NewInMemoryDB() *database.DBService {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	err = executeSQLFile(db, "../../schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	err = executeSQLFile(db, "../../mockdata.sql")
	if err != nil {
		log.Fatal(err)
	}

	return database.New(db)
}

func (s *Server) Atomic(
	ctx context.Context,
	opts *sql.TxOptions,
) (*database.AtomitcDBService, error) {
	a, err := s.Atomic(ctx, opts)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func NewDB() *database.DBService {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	return database.New(db)
}

type Server struct {
	port int

	db Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	fmt.Printf("opening on port %d", port)
	db := NewDB()

	NewServer := &Server{
		port: port,

		db: db,
	}
	atomicdb, err := db.Atomic(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	atomicdb.Rollback()
	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

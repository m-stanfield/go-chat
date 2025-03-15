package server

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"go-chat-react/internal/database"
)

type httpErrorInfo struct {
	StatusCode int
	Message    string
}

type serverVerification struct {
	Validated bool
	UserId    database.Id
	Server    database.Server
}

func getUserIdFromContext(r *http.Request) (database.Id, error) {
	val := r.Context().Value("userid")
	if val == nil {
		return database.Id(0), errors.New("unable to get userid from context")
	}
	userid, ok := val.(database.Id)
	if !ok {
		return database.Id(0), errors.New("unable to get userid from context")
	}
	return userid, nil
}

func (s *Server) validSession(userinfo database.UserLoginInfo, usertoken string) bool {
	// if no token has been set
	if userinfo.Token == "" {
		return false
	}
	// if the token has expired
	if time.Now().After(userinfo.TokenExpireTime) {
		return false
	}
	// if the token is not the same as the one in the database
	return userinfo.Token == usertoken
}

func (s *Server) getServerFromRequest(r *http.Request) (database.Server, error) {
	serverid_str := r.PathValue("serverid")
	serverid, err := strconv.Atoi(serverid_str)
	if err != nil {
		return database.Server{}, errors.New("invalid request: unable to parse server id")
	}
	if serverid <= 0 {
		return database.Server{}, errors.New("invalid request: invalid server id")
	}
	server, err := s.db.GetServer(database.Id(serverid))
	if err != nil {
		return database.Server{}, errors.New("error: unable to locate server")
	}
	return server, nil
}

func (s *Server) getChannelFromRequest(r *http.Request) (database.Channel, error) {
	channelid_str := r.PathValue("channelid")
	channelid, err := strconv.Atoi(channelid_str)
	if err != nil {
		return database.Channel{}, errors.New("invalid request: unable to parse server id")
	}
	if channelid <= 0 {
		return database.Channel{}, errors.New("invalid request: invalid server id")
	}
	channel, err := s.db.GetChannel(database.Id(channelid))
	if err != nil {
		return database.Channel{}, errors.New("error: unable to locate server")
	}
	return channel, nil
}

func (s *Server) getServerFromChannel(channelid database.Id) (database.Server, error) {
	channel, err := s.db.GetChannel(channelid)
	if err != nil {
		return database.Server{}, err
	}
	server, err := s.db.GetServer(channel.ServerId)
	if err != nil {
		return database.Server{}, err
	}
	return server, nil
}

package server

import (
	"errors"
	"fmt"
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

func parsePathFromID(r *http.Request, field string) (database.Id, error) {
	fieldIDStr := r.PathValue(field)
	fieldID, err := strconv.Atoi(fieldIDStr)
	if err != nil {
		return 0, fmt.Errorf("invalid request: unable to parse %s", field)
	}
	if fieldID <= 0 {
		return 0, fmt.Errorf("invalid request: valid %s id requires >=0", field)
	}
	return database.Id(fieldID), nil
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

func (s *Server) GetServerFromRequest(r *http.Request) (database.Server, error) {
	serverid, err := parsePathFromID(r, "serverid")
	if err != nil {
		return database.Server{}, errors.New("invalid request: unable to parse server id")
	}
	server, err := s.db.GetServer(database.Id(serverid))
	if err != nil {
		return database.Server{}, errors.New("error: unable to locate server")
	}
	return server, nil
}

func (s *Server) GetChannelFromRequest(r *http.Request) (database.Channel, error) {
	channelid, err := parsePathFromID(r, "channelid")
	if err != nil {
		return database.Channel{}, errors.New("invalid request: unable to parse server id")
	}
	channel, err := s.db.GetChannel(channelid)
	if err != nil {
		return database.Channel{}, errors.New("error: unable to locate server")
	}
	return channel, nil
}

func (s *Server) GetServerFromChannel(channelid database.Id) (database.Server, error) {
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

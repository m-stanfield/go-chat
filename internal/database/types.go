package database

import (
	"strconv"
	"time"
)

type Id = uint

func ParseIntToID(id int) (Id, error) {
	if id <= 0 {
		return Id(0), ErrUnsupportedNegativeValue
	}
	return Id(id), nil
}

func ParseStringToID(id string) (Id, error) {
	intid, err := strconv.Atoi(id)
	if err != nil {
		return Id(0), ErrParsingValue
	}
	return ParseIntToID(intid)
}

type User struct {
	UserId   Id
	UserName string
}

type UserLoginInfo struct {
	UserId          Id
	PasswordHash    string
	Salt            string
	Token           string
	TokenExpireTime time.Time
}

type UsernameLogEntry struct {
	UserId    Id
	Username  string
	Timestamp time.Time
}

type UserNicknameLogEntry struct {
	UserId    Id
	ServerId  Id
	Nickname  string
	Timestamp time.Time
}

type Server struct {
	ServerId   Id
	OwnerId    Id
	ServerName string
}

type Channel struct {
	ChannelId   Id
	ServerId    Id
	ChannelName string
	Timestamp   time.Time
}

type Message struct {
	MessageId        Id
	UserId           Id
	ChannelId        Id
	Contents         string
	Timestamp        time.Time
	Editted          *bool
	EdittedTimeStamp *time.Time
}

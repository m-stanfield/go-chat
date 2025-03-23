package database

import (
	"context"
	"time"
)

type UserService interface {
	GetUserIDFromUserName(username string) (Id, error)
	UpdateUserSessionToken(userid Id) (string, time.Time, error)
	GetUserLoginInfoFromToken(token string) (UserLoginInfo, error)
	GetUserLoginInfo(userid Id) (UserLoginInfo, error)
	ValidateUserLoginInfo(userid Id, password string) (bool, error)

	GetUser(userid Id) (User, error)
	CreateUser(username string, password string) (Id, error)
	UpdateUserName(userid Id, username string) error
	GetRecentUsernames(userid Id, number uint) ([]UsernameLogEntry, error)
}

type ServerService interface {
	GetUsersOfServer(serverid Id) ([]User, error)
	GetServersOfUser(userid Id) ([]Server, error)
	GetServer(serverid Id) (Server, error)
	CreateServer(ownerid Id, servername string) (Id, error)
	DeleteServer(serverid Id) error
	UpdateServerName(serverid Id, servername string) error
	IsUserInServer(userid Id, serverid Id) (bool, error)
}

type ChannelService interface {
	AddChannel(serverid Id, channelname string) (Id, error)
	DeleteChannel(channelid Id) error
	GetChannel(channelid Id) (Channel, error)
	GetChannelsOfServer(serverid Id) ([]Channel, error)
	UpdateChannel(channelid Id, username string) error
	AddUserToChannel(channelid Id, userid Id) error
	RemoveUserFromChannel(channelid Id, userid Id) error
	GetUsersInChannel(channelid Id) ([]User, error)
	IsUserInChannel(userid Id, channelid Id) (bool, error)
}

type MessageService interface {
	GetMessage(messageid Id) (Message, error)
	GetMessagesInChannel(channelid Id, number uint) ([]Message, error)
	AddMessage(channelid Id, userid Id, message string) (Id, error)
	UpdateMessage(messageid Id, message string) error
	DeleteMessage(messageid Id) error
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
		Atomic(ctx context.Context) (AtomicService, error)
		UserService
		ServerService
		ChannelService
		MessageService
		LifecycleService
	}
)

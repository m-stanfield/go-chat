package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func loadDB(filename string, schemafile string, datafile string) (*sql.DB, error) {
	_, schemaExistError := os.Stat(schemafile)
	schemaExists := !errors.Is(schemaExistError, os.ErrNotExist)
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	if !schemaExists {
		return nil, fmt.Errorf("can't find schema file: %w", err)
	}
	schemaString, err := os.ReadFile(schemafile)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(string(schemaString))
	if err != nil {
		return nil, err
	}
	dataString, err := os.ReadFile(datafile)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(string(dataString))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func setup() *DBService {
	wd, _ := os.Getwd()
	fmt.Println(wd)
	db, err := loadDB(":memory:", "../../schema.sql", "../../mockdata.sql")
	if err != nil {
		panic("unable to initialize database for tests")
	}
	return &DBService{db: db, conn: db}
}

func Test_DeleteMessage(t *testing.T) {
	db := setup()
	defer db.Close()
	message := "1111"
	id, err := db.AddMessage(1, 1, message)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	err = db.DeleteMessage(id)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	_, err = db.GetMessage(id)
	if !errors.Is(err, ErrRecordNotFound) {
		t.Fatalf("TestA: err: %v", err)
	}
}

func Test_GetMessage_InvalidId(t *testing.T) {
	db := setup()
	defer db.Close()
	_, err := db.GetMessage(100000000)
	if !errors.Is(err, ErrRecordNotFound) {
		t.Fatalf("TestA: err: %v", err)
	}
}

func Test_UpdateMessage(t *testing.T) {
	db := setup()
	defer db.Close()
	initialMessage := "1111"
	newMessage := "u2jklfdsa"
	start_time := time.Now()

	id, err := db.AddMessage(1, 1, initialMessage)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	time.Sleep(1 * time.Millisecond)
	message, err := db.GetMessage(id)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	err = db.UpdateMessage(message.MessageId, newMessage)
	newmessage, err := db.GetMessage(id)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if newmessage.Contents != newMessage {
		t.Fatalf(
			"TestA: invalmessage message expected: %s got: %s",
			newMessage,
			newmessage.Contents,
		)
	}
	if newmessage.MessageId != id {
		t.Fatalf("TestA: invalmessage message expected: %d got: %d", id, newmessage.MessageId)
	}
	if newmessage.Timestamp != message.Timestamp {
		t.Fatalf(
			"TestA: invalmessage message expected: %v got: %v",
			start_time,
			newmessage.Timestamp,
		)
	}
	if !newmessage.EdittedTimeStamp.After(message.Timestamp) {
		t.Fatalf(
			"TestA: invalmessage message expected: %v got: %v",
			start_time,
			newmessage.EdittedTimeStamp,
		)
	}
}

func Test_GetUser(t *testing.T) {
	db := setup()
	expectedUsername := "u1"
	expectedId := 1

	user, err := db.GetUser(uint(expectedId))
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if user.UserId != uint(expectedId) {
		t.Fatalf("TestA: invalid id. expected: %d got: %d", 1, user.UserId)
	}
	if user.UserName != expectedUsername {
		t.Fatalf("TestA: invalid username expected: %s got: %s", expectedUsername, user.UserName)
	}
}

func Test_AtomicFail(t *testing.T) {
	adb := setup()
	defer adb.Close()
	atomic, err := adb.Atomic(context.Background())
	db := atomic.Service()

	u, err := db.CreateUser("aaaaa", "password")
	if err != nil {
		t.Fatalf("database - AtomicFail: errored while creating user %v", err)
	}
	err = atomic.Rollback()
	if err != nil {
		t.Fatalf("database - AtomicFail: errored during rollback %v", err)
	}

	_, err = adb.GetUser(u)
	if err == nil {
		t.Fatalf("fail")
	}
}

func Test_AtomicPass(t *testing.T) {
	/*
		db := setup()
		defer db.Close()
		var userId *Id = nil
		err := db.Atomic(context.Background(), func(db *DBService) error {
			u, err := db.CreateUser("aaaaa", "password")
			if err != nil {
				return err
			}
			userId = &u
			return nil
		})
		if err != nil {
			t.Fatalf("fail")
		}
		_, err = db.GetUser(*userId)
		if err != nil {
			t.Fatalf("fail")
		}
	*/
}

func Test_CreateUser(t *testing.T) {
	db := setup()
	defer db.Close()
	expectedUsername := "u2jjklfdsa"
	expectedPassword := "SFioj*&*(0"

	id, err := db.CreateUser(expectedUsername, expectedPassword)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	user, err := db.GetUser(id)
	if user.UserId != id {
		t.Fatalf("TestA: invaluser user. expected: %d got: %d", id, user.UserId)
	}
	if user.UserName != expectedUsername {
		t.Fatalf("TestA: invaluser username expected: %s got: %s", expectedUsername, user.UserName)
	}
	valid_user, err := db.ValidateUserLoginInfo(id, expectedPassword)
	if err != nil {
		t.Fatalf("Test Create User: err: %v", err)
	}
	if !valid_user {
		t.Fatalf("Test Create User: invalid user")
	}
}

func Test_UpdateServerName(t *testing.T) {
	db := setup()
	defer db.Close()
	id := Id(1)
	newServerName := "u2jklfdsa"

	server, err := db.GetServer(id)
	if server.ServerName == newServerName {
		t.Fatalf("Test_UpdateUserName: invalid test data, test user name matches previous")
	}
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	err = db.UpdateServerName(server.ServerId, newServerName)
	newserver, err := db.GetServer(id)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if newserver.ServerName != newServerName {
		t.Fatalf("TestA: invaluser username expected: %s got: %s", newServerName, server.ServerName)
	}
}

func Test_UpdateUserName(t *testing.T) {
	db := setup()
	defer db.Close()
	id := Id(1)
	newUserName := "u2jklfdsa"

	user, err := db.GetUser(id)
	if user.UserName == newUserName {
		t.Fatalf("Test_UpdateUserName: invalid test data, test user name matches previous")
	}
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	err = db.UpdateUserName(user.UserId, newUserName)
	newuser, err := db.GetUser(id)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if newuser.UserName != newUserName {
		t.Fatalf("TestA: invaluser username expected: %s got: %s", newUserName, user.UserName)
	}
}

func Test_AutopopulateUserNameLogEntries(t *testing.T) {
	db := setup()
	defer db.Close()
	userid := Id(1)
	previousLogs, err := db.GetRecentUsernames(userid, 1)
	if err != nil {
		t.Fatalf("Autopop user name log: err: %v", err)
	}
	if previousLogs == nil {
		t.Fatalf("Autopop user name log: initial log function nil")
	}
	if len(previousLogs) == 0 {
		t.Fatalf("Autopop user name log: initial logs length zero")
	}
	newUserName := previousLogs[0].Username + "abcd"
	err = db.UpdateUserName(userid, newUserName)
	if err != nil {
		t.Fatalf("Autopop user name log: err: %v", err)
	}
	newLogs, err := db.GetRecentUsernames(userid, 2)
	if err != nil {
		t.Fatalf("Autopop user name log: err: %v", err)
	}
	if newLogs[0].Username != newUserName {
		t.Fatalf("Username recorded in log incorrect")
	}

	if newLogs[1].Username != previousLogs[0].Username {
		t.Fatalf("Autopop user name  log: err: previous name not correctly grabbed")
	}
}

func Test_GetServersOfUser(t *testing.T) {
	db := setup()
	defer db.Close()
	userid := Id(1)

	server, err := db.GetServersOfUser(userid)
	if err != nil {
		t.Fatalf("Autopop user name log: err: %v", err)
	}
	if len(server) != 2 {
		t.Fatalf(
			"GetServersOfUser failed: incorrect number of servers expected: %d got: %d",
			len(server),
			2,
		)
	}
	found1, found2 := false, false
	for _, s := range server {
		if s.ServerId == 1 {
			found1 = true
		} else if s.ServerId == 2 {
			found2 = true
		}
	}
	if !(found1 && found2) {
		t.Fatalf("GetServersOfUser failed: invalid server ids")
	}
}

func Test_GetUsersOfServer(t *testing.T) {
	db := setup()
	defer db.Close()
	serverid := Id(1)

	users, err := db.GetUsersOfServer(serverid)
	if err != nil {
		t.Fatalf("Autopop user name log: err: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf(
			"GetUsersOfServer failed: incorrect number of servers expected: %d got: %d",
			len(users),
			2,
		)
	}
	found1, found2 := false, false
	for _, s := range users {
		if s.UserId == 1 {
			found1 = true
		} else if s.UserId == 3 {
			found2 = true
		}
	}
	if !(found1 && found2) {
		t.Fatalf("GetUsersOfServer failed: invalid server ids")
	}
}

func Test_GetChannelsOfServer(t *testing.T) {
	db := setup()
	defer db.Close()
	serverid := Id(1)

	channels, err := db.GetChannelsOfServer(serverid)
	if err != nil {
		t.Fatalf("channel of server log: err: %v", err)
	}
	if len(channels) != 2 {
		t.Fatalf(
			"GetUsersOfServer failed: incorrect number of servers expected: %d got: %d",
			len(channels),
			2,
		)
	}
	found1, found2 := false, false
	for _, s := range channels {
		if s.ChannelId == 1 && s.ChannelName == "a" {
			found1 = true
		} else if s.ChannelId == 2 && s.ChannelName == "b" {
			found2 = true
		}
	}
	if !(found1 && found2) {
		t.Fatalf("GetUsersOfServer failed: invalid server ids")
	}
}

func Test_GetMessage(t *testing.T) {
	db := setup()
	defer db.Close()
	channelid := Id(1)
	userid := Id(1)
	messageid := Id(1)
	message := "1111"
	user, err := db.GetMessage(messageid)
	if err != nil {
		t.Fatalf("TestA: invalid error %e", err)
	}
	if user.UserId != userid {
		t.Fatalf("TestA: invalid user id expected: %d got: %d", userid, user.UserId)
	}
	if user.MessageId != messageid {
		t.Fatalf("TestA: invalid message id expected: %d got: %d", messageid, user.MessageId)
	}
	if user.ChannelId != channelid {
		t.Fatalf("TestA: invalid channel id expected: %d got: %d", channelid, user.ChannelId)
	}
	if user.Contents != message {
		t.Fatalf("TestA: invalid message contents expected: %s got: %s", message, user.Contents)
	}
}

func Test_AddMessage(t *testing.T) {
	db := setup()
	defer db.Close()
	channelid := Id(1)
	userid := Id(1)
	message := "sdfjal"

	messageid, err := db.AddMessage(channelid, userid, message)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	user, err := db.GetMessage(messageid)
	if user.MessageId != messageid {
		t.Fatalf("TestA: invalid message id expected: %d got: %d", messageid, user.MessageId)
	}
	if user.ChannelId != channelid {
		t.Fatalf("TestA: invalid channel id expected: %d got: %d", channelid, user.ChannelId)
	}
	if user.Contents != message {
		t.Fatalf("TestA: invalid message contents expected: %s got: %s", message, user.Contents)
	}
}

func Test_AddChannel(t *testing.T) {
	db := setup()
	defer db.Close()
	serverid := Id(2)
	expectedChannelname := "u2jklfdsa"

	id, err := db.AddChannel(serverid, expectedChannelname)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	user, err := db.GetChannel(id)
	if user.ChannelId != id {
		t.Fatalf("TestA: invaluser user. expected: %d got: %d", id, user.ChannelId)
	}
	if user.ChannelName != expectedChannelname {
		t.Fatalf(
			"TestA: invaluser username expected: %s got: %s",
			expectedChannelname,
			user.ChannelName,
		)
	}
}

func Test_GetChannel(t *testing.T) {
	db := setup()
	expectedChannelname := "b"
	expectedId := 2

	channel, err := db.GetChannel(uint(expectedId))
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if channel.ChannelId != uint(expectedId) {
		t.Fatalf("TestA: invalid id. expected: %d got: %d", 1, channel.ChannelId)
	}
	if channel.ChannelName != expectedChannelname {
		t.Fatalf(
			"TestA: invalid username expected: %s got: %s",
			expectedChannelname,
			channel.ChannelName,
		)
	}
}

func Test_UpdateChannelName(t *testing.T) {
	db := setup()
	defer db.Close()
	id := Id(1)
	newChannelName := "u2jklfdsa"

	user, err := db.GetChannel(id)
	if user.ChannelName == newChannelName {
		t.Fatalf("Test_UpdateChannelName: invalid test data, test user name matches previous")
	}
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	err = db.UpdateChannel(user.ChannelId, newChannelName)
	newuser, err := db.GetChannel(id)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if newuser.ChannelName != newChannelName {
		t.Fatalf("TestA: invaluser username expected: %s got: %s", newChannelName, user.ChannelName)
	}
}

func Test_GetMessageInChannel_Two(t *testing.T) {
	db := setup()
	defer db.Close()
	var channelid Id = 1
	var number int = 2

	messages, err := db.GetMessagesInChannel(channelid, uint(number))
	if err != nil {
		t.Fatalf("GetMessageInChannel Number error: %v", err)
	}

	if len(messages) != number {
		t.Fatalf("incorrect number of messages: got %d expected %d", len(messages), number)
	}
}

func Test_GetMessageInChannel_Three(t *testing.T) {
	db := setup()
	defer db.Close()
	var channelid Id = 1
	var number int = 3

	messages, err := db.GetMessagesInChannel(channelid, uint(number))
	if err != nil {
		t.Fatalf("GetMessageInChannel Number error: %v", err)
	}

	if len(messages) != number {
		t.Fatalf("incorrect number of messages: got %d expected %d", len(messages), number)
	}
}

func Test_GetMessageInChannel_NoMessages(t *testing.T) {
	db := setup()
	defer db.Close()
	var channelid Id = 1
	var number int = 0

	messages, err := db.GetMessagesInChannel(channelid, uint(number))
	if err != nil {
		t.Fatalf("GetMessageInChannel Number error: %v", err)
	}

	if len(messages) != number {
		t.Fatalf("incorrect number of messages: got %d expected %d", len(messages), number)
	}
}

func Test_GetMessageInChannel_InvalidChannel(t *testing.T) {
	db := setup()
	defer db.Close()
	var channelid Id = 10007183090
	var number int = 3

	messages, err := db.GetMessagesInChannel(channelid, uint(number))
	if err != nil {
		t.Fatalf("GetMessageInChannel Number error: %v", err)
	}

	if len(messages) != 0 {
		t.Fatalf("incorrect number of messages: got %d expected %d", len(messages), number)
	}
}

func Test_CreateServer(t *testing.T) {
	db := setup()
	defer db.Close()
	expectedServerName := "u2jjklfdsa"
	ownerid := Id(118)

	id, err := db.CreateServer(ownerid, expectedServerName)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	server, err := db.GetServer(id)
	if server.ServerId != id {
		t.Fatalf("TestA: invalid id. expected: %d got: %d", id, server.ServerId)
	}
	if server.ServerName != expectedServerName {
		t.Fatalf(
			"TestA: invalid username expected: %s got: %s",
			expectedServerName,
			server.ServerName,
		)
	}
	if server.OwnerId != ownerid {
		t.Fatalf("TestA: invalid owner id expected: %d got: %d", ownerid, server.OwnerId)
	}
}

func Test_GetUserIDFromUserName(t *testing.T) {
	db := setup()
	defer db.Close()
	username := "u1"
	id, err := db.GetUserIDFromUserName(username)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if id != 1 {
		t.Fatalf("TestA: invalid id. expected: %d got: %d", 1, id)
	}
}

func Test_UpdateUserSessionToken(t *testing.T) {
	db := setup()
	defer db.Close()
	id := Id(1)
	start_time := time.Now()
	val, expiretime, err := db.UpdateUserSessionToken(id)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if val == "" {
		t.Fatalf("TestA: invalid token value")
	}
	if expiretime.Before(start_time) {
		t.Fatalf("TestA: invalid expire time")
	}
}

func Test_GetUserLoginInfoFromToken(t *testing.T) {
	db := setup()
	defer db.Close()
	id := Id(1)
	token, _, err := db.UpdateUserSessionToken(id)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	userid, err := db.GetUserLoginInfoFromToken(token)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if userid.UserId != id {
		t.Fatalf("TestA: invalid id. expected: %d got: %d", id, userid.UserId)
	}
}

func Test_IsUserInServer_Valid(t *testing.T) {
	db := setup()
	defer db.Close()
	userid := Id(121)
	serverid := Id(1)
	inserver, err := db.IsUserInServer(userid, serverid)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if inserver {
		t.Fatalf("TestA: invalid in server")
	}
}

func Test_IsUserInServer_InValid(t *testing.T) {
	db := setup()
	defer db.Close()
	userid := Id(1)
	serverid := Id(1)
	_, err := db.IsUserInServer(userid, serverid)
	inserver, err := db.IsUserInServer(userid, serverid)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if !inserver {
		t.Fatalf("TestA: valid in server should be invalid")
	}
}

func Test_AddUserToChannel_Valid(t *testing.T) {
	db := setup()
	defer db.Close()
	userid := Id(3)
	channelid := Id(1)
	inserver, err := db.IsUserInChannel(userid, channelid)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if inserver {
		t.Fatalf("TestA: invalid in server")
	}

	err = db.AddUserToChannel(userid, channelid)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	inserver, err = db.IsUserInChannel(userid, channelid)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if !inserver {
		t.Fatalf("TestA: valid in server should be invalid")
	}
}

func Test_IsUserInChannel_Valid(t *testing.T) {
	db := setup()
	defer db.Close()
	userid := Id(1)
	channelid := Id(1)
	inserver, err := db.IsUserInChannel(userid, channelid)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if !inserver {
		t.Fatalf("TestA: valid in server should be invalid")
	}
}

func Test_IsUserInChannel_Invalid(t *testing.T) {
	db := setup()
	defer db.Close()
	userid := Id(3)
	channelid := Id(1)
	inserver, err := db.IsUserInChannel(userid, channelid)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if inserver {
		t.Fatalf("TestA: invalid in server")
	}
}

func Test_GetUsersInChannel_Valid(t *testing.T) {
	db := setup()
	defer db.Close()
	serverid := Id(1)
	users, err := db.GetUsersInChannel(serverid)
	if err != nil {
		t.Fatalf("TestA: err: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("TestA: invalid number of users")
	}
	if users[0].UserId != 1 {
		t.Fatalf("TestA: invalid user id")
	}
	if users[0].UserName != "u1" {
		t.Fatalf("TestA: invalid user name")
	}

	if users[1].UserId != 2 {
		t.Fatalf("TestA: invalid user id")
	}
	if users[1].UserName != "u2" {
		t.Fatalf("TestA: invalid user name")
	}
}

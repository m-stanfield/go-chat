package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

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

func setup() *service {
	wd, _ := os.Getwd()
	fmt.Println(wd)
	db, err := loadDB(":memory:", "../../schema.sql", "../../mockdata.sql")
	if err != nil {
		panic("unable to initialize database for tests")
	}
	return &service{db: db, conn: db}
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
	db := setup()
	defer db.Close()
	var userId *Id = nil
	err := db.Atomic(context.Background(), func(db Service) error {
		u, err := db.AddUser("aaaaa")
		if err != nil {
			return err
		}
		userId = &u
		return errors.New("mock error for testing rollback")
	})
	if err == nil {
		t.Fatalf("fail")
	}
	_, err = db.GetUser(*userId)
	if err == nil {
		t.Fatalf("fail")
	}

}

func Test_AtomicPass(t *testing.T) {
	db := setup()
	defer db.Close()
	var userId *Id = nil
	err := db.Atomic(context.Background(), func(db Service) error {
		u, err := db.AddUser("aaaaa")
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

}

func Test_AddUser(t *testing.T) {
	db := setup()
	defer db.Close()
	expectedUsername := "u2jklfdsa"

	id, err := db.AddUser(expectedUsername)
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
		t.Fatalf("GetServersOfUser failed: incorrect number of servers expected: %d got: %d", len(server), 2)
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
		t.Fatalf("GetUsersOfServer failed: incorrect number of servers expected: %d got: %d", len(users), 2)
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
		t.Fatalf("GetUsersOfServer failed: incorrect number of servers expected: %d got: %d", len(channels), 2)
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
		t.Fatalf("TestA: invaluser username expected: %s got: %s", expectedChannelname, user.ChannelName)
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
		t.Fatalf("TestA: invalid username expected: %s got: %s", expectedChannelname, channel.ChannelName)
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
	err = db.UpdateChannelName(user.ChannelId, newChannelName)
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

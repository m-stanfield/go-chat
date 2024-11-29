package database

import (
	"fmt"
	"strconv"
	"time"
)

func (r *Database) GetUserIDFromUserName(username string) (Id, error) {
	rows, err := r.conn.Query("SELECT userid FROM UserTable WHERE username = ?", username)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	count := 0
	var userid Id
	for rows.Next() {
		count += 1
		if count > 1 {
			return 0, ErrMultipleRecords
		}
		err := rows.Scan(&userid)
		if err != nil {
			return 0, err
		}
	}
	if count == 0 {
		return 0, ErrNoRecord
	}
	return userid, nil
}

func (r *Database) UpdateUserSessionToken(userid Id) (string, time.Time, error) {
	token := "token" + strconv.FormatUint(uint64(userid), 10)
	expire := time.Now().Add(24 * time.Hour)
	_, err := r.conn.Exec("UPDATE UserLoginTable SET token = ?, token_expire_time = ? WHERE userid=? ", token, expire, userid)
	if err != nil {
		return "", expire, err
	}
	return token, expire, nil
}

func (r *Database) GetUserLoginInfo(userid Id) (UserLoginInfo, error) {
	rows, err := r.conn.Query("SELECT userid, passwordhash, salt, token, token_expire_time FROM UserLoginTable WHERE userid = ?", userid)
	if err != nil {
		return UserLoginInfo{}, err
	}
	defer rows.Close()
	count := 0
	var user UserLoginInfo
	for rows.Next() {
		count += 1
		if count > 1 {
			return UserLoginInfo{}, ErrMultipleRecords
		}
		err := rows.Scan(&user.UserId, &user.PasswordHash, &user.Salt, &user.Token, &user.TokenExpireTime)
		if err != nil {
			return UserLoginInfo{}, err
		}
	}
	if count == 0 {
		return UserLoginInfo{}, ErrNoRecord
	}
	return user, nil
}

func (r *Database) GetUser(userid Id) (User, error) {
	rows, err := r.conn.Query("SELECT userid, username FROM UserTable WHERE userid = ?", userid)
	if err != nil {
		return User{}, err
	}
	defer rows.Close()
	count := 0
	var user User
	for rows.Next() {
		count += 1
		if count > 1 {
			return User{}, ErrMultipleRecords
		}
		err := rows.Scan(&user.UserId, &user.UserName)
		if err != nil {
			return User{}, err
		}
	}
	if count == 0 {
		return User{}, ErrNoRecord
	}
	return user, nil
}

func (r *Database) AddUser(username string) (Id, error) {
	d, err := r.conn.Exec("INSERT INTO UserTable ( username) VALUES ( ?)", username)
	if err != nil {
		return 0, fmt.Errorf("add user - username: %s err: %w", username, err)
	}
	id, err := d.LastInsertId()
	if err != nil {
		return 0, err
	}
	if id < 0 {
		return 0, ErrNegativeRowIndex
	}
	return Id(id), nil
}

func (r *Database) UpdateUserName(userid Id, username string) error {
	_, err := r.conn.Exec("UPDATE UserTable SET username = ? WHERE userid=? ", username, userid)
	return err
}

func (r *Database) GetRecentUsernames(userid Id, number uint) ([]UsernameLogEntry, error) {
	rows, err := r.conn.Query("SELECT userid, username, timestamp FROM UserNameLogTable WHERE userid = ? ORDER BY timestamp DESC LIMIT ?", userid, number)
	if err != nil {
		return []UsernameLogEntry{}, err
	}
	defer rows.Close()
	var names []UsernameLogEntry
	for rows.Next() {
		var name UsernameLogEntry
		err := rows.Scan(&name.UserId, &name.Username, &name.Timestamp)
		if err != nil {
			return []UsernameLogEntry{}, err
		}
		names = append(names, name)
	}
	return names, nil
}

func (r *Database) GetUsersOfServer(serverid Id) ([]User, error) {
	rows, err := r.conn.Query("SELECT U.userid, U.username FROM UsersServerTable as US INNER JOIN UserTable as U ON US.userid = U.userid WHERE US.serverid = ?", serverid)
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()
	var names []User
	for rows.Next() {
		var name User
		err := rows.Scan(&name.UserId, &name.UserName)
		if err != nil {
			return []User{}, err
		}
		names = append(names, name)
	}
	return names, nil
}

func (r *Database) GetServersOfUser(userid Id) ([]Server, error) {
	rows, err := r.conn.Query("SELECT S.serverid, S.ownerid, S.servername FROM UsersServerTable as U INNER JOIN ServerTable as S ON U.serverid = S.serverid WHERE U.userid = ?", userid)
	if err != nil {
		return []Server{}, err
	}
	defer rows.Close()
	var servers []Server
	for rows.Next() {
		var s Server
		err := rows.Scan(&s.ServerId, &s.OwnerId, &s.ServerName)
		if err != nil {
			return []Server{}, err
		}
		servers = append(servers, s)
	}
	return servers, nil
}

func (r *Database) GetChannelsOfServer(serverid Id) ([]Channel, error) {
	rows, err := r.conn.Query("SELECT channelid, serverid, channelname, timestamp FROM ChannelTable WHERE serverid = ?", serverid)
	if err != nil {
		return []Channel{}, err
	}
	defer rows.Close()
	var servers []Channel
	for rows.Next() {
		var s Channel
		err := rows.Scan(&s.ChannelId, &s.ServerId, &s.ChannelName, &s.Timestamp)
		if err != nil {
			return []Channel{}, err
		}
		servers = append(servers, s)
	}
	return servers, nil

}
func (r *Database) AddChannel(serverid Id, channelname string) (Id, error) {
	d, err := r.conn.Exec("INSERT INTO ChannelTable ( serverid, channelname) VALUES ( ?, ?)", serverid, channelname)
	if err != nil {
		return 0, fmt.Errorf("add user - username: %s err: %w", channelname, err)
	}
	id, err := d.LastInsertId()
	if err != nil {
		return 0, err
	}
	if id < 0 {
		return 0, ErrNegativeRowIndex
	}
	return Id(id), nil
}

func (r *Database) GetChannel(channelid Id) (Channel, error) {
	rows, err := r.conn.Query("SELECT channelid, channelname, serverid, timestamp FROM ChannelTable WHERE channelid = ?", channelid)
	if err != nil {
		return Channel{}, err
	}
	defer rows.Close()
	count := 0
	var channel Channel
	for rows.Next() {
		count += 1
		if count > 1 {
			return Channel{}, ErrMultipleRecords
		}
		err := rows.Scan(&channel.ChannelId, &channel.ChannelName, &channel.ServerId, &channel.Timestamp)
		if err != nil {
			return Channel{}, err
		}
	}
	if count == 0 {
		return Channel{}, ErrNoRecord
	}
	return channel, nil
}
func (r *Database) UpdateChannelName(userid Id, username string) error {
	_, err := r.conn.Exec("UPDATE ChannelTable SET channelname = ? WHERE channelid=? ", username, userid)
	return err
}

func (r *Database) GetMessage(messageid Id) (Message, error) {
	rows, err := r.conn.Query("SELECT messageid, channelid, userid, contents, timestamp, editted, edittimestamp FROM ChannelMessageTable WHERE messageid = ? ", messageid)
	if err != nil {
		return Message{}, err
	}
	defer rows.Close()
	var message Message
	for rows.Next() {
		err := rows.Scan(&message.MessageId, &message.ChannelId, &message.UserId, &message.Contents, &message.Timestamp, &message.Editted, &message.EdittedTimeStamp)
		if err != nil {
			return Message{}, err
		}
	}
	return message, nil
}

func (r *Database) AddMessage(channelid Id, userid Id, message string) (Id, error) {
	if userid == 0 || channelid == 0 {
		return 0, fmt.Errorf("add message - zero userid or channel id")
	}
	d, err := r.conn.Exec("INSERT INTO ChannelMessageTable (userid, channelid, contents) VALUES ( ?, ?, ?)", userid, channelid, message)
	if err != nil {
		return 0, fmt.Errorf("add user - userid: %d err: %w", userid, err)
	}
	id, err := d.LastInsertId()
	if err != nil {
		return 0, err
	}
	if id < 0 {
		return 0, ErrNegativeRowIndex
	}
	return Id(id), nil
}

func (r *Database) GetMessagesInChannel(channelid Id, number uint) ([]Message, error) {
	rows, err := r.conn.Query("SELECT messageid, channelid, userid, contents, timestamp, editted, edittimestamp FROM ChannelMessageTable WHERE channelid = ? ORDER BY timestamp DESC LIMIT ?", channelid, number)
	if err != nil {
		return []Message{}, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.MessageId, &message.ChannelId, &message.UserId, &message.Contents, &message.Timestamp, &message.Editted, &message.EdittedTimeStamp)
		if err != nil {
			return []Message{}, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

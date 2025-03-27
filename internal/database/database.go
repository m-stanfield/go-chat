package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type dbConn interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type AtomitcDBService struct {
	service  *DBService
	commit   func() error
	rollback func() error
}

func (a *AtomitcDBService) Service() *DBService {
	return a.service
}

func (a *AtomitcDBService) Commit() error {
	return a.commit()
}

func (a *AtomitcDBService) Rollback() error {
	return a.rollback()
}

type DBService struct {
	db   *sql.DB
	conn dbConn
}

func New(db *sql.DB) *DBService {
	return &DBService{db: db, conn: db}
}

func (r *DBService) ValidateUserLoginInfo(userid Id, password string) (bool, error) {
	user, err := r.GetUserLoginInfo(userid)
	if err != nil {
		return false, err
	}
	return comparePassword(user, password), nil
}

func (db *DBService) withTx(tx *sql.Tx) *DBService {
	return &DBService{db: db.db, conn: tx}
}

func (r *DBService) Atomic(ctx context.Context, opts *sql.TxOptions) (*AtomitcDBService, error) {
	tx, err := r.db.BeginTx(ctx, opts)
	if err != nil {
		return &AtomitcDBService{}, err
	}
	commit := func() error {
		return tx.Commit()
	}

	rollback := func() error {
		return tx.Rollback()
	}
	a := r.withTx(tx)
	return &AtomitcDBService{service: a, commit: commit, rollback: rollback}, nil
}

func hashPassword(password string, salt string) string {
	return (password + salt)
}

func comparePassword(userinfo UserLoginInfo, password string) bool {
	return hashPassword(password, userinfo.Salt) == userinfo.PasswordHash
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *DBService) Close() error {
	return s.db.Close()
}

func (r *DBService) CreateServer(ownerid Id, servername string) (Id, error) {
	d, err := r.conn.Exec(
		"INSERT INTO ServerTable (servername, ownerid) VALUES (?, ?)",
		servername,
		ownerid,
	)
	if err != nil {
		return 0, fmt.Errorf(
			"add server - servername: %s ownerid: %d err: %w",
			servername,
			ownerid,
			err,
		)
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

func (r *DBService) DeleteMessage(messageid Id) error {
	a, err := r.conn.Exec("DELETE FROM ChannelMessageTable WHERE messageid = ?", messageid)
	if err != nil {
		return err
	}
	rowsAffected, err := a.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (r *DBService) UpdateMessage(messageid Id, message string) error {
	_, err := r.conn.Exec(
		"UPDATE ChannelMessageTable SET contents = ? WHERE messageid=? ",
		message,
		messageid,
	)
	return err
}

func (r *DBService) GetUserIDFromUserName(username string) (Id, error) {
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
		return 0, ErrRecordNotFound
	}
	return userid, nil
}

func (r *DBService) UpdateUserSessionToken(userid Id) (string, time.Time, error) {
	token := "token" + strconv.FormatUint(uint64(userid), 10)
	expire := time.Now().Add(24 * time.Hour)
	_, err := r.conn.Exec(
		"UPDATE UserLoginTable SET token = ?, token_expire_time = ? WHERE userid=? ",
		token,
		expire,
		userid,
	)
	if err != nil {
		return "", expire, err
	}
	return token, expire, nil
}

func (r *DBService) GetUserLoginInfo(userid Id) (UserLoginInfo, error) {
	rows, err := r.conn.Query(
		"SELECT userid, passwordhash, salt, token, token_expire_time FROM UserLoginTable WHERE userid = ?",
		userid,
	)
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
		err := rows.Scan(
			&user.UserId,
			&user.PasswordHash,
			&user.Salt,
			&user.Token,
			&user.TokenExpireTime,
		)
		if err != nil {
			return UserLoginInfo{}, err
		}
	}
	if count == 0 {
		return UserLoginInfo{}, ErrRecordNotFound
	}
	return user, nil
}

func (r *DBService) GetUserLoginInfoFromToken(token string) (UserLoginInfo, error) {
	rows, err := r.conn.Query(
		"SELECT userid, passwordhash, salt, token, token_expire_time FROM UserLoginTable WHERE token = ?",
		token,
	)
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
		err := rows.Scan(
			&user.UserId,
			&user.PasswordHash,
			&user.Salt,
			&user.Token,
			&user.TokenExpireTime,
		)
		if err != nil {
			return UserLoginInfo{}, err
		}
	}
	if count == 0 {
		return UserLoginInfo{}, ErrRecordNotFound
	}
	return user, nil
}

func (r *DBService) GetUser(userid Id) (User, error) {
	rows, err := r.conn.Query("SELECT userid, username FROM UserTable WHERE userid = ?", userid)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrRecordNotFound
	}
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
		return User{}, ErrRecordNotFound
	}
	return user, nil
}

func (r *DBService) CreateUser(username string, password string) (Id, error) {
	d, err := r.conn.Exec("INSERT INTO UserTable (username) VALUES (?)", username)
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint &&
		sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return 0, ErrRecordAlreadyExists
	}
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
	random_salt := "salt" + strconv.FormatUint(uint64(id), 10)
	hashed_password := hashPassword(password, random_salt)
	_, err = r.conn.Exec(
		"INSERT INTO UserLoginTable (userid, passwordhash, salt, token) VALUES ( ?, ?, ?, ?)",
		id,
		hashed_password,
		random_salt,
		"",
	)
	if err != nil {
		return 0, err
	}

	return Id(id), nil
}

func (r *DBService) UpdateUserName(userid Id, username string) error {
	_, err := r.conn.Exec("UPDATE UserTable SET username = ? WHERE userid=? ", username, userid)
	return err
}

func (r *DBService) UpdateServerName(serverid Id, servername string) error {
	_, err := r.conn.Exec(
		"UPDATE ServerTable SET servername = ? WHERE serverid=? ",
		servername,
		serverid,
	)
	return err
}

func (r *DBService) GetRecentUsernames(userid Id, number uint) ([]UsernameLogEntry, error) {
	rows, err := r.conn.Query(
		"SELECT userid, username, timestamp FROM UserNameLogTable WHERE userid = ? ORDER BY timestamp DESC LIMIT ?",
		userid,
		number,
	)
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

func (r *DBService) GetUsersOfServer(serverid Id) ([]User, error) {
	rows, err := r.conn.Query(
		"SELECT U.userid, U.username FROM UsersServerTable as US INNER JOIN UserTable as U ON US.userid = U.userid WHERE US.serverid = ?",
		serverid,
	)
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

func (r *DBService) DeleteServer(serverid Id) error {
	_, err := r.conn.Exec("DELETE FROM ServerTable WHERE serverid = ?", serverid)
	return err
}

func (r *DBService) GetServersOfUser(userid Id) ([]Server, error) {
	rows, err := r.conn.Query(
		"SELECT S.serverid, S.ownerid, S.servername FROM UsersServerTable as U INNER JOIN ServerTable as S ON U.serverid = S.serverid WHERE U.userid = ?",
		userid,
	)
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

func (r *DBService) GetChannelsOfServer(serverid Id) ([]Channel, error) {
	rows, err := r.conn.Query(
		"SELECT channelid, serverid, channelname, timestamp FROM ChannelTable WHERE serverid = ?",
		serverid,
	)
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

func (r *DBService) IsUserInChannel(userid Id, channelid Id) (bool, error) {
	query := `SELECT COUNT(1) FROM UsersChannelTable WHERE channelid = ? AND userid = ?`
	var count int
	err := r.db.QueryRow(query, channelid, userid).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *DBService) AddUserToChannel(userid Id, channelid Id) error {
	d, err := r.conn.Exec(
		"INSERT INTO UsersChannelTable ( userid, channelid) VALUES ( ?, ?)",
		userid,
		channelid,
	)
	if err != nil {
		return fmt.Errorf("add user - userid: %d err: %w", userid, err)
	}
	id, err := d.LastInsertId()
	if err != nil {
		return err
	}
	if id < 0 {
		return ErrNegativeRowIndex
	}
	return nil
}

func (r *DBService) AddChannel(serverid Id, channelname string) (Id, error) {
	d, err := r.conn.Exec(
		"INSERT INTO ChannelTable ( serverid, channelname) VALUES ( ?, ?)",
		serverid,
		channelname,
	)
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

func (r *DBService) IsUserInServer(userid Id, serverid Id) (bool, error) {
	query := `SELECT COUNT(1) FROM UsersServerTable WHERE serverid = ? AND userid = ?`
	var count int
	err := r.db.QueryRow(query, serverid, userid).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *DBService) GetChannel(channelid Id) (Channel, error) {
	rows, err := r.conn.Query(
		"SELECT channelid, channelname, serverid, timestamp FROM ChannelTable WHERE channelid = ?",
		channelid,
	)
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
		err := rows.Scan(
			&channel.ChannelId,
			&channel.ChannelName,
			&channel.ServerId,
			&channel.Timestamp,
		)
		if err != nil {
			return Channel{}, err
		}
	}
	if count == 0 {
		return Channel{}, ErrRecordNotFound
	}
	return channel, nil
}

func (r *DBService) UpdateChannel(channelid Id, new_server_name string) error {
	_, err := r.conn.Exec(
		"UPDATE ChannelTable SET channelname = ? WHERE channelid = ? ",
		new_server_name,
		channelid,
	)
	return err
}

func (r *DBService) GetServer(serverid Id) (Server, error) {
	rows, err := r.conn.Query(
		"SELECT serverid, ownerid, servername FROM ServerTable WHERE serverid = ? ",
		serverid,
	)
	if err != nil {
		return Server{}, err
	}
	defer rows.Close()
	var server Server
	server_found := false
	for rows.Next() {
		server_found = true
		err := rows.Scan(
			&server.ServerId,
			&server.OwnerId,
			&server.ServerName,
		)
		if err != nil {
			return Server{}, err
		}
	}
	if !server_found {
		return Server{}, ErrRecordNotFound
	}
	return server, nil
}

func (r *DBService) GetMessage(messageid Id) (Message, error) {
	rows, err := r.conn.Query(
		"SELECT messageid, channelid, userid, contents, timestamp, editted, edittimestamp FROM ChannelMessageTable WHERE messageid = ? ",
		messageid,
	)
	if err != nil {
		return Message{}, err
	}
	defer rows.Close()
	count := 0
	var message Message
	for rows.Next() {
		count += 1
		err := rows.Scan(
			&message.MessageId,
			&message.ChannelId,
			&message.UserId,
			&message.Contents,
			&message.Timestamp,
			&message.Editted,
			&message.EdittedTimeStamp,
		)
		if err != nil {
			return Message{}, err
		}
	}
	if count == 0 {
		return Message{}, ErrRecordNotFound
	}
	return message, nil
}

func (r *DBService) AddMessage(channelid Id, userid Id, message string) (Id, error) {
	if userid == 0 || channelid == 0 {
		return 0, fmt.Errorf("add message - zero userid or channel id")
	}
	d, err := r.conn.Exec(
		"INSERT INTO ChannelMessageTable (userid, channelid, contents) VALUES ( ?, ?, ?)",
		userid,
		channelid,
		message,
	)
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

func (r *DBService) GetMessagesInChannel(channelid Id, number uint) ([]Message, error) {
	rows, err := r.conn.Query(
		"SELECT messageid, channelid, userid, contents, timestamp, editted, edittimestamp FROM ChannelMessageTable WHERE channelid = ? ORDER BY timestamp DESC LIMIT ?",
		channelid,
		number,
	)
	if err != nil {
		return []Message{}, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(
			&message.MessageId,
			&message.ChannelId,
			&message.UserId,
			&message.Contents,
			&message.Timestamp,
			&message.Editted,
			&message.EdittedTimeStamp,
		)
		if err != nil {
			return []Message{}, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (r *DBService) GetUsersInChannel(channelid Id) ([]User, error) {
	rows, err := r.conn.Query(
		"SELECT U.userid, U.username FROM UsersChannelTable as UC INNER JOIN UserTable as U ON UC.userid = U.userid WHERE UC.channelid = ?",
		channelid,
	)
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

func (r *DBService) RemoveUserFromChannel(channelid Id, userid Id) error {
	result, err := r.conn.Exec(
		"DELETE FROM UsersChannelTable WHERE channelid = ? AND userid = ?",
		channelid,
		userid,
	)
	// check and ensure that at least one row was deleted
	if err != nil {
		return fmt.Errorf(
			"remove user from channel - channelid: %d userid: %d err: %w",
			channelid,
			userid,
			err,
		)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (r *DBService) DeleteChannel(channelid Id) error {
	_, err := r.conn.Exec("DELETE FROM ChannelTable WHERE channelid = ?", channelid)
	return err
}

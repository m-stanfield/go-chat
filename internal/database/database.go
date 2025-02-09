package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string
	Atomic(context.Context, AtomicCallback) error
	GetUserIDFromUserName(username string) (Id, error)
	UpdateUserSessionToken(userid Id) (string, time.Time, error)
	GetUserLoginInfoFromToken(token string) (UserLoginInfo, error)
	GetUserLoginInfo(userid Id) (UserLoginInfo, error)
	ValidateUserLoginInfo(userid Id, password string) (bool, error)
	GetUser(userid Id) (User, error)
	CreateUser(username string, password string) (Id, error)
	UpdateUserName(userid Id, username string) error
	GetRecentUsernames(userid Id, number uint) ([]UsernameLogEntry, error)
	GetUsersOfServer(serverid Id) ([]User, error)
	GetServersOfUser(userid Id) ([]Server, error)
	GetChannelsOfServer(serverid Id) ([]Channel, error)
	AddChannel(serverid Id, channelname string) (Id, error)
	GetChannel(channelid Id) (Channel, error)
	UpdateChannelName(userid Id, username string) error
	GetMessage(messageid Id) (Message, error)
	AddMessage(channelid Id, userid Id, message string) (Id, error)
	GetMessagesInChannel(channelid Id, number uint) ([]Message, error)
	GetServer(serverid Id) (Server, error)
	CreateServer(ownerid Id, servername string) (Id, error)
	IsUserInServer(userid Id, serverid Id) (bool, error)

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}
type DBConn interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type AtomicCallback = func(r Service) error

type service struct {
	db   *sql.DB
	conn DBConn
}

var (
	dburl      = os.Getenv("BLUEPRINT_DB_URL")
	dbInstance *service
)

func (r *service) ValidateUserLoginInfo(userid Id, password string) (bool, error) {
	user, err := r.GetUserLoginInfo(userid)
	if err != nil {
		return false, err
	}
	return comparePassword(user, password), nil
}

func (db *service) withTx(tx *sql.Tx) *service {
	return &service{db: db.db, conn: tx}
}

func (r *service) Atomic(ctx context.Context, cb func(ds Service) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("tx err: %w, rb err: %w", err, rbErr)
			}
		} else {
			err = tx.Commit()
		}
	}()
	dbTx := r.withTx(tx)
	err = cb(dbTx)
	return err
}

func New() Service {
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

	dbInstance = &service{
		db:   db,
		conn: db,
	}
	return dbInstance
}

func hashPassword(password string, salt string) string {
	return (password + salt)
}

func comparePassword(userinfo UserLoginInfo, password string) bool {
	return hashPassword(password, userinfo.Salt) == userinfo.PasswordHash
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dburl)
	return s.db.Close()
}

func (r *service) CreateServer(ownerid Id, servername string) (Id, error) {
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

func (r *service) GetUserIDFromUserName(username string) (Id, error) {
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

func (r *service) UpdateUserSessionToken(userid Id) (string, time.Time, error) {
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

func (r *service) GetUserLoginInfo(userid Id) (UserLoginInfo, error) {
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
		return UserLoginInfo{}, ErrNoRecord
	}
	return user, nil
}

func (r *service) GetUserLoginInfoFromToken(token string) (UserLoginInfo, error) {
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
		return UserLoginInfo{}, ErrNoRecord
	}
	return user, nil
}

func (r *service) GetUser(userid Id) (User, error) {
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

func (r *service) CreateUser(username string, password string) (Id, error) {
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

func (r *service) UpdateUserName(userid Id, username string) error {
	_, err := r.conn.Exec("UPDATE UserTable SET username = ? WHERE userid=? ", username, userid)
	return err
}

func (r *service) GetRecentUsernames(userid Id, number uint) ([]UsernameLogEntry, error) {
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

func (r *service) GetUsersOfServer(serverid Id) ([]User, error) {
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

func (r *service) GetServersOfUser(userid Id) ([]Server, error) {
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

func (r *service) GetChannelsOfServer(serverid Id) ([]Channel, error) {
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

func (r *service) AddChannel(serverid Id, channelname string) (Id, error) {
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

func (r *service) IsUserInServer(userid Id, serverid Id) (bool, error) {
	query := `SELECT COUNT(1) FROM UsersServerTable WHERE serverid = ? AND userid = ?`
	var count int
	err := r.db.QueryRow(query, serverid, userid).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *service) GetChannel(channelid Id) (Channel, error) {
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
		return Channel{}, ErrNoRecord
	}
	return channel, nil
}

func (r *service) UpdateChannelName(userid Id, username string) error {
	_, err := r.conn.Exec(
		"UPDATE ChannelTable SET channelname = ? WHERE channelid=? ",
		username,
		userid,
	)
	return err
}

func (r *service) GetServer(serverid Id) (Server, error) {
	rows, err := r.conn.Query(
		"SELECT serverid, ownerid, servername FROM ServerTable WHERE serverid = ? ",
		serverid,
	)
	if err != nil {
		return Server{}, err
	}
	defer rows.Close()
	var server Server
	for rows.Next() {
		err := rows.Scan(
			&server.ServerId,
			&server.OwnerId,
			&server.ServerName,
		)
		if err != nil {
			return Server{}, err
		}
	}
	return server, nil
}

func (r *service) GetMessage(messageid Id) (Message, error) {
	rows, err := r.conn.Query(
		"SELECT messageid, channelid, userid, contents, timestamp, editted, edittimestamp FROM ChannelMessageTable WHERE messageid = ? ",
		messageid,
	)
	if err != nil {
		return Message{}, err
	}
	defer rows.Close()
	var message Message
	for rows.Next() {
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
	return message, nil
}

func (r *service) AddMessage(channelid Id, userid Id, message string) (Id, error) {
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

func (r *service) GetMessagesInChannel(channelid Id, number uint) ([]Message, error) {
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

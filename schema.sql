BEGIN TRANSACTION;
DROP TABLE IF EXISTS "ChannelTable";
CREATE TABLE IF NOT EXISTS "ChannelTable" (
	"channelid"	INTEGER NOT NULL UNIQUE,
	"serverid"	INTEGER NOT NULL,
	"channelname"	TEXT NOT NULL,
	"timestamp"	DATETIME DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
	PRIMARY KEY("channelid" AUTOINCREMENT),
	UNIQUE("serverid","channelname"),
	FOREIGN KEY("serverid") REFERENCES "ServerTable"("serverid")
);
DROP TABLE IF EXISTS "ServerTable";
CREATE TABLE IF NOT EXISTS "ServerTable" (
	"serverid"	INTEGER NOT NULL,
	"ownerid"	INTEGER NOT NULL,
	"servername"	TEXT NOT NULL,
	FOREIGN KEY("ownerid") REFERENCES "UserTable"("userid"),
	PRIMARY KEY("serverid" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "UserNameLogTable";
CREATE TABLE IF NOT EXISTS "UserNameLogTable" (
	"id"	INTEGER NOT NULL UNIQUE,
	"userid"	INTEGER NOT NULL,
	"username"	TEXT NOT NULL,
	"timestamp"	DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("userid") REFERENCES "UserTable"("userid")
);
DROP TABLE IF EXISTS "UserNicknameLogTable";
CREATE TABLE IF NOT EXISTS "UserNicknameLogTable" (
	"id"	INTEGER NOT NULL UNIQUE,
	"userid"	INTEGER NOT NULL,
	"serverid"	INTEGER NOT NULL,
	"nickname"	TEXT NOT NULL,
	"timestamp"	DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
	FOREIGN KEY("userid") REFERENCES "UserTable"("userid"),
	FOREIGN KEY("serverid") REFERENCES "ServerTable"("serverid"),
	PRIMARY KEY("id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "UserTable";
CREATE TABLE IF NOT EXISTS "UserTable" (
	"userid"	INTEGER NOT NULL,
	"username"	TEXT NOT NULL,
	UNIQUE("username"),
	PRIMARY KEY("userid" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "UsersServerTable";
CREATE TABLE IF NOT EXISTS "UsersServerTable" (
	"userid"	INTEGER NOT NULL,
	"serverid"	INTEGER NOT NULL,
	"nickname"	TEXT NOT NULL,
	"timestamp"	TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
	FOREIGN KEY("serverid") REFERENCES "ServerTable"("serverid"),
	PRIMARY KEY("userid","serverid"),
	FOREIGN KEY("userid") REFERENCES "UserTable"("userid")
);
DROP TABLE IF EXISTS "UserLoginTable";
CREATE TABLE IF NOT EXISTS "UserLoginTable" (
	"userid"	INTEGER NOT NULL UNIQUE,
	"passwordhash"	TEXT NOT NULL,
	"salt"	TEXT NOT NULL,
	"token"	TEXT NOT NULL,
	"token_expire_time"	DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
	FOREIGN KEY("userid") REFERENCES "UserTable"("userid")
);
DROP TABLE IF EXISTS "ChannelMessageTable";
CREATE TABLE IF NOT EXISTS "ChannelMessageTable" (
	"messageid"	INTEGER NOT NULL UNIQUE,
	"channelid"	INTEGER NOT NULL,
	"userid"	INTEGER NOT NULL,
	"contents"	TEXT NOT NULL,
	"timestamp"	DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
	"editted"	INTEGER,
	"edittimestamp"	DATETIME,
	FOREIGN KEY("userid") REFERENCES "UsersChannelTable"("userid"),
	FOREIGN KEY("channelid") REFERENCES "UsersChannelTable"("channelid"),
	PRIMARY KEY("messageid" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "UsersChannelTable";
CREATE TABLE IF NOT EXISTS "UsersChannelTable" (
	"userid"	INTEGER NOT NULL,
	"channelid"	INTEGER NOT NULL,
	FOREIGN KEY("userid") REFERENCES "UserTable"("userid"),
	FOREIGN KEY("channelid") REFERENCES "ChannelTable"("channelid"),
	PRIMARY KEY("userid","channelid")
);
DROP TRIGGER IF EXISTS "UpdateUserNameLog";
CREATE TRIGGER UpdateUserNameLog AFTER UPDATE OF username ON UserTable 
BEGIN
	INSERT INTO UserNameLogTable (userid, username) VALUES (new.userid, new.username);
END;
DROP TRIGGER IF EXISTS "UpdateUserNicknameLog";
CREATE TRIGGER UpdateUserNicknameLog AFTER UPDATE OF nickname ON UsersServerTable 
BEGIN
	INSERT INTO UserNicknameLogTable (userid, serverid, nickname) VALUES (new.userid, new.serverid, new.nickname);
END;

CREATE TRIGGER UpdateMessageLog AFTER UPDATE OF contents ON ChannelMessageTable 
BEGIN
	UPDATE ChannelMessageTable SET editted = 1, edittimestamp = strftime('%Y-%m-%d %H:%M:%f', 'now') WHERE messageid = old.messageid;
END;
COMMIT;

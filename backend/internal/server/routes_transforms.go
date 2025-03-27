package server

import (
	"time"

	"go-chat-react/internal/database"
)

func fromDBMessageToSeverMessage(message database.Message) ServerMessage {
	return ServerMessage{
		UserId:    message.UserId,
		ChannelId: message.ChannelId,
		MessageID: message.MessageId,
		Message:   message.Contents,
		Date:      message.Timestamp.Format(time.UnixDate),
	}
}

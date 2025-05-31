package websocket

import (
	"context"

	"github.com/coder/websocket"
)

type CoderWebSocketConnection struct {
	conn *websocket.Conn
}

func NewCoderWebSocketConnection(conn *websocket.Conn) *CoderWebSocketConnection {
	return &CoderWebSocketConnection{conn: conn}
}

func (w *CoderWebSocketConnection) Close(code StatusCode, message string) error {
	return w.conn.Close(websocket.StatusCode(code), message)
}

func (w *CoderWebSocketConnection) Read(ctx context.Context) (MessageType, []byte, error) {
	msgType, data, err := w.conn.Read(ctx)
	return MessageType(msgType), data, err
}

func (w *CoderWebSocketConnection) Write(
	ctx context.Context,
	msgType MessageType,
	data []byte,
) error {
	return w.conn.Write(ctx, websocket.MessageType(msgType), data)
}

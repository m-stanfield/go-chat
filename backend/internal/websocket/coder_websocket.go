package websocket

import (
	"context"
	"errors"
	"net/http"

	"github.com/coder/websocket"
)

type CoderWebSocketConnection struct {
	conn *websocket.Conn
}

func NewCoderWebSocketConnection(
	w http.ResponseWriter,
	r *http.Request,
) (*CoderWebSocketConnection, error) {
	opts := websocket.AcceptOptions{InsecureSkipVerify: true}
	conn, err := websocket.Accept(w, r, &opts)
	if err != nil {
		return nil, errors.New("failed to open websocket connection: " + err.Error())
	}
	return &CoderWebSocketConnection{conn: conn}, nil
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

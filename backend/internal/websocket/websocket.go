package websocket

import (
	"context"
	"errors"
	"log"
	"sync"
)

type IncomingMessage struct {
	Id      uint64
	Payload []byte
}

type (
	WebSocketConnection interface {
		Close(StatusCode, string) error
		Read(context.Context) (MessageType, []byte, error)
		Write(context.Context, MessageType, []byte) error
	}
	webSocketClient struct {
		ID      uint64
		conn    WebSocketConnection
		receive chan IncomingMessage
		send    chan []byte
		cancel  context.CancelFunc
		closed  bool
	}
)

func newWebSocketClient(
	Id uint64,
	conn WebSocketConnection,
	incoming chan IncomingMessage,
) *webSocketClient {
	ctx, cancel := context.WithCancel(context.Background())

	send := make(chan []byte)
	client := webSocketClient{
		ID:      Id,
		conn:    conn,
		cancel:  cancel,
		send:    send,
		receive: incoming,
		closed:  false,
	}
	go client.read(ctx)
	go client.write(ctx)
	return &client
}

func (c *webSocketClient) close(status StatusCode) error {
	if c.closed {
		return errors.New("client already closed")
	}
	err := c.conn.Close(status, "") // TODO: determine proper status
	if err != nil {
		return err
	}
	c.cancel()
	close(c.receive)
	return nil
}

func (c *webSocketClient) read(ctx context.Context) {
	for {
		messageType, message, err := c.conn.Read(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Printf("Client %d read cancelled", c.ID)
				return
			} else {
				log.Printf("Client %d read error: %v", c.ID, err)
				c.close(StatusAbnormalClosure)
			}
			return
		}

		if messageType == MessageBinary || messageType == MessageText {
			log.Printf("Received from client %d (%d bytes): %s", c.ID, len(message), message)
			msg := IncomingMessage{
				Id:      c.ID,
				Payload: message,
			}
			c.receive <- msg
		}
	}
}

func (c *webSocketClient) write(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.send:
			if !ok {
				c.close(StatusNormalClosure)
				return
			}
			err := c.conn.Write(ctx, MessageText, msg)
			if err != nil {
				log.Printf("Client %d write error: %v", c.ID, err)
				c.close(StatusAbnormalClosure)
				return
			}
		}
	}
}

type WebSocketManager struct {
	clients map[uint64]*webSocketClient
	mutex   sync.RWMutex
}

func (m *WebSocketManager) NewConnection(
	Id uint64,
	conn WebSocketConnection,
) chan IncomingMessage {
	incoming := make(chan IncomingMessage, 100) // Buffered channel for incoming messages
	client := newWebSocketClient(Id, conn, incoming)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.clients[client.ID] = client
	log.Printf("Client %d registered. Total clients: %d", client.ID, len(m.clients))
	return incoming
}

func (m *WebSocketManager) CloseConnection(id uint64) {
	// remove later, redundant with Deregistre
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, ok := m.clients[id]
	if !ok {
		return
	}
	delete(m.clients, id)
	close(client.send)
	log.Printf("Client %d unregistered. Total clients: %d", client.ID, len(m.clients))
}

func (m *WebSocketManager) SendToClient(Id uint64, message []byte) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	client, ok := m.clients[Id]
	if !ok {
		log.Printf("Client %d not found.", Id)
		return false
	}
	select {
	case client.send <- message:
		return true
	default:
		log.Printf("Client %d send channel full, dropping message.", Id)
		return false
	}
}

package websocket

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/google/uuid"
)

type IncomingMessage struct {
	Payload []byte
}

type (
	WebSocketConnection interface {
		Close(StatusCode, string) error
		Read(context.Context) (MessageType, []byte, error)
		Write(context.Context, MessageType, []byte) error
	}
	webSocketClient struct {
		ID      string
		conn    WebSocketConnection
		receive chan IncomingMessage
		send    chan []byte
		cancel  context.CancelFunc
		closed  bool
	}
)

func newWebSocketClient(
	Id string,
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
	// TODO: find correct ordering of close operations. how to handle closing channel vs connetion
	if c.closed {
		return errors.New("client already closed")
	}

	c.cancel()
	close(c.receive)
	log.Printf("Client %s closed with status %d", c.ID, status)
	err := c.conn.Close(status, "") // TODO: determine proper status
	if err != nil {
		return err
	}
	return nil
}

func (c *webSocketClient) read(ctx context.Context) {
	for {
		messageType, message, err := c.conn.Read(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Printf("Client %s read cancelled", c.ID)
				c.close(StatusNormalClosure)
			} else {
				log.Printf("Client %s read error: %v", c.ID, err)
				c.close(StatusAbnormalClosure)
			}
			return
		}

		if messageType == MessageBinary || messageType == MessageText {
			log.Printf("Received from client %s (%d bytes): %s", c.ID, len(message), message)
			msg := IncomingMessage{
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
				log.Printf("Client %s write error: %v", c.ID, err)
				c.close(StatusAbnormalClosure)
				return
			}
		}
	}
}

type WebSocketManager struct {
	clients map[string]*webSocketClient
	mutex   sync.RWMutex
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients: make(map[string]*webSocketClient),
		mutex:   sync.RWMutex{},
	}
}

func (m *WebSocketManager) NewConnection(
	conn WebSocketConnection,
) (string, chan IncomingMessage) {
	Id := uuid.New().String()
	incoming := make(chan IncomingMessage, 100) // Buffered channel for incoming messages
	client := newWebSocketClient(Id, conn, incoming)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.clients[client.ID] = client
	log.Printf("Client %s registered. Total clients: %d", client.ID, len(m.clients))
	return Id, incoming
}

func (m *WebSocketManager) CloseConnection(id string) {
	// remove later, redundant with Deregistre
	m.mutex.Lock()
	defer m.mutex.Unlock()
	client, ok := m.clients[id]
	if !ok {
		return
	}
	delete(m.clients, id)
	close(client.send)
	log.Printf("Client %s unregistered. Total clients: %d", client.ID, len(m.clients))
}

func (m *WebSocketManager) SendToClient(Id string, message []byte) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	client, ok := m.clients[Id]
	if !ok {
		log.Printf("Client %s not found.", Id)
		return false
	}
	select {
	case client.send <- message:
		return true
	default:
		log.Printf("Client %s send channel full, dropping message.", Id)
		return false
	}
}

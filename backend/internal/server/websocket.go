package server

import (
	"context"
	"log"
	"sync"

	"github.com/coder/websocket"
)

type IncomingMessage struct {
	Id      uint64
	Payload []byte
}

type webSocketClient struct {
	ID      uint64
	conn    *websocket.Conn
	receive chan IncomingMessage
	send    chan []byte
	cancel  context.CancelFunc
	onClose func(*webSocketClient)
}

func newWebSocketClient(Id uint64, conn *websocket.Conn, onClose func(*webSocketClient),
) *webSocketClient {
	ctx, cancel := context.WithCancel(context.Background())

	send := make(chan []byte)
	receive := make(chan IncomingMessage)
	client := webSocketClient{
		ID:      Id,
		conn:    conn,
		cancel:  cancel,
		send:    send,
		receive: receive,
		onClose: onClose,
	}
	go client.read(ctx)
	go client.write(ctx)
	return &client
}

func (c *webSocketClient) close() {
	c.conn.Close(websocket.StatusNormalClosure, "") // TODO: determine proper status
	c.cancel()
	close(c.receive)
	c.onClose(c)
}

func (c *webSocketClient) read(ctx context.Context) {
	defer c.close()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			messageType, message, err := c.conn.Read(ctx)
			if err != nil {
				log.Printf("Client %d read error: %v", c.ID, err)
				return
			}

			if messageType == websocket.MessageBinary || messageType == websocket.MessageText {
				log.Printf("Received from client %d (%d bytes): %s", c.ID, len(message), message)
				msg := IncomingMessage{
					Id:      c.ID,
					Payload: message,
				}
				c.receive <- msg
			}
		}
	}
}

func (c *webSocketClient) write(ctx context.Context) {
	defer c.close()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			err := c.conn.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				log.Printf("Client %d write error: %v", c.ID, err)
				return
			}
		}
	}
}

type WebSocketManager struct {
	clients    map[uint64]*webSocketClient
	register   chan *webSocketClient
	deregister chan uint64
	mutex      sync.RWMutex
}

func (m *WebSocketManager) NewConnection(Id uint64, conn *websocket.Conn) {
	onClose := func(wsc *webSocketClient) {
		m.deregister <- wsc.ID
	}
	client := newWebSocketClient(Id, conn, onClose)
	m.register <- client
}

func (m *WebSocketManager) CloseConnection(id uint64) bool {
	m.mutex.Lock()
	client, ok := m.clients[id]
	if ok {
		delete(m.clients, id)
	}
	m.mutex.Unlock()
	if ok {
		// splitting into second to allow lock to be held as short as possible
		client.close()
	}
	return ok
}

func (m *WebSocketManager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mutex.Lock()
			m.clients[client.ID] = client
			m.mutex.Unlock()
			log.Printf("Client %d registered. Total clients: %d", client.ID, len(m.clients))
		case removeId := <-m.deregister:
			m.mutex.Lock()
			client, ok := m.clients[removeId]
			if !ok {
				continue
			}
			delete(m.clients, removeId)
			m.mutex.Unlock()

			close(client.send)
			client.close()
			log.Printf("Client %d unregistered. Total clients: %d", client.ID, len(m.clients))
		}
	}
}

func (m *WebSocketManager) SendToClient(Id uint64, message []byte) bool {
	m.mutex.RLock()
	client, ok := m.clients[Id]
	m.mutex.RUnlock()
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

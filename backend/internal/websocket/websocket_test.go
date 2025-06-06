package websocket

import (
	"context"
	"errors"
	"testing"
	"time"
)

// mockWebSocketConnection mocks WebSocketConnection for testing.
type mockWebSocketConnection struct {
	readChan  chan []byte
	writeChan chan []byte
	closeErr  error
	readErr   error
	writeErr  error
	closed    bool
}

func (m *mockWebSocketConnection) Close(code StatusCode, reason string) error {
	m.closed = true
	return m.closeErr
}

func (m *mockWebSocketConnection) Read(ctx context.Context) (MessageType, []byte, error) {
	select {
	case msg := <-m.readChan:
		return MessageText, msg, m.readErr
	case <-ctx.Done():
		return 0, nil, ctx.Err()
	}
}

func (m *mockWebSocketConnection) Write(
	ctx context.Context,
	msgType MessageType,
	msg []byte,
) error {
	m.writeChan <- msg
	return m.writeErr
}

// Test registration and deregistration in WebSocketManager
func TestWebSocketManager_RegisterAndDeregister(t *testing.T) {
	manager := &WebSocketManager{
		clients: make(map[string]*webSocketClient),
	}

	mockConn := &mockWebSocketConnection{
		readChan:  make(chan []byte),
		writeChan: make(chan []byte),
	}

	clientID := "test-client-1"
	manager.NewConnection(clientID, mockConn)

	// Wait for registration
	time.Sleep(50 * time.Millisecond)

	manager.mutex.RLock()
	if _, ok := manager.clients[clientID]; !ok {
		t.Fatal("Client was not registered")
	}
	manager.mutex.RUnlock()

	// Close connection and deregister
	manager.CloseConnection(clientID)

	time.Sleep(50 * time.Millisecond)

	manager.mutex.RLock()
	if _, ok := manager.clients[clientID]; ok {
		t.Fatal("Client was not deregistered")
	}
	manager.mutex.RUnlock()
}

// Test SendToClient logic
func TestWebSocketManager_SendToClient(t *testing.T) {
	manager := &WebSocketManager{
		clients: make(map[string]*webSocketClient),
	}

	mockConn := &mockWebSocketConnection{
		readChan:  make(chan []byte),
		writeChan: make(chan []byte, 1), // buffered to avoid deadlock
	}

	clientID := "test-client-2"
	manager.NewConnection(clientID, mockConn)

	time.Sleep(50 * time.Millisecond)

	// Send message
	msg := []byte("hello")
	ok := manager.SendToClient(clientID, msg)
	if !ok {
		t.Fatal("Failed to send message to client")
	}

	// Verify message sent
	select {
	case sentMsg := <-mockConn.writeChan:
		if string(sentMsg) != "hello" {
			t.Fatalf("Unexpected message sent: %s", sentMsg)
		}
	case <-time.After(time.Second):
		t.Fatal("Message was not written by client")
	}
}

// Test read and write loops of webSocketClient
func TestWebSocketClient_ReadWrite(t *testing.T) {
	mockConn := &mockWebSocketConnection{
		readChan:  make(chan []byte, 1),
		writeChan: make(chan []byte, 1),
	}

	incoming := make(chan IncomingMessage, 1)
	client := newWebSocketClient("42", mockConn, incoming)

	// Simulate sending a message
	mockConn.readChan <- []byte("incoming")
	select {
	case msg := <-client.receive:
		if string(msg.Payload) != "incoming" {
			t.Fatalf("Unexpected payload: %s", msg.Payload)
		}
	case <-time.After(time.Second):
		t.Fatal("Read did not send to receive channel")
	}

	// Simulate writing a message
	outMsg := []byte("outgoing")
	client.send <- outMsg

	select {
	case written := <-mockConn.writeChan:
		if string(written) != "outgoing" {
			t.Fatalf("Unexpected written message: %s", written)
		}
	case <-time.After(time.Second):
		t.Fatal("Write did not occur")
	}

	// Simulate error on read to trigger close
	mockConn.readErr = errors.New("read error")
	mockConn.readChan <- []byte("trigger close")

	time.Sleep(50 * time.Millisecond)
	if !mockConn.closed {
		t.Fatal("Connection was not closed")
	}
}

// Test SendToClient with missing client
func TestWebSocketManager_SendToMissingClient(t *testing.T) {
	manager := &WebSocketManager{
		clients: make(map[string]*webSocketClient),
	}
	ok := manager.SendToClient("999", []byte("test"))
	if ok {
		t.Fatal("Expected SendToClient to return false for missing client")
	}
}

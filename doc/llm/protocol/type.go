package handlers

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	Data string
}

type Session struct {
	ID        string
	MessageCh chan Message
	CreatedAt time.Time
}

type Manager struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
	}
}

func (m *Manager) CreateSession() *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	sessionID := fmt.Sprintf("%d", time.Now().UnixNano())
	session := &Session{
		ID:        sessionID,
		MessageCh: make(chan Message, 100),
		CreatedAt: time.Now(),
	}

	m.sessions[sessionID] = session
	return session
}

func (m *Manager) GetSession(sessionID string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[sessionID]
	return session, exists
}

func (m *Manager) RemoveSession(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.sessions[sessionID]; exists {
		close(session.MessageCh)
		delete(m.sessions, sessionID)
	}
	fmt.Println("delete session is ok")
}

func (m *Manager) SendMessage(sessionID string, message Message) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if session, exists := m.sessions[sessionID]; exists {
		select {
		case session.MessageCh <- message:
			return true
		default:
			return false
		}
	}
	return false
}

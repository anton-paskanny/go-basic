package storage

import (
	"errors"
	"sync"
	"time"

	"order-api-auth/models"
)

// InMemoryStorage represents an in-memory storage
type InMemoryStorage struct {
	users    map[string]*models.User
	sessions map[string]*models.Session
	mu       sync.RWMutex
}

// NewInMemoryStorage creates a new in-memory storage
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		users:    make(map[string]*models.User),
		sessions: make(map[string]*models.Session),
	}
}

// CreateUser creates a new user
func (s *InMemoryStorage) CreateUser(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if user with this phone already exists
	for _, existingUser := range s.users {
		if existingUser.Phone == user.Phone {
			return errors.New("user with this phone already exists")
		}
	}

	s.users[user.ID] = user
	return nil
}

// GetUserByPhone gets user by phone number
func (s *InMemoryStorage) GetUserByPhone(phone string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Phone == phone {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

// GetUserByID gets user by ID
func (s *InMemoryStorage) GetUserByID(id string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// CreateSession creates a new session
func (s *InMemoryStorage) CreateSession(session *models.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[session.ID] = session
	return nil
}

// GetSession gets session by ID
func (s *InMemoryStorage) GetSession(sessionID string) (*models.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, errors.New("session not found")
	}

	return session, nil
}

// MarkSessionAsUsed marks session as used
func (s *InMemoryStorage) MarkSessionAsUsed(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return errors.New("session not found")
	}

	session.IsUsed = true
	return nil
}

// CleanupExpiredSessions removes expired sessions
func (s *InMemoryStorage) CleanupExpiredSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for sessionID, session := range s.sessions {
		if session.ExpiresAt.Before(now) {
			delete(s.sessions, sessionID)
		}
	}
}

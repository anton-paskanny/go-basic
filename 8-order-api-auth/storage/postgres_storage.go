package storage

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"order-api-auth/models"
)

// UserStorage interface for working with users
type UserStorage interface {
	CreateUser(user *models.User) error
	GetUserByPhone(phone string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
}

// SessionStorage interface for working with sessions
type SessionStorage interface {
	CreateSession(session *models.Session) error
	GetSession(sessionID string) (*models.Session, error)
	MarkSessionAsUsed(sessionID string) error
	CleanupExpiredSessions()
}

// PostgreSQLStorage represents a PostgreSQL storage implementation
type PostgreSQLStorage struct {
	db *gorm.DB
}

// NewPostgreSQLStorage creates a new PostgreSQL storage
func NewPostgreSQLStorage(db *gorm.DB) *PostgreSQLStorage {
	return &PostgreSQLStorage{
		db: db,
	}
}

// CreateUser creates a new user
func (s *PostgreSQLStorage) CreateUser(user *models.User) error {
	result := s.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetUserByPhone gets user by phone number
func (s *PostgreSQLStorage) GetUserByPhone(phone string) (*models.User, error) {
	var user models.User
	result := s.db.Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByID gets user by ID
func (s *PostgreSQLStorage) GetUserByID(id string) (*models.User, error) {
	var user models.User
	result := s.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// CreateSession creates a new session
func (s *PostgreSQLStorage) CreateSession(session *models.Session) error {
	result := s.db.Create(session)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetSession gets session by ID
func (s *PostgreSQLStorage) GetSession(sessionID string) (*models.Session, error) {
	var session models.Session
	result := s.db.Where("id = ?", sessionID).First(&session)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, result.Error
	}
	return &session, nil
}

// MarkSessionAsUsed marks session as used
func (s *PostgreSQLStorage) MarkSessionAsUsed(sessionID string) error {
	result := s.db.Model(&models.Session{}).Where("id = ?", sessionID).Update("is_used", true)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CleanupExpiredSessions removes expired sessions
func (s *PostgreSQLStorage) CleanupExpiredSessions() {
	now := time.Now()
	s.db.Where("expires_at < ?", now).Delete(&models.Session{})
}

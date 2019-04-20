package modules

import (
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"github.com/in2store/service-in2-auth/database"
	"github.com/johnnyeven/libtools/sqlx"
	"github.com/johnnyeven/libtools/sqlx/builder"
	"github.com/johnnyeven/libtools/timelib"
	"strings"
)

// 根据sessionID获取session
func GetSessionBySessionID(sessionID string, db *sqlx.DB) (*database.Session, error) {
	session := &database.Session{
		SessionID: sessionID,
	}
	err := session.FetchBySessionID(db)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// 根据用户ID获取可用的session
func GetSessionByUserID(userID uint64, db *sqlx.DB) (*database.Session, error) {
	s := &database.Session{}
	t := s.T()
	condition := builder.And(t.F("UserID").Eq(userID), t.F("ExpireTime").Gt(timelib.Now()))
	sessions, count, err := s.FetchList(db, -1, 0, condition)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}

	return &sessions[0], nil
}

// 创建session
func CreateSession(userID uint64, db *sqlx.DB) (*database.Session, error) {
	session := &database.Session{
		SessionID: NewSessionID(userID),
		UserID:    userID,
	}
	err := session.Create(db)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// 根据提供的sessionID刷新session
func RefreshSessionID(sessionID string, db *sqlx.DB) (*database.Session, error) {
	session := &database.Session{
		SessionID: sessionID,
	}
	err := session.FetchBySessionID(db)
	if err != nil {
		return nil, err
	}

	refreshSession := &database.Session{
		SessionID: NewSessionID(session.UserID),
		UserID:    session.UserID,
	}
	err = refreshSession.Create(db)
	if err != nil {
		return nil, err
	}

	return refreshSession, nil
}

func NewSessionID(userID uint64) string {
	guid := strings.Replace(uuid.New().String(), "-", "", -1)
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s%d", guid, userID)))
	return base58.Encode(hash[:])
}

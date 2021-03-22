package sessions

import (
	"time"

	"github.com/google/uuid"
)

type SessionInfo struct {
	Name   string
	Expire time.Time
}

type Sessions struct {
	Keys map[string]SessionInfo
}

var GSession Sessions = New()

func New() Sessions {
	return Sessions{
		Keys: make(map[string]SessionInfo),
	}
}

func (s *Sessions) NewSession(username string) (string, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	uuidStr := uuid.String()

	var newSessionInfo SessionInfo
	newSessionInfo.Name = username
	newSessionInfo.Expire = time.Now().Add(24 * time.Hour)
	s.Keys[uuidStr] = newSessionInfo

	return uuidStr, nil
}

func (s *Sessions) ValidateSession(uuid string) bool {
	timeNow := time.Now()
	info, ok := s.Keys[uuid]
	if ok {
		return timeNow.Before(info.Expire)
	} else {
		return false
	}
}

func (s *Sessions) Get(uuid string) (SessionInfo, bool) {
	sessionInfo, ok := s.Keys[uuid]
	return sessionInfo, ok
}

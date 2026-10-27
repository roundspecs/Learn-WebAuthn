package main

import (
	"github.com/couchbase/gocb/v2"
	"github.com/go-webauthn/webauthn/webauthn"
)

type SessionModel struct{ scope *gocb.Scope }

func (m *SessionModel) Insert(key string, s *webauthn.SessionData) {
	m.scope.Collection("sessions").Insert(key, s, nil)
}

func (m *SessionModel) Get(key string) (s webauthn.SessionData) {
	doc, err := m.scope.Collection("sessions").Get(key, nil)
	if err != nil || doc == nil {
		return s
	}
	_ = doc.Content(&s)
	return s
}

func (m *SessionModel) Delete(key string) {
	m.scope.Collection("sessions").Remove(key, nil)
}

type UserSession struct {
	UserID   []byte `json:"userId"`
	Username string `json:"username"`
}

func (m *SessionModel) SaveUserSession(key string, s *UserSession) error {
	_, err := m.scope.Collection("sessions").Upsert("session:"+key, s, nil)
	return err
}

func (m *SessionModel) GetUserSession(key string) (*UserSession, error) {
	doc, err := m.scope.Collection("sessions").Get("session:"+key, nil)
	if err != nil || doc == nil {
		return nil, err
	}
	var s UserSession
	err = doc.Content(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (m *SessionModel) DeleteUserSession(key string) error {
	_, err := m.scope.Collection("sessions").Remove("session:"+key, nil)
	return err
}

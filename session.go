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

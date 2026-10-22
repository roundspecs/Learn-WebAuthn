package main

import (
	"encoding/base64"

	"github.com/couchbase/gocb/v2"
	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	Id          []byte                `json:"id"`
	Name        string                `json:"name"`
	DisplayName string                `json:"displayName"`
	Credentials []webauthn.Credential `json:"credentials"`
}

func (u *User) WebAuthnID() []byte                         { return u.Id }
func (u *User) WebAuthnName() string                       { return u.Name }
func (u *User) WebAuthnDisplayName() string                { return u.DisplayName }
func (u *User) WebAuthnCredentials() []webauthn.Credential { return u.Credentials }

type UserModel struct{ scope *gocb.Scope }

func (m *UserModel) Upsert(u *User) {
	id := base64.StdEncoding.EncodeToString(u.Id)
	m.scope.Collection("users").Upsert(id, u, nil)
}

func (m *UserModel) Get(userID []byte) (*User, error) {
	var u User
	id := base64.StdEncoding.EncodeToString(userID)
	doc, err := m.scope.Collection("users").Get(id, nil)
	if err != nil {
		return nil, err
	}
	err = doc.Content(&u)
	return &u, err
}

func (m *UserModel) GetByUsername(username string) (*User, error) {
	q := `select u.* from users u where u.name == $username`
	result, err := m.scope.Query(q, &gocb.QueryOptions{
		NamedParameters: map[string]any{"username": username},
		ScanConsistency: gocb.QueryScanConsistencyRequestPlus,
	})
	if err != nil {
		return nil, err
	}
	var u User
	err = result.One(&u)
	if err != nil {
		return nil, err
	}
	return &u, err
}

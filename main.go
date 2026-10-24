package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/go-webauthn/webauthn/webauthn"
)

type UsernameRequest struct {
	Username string `json:"username"`
}

var wa, _ = webauthn.New(&webauthn.Config{
	RPID:          "localhost",
	RPDisplayName: "Learn-WebAuthn",
	RPOrigins:     []string{"http://localhost:8000/"},
})

var userModel = &UserModel{}
var sessionModel = &SessionModel{}

func init() {
	scope := cbScope()
	userModel.scope = scope
	sessionModel.scope = scope
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("POST /register/begin", registerBegin)
	mux.HandleFunc("POST /register/finish", registerFinish)
	mux.HandleFunc("POST /login/begin", loginBegin)
	mux.HandleFunc("POST /login/finish", loginFinish)
	srv := http.Server{Addr: ":8000", Handler: mux}
	fmt.Print(srv.ListenAndServe())
}

func cbScope() *gocb.Scope {
	cluster, err := gocb.Connect("couchbase://localhost", gocb.ClusterOptions{Username: "root", Password: "secret"})
	if err != nil {
		fmt.Println(err)
	}
	bucket := cluster.Bucket("leadgen")
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		fmt.Println(err)
	}
	scope := bucket.Scope("learn-webauthn")
	return scope
}

func registerBegin(w http.ResponseWriter, r *http.Request) {
	var req UsernameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	username := req.Username

	if u, _ := userModel.GetByUsername(username); u != nil {
		http.Error(w, "Username exists", http.StatusBadRequest)
		return
	}
	u := &User{
		Id:          make([]byte, 64),
		Name:        username,
		DisplayName: username,
	}
	_, _ = rand.Read(u.Id)
	userModel.Upsert(u)

	creation, session, err := wa.BeginRegistration(u)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	key := getOrSetCookie(w, r, "webauthn_ceremony")
	sessionModel.Insert(key, session)

	json.NewEncoder(w).Encode(creation)
}

func registerFinish(w http.ResponseWriter, r *http.Request) {
	key := getOrSetCookie(w, r, "webauthn_ceremony")
	session := sessionModel.Get(key)
	deleteCookie(w, key)
	sessionModel.Delete(key)
	u, _ := userModel.Get(session.UserID)
	if u == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	credentials, err := wa.FinishRegistration(u, session, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	u.Credentials = append(u.Credentials, *credentials)
	userModel.Upsert(u)
	w.WriteHeader(http.StatusCreated)
}

func loginBegin(w http.ResponseWriter, r *http.Request) {
	var req UsernameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	username := req.Username

	u, err := userModel.GetByUsername(username)
	if err != nil || u == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	assertion, session, err := wa.BeginLogin(u)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	key := getOrSetCookie(w, r, "webauthn_ceremony")
	sessionModel.Insert(key, session)

	json.NewEncoder(w).Encode(assertion)
}

func loginFinish(w http.ResponseWriter, r *http.Request) {
	key := getOrSetCookie(w, r, "webauthn_ceremony")
	session := sessionModel.Get(key)
	deleteCookie(w, key)
	sessionModel.Delete(key)
	u, _ := userModel.Get(session.UserID)
	if u == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	credential, err := wa.FinishLogin(u, session, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	for i, c := range u.Credentials {
		if bytes.Equal(c.ID, credential.ID) {
			u.Credentials[i] = *credential
			break
		}
	}
	userModel.Upsert(u)
	w.WriteHeader(http.StatusOK)
}

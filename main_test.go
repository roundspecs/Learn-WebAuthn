package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"uuid"

	"github.com/go-webauthn/webauthn/webauthn"
)

func TestRegisterBegin(t *testing.T) {
	testUsername := uuid.New().String()
	tests := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
		expectCookie bool
		waitDuration time.Duration
	}{
		{
			"No Body",
			``,
			http.StatusBadRequest,
			"Username is required\n",
			false,
			0,
		},
		{
			"Empty Username",
			`{"username":""}`,
			http.StatusBadRequest,
			"Username is required\n",
			false,
			0,
		},
		{
			"Valid Body",
			`{"username":"` + testUsername + `"}`,
			http.StatusOK,
			"",
			true,
			0,
		},
		{
			"Existing Username",
			`{"username":"` + testUsername + `"}`,
			http.StatusBadRequest,
			"Username exists\n",
			false,
			time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(tt.waitDuration)
			r := httptest.NewRequest(http.MethodPost, "/register/begin", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			registerBegin(w, r)
			b, _ := io.ReadAll(w.Body)
			body := string(b)

			if w.Code != tt.expectedCode {
				t.Errorf("Code: expected %v, got %v\nBody: %v", tt.expectedCode, w.Code, body)
			}
			if tt.expectedBody != "" {
				if string(body) != tt.expectedBody {
					t.Errorf("Body: expected %q, got %q", tt.expectedBody, body)
				}
			}
			if tt.expectCookie {
				foundCookie := false
				cookies := w.Result().Cookies()
				for _, cookie := range cookies {
					if cookie.Name == "webauthn_ceremony" {
						foundCookie = true
					}
				}
				if !foundCookie {
					t.Errorf("Expected cookie with name \"webauthn_ceremony\" but couldn't find it\nBody: %v", body)
				}
			}
		})
	}
}

func TestRegisterFinish(t *testing.T) {
	testUsername := uuid.New().String()
	r := httptest.NewRequest(http.MethodPost, "/register/begin", strings.NewReader(`{"username":"`+testUsername+`"}`))
	w := httptest.NewRecorder()
	registerBegin(w, r)

	var ceremonyCookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == "webauthn_ceremony" {
			ceremonyCookie = c
			break
		}
	}

	tests := []struct {
		name         string
		cookie       *http.Cookie
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "No Cookie",
			cookie:       nil,
			body:         ``,
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Unauthorized\n",
		},
		{
			name:         "Invalid Cookie",
			cookie:       &http.Cookie{Name: "webauthn_ceremony", Value: "invalid-key"},
			body:         ``,
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Unauthorized\n",
		},
		{
			name:         "Empty Body",
			cookie:       ceremonyCookie,
			body:         ``,
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Unauthorized\n",
		},
		{
			name:         "Invalid JSON Body",
			cookie:       ceremonyCookie,
			body:         `{"bad":"data"}`,
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Unauthorized\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/register/finish", strings.NewReader(tt.body))
			if tt.cookie != nil {
				r.AddCookie(tt.cookie)
			}
			w := httptest.NewRecorder()

			registerFinish(w, r)
			b, _ := io.ReadAll(w.Body)
			body := string(b)

			if w.Code != tt.expectedCode {
				t.Errorf("Code: expected %v, got %v\nBody: %v", tt.expectedCode, w.Code, body)
			}
			if tt.expectedBody != "" && body != tt.expectedBody {
				t.Errorf("Body: expected %q, got %q", tt.expectedBody, body)
			}
		})
	}
}

func TestLoginBegin(t *testing.T) {
	testUsername := uuid.New().String()
	u := &User{
		Id:          make([]byte, 64),
		Name:        testUsername,
		DisplayName: testUsername,
		Credentials: []webauthn.Credential{
			{
				ID:              []byte("cred-test-id"),
				PublicKey:       []byte("cred-test-pubkey"),
				AttestationType: "none",
			},
		},
	}
	userModel.Upsert(u)

	tests := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
		expectCookie bool
		waitDuration time.Duration
	}{
		{
			"No Body",
			``,
			http.StatusBadRequest,
			"Username is required\n",
			false,
			0,
		},
		{
			"Empty Username",
			`{"username":""}`,
			http.StatusBadRequest,
			"Username is required\n",
			false,
			0,
		},
		{
			"Non-existent User",
			`{"username":"non-existent-` + uuid.New().String() + `"}`,
			http.StatusNotFound,
			"User not found\n",
			false,
			0,
		},
		{
			"Valid Body",
			`{"username":"` + testUsername + `"}`,
			http.StatusOK,
			"",
			true,
			time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(tt.waitDuration)
			r := httptest.NewRequest(http.MethodPost, "/login/begin", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			loginBegin(w, r)
			b, _ := io.ReadAll(w.Body)
			body := string(b)

			if w.Code != tt.expectedCode {
				t.Errorf("Code: expected %v, got %v\nBody: %v", tt.expectedCode, w.Code, body)
			}
			if tt.expectedBody != "" {
				if string(body) != tt.expectedBody {
					t.Errorf("Body: expected %q, got %q", tt.expectedBody, body)
				}
			}
			if tt.expectCookie {
				foundCookie := false
				cookies := w.Result().Cookies()
				for _, cookie := range cookies {
					if cookie.Name == "webauthn_ceremony" {
						foundCookie = true
					}
				}
				if !foundCookie {
					t.Errorf("Expected cookie with name \"webauthn_ceremony\" but couldn't find it\nBody: %v", body)
				}
			}
		})
	}
}

func TestLoginFinish(t *testing.T) {
	testUsername := uuid.New().String()
	u := &User{
		Id:          make([]byte, 64),
		Name:        testUsername,
		DisplayName: testUsername,
		Credentials: []webauthn.Credential{
			{
				ID:              []byte("cred-test-id"),
				PublicKey:       []byte("cred-test-pubkey"),
				AttestationType: "none",
			},
		},
	}
	userModel.Upsert(u)
	time.Sleep(time.Second)

	r := httptest.NewRequest(http.MethodPost, "/login/begin", strings.NewReader(`{"username":"`+testUsername+`"}`))
	w := httptest.NewRecorder()
	loginBegin(w, r)

	var ceremonyCookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == "webauthn_ceremony" {
			ceremonyCookie = c
			break
		}
	}

	tests := []struct {
		name         string
		cookie       *http.Cookie
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "No Cookie",
			cookie:       nil,
			body:         ``,
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Unauthorized\n",
		},
		{
			name:         "Invalid Cookie",
			cookie:       &http.Cookie{Name: "webauthn_ceremony", Value: "invalid-key"},
			body:         ``,
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Unauthorized\n",
		},
		{
			name:         "Empty Body",
			cookie:       ceremonyCookie,
			body:         ``,
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Unauthorized\n",
		},
		{
			name:         "Invalid JSON Body",
			cookie:       ceremonyCookie,
			body:         `{"bad":"data"}`,
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Unauthorized\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/login/finish", strings.NewReader(tt.body))
			if tt.cookie != nil {
				r.AddCookie(tt.cookie)
			}
			w := httptest.NewRecorder()

			loginFinish(w, r)
			b, _ := io.ReadAll(w.Body)
			body := string(b)

			if w.Code != tt.expectedCode {
				t.Errorf("Code: expected %v, got %v\nBody: %v", tt.expectedCode, w.Code, body)
			}
			if tt.expectedBody != "" && body != tt.expectedBody {
				t.Errorf("Body: expected %q, got %q", tt.expectedBody, body)
			}
		})
	}
}

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

func TestGetSession(t *testing.T) {
	t.Run("No Cookie", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/session", nil)
		w := httptest.NewRecorder()

		getSession(w, r)
		b, _ := io.ReadAll(w.Body)
		body := string(b)

		if w.Code != http.StatusOK {
			t.Errorf("Code: expected %v, got %v\nBody: %v", http.StatusOK, w.Code, body)
		}
		expectedBody := "{\"loggedIn\":false}\n"
		if body != expectedBody {
			t.Errorf("Body: expected %q, got %q", expectedBody, body)
		}
	})

	t.Run("Invalid Cookie", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/session", nil)
		r.AddCookie(&http.Cookie{Name: "session", Value: "non-existent-token"})
		w := httptest.NewRecorder()

		getSession(w, r)
		b, _ := io.ReadAll(w.Body)
		body := string(b)

		if w.Code != http.StatusOK {
			t.Errorf("Code: expected %v, got %v\nBody: %v", http.StatusOK, w.Code, body)
		}
		expectedBody := "{\"loggedIn\":false}\n"
		if body != expectedBody {
			t.Errorf("Body: expected %q, got %q", expectedBody, body)
		}
	})

	t.Run("Valid Session", func(t *testing.T) {
		testUsername := uuid.New().String()
		u := &User{
			Id:          make([]byte, 64),
			Name:        testUsername,
			DisplayName: testUsername,
		}
		userModel.Upsert(u)
		time.Sleep(time.Second)

		sessionToken := generateToken()
		_ = sessionModel.SaveUserSession(sessionToken, &UserSession{
			UserID:   u.Id,
			Username: u.Name,
		})

		r := httptest.NewRequest(http.MethodGet, "/session", nil)
		r.AddCookie(&http.Cookie{Name: "session", Value: sessionToken})
		w := httptest.NewRecorder()

		getSession(w, r)
		b, _ := io.ReadAll(w.Body)
		body := string(b)

		if w.Code != http.StatusOK {
			t.Errorf("Code: expected %v, got %v\nBody: %v", http.StatusOK, w.Code, body)
		}
		expectedBody := `{"loggedIn":true,"username":"` + testUsername + `"}` + "\n"
		if body != expectedBody {
			t.Errorf("Body: expected %q, got %q", expectedBody, body)
		}
	})
}

func TestLogout(t *testing.T) {
	sessionToken := generateToken()
	_ = sessionModel.SaveUserSession(sessionToken, &UserSession{
		UserID:   []byte("test-user-id"),
		Username: "test-user",
	})

	r := httptest.NewRequest(http.MethodPost, "/logout", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: sessionToken})
	w := httptest.NewRecorder()

	logout(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Code: expected %v, got %v", http.StatusOK, w.Code)
	}

	foundDeletedCookie := false
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == "session" && cookie.MaxAge < 0 {
			foundDeletedCookie = true
		}
	}
	if !foundDeletedCookie {
		t.Errorf("Expected session cookie to be deleted")
	}

	sess, _ := sessionModel.GetUserSession(sessionToken)
	if sess != nil {
		t.Errorf("Expected user session to be deleted from Couchbase")
	}
}

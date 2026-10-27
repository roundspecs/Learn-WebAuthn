package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func setCookie(w http.ResponseWriter, name string, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   30 * 24 * 60 * 60, // 30 days
		SameSite: http.SameSiteLaxMode,
	})
}

func getOrSetCookie(w http.ResponseWriter, r *http.Request, name string) (value string) {
	if c, err := r.Cookie(name); err == nil && c.Value != "" {
		return c.Value
	}
	b := make([]byte, 16)
	rand.Read(b)
	value = base64.StdEncoding.EncodeToString(b)
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	})
	return value
}

func deleteCookie(w http.ResponseWriter, name string) (value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	return value
}

package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

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

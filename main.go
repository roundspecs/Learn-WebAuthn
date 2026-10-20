package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("POST /register/begin", registerBegin)
	mux.HandleFunc("POST /register/finish", registerFinish)
	mux.HandleFunc("POST /login/begin", loginBegin)
	mux.HandleFunc("POST /login/finish", loginFinish)

	srv := http.Server{
		Addr: ":8000",
		Handler: mux,
	}

	srv.ListenAndServe()
}

func registerBegin(w http.ResponseWriter, r *http.Request) {}
func registerFinish(w http.ResponseWriter, r *http.Request) {}
func loginBegin(w http.ResponseWriter, r *http.Request) {}
func loginFinish(w http.ResponseWriter, r *http.Request) {}

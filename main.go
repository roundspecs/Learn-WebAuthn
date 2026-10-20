package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /", http.FileServer(http.Dir("./static")))

	srv := http.Server{
		Addr: ":8000",
		Handler: mux,
	}

	srv.ListenAndServe()
}

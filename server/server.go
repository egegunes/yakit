package server

import (
	"net/http"
	"time"
)

func New(h http.Handler, serverAddress string) *http.Server {
	return &http.Server{
		Addr:         serverAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      h,
	}
}

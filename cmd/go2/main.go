package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"runtime/debug"

	"github.com/atlanssia/go2/pkg/httpproxy"
	"github.com/atlanssia/go2/pkg/httpserver"
	"github.com/atlanssia/go2/pkg/log"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			ctx := context.Background()
			log.Error(ctx, "panic: %v", string(debug.Stack()))
		}
		log.Sync()
	}()

	// defer cancelWeb()

	ctx := context.Background()
	log.Info(ctx, "initializing...")

	s := newHttpServer()
	log.Error(ctx, "http fatal: %v", s.ListenAndServe())
}

func newHttpServer() *http.Server {
	server := &http.Server{
		Addr: ":6800",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				httpproxy.HandleTunneling(w, r)
			} else {
				httpserver.HandleHTTP(w, r)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	return server
}

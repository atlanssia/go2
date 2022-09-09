package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/atlanssia/go2/pkg/httpproxy"
	"github.com/atlanssia/go2/pkg/httpserver"
	"github.com/atlanssia/go2/pkg/log"
	"github.com/atlanssia/go2/pkg/utils"
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
	s := newHttpServer(6800)
	log.Error(ctx, "http fatal: %v", s.ListenAndServe())
}

func newHttpServer(port int) *http.Server {
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lrw := utils.NewLoggingResponseWriter(w)
			defer func() {
				log.Access(log.AccessLog{
					RemoteAddr:           r.RemoteAddr,
					Method:               r.Method,
					Proto:                r.Proto,
					RequestContentLength: r.ContentLength,
					Host:                 r.Host,
					RequestURI:           r.RequestURI,
					Status:               lrw.StatusCode(),
					Url:                  r.URL.String(),
					UserAgent:            r.Header.Get("User-Agent"),
					RequestTime:          int64(time.Since(start)),
				})
			}()
			if r.Method == http.MethodConnect {
				httpproxy.HandleTunneling(lrw, r)
			} else {
				httpserver.HandleHTTP(lrw, r)
			}

		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	return server
}

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"runtime/debug"
	"sync"
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
	log.Info(nil, "system initializing...")
	// defer s.Shutdown()
	wg := &sync.WaitGroup{}
	go listenHttp(wg)
	go listenHttps(wg)
	log.Info(nil, "system running...")
	time.Sleep(time.Second)
	wg.Wait()
	log.Info(nil, "system exit")
}

func listenHttp(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	startListening(6800, "", "")
}

func listenHttps(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	cert := "ca/server.pem"
	key := "ca/server.key"
	startListening(6801, cert, key)
}

func startListening(port int, cert string, key string) {
	s := newHttpServer(port)
	var err error
	if cert != "" && key != "" {
		err = s.ListenAndServeTLS(cert, key)
	} else {
		err = s.ListenAndServe()
	}
	if err != nil {
		log.Error(nil, "system exit: %v", err)
	}
}

func newHttpServer(port int) *http.Server {
	log.Info(nil, "initializing listening on port %d", port)
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

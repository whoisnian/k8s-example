package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/logger"
)

const listenAddr = "0.0.0.0:8080"

func pongHandler(store *httpd.Store) {
	logger.Debug("Recerved ping request.")
	store.Respond200([]byte("pong\n"))
}

func main() {
	logger.SetDebug(true)

	go func() {
		mux := httpd.NewMux()
		mux.Handle("/ping", "GET", pongHandler)

		logger.Info("Service started. (", os.Getpid(), ")")
		if err := http.ListenAndServe(listenAddr, logger.Req(logger.Recovery(mux))); err != nil {
			logger.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	rec := <-c
	logger.Warn("Got signal: ", rec)

	time.Sleep(time.Second * 20)
}

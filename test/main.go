package main

import (
	"bytes"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/logger"
)

var (
	listenAddr = "0.0.0.0:8080"

	MiB = bytes.Repeat([]byte("0123456789abcdef"), 65536)
)

func pingHandler(store *httpd.Store) {
	logger.Debug("Received ping request.")
	store.Respond200([]byte("pong\n"))
}

func memHandler(store *httpd.Store) {
	cnt, err := strconv.Atoi(store.R.FormValue("cnt"))
	if err != nil {
		logger.Panic(err)
	}
	logger.Debug("Received mem (", cnt, ") request.")

	buf := [][]byte{}
	for i := 0; i < cnt; i++ {
		tmp := make([]byte, len(MiB))
		copy(tmp, MiB)
		buf = append(buf, tmp)
	}

	time.Sleep(time.Second * 10)
	res := strconv.Itoa(len(buf)*len(buf[0])/1024/1024) + " MiB\n"
	store.Respond200([]byte(res))
}

func main() {
	logger.SetDebug(true)

	mux := httpd.NewMux()
	mux.Handle("/ping", "GET", pingHandler)
	mux.Handle("/mem", "GET", memHandler)

	go func() {
		logger.Info("Service started. (", os.Getpid(), ")")
		if err := http.ListenAndServe(listenAddr, logger.Req(logger.Recovery(mux))); err != nil {
			logger.Fatal(err)
		}
	}()

	done := make(chan struct{})
	once := sync.Once{}
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		for rec := range ch {
			logger.Warn("Signal received: ", rec)
			once.Do(func() { close(done) })
		}
	}()

	<-done
	time.Sleep(time.Second * 20)
	logger.Info("Service stopped.")
}

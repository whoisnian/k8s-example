package main

import (
	"bytes"
	"io"
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
	// config env
	listenAddr = "0.0.0.0:8080"
	upstream   = ""
	modeFile   = "./mode"

	// runtime flag
	modeList   = []string{"normal", "reject", "panic"}
	serverMode = "normal"

	// constant
	MiB = bytes.Repeat([]byte("0123456789abcdef"), 65536)
)

func init() {
	if val, ok := os.LookupEnv("LISTEN_ADDR"); ok {
		listenAddr = val
	}
	if val, ok := os.LookupEnv("UPSTREAM"); ok {
		upstream = val
	}
	if val, ok := os.LookupEnv("MODE_FILE"); ok {
		modeFile = val
	}
	if f, err := os.Open(modeFile); err == nil {
		defer f.Close()
		if fbuf, err := io.ReadAll(f); err == nil {
			fstr := string(bytes.TrimSpace(fbuf))
			for _, mode := range modeList {
				if fstr == mode {
					serverMode = mode
					break
				}
			}
		}
	}
}

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

func upstreamHandler(store *httpd.Store) {
	logger.Debug("Received upstream request.")
	if upstream == "" {
		store.Respond200([]byte("success\n"))
	} else {
		resp, err := http.Get(upstream)
		if err != nil {
			logger.Panic(err)
		}
		defer resp.Body.Close()
		io.Copy(store.W, resp.Body)
	}
}

func main() {
	logger.SetDebug(true)
	if serverMode == "panic" {
		logger.Panic("panic mode")
	}

	mux := httpd.NewMux()
	mux.Handle("/ping", "GET", pingHandler)
	mux.Handle("/mem", "GET", memHandler)
	mux.Handle("/upstream", "GET", upstreamHandler)

	logger.Info("Service started. (", os.Getpid(), ")")
	if serverMode != "reject" {
		go func() {
			if err := http.ListenAndServe(listenAddr, logger.Req(logger.Recovery(mux))); err != nil {
				logger.Fatal(err)
			}
		}()
	}

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

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

	"github.com/whoisnian/glb/config"
	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/logger"
)

var (
	// config env
	CFG struct {
		ListenAddr string `flag:"l,0.0.0.0:8080,Server listen addr"`
		Upstream   string `flag:"u,,Reverse proxy target"`
		ModeFile   string `flag:"m,./mode,Path of modefile"`
		PodName    string `flag:"n,unknown,Present instance name"`
	}

	// runtime flag
	modeList = []string{"normal", "reject", "panic"}

	// constant
	MiB = bytes.Repeat([]byte("0123456789abcdef"), 65536)
)

func parseMode(modeFile string) string {
	if f, err := os.Open(modeFile); err == nil {
		defer f.Close()
		if fbuf, err := io.ReadAll(f); err == nil {
			fstr := string(bytes.TrimSpace(fbuf))
			for _, mode := range modeList {
				if fstr == mode {
					return mode
				}
			}
		}
	}
	return "normal"
}

func pingHandler(store *httpd.Store) {
	logger.Info("Received ping request.")
	store.Respond200([]byte("pong\n"))
}

func memHandler(store *httpd.Store) {
	cnt, err := strconv.Atoi(store.R.FormValue("cnt"))
	if err != nil {
		logger.Panic(err)
	}
	logger.Info("Received mem (", cnt, ") request.")

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
	logger.Info("Received upstream request.")
	if CFG.Upstream == "" {
		store.Respond200([]byte("success\n"))
	} else {
		resp, err := http.Get(CFG.Upstream)
		if err != nil {
			logger.Panic(err)
		}
		defer resp.Body.Close()
		io.Copy(store.W, resp.Body)
	}
}

func podnameHandler(store *httpd.Store) {
	logger.Info("Received podname request.")
	store.Respond200([]byte(CFG.PodName + "\n"))
}

func main() {
	if err := config.FromCommandLine(&CFG); err != nil {
		logger.Panic(err)
	}

	serverMode := parseMode(CFG.ModeFile)
	if serverMode == "panic" {
		logger.Panic("panic mode")
	}

	mux := httpd.NewMux()
	mux.Handle("/ping", "GET", pingHandler)
	mux.Handle("/mem", "GET", memHandler)
	mux.Handle("/upstream", "GET", upstreamHandler)
	mux.Handle("/podname", "GET", podnameHandler)

	logger.Info("Service started. (", os.Getpid(), ")")
	if serverMode != "reject" {
		go func() {
			if err := http.ListenAndServe(CFG.ListenAddr, logger.Req(logger.Recovery(mux))); err != nil {
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
	time.Sleep(time.Second * 5)
	logger.Info("Service stopped.")
}

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/whoisnian/k8s-example/src/file/global"
	"github.com/whoisnian/k8s-example/src/file/router"
)

func main() {
	global.SetupConfig()

	if global.CFG.Version {
		fmt.Printf("%s %s(%s)\n", global.AppName, global.Version, global.BuildTime)
		return
	}

	server := &http.Server{
		Addr:              global.CFG.ListenAddr,
		Handler:           router.Setup(),
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 180,
		MaxHeaderBytes:    http.DefaultMaxHeaderBytes,
	}
	go func() {
		log.Printf("service is starting: http://%s", global.CFG.ListenAddr)
		if err := server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			log.Printf("service is shutting down")
		} else if err != nil {
			log.Fatalln(err)
		}
	}()

	waitFor(syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println(err)
	}
	log.Println("service has been shut down")
}

func waitFor(signals ...os.Signal) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)
	defer signal.Stop(c)

	<-c
}

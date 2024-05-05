package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/whoisnian/k8s-example/src/user/global"
	"github.com/whoisnian/k8s-example/src/user/model"
	"github.com/whoisnian/k8s-example/src/user/router"
	"go.uber.org/zap"
)

func main() {
	global.SetupConfig()
	global.SetupLogger()
	global.LOG.Info("setup config successfully", zap.Any("CFG", global.CFG))

	if global.CFG.Version {
		fmt.Printf("%s %s(%s)\n", global.AppName, global.Version, global.BuildTime)
		return
	}

	global.SetupRedis()
	global.LOG.Info("setup redis successfully")
	global.SetupMySQL()
	global.LOG.Info("setup mysql successfully")

	if global.CFG.AutoMigrate {
		model.SetupAutoMigrate(global.DB)
		global.LOG.Info("setup auto-migrate successfully")
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
		global.LOG.Info("service is starting", zap.String("addr", global.CFG.ListenAddr))
		if err := server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			global.LOG.Warn("service is shutting down")
		} else if err != nil {
			global.LOG.Fatal(err.Error())
		}
	}()

	waitFor(syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		global.LOG.Warn(err.Error())
	}
	global.LOG.Info("service has been shut down")
}

func waitFor(signals ...os.Signal) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)
	defer signal.Stop(c)

	<-c
}
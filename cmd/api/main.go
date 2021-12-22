// cmd/api/main.go

package main

import (
	"github.com/joho/godotenv"
	"github.com/shadowbane/golang-microservice-sekolah/cmd/api/router"
	"github.com/shadowbane/golang-microservice-sekolah/pkg/application"
	"github.com/shadowbane/golang-microservice-sekolah/pkg/exithandler"
	"github.com/shadowbane/golang-microservice-sekolah/pkg/logger"
	"github.com/shadowbane/golang-microservice-sekolah/pkg/server"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	zap.S().Info("Starting Application")

	if err := godotenv.Load(); err != nil {
		zap.S().Warnf("Failed to load env vars!")
	}

	app, err := application.Start()
	if err != nil {
		zap.S().Fatal(err.Error())
	}

	srv := server.
		Get().
		WithAddr(app.Cfg.GetAPIPort()).
		WithRouter(router.Get(app)).
		WithErrLogger(zap.S())

	go func() {
		zap.S().Info("starting server at ", app.Cfg.GetAPIPort())

		if err := srv.Start(); err != nil {
			zap.S().Fatal(err.Error())
		}
	}()

	exithandler.Init(func() {
		if err := srv.Close(); err != nil {
			zap.S().Error(err.Error())
		}
	})
}

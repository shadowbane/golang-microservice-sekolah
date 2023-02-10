package application

import (
	"github.com/jinzhu/gorm"
	"github.com/shadowbane/go-logger"
	"github.com/shadowbane/golang-microservice-sekolah/pkg/config"
)

type Application struct {
	Cfg *config.Config
	DB  *gorm.DB
}

func Start() (*Application, error) {
	cfg := config.Get()
	db := cfg.ConnectToDatabase()

	logger.Init(cfg.LogConfig)

	return &Application{
		Cfg: cfg,
		DB:  db,
	}, nil
}

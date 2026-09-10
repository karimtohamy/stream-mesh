package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"stream-mesh/streaming/internal/bootstrap"
	"stream-mesh/streaming/internal/broker"
	"stream-mesh/streaming/internal/config"
	"stream-mesh/streaming/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	cfg      *config.Config
	rabbitMQ *broker.RabbitClient
	db       *gorm.DB
	server   *http.Server
}

func NewApp(cfg *config.Config) (*App, error) {
	db, err := gorm.Open(mysql.Open(cfg.MySQL.DSN()), &gorm.Config{})
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	if cfg.App.AutoMigrate {
		err = models.InitModels(db, &cfg.App)
		if err != nil {
			fmt.Print(err.Error())
			return nil, err
		}
	}
	rabbitClient, err := broker.NewRabbitClient(cfg)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return &App{
		cfg:      cfg,
		rabbitMQ: rabbitClient,
		db:       db,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	router := gin.Default()
	brokerListener := broker.NewListener(a.rabbitMQ, a.cfg)
	bootstrap.Init(ctx, a.db, router, brokerListener)
	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.App.Port),
		Handler: router,
	}
	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v", err)
		}
	}()
	return nil
}

func (a *App) ShutDown(ctx context.Context) {
	if err := a.server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	a.rabbitMQ.Close()
}

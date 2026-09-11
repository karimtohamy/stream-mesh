package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"stream-mesh/streaming/internal/broker"
	"stream-mesh/streaming/internal/config"
	"stream-mesh/streaming/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	Cfg      *config.Config
	RabbitMQ *broker.RabbitClient
	Listener *broker.Listener
	Db       *gorm.DB
	Server   *http.Server
	Redis    *redis.Client
	Router   *gin.Engine
}

func NewApp() (*App, error) {
	//load config
	cfg, err := config.Load()
	if err != nil {
		log.Print("failed to load config")
		return nil, err
	}
	ctx := context.Background()
	//db connection
	db, err := gorm.Open(postgres.Open(cfg.Postgres.DSN()), &gorm.Config{})
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
	log.Println("connected to Postgres and migrations ran")

	//rabbitMQ conn
	log.Println("connected to RabbitMQ")
	rabbitClient, err := broker.NewRabbitClient(cfg)

	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	//register listener
	brokerListener := broker.NewListener(rabbitClient, cfg)
	//redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr(),
		DB:   cfg.Redis.DB,
	})
	fmt.Print(redisClient)
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	log.Println("connected to Redis")
	return &App{
		Cfg:      cfg,
		RabbitMQ: rabbitClient,
		Listener: brokerListener,
		Db:       db,
		Redis:    redisClient,
		Router:   gin.Default(),
	}, nil
}

func (a *App) Start(ctx context.Context) error {

	a.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Cfg.App.Port),
		Handler: a.Router,
	}
	go func() {
		if err := a.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server error: %v", err)
		}
	}()
	return nil
}

func (a *App) ShutDown(ctx context.Context) {
	if err := a.Server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	a.RabbitMQ.Close()
}

package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	RabbitMQ RabbitMqConfig
	Postgres PostgresConfig
	Redis    RedisConfig
}

type AppConfig struct {
	Port        int
	AutoMigrate bool
	Secret      string
}

type RabbitMqConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Exchange        string
	TranscodeCmpKey string
	TranscodeQueue  string
	UsageTickKey    string
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	Host string
	Port int
	DB   int
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	val := getEnv(key, "")
	if val == "" {
		return fallback
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return n
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		App: AppConfig{
			Port:        getEnvAsInt("STREAMING_PORT", 8086),
			AutoMigrate: getEnv("AUTO_MIGRATE", "false") == "true",
			Secret:      getEnv("JWT_SECRET", ""),
		},
		RabbitMQ: RabbitMqConfig{
			Host:            getEnv("RABBITMQ_HOST", "localhost"),
			Port:            getEnvAsInt("RABBITMQ_PORT", 5672),
			User:            getEnv("RABBITMQ_USER", "guest"),
			Password:        getEnv("RABBITMQ_PASSWORD", "guest"),
			Exchange:        getEnv("RABBIT_EXCHANGE", "stream.events"),
			TranscodeCmpKey: getEnv("RABBIT_TRANSCODE_CMP_KEY", "media.transcode.completed"),
			TranscodeQueue:  getEnv("RABBIT_STREAMING_QUEUE", "streaming.transcode.results"),
			UsageTickKey:    getEnv("RABBIT_USAGE_TICK_KEY", "stream.usage.tick"),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvAsInt("POSTGRES_PORT", 5432),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", ""),
			Database: getEnv("POSTGRES_DATABASE", "stream_db"),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: getEnvAsInt("REDIS_PORT", 6379),
			DB:   getEnvAsInt("REDIS_DB", 0),
		},
	}, nil
}

func (r *RabbitMqConfig) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", r.User, r.Password, r.Host, r.Port)
}

func (p *PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		p.Host, p.User, p.Password, p.Database, p.Port)
}
func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

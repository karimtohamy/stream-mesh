package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	RabbitMQ RabbitConfig
	MinIO    MinIOConfig
	MySQL    MySQLConfig
	Redis    RedisConfig
}

type AppConfig struct {
	Port            int
	TransmuxWorkers int
}

type RabbitConfig struct {
	Host             string
	Port             int
	User             string
	Password         string
	Exchange         string
	TranscodeQueue   string
	TranscodeReqKey  string
	TranscodeCmpKey  string
	UsageTickKey     string
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

type RedisConfig struct {
	Host string
	Port int
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}
	return val
}

func getEnvAsBool(key string, fallback bool) bool {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return fallback
	}
	return val
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		App: AppConfig{
			Port:            getEnvAsInt("MEDIA_SYNC_PORT", 8085),
			TransmuxWorkers: getEnvAsInt("TRANSMUX_WORKERS", 3),
		},
		RabbitMQ: RabbitConfig{
			Host:            getEnv("RABBITMQ_HOST", "localhost"),
			Port:            getEnvAsInt("RABBITMQ_PORT", 5672),
			User:            getEnv("RABBITMQ_USER", "guest"),
			Password:        getEnv("RABBITMQ_PASSWORD", "guest"),
			Exchange:        getEnv("RABBIT_EXCHANGE", "stream.events"),
			TranscodeQueue:  getEnv("RABBIT_TRANSCODE_QUEUE", "media.transcode.events"),
			TranscodeReqKey: getEnv("RABBIT_TRANSCODE_REQ_KEY", "media.transcode.request"),
			TranscodeCmpKey: getEnv("RABBIT_TRANSCODE_CMP_KEY", "media.transcode.completed"),
			UsageTickKey:    getEnv("RABBIT_USAGE_TICK_KEY", "stream.usage.tick"),
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ROOT_USER", "minio_admin"),
			SecretKey: getEnv("MINIO_ROOT_PASSWORD", "minio_secret"),
			UseSSL:    getEnvAsBool("MINIO_USE_SSL", false),
		},
		MySQL: MySQLConfig{
			Host:     getEnv("MYSQL_HOST", "localhost"),
			Port:     getEnvAsInt("MYSQL_PORT", 3306),
			User:     getEnv("MYSQL_USER", "stream_admin"),
			Password: getEnv("MYSQL_PASSWORD", "stream_secret"),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: getEnvAsInt("REDIS_PORT", 6379),
		},
	}, nil
}

func (r *RabbitConfig) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", r.User, r.Password, r.Host, r.Port)
}

func (m *MySQLConfig) DSN(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", m.User, m.Password, m.Host, m.Port, dbName)
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
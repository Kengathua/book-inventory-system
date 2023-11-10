package config

import (
	"os"
	"time"
)

func GetEnv(key string, defaultValue string) string {
	if os.Getenv(key) == "" {
		return defaultValue
	}
	return os.Getenv(key)
}

type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"POSTGRESQL_URL"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	DatabaseURL          string        `mapstructure:"DATABASE_URL"`
	TestDatabaseURL      string        `mapstructure:"TEST_DATABASE_URL"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	PaginationPageSize   string        `mapstructure:"PAGINATION_PAGE_SIZE"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EmailSender          string        `mapstructure:"EMAIL_SENDER"`
	AuthEmailSender      string        `mapstructure:"AUTH_EMAIL_SENDER"`
	BaseDir              string        `mapstructure:"BASE_DIR"`
	CallbackBaseUrl      string        `mapstructure:"CALLBACK_BASE_URL"`
	TaskSchedulerLock    string        `mapstructure:"TASK_SCHEDULER_LOCK"`
}

// LoadConfig reads configuration from environment file or environment variables.
func LoadConfig() (config Config, err error) {
	config.Environment = GetEnv("ENVIRONMENT", "")
	config.DBDriver = GetEnv("DB_DRIVER", "")
	config.MigrationURL = GetEnv("MIGRATION_URL", "")
	config.DatabaseURL = GetEnv("DATABASE_URL", "postgres://demo_user:demo_pass@localhost:5432/demo_test")
	config.DBSource = GetEnv("POSTGRESQL_URL", "")
	config.PaginationPageSize = GetEnv("PAGINATION_PAGE_SIZE", "50")
	config.TestDatabaseURL = GetEnv("TEST_DATABASE_URL", "postgres://demo_user:demo_pass@localhost:5432/demo_test")
	config.EmailSender = GetEnv("EMAIL_SENDER", "")
	config.AuthEmailSender = GetEnv("AUTH_EMAIL_SENDER", "")
	config.BaseDir = GetEnv("BASE_DIR", "")
	config.CallbackBaseUrl = GetEnv("CALLBACK_BASE_URL", "")
	config.TaskSchedulerLock = GetEnv("TASK_SCHEDULER_LOCK", "task_scheduler_lock")
	return config, nil
}

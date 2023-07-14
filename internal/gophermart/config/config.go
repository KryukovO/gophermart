package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	address        = ":8081"
	dsn            = ""
	accrualAddress = ""

	secretKey         = ""
	userTokenTTL      = 30 * time.Minute
	shutdownTimeout   = 10 * time.Second
	repositoryTimeout = 3 * time.Second
	migrations        = "sql/migrations"
	accrualWorkers    = 3
	accrualInterval   = 3 * time.Second
	accrualShutdown   = 3 * time.Second
)

type Config struct {
	Address        string // Адрес эндпоинта сервера (host:port)
	DSN            string // Адрес подключения к БД
	AccrualAddress string // Адрес системы расчёта начислений

	SecretKey          string        // Ключ шифрования токена авторизации
	UserTokenTTL       time.Duration // Время жизни токена пользователя
	ShutdownTimeout    time.Duration // Таймаут для graceful shutdown сервера
	RepositioryTimeout time.Duration // Таймаут соединения с хранилищем
	Migrations         string        // Путь до директории с файлами миграции
	AccrualWorkers     uint          // Количество одновременно исходящих запросов к сервису Accrual
	AccrualInterval    time.Duration // Интервал генерации новой партии запросов к сервису Accrual
	AccrualShutdown    time.Duration // Таймаут для завершения соединения с Accrual
}

func NewConfig() *Config {
	vpr := viper.New()

	vpr.AllowEmptyEnv(false)

	vpr.BindEnv("run_address")
	vpr.BindEnv("database_uri")
	vpr.BindEnv("accrual_system_address")
	vpr.BindEnv("jwt_secret")
	vpr.BindEnv("jwt_ttl")
	vpr.BindEnv("server_shutdown")
	vpr.BindEnv("repository_timeout")
	vpr.BindEnv("database_migrations")
	vpr.BindEnv("accrual_connector_workers")
	vpr.BindEnv("accrual_connector_interval")
	vpr.BindEnv("accrual_connector_shutdown")

	vpr.SetDefault("run_address", address)
	vpr.SetDefault("database_uri", dsn)
	vpr.SetDefault("accrual_system_address", accrualAddress)
	vpr.SetDefault("jwt_secret", secretKey)
	vpr.SetDefault("jwt_ttl", userTokenTTL)
	vpr.SetDefault("server_shutdown", shutdownTimeout)
	vpr.SetDefault("repository_timeout", repositoryTimeout)
	vpr.SetDefault("database_migrations", migrations)
	vpr.SetDefault("accrual_connector_workers", accrualWorkers)
	vpr.SetDefault("accrual_connector_interval", accrualInterval)
	vpr.SetDefault("accrual_connector_shutdown", accrualShutdown)

	return &Config{
		Address:            vpr.GetString("run_address"),
		DSN:                vpr.GetString("database_uri"),
		AccrualAddress:     vpr.GetString("accrual_system_address"),
		SecretKey:          vpr.GetString("jwt_secret"),
		UserTokenTTL:       vpr.GetDuration("jwt_ttl"),
		ShutdownTimeout:    vpr.GetDuration("server_shutdown"),
		RepositioryTimeout: vpr.GetDuration("repository_timeout"),
		Migrations:         vpr.GetString("database_migrations"),
		AccrualWorkers:     vpr.GetUint("accrual_connector_workers"),
		AccrualInterval:    vpr.GetDuration("accrual_connector_interval"),
		AccrualShutdown:    vpr.GetDuration("accrual_connector_shutdown"),
	}
}

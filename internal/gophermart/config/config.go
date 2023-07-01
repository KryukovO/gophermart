package config

type Config struct {
	Address        string `env:"RUN_ADDRESS"`  // Адрес эндпоинта сервера (host:port)
	DSN            string `env:"DATABASE_URI"` // Адрес подключения к БД
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`

	ShutdownTimeout uint   // Таймаут для graceful shutdown сервера
	Migrations      string // Путь до директории с файлами миграции
}

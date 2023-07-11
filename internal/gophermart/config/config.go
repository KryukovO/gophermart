package config

type Config struct {
	Address        string `env:"RUN_ADDRESS"`            // Адрес эндпоинта сервера (host:port)
	DSN            string `env:"DATABASE_URI"`           // Адрес подключения к БД
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"` // Адрес системы расчёта начислений

	SecretKey          string // Ключ шифрования токена авторизации
	RepositioryTimeout uint   // Таймаут соединения с хранилищем, сек
	ShutdownTimeout    uint   // Таймаут для graceful shutdown сервера, сек
	Migrations         string // Путь до директории с файлами миграции
	UserTokenTTL       uint   // Время жизни токена пользователя, мин
	AccrualWorkers     uint   // Количество одновременно исходящих запросов к сервису Accrual
	AccrualInterval    uint   // Интервал генерации новой партии запросов к сервису Accrual, сек
	AccrualShutdown    uint   // Таймаут для завершения соединения с Accrual, сек
}

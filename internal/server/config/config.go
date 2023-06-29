package config

type Config struct {
	Address        string `env:"RUN_ADDRESS"`
	DSN            string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

# Loyalty points service

Сервис представляет собой HTTP API со следующей функциональностью:
- Регистрация, аутентификация и авторизация пользователей
- Приём номеров заказов от зарегистрированных пользователей
- Учёт и ведение списка переданных номеров заказов зарегистрированного пользователя
- Учёт и ведение накопительного счёта зарегистрированного пользователя
- Проверка принятых номеров заказов через систему расчёта баллов лояльности
- Начисление за каждый подходящий номер заказа положенного вознаграждения на счёт лояльности пользователя

## Usage

Сборка бинарного файла сервиса осуществляется командой `make build`. Скомпилированный бинарный файл размещается по пути `cmd/gophermart/gophermart`.

При запуске сервиса считываются значения следующих условных переменных:
- `RUN_ADDRESS` - Адрес и порт запуска сервиса (host:port)
- `DATABASE_URI` - Адрес подключения к БД
- `ACCRUAL_SYSTEM_ADDRESS` - Адрес сервиса расчета баллов лояльности
- `JWT_SECRET` - Ключ шифрования токена авторизации
- `JWT_TTL` - Время жизни токена пользователя
- `SERVER_SHUTDOWN` - Таймаут для graceful shutdown сервера
- `REPOSITORY_TIMEOUT` - Таймаут соединения с хранилищем
- `DATABASE_MIGRATIONS` - Путь до директории с файлами миграции
- `ACCRUAL_CONNECTOR_WORKERS` - Количество одновременно исходящих запросов к сервису расчета баллов лояльности
- `ACCRUAL_CONNECTOR_INTERVAL` - Интервал генерации новой партии запросов к сервису расчета баллов лояльности
- `ACCRUAL_CONNECTOR_SHUTDOWN` - Таймаут для завершения соединения с сервисом расчета баллов лояльности

В случае отсутствия переменной окружения в системе используется значение по умолчанию, кроме того поддерживается следующие флаги запуска, перекрывающие соответствующие значения переменных окружения:
```
-r, --accrual string     Accrual system address
--accshutdown duration   Accrual connector shutdown timeout (default 3s)
-a, --address string     Address to run HTTP server (default ":8081")
-d, --dsn string         URI to database
-h, --help               Shows gophermart usage
--interval duration      Interval for generating requests to Accrual (default 3s)
--migrations string      Directory of database migration files (default "sql/migrations")
--secret string          Authorization token encryption key
--shutdown duration      Server shutdown timeout (default 10s)
--timeout duration       Repository connection timeout (default 3s)
--userttl duration       User token lifetime (default 30m0s)
--workers uint           Number of concurrent requests to Accrual (default 3)
```
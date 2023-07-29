# Cumulative loyalty system "Gophermart"

[![gophermart](https://github.com/KryukovO/gophermart/actions/workflows/gophermart.yml/badge.svg)](https://github.com/KryukovO/gophermart/actions/workflows/gophermart.yml) [![codecov](https://codecov.io/gh/KryukovO/gophermart/branch/master/graph/badge.svg?token=KWJK8NWS4V)](https://codecov.io/gh/KryukovO/gophermart)

Система преднаначена для расчета бонусных начислений и ведения накопительного бонусного счета пользователя и состоит из следующих сервисов:
- Сервис расчета баллов лояльности
- Сервис начисления баллов лояльности

Для взаимодействия с сервисами предоставляется [HTTP API](./docs/api.md).

Используемые технологии:
- PostgreSQL (в качестве хранилища данных)
- Docker (для запуска сервиса)
- Swagger (для документации API)
- Echo (веб фреймворк)
- golang-migrate/migrate (для миграций БД)
- pgx (драйвер для работы с PostgreSQL)
- golang/mock, testify (для тестирования)

# Getting started

Для запуска системы необходимо создать файл `.env`, используя шаблон `.env.example`, и заполнить недостающие значения.

# Usage

Запуск системы можно осуществить в форме docker-контейнеров с помощью команды `make docker-run`. Остановка системы осуществляется командой `make docker-stop`. Сервис расчета баллов лояльности после запуска будет доступен по адресу http://127.0.0.1:8080, Сервис начисления баллов лояльности - http://127.0.0.1:8081.

Альтернативой запуску контенеров является ручной запуск сервисов, подробнее об этом [здесь](./cmd/accrual/README.md) и [здесь](./cmd/gophermart/README.md).

Документацию по сервису начисления баллов лояльности после запуска можно посмотреть по адресу http://127.0.0.1:8081/swagger/index.html.

Unit-тестирование можно провести посредством трёх команд:
1. `make test` - простое тестирование
2. `make cover` - тестирование с покрытием
3. `make cover-html` - тестирование в результате которого будет сгенерирован файл `cover.html` с подробным отчётом о покрытии

Запуск линтера осуществляется командой `make lint`.

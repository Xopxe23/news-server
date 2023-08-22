# REST API Для Создания News на Go

## Реализованы вторизация и добавление в закладки

## В курсе разобранны следующие концепции:
- Разработка Веб-Приложений на Go, следуя дизайну REST API.
- Подход Чистой Архитектуры в построении структуры приложения. Техника внедрения зависимости.
- Работа с БД Postgres. Генерация файлов миграций. 
- Конфигурация приложения с помощью библиотек <a href="https://github.com/spf13/viper">spf13/viper</a> и <a href="github.com/kelseyhightower/envconfig">kelseyhightower/envconfig</a>. Работа с переменными окружения.
- Работа с БД, используя библиотеку sql</a>.
- Регистрация и аутентификация. Работа с JWT. Middleware.
- Написание SQL запросов.
- Graceful Shutdown

### Для запуска приложения:

```
go run cmd/main.go 
```

Если приложение запускается впервые, необходимо применить миграции к базе данных указав бд в .env:

```
migrate -path ./schema -database "postgresql://user:password@host:port/db_name?sslmode=..." up
```

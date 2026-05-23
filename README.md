# Go CRUD API

pet проект REST API для управления пользователями, продуктами и корзинами.

## Технологии

| Технология | Назначение |
|---|---|
| [Go 1.26.3](https://go.dev/) | Язык программирования |
| [chi v5](https://github.com/go-chi/chi) | HTTP-роутер |
| [pgx v5](https://github.com/jackc/pgx) | PostgreSQL-драйвер |
| [sqlc](https://sqlc.dev/) | Генерация типизированного Go-кода из SQL |
| [golang-migrate](https://github.com/golang-migrate/migrate) | Миграции базы данных |
| [swaggo/swag](https://github.com/swaggo/swag) | Генерация OpenAPI / Swagger UI |
| [godotenv](https://github.com/joho/godotenv) | Загрузка `.env` файлов |
| [google/uuid](https://github.com/google/uuid) | Генерация UUIDv7 |
| [PostgreSQL 16](https://www.postgresql.org/) | База данных |
| [Docker / Compose](https://docs.docker.com/compose/) | Контейнеризация |

## Структура проекта

```
crud/
├── docker-compose.yml
└── src/
    ├── Dockerfile
    ├── cmd/
    │   ├── main.go          # точка входа
    │   └── seed/main.go     # наполнение БД тестовыми данными
    ├── internal/
    │   ├── db/              # sqlc-сгенерированный код
    │   └── handler/         # HTTP-хендлеры
    ├── migrations/          # SQL-миграции (up/down)
    ├── query/               # SQL-запросы для sqlc
    └── sqlc.yaml
```

## Эндпоинты

### Пользователи
| Метод | URL | Описание |
|---|---|---|
| `GET` | `/users` | Список пользователей |
| `POST` | `/users` | Создать пользователя |
| `GET` | `/users/{id}` | Получить пользователя |
| `PUT` | `/users/{id}` | Обновить пользователя |
| `DELETE` | `/users/{id}` | Удалить пользователя |

### Продукты
| Метод | URL | Описание |
|---|---|---|
| `GET` | `/products` | Список продуктов |
| `POST` | `/products` | Создать продукт |
| `GET` | `/products/{id}` | Получить продукт |
| `PUT` | `/products/{id}` | Обновить продукт |
| `DELETE` | `/products/{id}` | Удалить продукт |

### Корзины
| Метод | URL | Описание |
|---|---|---|
| `GET` | `/carts/user/{userID}` | Корзина пользователя |
| `POST` | `/carts` | Добавить товар в корзину |
| `PUT` | `/carts/{id}` | Изменить количество |
| `DELETE` | `/carts/{id}` | Удалить из корзины |

## Запуск

### Docker

```bash
docker compose up --build
```

### Локально

```bash
cp src/.env.example src/.env.example
cd src
go run ./cmd/main.go
```

### Seed (тестовые данные)

```bash
cd src
go run ./cmd/seed/main.go
```

### Swagger UI

```
http://localhost:5000/swagger/index.html
```

### Регенерация sqlc

```bash
cd src
sqlc generate
```

### Регенерация swagger

```bash
cd src
swag init -g cmd/main.go --output docs
```

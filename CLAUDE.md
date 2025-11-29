# LabGuard

Система лицензирования для продажи лабораторных работ и курсовых.

## Стек

- Go 1.25
- PostgreSQL (pgx, чистый SQL)
- Telegram Bot (gopkg.in/telebot.v4)
- chi router
- golang-migrate для миграций
- yaml конфиги

## Архитектурные принципы

### Dependency Injection

Все зависимости передаются через конструкторы. Никаких глобальных переменных.

```go
// Сервис зависит от интерфейса, не от реализации
type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

### Интерфейсы

Следуем идиоматичному Go-подходу: **интерфейсы определяются на стороне потребителя** (consumer-side), а не производителя.

```go
// ✅ ПРАВИЛЬНО: интерфейс определён в пакете, который его использует
// internal/bot/handlers/start.go
type RegisterAPIClient interface {
    CheckUserExists(telegramID int64) (bool, error)
    RegisterUser(telegramID int64, name, group string) (string, error)
}

func NewStartHandler(apiClient RegisterAPIClient, logger *slog.Logger) *StartHandler {
    return &StartHandler{client: apiClient, logger: logger}
}

// internal/bot/services/api/client.go
type HttpClient struct { ... }

// HttpClient реализует RegisterAPIClient, но НЕ импортирует интерфейс
func (c *HttpClient) CheckUserExists(telegramID int64) (bool, error) { ... }
func (c *HttpClient) RegisterUser(telegramID int64, name, group string) (string, error) { ... }
```

**Преимущества:**
- Минимальные интерфейсы — только нужные методы
- Нет циклических зависимостей
- Легко тестировать — моки создаются прямо в тестах
- Меньше coupling между пакетами

**Не создаём отдельную папку `internal/interfaces/`** — это anti-pattern для Go.

### Слои приложения

```
handler → service → repository (storage)
    ↓         ↓           ↓
 валидация  логика      БД
```

- **handler** — принимает запросы, валидирует вход, вызывает service
- **service** — бизнес-логика, работает с интерфейсами repository
- **storage** — реализация работы с БД

### Валидация

Валидация входных данных выносится в отдельный пакет `internal/bot/validators/`. Каждый файл содержит валидаторы для конкретной области (регистрация, покупки, продукты). Валидаторы возвращают предопределённые ошибки для единообразия сообщений пользователю.

## Структура проекта

```
labguard/
├── cmd/
│   ├── server/main.go              # Только app.Run()
│   ├── bot/main.go                 # Только app.Run()
│   └── client/main.go
│
├── internal/
│   ├── server/
│   │   ├── app/
│   │   │   └── app.go              # DI, сборка зависимостей, запуск
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── handler/
│   │   │   ├── handler.go          # Регистрация роутов
│   │   │   ├── user.go
│   │   │   ├── product.go
│   │   │   ├── purchase.go
│   │   │   ├── device.go
│   │   │   └── verify.go
│   │   ├── service/
│   │   │   ├── user/
│   │   │   │   └── user.go
│   │   │   ├── product/
│   │   │   │   └── product.go
│   │   │   ├── purchase/
│   │   │   │   └── purchase.go
│   │   │   ├── device/
│   │   │   │   └── device.go
│   │   │   └── license/
│   │   │       └── license.go
│   │   └── middleware/
│   │       └── ...
│   │
│   ├── bot/
│   │   ├── app/
│   │   │   └── app.go              # DI, сборка зависимостей, запуск, регистрация хендлеров
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── handlers/
│   │   │   ├── base.go             # Базовый обработчик для переиспользования
│   │   │   ├── start.go            # Регистрация пользователя (/start)
│   │   │   ├── products.go         # Список продуктов (/products)
│   │   │   ├── buy.go              # Покупка продукта (/buy)
│   │   │   ├── my.go               # Личный кабинет (/my)
│   │   │   └── devices.go          # Сброс fingerprint (/devices)
│   │   ├── services/
│   │   │   └── api/
│   │   │       └── client.go       # Реализует интерфейсы из handlers
│   │   ├── validators/
│   │   │   ├── registration.go     # Валидация ФИО и группы
│   │   │   ├── product.go          # Валидация product_slug
│   │   │   └── purchase.go         # Валидация для покупок
│   │   ├── keyboards/
│   │   │   ├── start.go            # Клавиатура подтверждения регистрации
│   │   │   └── products.go         # Inline-клавиатура списка продуктов
│   │   ├── models/
│   │   │   └── products.go         # Модель продукта для бота
│   │   └── middleware/
│   │       └── loggers/
│   │           └── message.go
│   │
│   ├── client/
│   │   ├── app/
│   │   │   └── app.go
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── fingerprint/
│   │   │   └── fingerprint.go
│   │   └── api/
│   │       └── api.go              # HTTP клиент к server
│   │
│   ├── storage/
│   │   ├── storage.go
│   │   └── postgres/
│   │       ├── postgres.go         # Подключение, транзакции
│   │       ├── user.go
│   │       ├── product.go
│   │       └── purchase.go
│   │
│   └── model/
│       ├── user.go                 # Модели для server
│       ├── product.go
│       └── purchase.go
│
├── pkg/
│   └── cache/
│       └── cacheTTL.go             # Generic кеш с TTL
│
├── assets/
│   └── labguard.exe                # Скомпилированный клиент
│
├── migrations/
│   └── 001_init.sql
│
├── configs/
│   ├── bot.yaml
│   └── server.yaml
│
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

## Модель данных

```sql
-- users
id SERIAL PRIMARY KEY
telegram_id BIGINT UNIQUE NOT NULL
token VARCHAR(64) UNIQUE NOT NULL
fingerprint VARCHAR(255)
fingerprint_updated_at TIMESTAMP

-- products
id SERIAL PRIMARY KEY
slug VARCHAR(50) UNIQUE NOT NULL
name VARCHAR(255) NOT NULL
price DECIMAL(10,2) NOT NULL
github_repo VARCHAR(255)

-- purchases
id SERIAL PRIMARY KEY
user_id INT REFERENCES users(id)
product_id INT REFERENCES products(id)
purchased_at TIMESTAMP DEFAULT NOW()
UNIQUE(user_id, product_id)
```

**Ограничения:**
- Один fingerprint на пользователя
- Сброс fingerprint: раз в 30 дней через бота
- Одна покупка = один продукт для юзера (UNIQUE constraint)

## API Эндпоинты

### Для клиента (exe)

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| POST | /api/v1/verify | Проверка лицензии |

**Запрос:**
```json
{
    "token": "abc-123-xyz",
    "product_id": "lab1",
    "fingerprint": "fp_xxxxx"
}
```

**Ответ:**
```json
{"valid": true}
// или
{"valid": false, "reason": "device_mismatch | not_purchased | invalid_token"}
```

### Для бота

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| POST | /api/v1/users | Регистрация пользователя (/start) |
| POST | /api/v1/purchases | Создание покупки |
| POST | /api/v1/device/reset | Сброс fingerprint |
| GET | /api/v1/users/:telegram_id/purchases | Мои покупки (личный кабинет) |
| GET | /api/v1/products?telegram_id=123 | Все продукты с флагом purchased |

**GET /api/v1/users/:telegram_id/purchases**
```json
{
    "token": "abc-123-xyz",
    "purchases": [
        {
            "product_slug": "lab1",
            "name": "Лабораторная работа №1",
            "purchased_at": "2025-01-15T10:00:00Z"
        }
    ]
}
```

**GET /api/v1/products?telegram_id=123**
```json
{
    "products": [
        {
            "slug": "lab1",
            "name": "Лабораторная работа №1",
            "price": 500,
            "purchased": true
        },
        {
            "slug": "lab2",
            "name": "Лабораторная работа №2",
            "price": 500,
            "purchased": false
        }
    ]
}
```

## Сервисы

### Bot

Telegram бот для авторизации и покупки. Общается с сервером через HTTP API.

**Команды:**
- `/start` — регистрация, создание токена
- `/products` — список доступных продуктов
- `/buy <product>` — начать покупку
- `/my` — мои покупки и токен
- `/devices` — сброс fingerprint

**Флоу покупки (MVP):**
1. Пользователь: `/buy lab1`
2. Бот: показывает цену и реквизиты карты
3. Пользователь: оплачивает, нажимает "Я оплатил"
4. Бот: создаёт заявку, уведомляет админа
5. Админ: проверяет оплату, подтверждает в боте
6. Бот: создаёт purchase через API, отправляет `labguard.exe` + `labguard.key`
7. Доступ к GitHub репозиторию: выдаётся вручную

### Server

REST API сервер. Обрабатывает запросы от бота и клиента.

Эндпоинты описаны в секции "API Эндпоинты".

### Client

Exe-файл, лежит в корне продукта.

**Файлы:**
- `labguard.exe` — клиент проверки
- `labguard.key` — файл с токеном пользователя

**Логика:**
1. Читает токен из `labguard.key`
2. Генерирует fingerprint машины (CPU ID, MAC, hostname)
3. Отправляет POST /api/v1/verify
4. Exit code: 0 = доступ есть, 1 = нет доступа

## Конфигурация

```yaml
# configs/bot.yaml
telegram:
  token: "BOT_TOKEN"
  admin_id: 123456789

database:
  dsn: "postgres://user:pass@localhost:5432/labguard?sslmode=disable"

server:
  url: "https://api.labguard.ru"
```

## Разработка

Проект пишется вручную. Claude Code используется для ревью кода.

### Git коммиты

Используются [Conventional Commits](https://www.conventionalcommits.org/ru/v1.0.0/):

- `feat:` — новая функциональность
- `fix:` — исправление бага
- `refactor:` — рефакторинг без изменения функциональности
- `docs:` — изменения в документации
- `test:` — добавление/изменение тестов
- `chore:` — рутинные задачи (обновление зависимостей, конфиги)
- `style:` — форматирование кода (без изменения логики)

**Формат сообщения:**
```
<тип>: <краткое описание>

- Значимое изменение 1
- Значимое изменение 2
- Значимое изменение 3
```

**Требования:**
- Заголовок: до 72 символов, описывает суть изменения
- Тело (опционально): список ключевых изменений, если их несколько
- Сообщение должно быть лаконичным, без лишних деталей
- **Перед составлением сообщения:** всегда проверяй изменения через `git diff` или `git diff <commit>`, чтобы точно описать все изменения
- Всегда обращай внимание на все изменения. В том числе новые файлы. А не только, файлы где изменения.

**Примеры:**
```
feat: добавил команду /start в бота

fix: исправил обработку ошибок в UserService

refactor: прокинул логгер через DI

- Добавил поле Logger в BotApp
- Передаю logger через конструктор NewBot
- Удалил глобальные переменные

docs: обновил CLAUDE.md с правилами коммитов
```

### Принципы ревью

Поскольку проект пишется вручную и это MVP, некоторые архитектурные принципы могут не соблюдаться сразу. Claude Code должен:

- **Мягко напоминать** о несоответствиях CLAUDE.md, но не блокировать работу
- **Спрашивать**, а не требовать: "Вижу что здесь нет конфига, хочешь добавить или пока оставим так?"
- **Предлагать обновить CLAUDE.md**, если решение отличается от запланированного: "Ты решил сделать X вместо Y. Обновить CLAUDE.md?"
- **Не считать ошибкой** отступления от плана — это нормально для MVP
- **Фокусироваться на реальных багах**, а не на архитектурных несоответствиях

Примеры допустимых отклонений в MVP:
- Отсутствие некоторых слоёв (handler → storage напрямую)
- Хардкод вместо конфига
- Упрощённая обработка ошибок
- Отсутствие middleware

**Запуск:**
```bash
# БД
docker-compose up -d postgres

# Миграции
migrate -path migrations -database "postgres://..." up

# Бот
go run cmd/bot/main.go -config configs/bot.yaml

# Сервер
go run cmd/server/main.go -config configs/server.yaml
```

**Сборка клиента:**
```bash
GOOS=windows GOARCH=amd64 go build -o assets/labguard.exe ./cmd/client
```

**Порядок разработки:**
1. Bot — команды, регистрация, покупка
2. Server — эндпоинты API
3. Client — проверка лицензии

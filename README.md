# Микросервисная система управления грузоперевозками (Go)

Данный проект представляет собой полнофункциональную систему для организации грузоперевозок, построенную на микросервисной архитектуре на языке Go. Система включает аутентификацию, управление грузами, маршрутами, заказами, уведомлениями и платежами. Все микросервисы контейнеризированы и могут быть легко развёрнуты с помощью Docker Compose.

## Основные возможности

- Регистрация и аутентификация пользователей (JWT).
- Управление грузами (CRUD, фильтрация по владельцу).
- Построение маршрутов и добавление грузов с автоматическим созданием точек погрузки/разгрузки.
- Оформление заказов с привязкой к грузу и маршруту, генерация PDF-договора.
- Эмуляция отправки уведомлений (email/SMS) с сохранением в базу.
- Управление платежами и имитация оплаты.
- Единый API Gateway с проверкой JWT и маршрутизацией запросов к микросервисам.
- Автоматические миграции базы данных при запуске сервисов.
- Полная контейнеризация через Docker Compose.

## Технологический стек

- **Язык:** Go 1.24
- **Веб-фреймворк:** - в проекте используется `chi` для маршрутизации.
- **Базы данных:** PostgreSQL (отдельная база на сервис)
- **Аутентификация:** JWT (golang-jwt), bcrypt
- **Миграции:** golang-migrate
- **Контейнеризация:** Docker, Docker Compose
- **Межсервисное взаимодействие:** HTTP (c использованием `httpx` в Go – `net/http`)

## Состав микросервисов

| Сервис | Порт | База данных | Назначение |
|--------|------|-------------|------------|
| User Service | 8000 | `user_db` | Регистрация, аутентификация, JWT |
| Cargo Service | 8001 | `cargo_db` | CRUD грузов |
| Route Service | 8002 | `route_db` | Управление маршрутами и точками |
| Order Service | 8003 | `order_db` | Заказы, PDF, интеграция |
| Notification Service | 8006 | `notification_db` | Уведомления (эмуляция) |
| Payment Service | 8007 | `payment_db` | Платежи |
| API Gateway | 8005 | — | Единая точка входа, аутентификация, прокси |

## Запуск проекта

### Предварительные требования

- Docker и Docker Compose (версия 3.8+)
- Git

### Инструкция по запуску

1. **Клонируйте репозиторий**
   ```bash
   git clone <repository-url>
   cd go_project_transfer
   ```

2. **(Опционально) Измените переменные окружения**  
   Все настройки уже заданы в `docker-compose.yml` и в коде. При необходимости отредактируйте файл `scripts/init-databases.sql` для создания баз данных.

3. **Запустите все сервисы**
   ```bash
   docker-compose up --build
   ```

4. **Проверьте работоспособность**
   - API Gateway доступен на `http://localhost:8005`
   - Документация (Swagger) пока не реализована, но можно тестировать через curl.

### Примеры запросов

- **Регистрация**
  ```bash
  curl -X POST http://localhost:8005/api/register \
    -H "Content-Type: application/json" \
    -d '{"email":"user@example.com","password":"123","full_name":"User","role":"driver"}'
  ```

- **Логин**
  ```bash
  curl -X POST http://localhost:8005/api/login \
    -H "Content-Type: application/json" \
    -d '{"email":"user@example.com","password":"123"}'
  ```

- **Создание груза** (требуется токен)
  ```bash
  TOKEN=<your_jwt_token>
  curl -X POST http://localhost:8005/api/cargo \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"title":"Груз","weight":100,"pickup_location":"A","delivery_location":"B","pickup_date":"2026-05-10T09:00:00Z","delivery_date":"2026-05-15T18:00:00Z","owner_id":1}'
  ```

- **Создание маршрута и добавление груза**
  ```bash
  curl -X POST http://localhost:8005/api/routes -H "Authorization: Bearer $TOKEN"
  # Запомните route_id
  curl -X POST http://localhost:8005/api/routes/1/cargo/1 -H "Authorization: Bearer $TOKEN"
  ```

- **Создание заказа**
  ```bash
  curl -X POST http://localhost:8005/api/orders \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"cargo_id":1,"customer_id":1}'
  ```

- **Подтверждение заказа** (создаёт уведомление и платёж)
  ```bash
  curl -X POST http://localhost:8005/api/orders/1/confirm -H "Authorization: Bearer $TOKEN"
  ```

- **Просмотр платежей**
  ```bash
  curl -X GET http://localhost:8005/api/payments -H "Authorization: Bearer $TOKEN"
  ```

## Остановка системы

```bash
docker-compose down -v   # -v удаляет тома базы данных
```

## Структура проекта

```
.
├── api-gateway/
├── cargo-service/
├── route-service/
├── order-service/
├── notification-service/
├── payment-service/
├── user-service/
├── pkg/
│   └── migrate/             # общий пакет миграций
├── scripts/
│   └── init-databases.sql   # создание баз данных
├── docker-compose.yml
└── README.md
```

## Дальнейшее развитие

- Добавить сервисы: Parser, Cart, Analytics.
- Реализовать автоматическую генерацию PDF с реальными данными.
- Настроить централизованное логирование и метрики.
- Написать сквозные тесты.

## Лицензия

MIT
```

Это README описывает проект, его возможности, запуск и примеры использования. Следующим файлом я предоставлю **ARCHITECTURE.md** с диаграммой и описанием взаимодействий.
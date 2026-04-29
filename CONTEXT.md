## План переписывания `transfer-project` на Go (от новичка до микросервисов)
### Этап 1. Шаблон микросервиса на Go + HTTP + PostgreSQL

**Цель:** Создать базовый шаблон, который можно копировать для каждого микросервиса.

- Использовать стандартный пакет `net/http` для REST API.
- Подключение к PostgreSQL через `database/sql` + драйвер `pgx` (или `lib/pq`).
- Простая миграция схем (например, `golang-migrate` или самописный скрипт).
- Обработка конфигурации через переменные окружения (`os.Getenv`).
- Структура папок шаблона:
  ```
  template/
    cmd/          (точка входа)
    internal/
      config/
      db/
      handlers/
      models/
      services/
      utils/
    go.mod
  ```
- Реализовать health check (`/health`) и простой эндпоинт `GET /ping`.
- Написать Dockerfile для сервиса.
- Протестировать локально с PostgreSQL в Docker.

**Результат:** Готовый шаблон, который можно копировать для следующих сервисов.

### Этап 2. User Service 

**Цель:** Реализовать аутентификацию и управление пользователями.

- Модель `User` (структура) с полями: ID (uint), Email, HashedPassword, Role, FullName, CreatedAt.
- Роуты:
  - `POST /register` – регистрация, хеширование пароля (bcrypt).
  - `POST /login` – выдача JWT (использовать `golang-jwt/jwt`).
  - `GET /me` – получение профиля (защищённый).
- JWT middleware для проверки токена.
- Хранилище на PostgreSQL.
- Тесты: ручные через curl, написать простые unit-тесты.

**Результат:** Полностью рабочий User Service на Go.

### Этап 3. Cargo Service

**Цель:** CRUD для грузов.

- Модель `Cargo`: ID, Title, Description, Weight, PickupLocation, DeliveryLocation, PickupDate, DeliveryDate, OwnerID, Status, CreatedAt.
- Эндпоинты: `POST /`, `GET /`, `GET /{id}`, `PUT /{id}`, `DELETE /{id}`.
- Аутентификация через заголовки `X-User-ID` и `X-User-Role`, извлечённые из JWT в Gateway (пока нет Gateway, при тестировании просто передаём заголовки).
- Фильтрация через query-параметры (status, owner_id).
- Импорт тестовых грузов из JSON (эндпоинт `/import-test`).
- Тесты.

**Результат:** Cargo Service на Go.

### Этап 4. API Gateway

**Цель:** Единая точка входа, проверка JWT, проксирование запросов.

- Реализовать Gateway на Go с использованием `net/http/httputil.ReverseProxy` или вручную через `http.Client`.
- Middleware для проверки JWT и добавления заголовков `X-User-ID`, `X-User-Role`.
- Маршрутизация по префиксам: `/api/users/*` → User Service, `/api/cargo/*` → Cargo Service.
- Проброс correlation ID (генерировать UUID).
- Тестирование через curl.

**Результат:** Gateway, который проксирует запросы к User и Cargo сервисам.

### Этап 5. Route Service 

**Цель:** Управление маршрутами и точками.

- Модели: `Route` и `RoutePoint`.
- Эндпоинты: создание маршрута, добавление точки, получение маршрута с точками, обновление точки.
- Интеграция с Cargo Service (клиент на Go через HTTP).
- Добавление груза в маршрут: запрос данных из Cargo, создание двух точек.
- Заглушка расчёта расстояний (фиксированное значение).
- Тесты с моками.

**Результат:** Route Service на Go, взаимодействующий с Cargo.

### Этап 6. Order Service 

**Цель:** Управление заказами, интеграция с Cargo, Route, Notification, Payment.

- Модель `Order`.
- Эндпоинт `POST /orders`: проверка груза через Cargo, создание маршрута через Route.
- Эндпоинт `PATCH /orders/{id}/status`: при подтверждении генерировать PDF (использовать пакет `gopdf` или вызов утилиты), отправлять уведомление (HTTP вызов Notification Service).
- Эндпоинт `POST /orders/{id}/confirm`.
- Клиенты для Cargo, Route, Notification, Payment (реализация HTTP-клиентов).
- Генерация PDF – отдельный сервис или функция.

**Результат:** Order Service на Go, с интеграцией.

### Этап 7. Notification Service

**Цель:** Эмуляция отправки уведомлений.

- Модель `Notification`.
- Эндпоинт `POST /notify`: сохраняет в БД и логирует отправку.
- GET `/notifications` – история (для админа).
- Можно добавить очередь (через каналы или Redis) для асинхронности.

**Результат:** Notification Service на Go.

### Этап 8. Payment Service

**Цель:** Управление платежами.

- Модель `Payment`.
- Эндпоинты: `POST /payments` (создать), `GET /payments`, `GET /payments/{id}`, `PATCH /payments/{id}/pay` (имитация оплаты).

**Результат:** Payment Service.

### Этап 9. Оставшиеся микросервисы (Parser, Cart, Analytics)

- **Parser Service:** парсинг JSON, импорт в Cargo.
- **Cart Service:** корзина, checkout.
- **Analytics Service:** отчёты (пока без Spark, просто агрегаты из БД). Можно добавить аналитику через SQL запросы.

### Этап 10. Интеграция, тестирование, документация

- Запуск всех сервисов через Docker Compose.
- Сквозные тесты (написать на Go с использованием `testing` и `httptest`).
- Документация: README, ARCHITECTURE, API.
- Настройка CI (GitHub Actions) для сборки и тестов.

### Дополнительные улучшения (по желанию)

- Перейти на gRPC для внутреннего общения (вместо HTTP).
- Добавить OpenTelemetry для трассировки.
- Использовать очереди (NATS, RabbitMQ) для асинхронных уведомлений.
- Реализовать кэширование через Redis.


### Рекомендации по обучению параллельно

- После **Этапа 1** - основы веб-сервера на Go.
- После **Этапа 2** - работ с БД, хеширование, JWT.
- После **Этапа 4** - проксирование, middleware.
- После **Этапа 5** - HTTP-клиенты, обработка JSON.
- После **Этапа 6** - интеграционное тестирование, генерация PDF.

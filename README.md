# ByteChat

**ByteChat** — это простой мессенджер на Go, использующий:
- **Gin** для HTTP-эндпоинтов (регистрация, логин, профиль, логаут),
- **PostgreSQL** для хранения пользователей и чёрного списка токенов,
- **JWT** для авторизации пользователей,
- **WebSocket** (Gorilla WebSocket) для чата в реальном времени.

Проект демонстрирует полный цикл от **регистрации** и **логина** до **онлайн-чата** и **сохранения** сообщений (опционально).

---

## Возможности
1. **Регистрация** (`POST /register`):
   - Создание нового пользователя с `username` и `password` (хешируется `bcrypt`).
2. **Авторизация** (`POST /login`):
   - Выдача **JWT-токена** при успешном логине.
3. **Профиль** (`GET /profile`):
   - Защищённый эндпоинт (требует `Authorization: Bearer <token>`).
4. **Логаут** (`POST /logout`):
   - Добавляет токен в «чёрный список», блокируя дальнейшее использование.
5. **Чат** (`GET /ws`):
   - Реальное время, клиенты подключаются по WebSocket.
   - Отправленные сообщения рассылаются всем онлайн-пользователям.
   - (Опционально) сохраняются в базе данных.

---

## Структура проекта
```
bytechat/
├── go.mod
├── main.go
└── internal/
    ├── handlers/
    │   ├── user_handler.go     # эндпоинты /register, /login, /logout, /profile
    │   └── chat_handler.go     # эндпоинт /ws для WebSocket
    ├── chat/
    │   ├── hub.go              # Hub - управление всеми клиентами
    │   ├── client.go           # Client - отдельное WebSocket-соединение
    │   └── message.go          # (опционально) модель сообщения
    ├── models/
    │   └── user.go             # модель пользователя
    ├── repositories/
    │   ├── user_repository.go  # хранение / поиск пользователей в БД
    │   └── message_repository.go # (опционально) хранение сообщений
    ├── services/
    │   ├── user_service.go     # бизнес-логика (регистрация, логин)
    │   └── chat_service.go     # (опционально) логика сохранения сообщений
    └── utils/
        ├── db.go               # подключение к PostgreSQL
        ├── auth.go             # JWT (генерация / проверка)
        ├── token_blacklist.go  # чёрный список токенов (logout)
        └── password.go         # bcrypt-хеширование паролей
```

---

## Установка и запуск

### 1. Клонировать репозиторий
```bash
git clone https://github.com/username/bytechat.git
cd bytechat
```

### 2. Настроить базу данных (PostgreSQL)
Создать базу данных `bytechat` и пользователя, например:
```sql
CREATE DATABASE bytechat;
CREATE USER bytechat_user WITH ENCRYPTED PASSWORD 'secret';
GRANT ALL PRIVILEGES ON DATABASE bytechat TO bytechat_user;
```
(При необходимости создать таблицы `users` и `token_blacklist`, а также `messages`, если нужен хранение сообщений).

### 3. Настроить строку подключения (dsn) в `internal/utils/db.go`
```go
dsn := "postgres://bytechat_user:secret@localhost:5432/bytechat?sslmode=disable"
```

### 4. Установить зависимости и запустить
```bash
go mod tidy
go run main.go
```
Сервер запустится на `:8080`.

---

## Использование

### 1. Регистрация (POST /register)
```json
POST http://localhost:8080/register
Content-Type: application/json

{
  "username": "alice",
  "password": "123456"
}
```
Возвращает JSON с полями `id, username`.

### 2. Авторизация (POST /login)
```json
POST http://localhost:8080/login
Content-Type: application/json

{
  "username": "alice",
  "password": "123456"
}
```
Возвращает `token` (JWT-токен).

### 3. Профиль (GET /profile)
Передаём токен:
```
GET http://localhost:8080/profile
Authorization: Bearer <token>
```
Возвращает данные пользователя. Если токен неверный или просрочен — `401 Unauthorized`.

### 4. Логаут (POST /logout)
```
POST http://localhost:8080/logout
Authorization: Bearer <token>
```
Добавляет токен в чёрный список. Любая попытка использовать этот токен вернёт `401`.

### 5. Чат (GET /ws)
Подключение по WebSocket:
```
ws://localhost:8080/ws
```
- Клиент отправляет сообщение — сервер рассылает его всем подключённым.
- (Опционально) передача токена `?token=...` или в заголовке.

---

## Тестирование

1. **Postman**: можно отправлять `POST /register`, `POST /login` и т.д.  
2. **WebSocket** в Postman или [websocket.org/echo.html](https://www.websocket.org/echo.html) для `/ws`.  
3. **`go test ./...`** (если есть юнит-тесты).

---

## Roadmap / Возможные улучшения

1. **Хранение сообщений**: 
   - Добавить таблицу `messages`, сохранять отправленные сообщения.
   - Загружать историю при подключении.  
2. **Комнаты**: 
   - Каждый клиент заходит в свою комнату (roomID).
   - Hub для каждой комнаты.  
3. **Docker** / **Docker-compose**:
   - Запускать Go + PostgreSQL вместе.  
4. **Деплой** в облако:
   - DigitalOcean, AWS, Render — чтобы развернуть приложение онлайн.

---


---

**Удачи в использовании ByteChat!**  
Если есть вопросы — задавайте в [Issues](https://github.com/username/bytechat/issues).

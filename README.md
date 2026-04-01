# miniBank

Простое банковское API на Go с аутентификацией и архитектурой Clean Code.

## 🚀 Быстрый старт

### 1. Настройка окружения

Скопируйте пример конфигурации:
```bash
cp .env.example .env
```

Отредактируйте `.env` файл под ваши настройки:
```bash
# Конфигурация базы данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=bank_db
DB_SSL_MODE=disable

# Конфигурация сервера
SERVER_PORT=7556

# Конфигурация аутентификации
JWT_SECRET=your-super-secret-256-bit-jwt-key
JWT_EXPIRE_HOURS=24

# Конфигурация приложения
APP_ENV=development
LOG_LEVEL=info
```

### 2. Настройка базы данных PostgreSQL

Создайте базу данных:
```sql
CREATE DATABASE bank_db;
```

### 3. Запуск приложения

```bash
go run main.go
```

Приложение автоматически:
- Подключится к базе данных
- Выполнит миграции (создаст таблицы)
- Запустит HTTP сервер

## 🔐 Аутентификация

### Регистрация пользователя
```bash
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "secure123",
  "first_name": "Иван",
  "last_name": "Иванов"
}
```

### Вход в систему
```bash
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "secure123"
}
```

Ответ содержит JWT токен:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "Иван",
    "last_name": "Иванов"
  }
}
```

### Использование токена

Добавляйте токен в заголовок всех защищенных запросов:
```
Authorization: Bearer <ваш-jwt-токен>
```

## 📋 API Endpoints

### Публичные (без авторизации)
- `GET /ping` - проверка работоспособности
- `GET /health` - детальная информация о состоянии сервиса
- `POST /auth/register` - регистрация нового пользователя  
- `POST /auth/login` - вход в систему

### Защищенные (требуют JWT токен)

#### Профиль пользователя
- `GET /profile` - получить профиль текущего пользователя
- `PUT /profile` - обновить профиль

#### Счета
- `GET /accounts` - получить все счета пользователя
- `POST /accounts` - создать новый счет
- `GET /accounts/search?name=<имя>` - поиск счетов по владельцу
- `POST /accounts/:id/top-up` - пополнить счет
- `POST /accounts/:id/withdraw` - снять со счета
- `POST /accounts/transfer` - перевод между счетами
- `GET /accounts/:id/transactions` - история транзакций по счету

#### Транзакции
- `GET /transactions` - получить все транзакции пользователя

## 🏗️ Структура проекта

```
├── main.go                    # Точка входа
├── internal/
│   ├── app/                   # Инициализация приложения
│   ├── config/                # Конфигурация
│   ├── controller/            # HTTP обработчики
│   │   ├── auth.go           # Аутентификация
│   │   ├── account.go        # Счета
│   │   ├── middleware.go     # JWT middleware
│   │   └── router.go         # Роутинг
│   ├── service/               # Бизнес-логика
│   │   ├── auth.go          # Аутентификация
│   │   ├── account.go       # Счета
│   │   └── transaction.go   # Транзакции
│   ├── repository/            # Работа с БД
│   │   ├── user.go          # Пользователи
│   │   ├── account.go       # Счета
│   │   └── transaction.go   # Транзакции
│   ├── models/                # Модели данных
│   │   ├── user.go          # Пользователь
│   │   ├── account.go       # Счет
│   │   └── transaction.go   # Транзакция
│   └── db/                    # Подключение к БД и миграции
├── .env                       # Конфигурация (не в git!)
├── .env.example               # Пример конфигурации
└── .gitignore                 # Исключения для git
```

## 🔧 Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DB_HOST` | Хост базы данных | `localhost` |
| `DB_PORT` | Порт базы данных | `5432` |
| `DB_USER` | Пользователь БД | `postgres` |
| `DB_PASSWORD` | Пароль БД | `postgres` |
| `DB_NAME` | Имя базы данных | `bank_db` |
| `DB_SSL_MODE` | Режим SSL | `disable` |
| `SERVER_PORT` | Порт HTTP сервера | `7556` |
| `JWT_SECRET` | Секретный ключ для JWT | обязательно! |
| `JWT_EXPIRE_HOURS` | Время жизни токена (часы) | `24` |
| `APP_ENV` | Окружение | `development` |
| `LOG_LEVEL` | Уровень логирования | `info` |

## 🎨 Примеры использования

### 1. Регистрация и создание счета

```bash
# Регистрируемся
curl -X POST http://localhost:7556/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"secure123","first_name":"John","last_name":"Doe"}'

# Входим и получаем токен
curl -X POST http://localhost:7556/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"secure123"}'

# Создаем счет (используя полученный токен)
curl -X POST http://localhost:7556/accounts \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Основной счет"}'
```

### 2. Пополнение счета

```bash
curl -X POST http://localhost:7556/accounts/1/top-up \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"amount":50000}'
```

### 3. Перевод между счетами

```bash
curl -X POST http://localhost:7556/accounts/transfer \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"from_id":1,"to_id":2,"amount":10000}'
```

## 🌟 Новые возможности

✅ **Полноценная система аутентификации** - JWT токены, bcrypt хэширование паролей  
✅ **Регистрация и авторизация пользователей** - безопасный вход в систему  
✅ **Управление профилем пользователя** - обновление информации  
✅ **Привязка счетов к пользователям** - каждый счет принадлежит конкретному пользователю  
✅ **JWT middleware** - автоматическая проверка токенов на защищенных endpoints  
✅ **Настраиваемые токены** - время жизни и секрет из конфигурации  
✅ **Структурированное логирование** - подробные логи с user_id  
✅ **Индексы БД** - оптимизированные запросы  

## 🔐 Безопасность

- ✅ JWT токены для аутентификации
- ✅ bcrypt хэширование паролей
- ✅ Проверка принадлежности счетов пользователям
- ✅ Валидация входных данных
- ✅ Секретные ключи в переменных окружения
- ✅ Автоматическое исключение .env из git

## 🚀 Следующие улучшения

- [ ] Graceful shutdown
- [ ] Rate limiting
- [ ] Двухфакторная аутентификация (2FA)
- [ ] RBAC (Role-Based Access Control)
- [ ] Unit тесты
- [ ] Swagger документация
- [ ] Мониторинг и метрики
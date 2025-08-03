# Pinstack User Service 👥

**Pinstack User Service** — микросервис для управления информацией о пользователях в системе **Pinstack**.

## Основные функции:
- CRUD-операции для пользователей (создание, чтение, обновление, удаление).
- Хранение данных пользователей.
- Взаимодействие с другими микросервисами через gRPC.

## Технологии:
- **Go** — основной язык разработки.
- **gRPC** — для межсервисной коммуникации.
- **Docker** — для контейнеризации.

## CI/CD Pipeline 🚀

### GitHub Actions
Проект использует GitHub Actions для автоматического тестирования при каждом push/PR.

**Этапы CI:**
1. **Unit Tests** — юнит-тесты с покрытием кода
2. **Integration Tests** — интеграционные тесты с полной инфраструктурой 
3. **Auto Cleanup** — автоматическая очистка Docker ресурсов

### Makefile команды 📋

#### Основные команды разработки:
```bash
# Проверка кода и тесты
make fmt                    # Форматирование кода (gofmt)
make lint                   # Проверка кода (go vet)
make test-unit              # Юнит-тесты с покрытием
make test-integration       # Интеграционные тесты (с Docker)
make test-all               # Все тесты: форматирование + линтер + юнит + интеграционные

# CI локально
make ci-local               # Полный CI процесс локально (имитация GitHub Actions)
```

#### Управление инфраструктурой:
```bash
# Настройка репозитория
make setup-system-tests        # Клонирует/обновляет pinstack-system-tests репозиторий

# Запуск инфраструктуры
make start-user-infrastructure  # Поднимает все Docker контейнеры для тестов
make check-services            # Проверяет готовность всех сервисов

# Интеграционные тесты
make test-user-integration     # Запускает только интеграционные тесты
make quick-test               # Быстрый запуск тестов без пересборки контейнеров

# Остановка и очистка
make stop-user-infrastructure  # Останавливает все тестовые контейнеры
make clean-user-infrastructure # Полная очистка (контейнеры + volumes + образы)
make clean                    # Полная очистка проекта + Docker
```

#### Логи и отладка:
```bash
# Просмотр логов сервисов
make logs-user              # Логи User Service
make logs-auth              # Логи Auth Service  
make logs-gateway           # Логи API Gateway
make logs-db                # Логи User Database
make logs-auth-db           # Логи Auth Database

# Экстренная очистка
make clean-docker-force     # Удаляет ВСЕ Docker ресурсы (с подтверждением)
```

### Зависимости для интеграционных тестов 🐳

Для интеграционных тестов автоматически поднимаются контейнеры:
- **user-db-test** — PostgreSQL для User Service
- **user-migrator-test** — миграции User Service  
- **user-service-test** — сам User Service
- **auth-db-test** — PostgreSQL для Auth Service
- **auth-migrator-test** — миграции Auth Service
- **auth-service-test** — Auth Service
- **api-gateway-test** — API Gateway

> 📍 **Требования:** Docker, docker-compose  
> 🚀 **Все сервисы собираются автоматически из Git репозиториев**  
> 🔄 **Репозиторий `pinstack-system-tests` клонируется автоматически при запуске тестов**

### Быстрый старт разработки ⚡

```bash
# 1. Проверить код
make fmt lint

# 2. Запустить юнит-тесты
make test-unit

# 3. Запустить интеграционные тесты
make test-integration

# 4. Или всё сразу
make ci-local

# 5. Очистка после работы
make clean
```

### Особенности 🔧

- **Отключение кеша тестов:** все тесты запускаются с флагом `-count=1`
- **Фокус на User Service:** интеграционные тесты тестируют только User endpoints
- **Автоочистка:** CI автоматически удаляет все Docker ресурсы после себя
- **Параллельность:** в CI юнит и интеграционные тесты запускаются последовательно

> ✅ Сервис готов к использованию.

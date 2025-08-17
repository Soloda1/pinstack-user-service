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
- **Prometheus** — для сбора метрик и мониторинга.
- **Grafana** — для визуализации метрик.
- **Loki** — для централизованного сбора логов.

## Архитектура

Проект построен на основе **гексагональной архитектуры (Hexagonal Architecture)** с четким разделением слоев:

### Структура проекта
```
├── cmd/                    # Точки входа приложения
│   ├── server/             # gRPC сервер
│   └── migrate/            # Миграции БД
├── internal/
│   ├── domain/             # Доменный слой
│   │   ├── models/         # Доменные модели
│   │   └── ports/          # Интерфейсы (порты)
│   │       ├── input/      # Входящие порты (use cases)
│   │       └── output/     # Исходящие порты (репозитории, кэш, метрики)
│   ├── application/        # Слой приложения
│   │   └── service/        # Бизнес-логика и сервисы
│   └── infrastructure/     # Инфраструктурный слой
│       ├── inbound/        # Входящие адаптеры (gRPC, HTTP)
│       └── outbound/       # Исходящие адаптеры (PostgreSQL, Redis, Prometheus)
├── migrations/             # SQL миграции
└── mocks/                 # Моки для тестирования
```

### Принципы архитектуры
- **Dependency Inversion**: Зависимости направлены к доменному слою
- **Clean Architecture**: Четкое разделение ответственности между слоями
- **Port & Adapter Pattern**: Интерфейсы определяются в domain, реализуются в infrastructure
- **Testability**: Легкое модульное тестирование благодаря dependency injection

### Мониторинг и метрики
Сервис включает полную интеграцию с системой мониторинга:
- **Prometheus метрики**: Автоматический сбор метрик gRPC, базы данных, кэша
- **Structured logging**: Интеграция с Loki для централизованного сбора логов
- **Health checks**: Проверки состояния всех компонентов
- **Performance monitoring**: Метрики времени ответа и throughput

## CI/CD Pipeline 🚀

### GitHub Actions
Проект использует GitHub Actions для автоматического тестирования при каждом push/PR.

**Этапы CI:**
1. **Unit Tests** — юнит-тесты с покрытием кода
2. **Integration Tests** — интеграционные тесты с полной инфраструктурой 
3. **Auto Cleanup** — автоматическая очистка Docker ресурсов

### Makefile команды 📋

#### Команды разработки

### Настройка и запуск
```bash
# Создание необходимых сетей
make create-networks

# Запуск легкой среды разработки (только база данных и кэш)
make start-dev-light

# Запуск полной среды разработки (с мониторингом)
make start-dev-full

# Остановка среды разработки
make stop-dev
```

### Мониторинг
```bash
# Запуск полного стека мониторинга (Prometheus, Grafana, Loki)
make start-monitoring

# Остановка мониторинга
make stop-monitoring

# Просмотр логов мониторинга
make logs-monitoring
```

### Доступ к сервисам мониторинга
- **Prometheus**: http://localhost:9090 - метрики и мониторинг
- **Grafana**: http://localhost:3000 - дашборды и визуализация
- **User Service Metrics**: http://localhost:9101/metrics - метрики пользовательского сервиса

### База данных
```bash
# Миграции
make migrate-up
make migrate-down

# Тестирование
make test
make test-cover
```

### Docker
```bash
# Сборка образа
make build

# Сборка и запуск в Docker
make docker-build
make docker-run
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

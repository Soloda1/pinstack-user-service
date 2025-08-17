.PHONY: test test-unit test-integration test-user-integration clean build run docker-build setup-system-tests setup-monitoring start-monitoring start-prometheus-stack start-elk-stack stop-monitoring clean-monitoring check-monitoring-health logs-prometheus logs-grafana logs-loki logs-elasticsearch logs-kibana start-dev-full stop-dev-full clean-dev-full start-dev-light

BINARY_NAME=user-service
DOCKER_IMAGE=pinstack-user-service:latest
GO_VERSION=1.24.2
SYSTEM_TESTS_DIR=../pinstack-system-tests
SYSTEM_TESTS_REPO=https://github.com/Soloda1/pinstack-system-tests.git
MONITORING_DIR=../pinstack-monitoring-service
MONITORING_REPO=https://github.com/Soloda1/pinstack-monitoring-service.git

# Проверка версии Go
check-go-version:
	@echo "🔍 Проверка версии Go..."
	@go version | grep -q "go$(GO_VERSION)" || (echo "❌ Требуется Go $(GO_VERSION)" && exit 1)
	@echo "✅ Go $(GO_VERSION) найден"

# Настройка monitoring репозитория
setup-monitoring:
	@echo "🔄 Проверка monitoring репозитория..."
	@if [ ! -d "$(MONITORING_DIR)" ]; then \
		echo "📥 Клонирование pinstack-monitoring-service..."; \
		git clone $(MONITORING_REPO) $(MONITORING_DIR); \
	else \
		echo "🔄 Обновление pinstack-monitoring-service..."; \
		cd $(MONITORING_DIR) && git pull origin main; \
	fi
	@echo "✅ Monitoring готов"

# Настройка system tests репозитория
setup-system-tests:
	@echo "🔄 Проверка system tests репозитория..."
	@if [ ! -d "$(SYSTEM_TESTS_DIR)" ]; then \
		echo "📥 Клонирование pinstack-system-tests..."; \
		git clone $(SYSTEM_TESTS_REPO) $(SYSTEM_TESTS_DIR); \
	else \
		echo "🔄 Обновление pinstack-system-tests..."; \
		cd $(SYSTEM_TESTS_DIR) && git pull origin main; \
	fi
	@echo "✅ System tests готовы"

# Форматирование и проверки
fmt: check-go-version
	gofmt -s -w .
	go fmt ./...

lint: check-go-version
	go vet ./...
	golangci-lint run

# Юнит тесты
test-unit: check-go-version
	go test -v -count=1 -race -coverprofile=coverage.txt ./...

# Запуск полной инфраструктуры для интеграционных тестов из существующего docker-compose
start-user-infrastructure: setup-system-tests
	@echo "🚀 Запуск полной инфраструктуры для интеграционных тестов..."
	@echo "🔍 Проверка и создание сетей..."
	@docker network create pinstack 2>/dev/null || true
	@docker network create pinstack-test 2>/dev/null || true
	cd $(SYSTEM_TESTS_DIR) && \
	USER_SERVICE_CONTEXT=../pinstack-user-service docker compose -f docker-compose.test.yml up -d \
		user-db-test \
		user-migrator-test \
		user-service-test \
		auth-db-test \
		auth-migrator-test \
		auth-service-test \
		api-gateway-test \
		redis
	@echo "⏳ Ожидание готовности сервисов..."
	@sleep 30

# Проверка готовности сервисов
check-services:
	@echo "🔍 Проверка готовности сервисов..."
	@docker exec pinstack-user-db-test pg_isready -U postgres || (echo "❌ User база данных не готова" && exit 1)
	@docker exec pinstack-auth-db-test pg_isready -U postgres || (echo "❌ Auth база данных не готова" && exit 1)
	@timeout 30 bash -c 'until docker exec pinstack-redis-test redis-cli ping | grep -q PONG; do echo "⏳ Ожидание Redis..."; sleep 2; done' || (echo "❌ Redis не готов" && exit 1)
	@echo "✅ Базы данных и Redis готовы"
	@echo "=== User Service logs ==="
	@docker logs pinstack-user-service-test --tail=10
	@echo "=== Auth Service logs ==="
	@docker logs pinstack-auth-service-test --tail=10
	@echo "=== API Gateway logs ==="
	@docker logs pinstack-api-gateway-test --tail=10
	@echo "=== Redis logs ==="
	@docker logs pinstack-redis-test --tail=5

# Интеграционные тесты только для user service
test-user-integration: check-services
	@echo "🧪 Запуск интеграционных тестов для User Service..."
	cd $(SYSTEM_TESTS_DIR) && \
	go test -v -count=1 -timeout=10m ./internal/scenarios/integration/gateway_user/...

# Остановка всех контейнеров
stop-user-infrastructure:
	@echo "🛑 Остановка всей инфраструктуры..."
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml stop \
		api-gateway-test \
		auth-service-test \
		auth-migrator-test \
		auth-db-test \
		user-service-test \
		user-migrator-test \
		user-db-test \
		redis
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml rm -f \
		api-gateway-test \
		auth-service-test \
		auth-migrator-test \
		auth-db-test \
		user-service-test \
		user-migrator-test \
		user-db-test \
		redis

# Полная очистка (включая volumes)
clean-user-infrastructure:
	@echo "🧹 Полная очистка всей инфраструктуры..."
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml down -v
	@echo "🧹 Очистка Docker контейнеров, образов и volumes..."
	docker container prune -f
	docker image prune -a -f
	docker volume prune -f
	docker network prune -f
	@echo "✅ Полная очистка завершена"

# Полные интеграционные тесты (с очисткой)
test-integration: start-user-infrastructure test-user-integration stop-user-infrastructure

# Все тесты
test-all: fmt lint test-unit test-integration

# Логи сервисов
logs-user:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f user-service-test

logs-auth:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f auth-service-test

logs-gateway:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f api-gateway-test

logs-db:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f user-db-test

logs-auth-db:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f auth-db-test

logs-redis:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f redis

# Redis утилиты для отладки
redis-cli:
	@echo "🔍 Подключение к Redis CLI..."
	docker exec -it pinstack-redis-test redis-cli

redis-info:
	@echo "📊 Информация о Redis..."
	docker exec pinstack-redis-test redis-cli info

redis-keys:
	@echo "🔑 Все ключи в Redis..."
	docker exec pinstack-redis-test redis-cli keys "*"

redis-flush:
	@echo "🧹 Очистка всех данных Redis..."
	@read -p "Очистить все данные Redis? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	docker exec pinstack-redis-test redis-cli flushall
	@echo "✅ Redis очищен"

# Быстрый тест с локальным user-service
quick-test-local: setup-system-tests
	@echo "⚡ Быстрый запуск тестов с локальным user-service..."
	@echo "🔍 Проверка и создание сетей..."
	@docker network create pinstack 2>/dev/null || true
	@docker network create pinstack-test 2>/dev/null || true
	cd $(SYSTEM_TESTS_DIR) && \
	USER_SERVICE_CONTEXT=../pinstack-user-service docker compose -f docker-compose.test.yml up -d \
		user-db-test user-migrator-test user-service-test \
		auth-db-test auth-migrator-test auth-service-test \
		api-gateway-test redis
	@echo "⏳ Ожидание готовности сервисов..."
	@sleep 30
	@timeout 30 bash -c 'until docker exec pinstack-redis-test redis-cli ping | grep -q PONG; do echo "⏳ Ожидание Redis..."; sleep 2; done'
	cd $(SYSTEM_TESTS_DIR) && \
	go test -v -count=1 -timeout=5m ./internal/scenarios/integration/gateway_user/...
	$(MAKE) stop-user-infrastructure

# Очистка
clean: clean-user-infrastructure
	go clean
	rm -f $(BINARY_NAME)
	@echo "🧹 Финальная очистка Docker системы..."
	docker system prune -a -f --volumes
	@echo "✅ Вся очистка завершена"

# Экстренная полная очистка Docker (если что-то пошло не так)
clean-docker-force:
	@echo "🚨 ЭКСТРЕННАЯ ПОЛНАЯ ОЧИСТКА DOCKER..."
	@echo "⚠️  Это удалит ВСЕ Docker контейнеры, образы, volumes и сети!"
	@read -p "Продолжить? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	docker stop $$(docker ps -aq) 2>/dev/null || true
	docker rm $$(docker ps -aq) 2>/dev/null || true
	docker rmi $$(docker images -q) 2>/dev/null || true
	docker volume rm $$(docker volume ls -q) 2>/dev/null || true
	docker network rm $$(docker network ls -q) 2>/dev/null || true
	docker system prune -a -f --volumes
	@echo "💥 Экстренная очистка завершена"

# CI локально (имитация GitHub Actions)
ci-local: test-all
	@echo "🎉 Локальный CI завершен успешно!"

# Быстрый тест (только запуск без пересборки)
quick-test: start-user-infrastructure
	@echo "⚡ Быстрый запуск тестов без пересборки..."
	cd $(SYSTEM_TESTS_DIR) && \
	go test -v -count=1 -timeout=5m ./internal/scenarios/integration/gateway_user/...
	$(MAKE) stop-user-infrastructure

######################
# Monitoring Stack   #
######################

# Запуск полного monitoring stack
start-monitoring: setup-monitoring
	@echo "📊 Запуск monitoring stack..."
	@echo "🔍 Проверка и создание сетей..."
	@docker network create pinstack 2>/dev/null || true
	@docker network create pinstack-test 2>/dev/null || true
	cd $(MONITORING_DIR) && \
	docker compose up -d
	@echo "🔗 Подключение monitoring к тестовым сетям..."
	@docker network connect pinstack-system-tests_pinstack-test pinstack-prometheus 2>/dev/null || true
	@docker network connect pinstack-system-tests_pinstack-test pinstack-grafana 2>/dev/null || true
	@docker network connect pinstack-system-tests_pinstack-test pinstack-loki 2>/dev/null || true
	@echo "⏳ Ожидание готовности monitoring сервисов..."
	@sleep 15
	@echo "✅ Monitoring stack запущен:"
	@echo "  📊 Prometheus: http://localhost:9090"
	@echo "  📈 Grafana: http://localhost:3000 (admin/admin)"
	@echo "  🔍 Loki: http://localhost:3100"
	@echo "  📋 Kibana: http://localhost:5601"
	@echo "  💾 Elasticsearch: http://localhost:9200"
	@echo "  🐧 PgAdmin: http://localhost:5050 (admin@admin.com/admin)"
	@echo "  🐛 Kafka UI: http://localhost:9091"

# Запуск только Prometheus stack (Prometheus + Grafana + Loki)
start-prometheus-stack: setup-monitoring
	@echo "📊 Запуск Prometheus stack..."
	@echo "🔍 Проверка и создание сетей..."
	@docker network create pinstack 2>/dev/null || true
	@docker network create pinstack-test 2>/dev/null || true
	cd $(MONITORING_DIR) && \
	docker compose up -d prometheus grafana loki promtail
	@echo "🔗 Подключение Prometheus stack к тестовым сетям..."
	@docker network connect pinstack-system-tests_pinstack-test pinstack-prometheus 2>/dev/null || true
	@docker network connect pinstack-system-tests_pinstack-test pinstack-grafana 2>/dev/null || true
	@docker network connect pinstack-system-tests_pinstack-test pinstack-loki 2>/dev/null || true
	@echo "⏳ Ожидание готовности Prometheus stack..."
	@sleep 10
	@echo "✅ Prometheus stack запущен:"
	@echo "  📊 Prometheus: http://localhost:9090"
	@echo "  📈 Grafana: http://localhost:3000 (admin/admin)"
	@echo "  🔍 Loki: http://localhost:3100"

# Запуск только ELK stack
start-elk-stack: setup-monitoring
	@echo "📊 Запуск ELK stack..."
	@echo "🔍 Проверка и создание сетей..."
	@docker network create pinstack 2>/dev/null || true
	@docker network create pinstack-test 2>/dev/null || true
	cd $(MONITORING_DIR) && \
	docker compose up -d elasticsearch logstash kibana filebeat
	@echo "⏳ Ожидание готовности ELK stack..."
	@sleep 30
	@echo "✅ ELK stack запущен:"
	@echo "  📋 Kibana: http://localhost:5601"
	@echo "  💾 Elasticsearch: http://localhost:9200"

# Проверка состояния monitoring сервисов
check-monitoring-health:
	@echo "🔍 Проверка состояния monitoring сервисов..."
	@echo "Prometheus:" && curl -s http://localhost:9090/-/healthy | head -1 || echo "❌ Prometheus недоступен"
	@echo "Grafana:" && curl -s http://localhost:3000/api/health | head -1 || echo "❌ Grafana недоступна"
	@echo "Loki:" && curl -s http://localhost:3100/ready | head -1 || echo "❌ Loki недоступен"
	@echo "Elasticsearch:" && curl -s http://localhost:9200/_cluster/health | head -1 || echo "❌ Elasticsearch недоступен"
	@echo "Kibana:" && curl -s http://localhost:5601/api/status | head -1 || echo "❌ Kibana недоступна"

# Остановка monitoring stack
stop-monitoring:
	@echo "🛑 Остановка monitoring stack..."
	@if [ -d "$(MONITORING_DIR)" ]; then \
		cd $(MONITORING_DIR) && docker compose stop; \
	else \
		echo "⚠️  Monitoring директория не найдена"; \
	fi

# Полная очистка monitoring stack
clean-monitoring:
	@echo "🧹 Очистка monitoring stack..."
	@if [ -d "$(MONITORING_DIR)" ]; then \
		cd $(MONITORING_DIR) && docker compose down -v; \
		echo "🧹 Очистка monitoring volumes..."; \
		docker volume rm pinstack-monitoring-service_elasticsearch_data 2>/dev/null || true; \
		docker volume rm pinstack-monitoring-service_filebeat_data 2>/dev/null || true; \
		echo "✅ Monitoring stack очищен"; \
	else \
		echo "⚠️  Monitoring директория не найдена"; \
	fi

# Логи monitoring сервисов
logs-prometheus:
	@if [ -d "$(MONITORING_DIR)" ]; then \
		cd $(MONITORING_DIR) && docker compose logs -f prometheus; \
	else \
		echo "⚠️  Monitoring директория не найдена"; \
	fi

logs-grafana:
	@if [ -d "$(MONITORING_DIR)" ]; then \
		cd $(MONITORING_DIR) && docker compose logs -f grafana; \
	else \
		echo "⚠️  Monitoring директория не найдена"; \
	fi

logs-loki:
	@if [ -d "$(MONITORING_DIR)" ]; then \
		cd $(MONITORING_DIR) && docker compose logs -f loki; \
	else \
		echo "⚠️  Monitoring директория не найдена"; \
	fi

logs-elasticsearch:
	@if [ -d "$(MONITORING_DIR)" ]; then \
		cd $(MONITORING_DIR) && docker compose logs -f elasticsearch; \
	else \
		echo "⚠️  Monitoring директория не найдена"; \
	fi

logs-kibana:
	@if [ -d "$(MONITORING_DIR)" ]; then \
		cd $(MONITORING_DIR) && docker compose logs -f kibana; \
	else \
		echo "⚠️  Monitoring директория не найдена"; \
	fi

# Комбинированные команды

# Полный development environment с мониторингом
start-dev-full: setup-monitoring start-monitoring start-user-infrastructure
	@echo "� Дополнительное подключение monitoring к тестовым сетям..."
	@docker network connect pinstack-system-tests_pinstack-test pinstack-prometheus 2>/dev/null || true
	@docker network connect pinstack-system-tests_pinstack-test pinstack-grafana 2>/dev/null || true
	@docker network connect pinstack-system-tests_pinstack-test pinstack-loki 2>/dev/null || true
	@echo "�🚀 Полная dev среда запущена!"
	@echo ""
	@echo "=== Приложения ==="
	@echo "  🔗 API Gateway: http://localhost:8080"
	@echo "  👤 User Service: http://localhost:8081"
	@echo "  🔐 Auth Service: http://localhost:8082"
	@echo ""
	@echo "=== Мониторинг ==="
	@echo "  📊 Prometheus: http://localhost:9090"
	@echo "  📈 Grafana: http://localhost:3000 (admin/admin)"
	@echo "  🔍 Loki: http://localhost:3100"
	@echo "  📋 Kibana: http://localhost:5601"
	@echo ""
	@echo "=== Базы данных ==="
	@echo "  🐧 PgAdmin: http://localhost:5050 (admin@admin.com/admin)"
	@echo "  🔴 Redis: localhost:6379"
	@echo "  🐛 Kafka UI: http://localhost:9091"

# Остановка всей dev среды
stop-dev-full: stop-monitoring stop-user-infrastructure
	@echo "🛑 Полная dev среда остановлена"

# Очистка всей dev среды
clean-dev-full: clean-monitoring clean-user-infrastructure
	@echo "🧹 Полная dev среда очищена"

# Запуск только с Prometheus stack (без ELK)
start-dev-light: setup-monitoring start-prometheus-stack start-user-infrastructure
	@echo "� Дополнительное подключение Prometheus stack к тестовым сетям..."
	@docker network connect pinstack-system-tests_pinstack-test pinstack-prometheus 2>/dev/null || true
	@docker network connect pinstack-system-tests_pinstack-test pinstack-grafana 2>/dev/null || true
	@docker network connect pinstack-system-tests_pinstack-test pinstack-loki 2>/dev/null || true
	@echo "�🚀 Легкая dev среда запущена (без ELK stack)!"
	@echo ""
	@echo "=== Приложения ==="
	@echo "  🔗 API Gateway"
	@echo "  👤 User Service"
	@echo "  🔐 Auth Service"
	@echo ""
	@echo "=== Мониторинг ==="
	@echo "  📊 Prometheus: http://localhost:9090"
	@echo "  📈 Grafana: http://localhost:3000 (admin/admin)"
	@echo "  🔍 Loki: http://localhost:3100"
### Установка и запуск

### Запуск локально
1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/Dobi-Vanish/L0
2. Перейдите в корень проекта и запустите через makefile поочерёдно команды:
   ```bash
   cd L0
   make build
   make docker-build
   make docker-up
   make migrate-up
3. Перезапустите контейнеры для применения всех настроек:
   ```bash
   make docker-down
   make docker-up

После успешного запуска доступны следующие endpoints:
 - Prometheus Metrics: `http://localhost:9090`
 - Статистика кэша: `http://localhost:8081/cache/stats`
 - Просмотр заказа по ID: `http://localhost:8081/`
 - Добавление заказа (через Postman): `http://localhost:8081/add_order`
Для подключения и просмотра логов через MongoDB подключиться через URI: `mongodb://localhost:27017`.

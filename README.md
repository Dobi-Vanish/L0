ПРИМЕЧАНИЕ: Изначальное видео оказывается записалось без смены окон, после того как перезаписал - отредактировать ответ уже было нельзя. Вот ссылка на подходящее видео - [https://youtu.be/9wddzxrZc4E?si=UcpKfekVMgTzUCU-](https://youtu.be/9wddzxrZc4E)

### Запуск проекта

1. Клонируйте репозиторий:
```bash
git clone [<repository-url>](https://github.com/Dobi-Vanish/L0)


Перейдите в корнь проекта. Далее поочерёдно введите команды:
```bash
make build
make docker-build
make docker-up
make migrate-up


Перезапустите контейнеры для применения всех настроек:
```bash
make docker-down
make docker-up

После успешного запуска доступны следующие endpoints:

    Prometheus Metrics: `http://localhost:9090`

    Статистика кэша: `http://localhost:8081/cache/stats`

    Просмотр заказа по ID: `http://localhost:8081/`

    Добавление заказа (через Postman): `http://localhost:8081/add_order`

Для подключения и просмотра логов через MongoDB подключиться через URI: `mongodb://localhost:27017`.

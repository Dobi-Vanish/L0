ПРИМЕЧАНИЕ: Изначальное видео оказывается записалось без смены окон, после того как перезаписал - отредактировать ответ уже было нельзя. Вот ссылка на подходящее видео - [https://youtu.be/9wddzxrZc4E?si=UcpKfekVMgTzUCU-](https://youtu.be/9wddzxrZc4E)

Для запуска перейдите в корень проекта и пропишите команды:
make build
make docker-build
make docker-up
make migrate-up

К сожалению, надо будет после этого завершить и заново запустить образы в Docker'e чтобы БД подтянулись:
make docker-down
make docker-up

После этого сервис будет полностью готов к работе.
http://localhost:9090 - Prometheus.
http://localhost:8081/cache/stats - статистика кэша.
http://localhost:8081/ - просмотр заказа по ID.
http://localhost:8081/add_order - через Postman добавить экземпляры в БД в PostgreSQL.
Для подключения и просмотра логов через MongoDB подключиться через URI: mongodb://localhost:27017 .

# User Service

Запуск локально через Docker Compose:

```bash
docker compose up --build
```

Сервисы:
- `postgres` — база на порту `5432`.
- `user-service` — gRPC на `:50051`, читает конфиг `/root/config.yaml`, миграции включены.

Конфигурация: `source/user-service/config.yaml` (монтируется в контейнер). При необходимости поменяйте строки подключения/секреты/порты.

Protobuf: `source/api/user-service/user/v1/user.proto`, сгенерированный код в `gen/go/user/v1`.

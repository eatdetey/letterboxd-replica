# API Gateway

HTTP шлюз на Fiber, проксирует запросы к `user-service` (gRPC).

## Запуск локально

1. Убедитесь, что `user-service` запущен и доступен по адресу из конфига.
2. Настройте конфиг `source/configs/api-gateway/config.yaml`:
   ```yaml
   http:
     addr: ":8080"
   user_service:
     address: "user-service:50051"
   ```
3. Запустите:
   ```bash
   cd source/api-gateway
   go run ./cmd/server
   ```

### Генерация gRPC клиента
```bash
cd source/api-gateway
make proto
```

## Маршруты
- `POST /api/v1/auth/register` — body: `{username, password, email}`
- `POST /api/v1/auth/login` — body: `{username, password}`
- `POST /api/v1/auth/refresh` — body: `{refresh_token}`
- `GET  /api/v1/users?ids=1,2&usernames=a,b&limit=10&offset=0`

Ответы содержат данные пользователя и/или пару токенов, ошибки пробрасываются с HTTP кодами, соответствующими gRPC статусам.

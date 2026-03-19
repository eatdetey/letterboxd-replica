# API Gateway

HTTP шлюз на Fiber, проксирует запросы к `user-service` и `movie-service` (gRPC).

## Запуск локально

1. Убедитесь, что `user-service` и `movie-service` запущены и доступны по адресам из конфига.
2. Настройте конфиг `source/configs/api-gateway/config.yaml`:
   ```yaml
   http:
     addr: ":8080"
   user_service:
     address: "user-service:50051"
   movie_service:
     address: "movie-service:50051"
   ```
3. Запустите:
   ```bash
   cd source/api-gateway
   go run ./cmd/server
   ```

### Генерация gRPC клиента
```bash
cd source
make generate
```

## Маршруты
- `POST /api/v1/auth/register` — body: `{username, password, email}`
- `POST /api/v1/auth/login` — body: `{username, password}`
- `POST /api/v1/auth/refresh` — body: `{refresh_token}`
- `GET  /api/v1/users?ids=1,2&usernames=a,b&limit=10&offset=0`
- `GET  /api/v1/movies?limit=20&offset=0&search=&genre=&playlist_id=&enrich_playlists=`
- `GET  /api/v1/playlists`
- `POST /api/v1/playlists` — body: `{name}`
- `PUT  /api/v1/playlists/{id}` — body: `{name}`
- `DELETE /api/v1/playlists/{id}`
- `POST /api/v1/playlists/{id}/movies` — body: `{movie_id}`
- `DELETE /api/v1/playlists/{id}/movies/{movie_id}`

Для `GET /api/v1/movies` заголовок `Authorization` прокидывается в `movie-service` через gRPC metadata (`authorization`).

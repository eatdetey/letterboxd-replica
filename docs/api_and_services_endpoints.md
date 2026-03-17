

# Техническая спецификация API и микросервисов (Final v2.0)

## 1. Архитектурный обзор

*   **Внешний интерфейс (Client ↔ Gateway):** REST API over HTTP. Формат: JSON.
    *   Клиент видит разделение: `GET /movies/{id}` и `GET /movies`.
*   **Внутренний интерфейс (Gateway ↔ Services):** gRPC over HTTP/2. Формат: Protobuf.
    *   **Movie Service** предоставляет единый универсальный метод `GetMovies`.
*   **Аутентификация:** JWT Bearer Token в заголовке `Authorization`.
*   **Роли:** `user`, `admin`.

---

## 2. Контракты внешнего API (HTTP/JSON)

Эти эндпоинты предоставляет **API Gateway**. Gateway транслирует разные HTTP запросы в единый gRPC вызов с разными параметрами.

### 2.1. User Service (Auth & Profile)

#### Регистрация
`POST /api/v1/auth/register`
```json
// Request
{
  "username": "string",
  "password": "string",
  "email": "string"
}

// Response 201 Created
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "role": "user",
  "tokens": {
    "access": "string",
    "refresh": "string"
  }
}
```

#### Вход (Login)
`POST /api/v1/auth/login`
```json
// Request
{
  "username": "string",
  "password": "string"
}

// Response 200 OK
{
  "id": "uuid",
  "username": "string",
  "role": "user",
  "tokens": {
    "access": "string",
    "refresh": "string"
  }
}
```

#### Обновление токена
`POST /api/v1/auth/refresh`
```json
// Request
{
  "refresh_token": "string"
}

// Response 200 OK
{
  "access": "string",
  "refresh": "string"
}
```

---

### 2.2. Movie Service (Каталог и Поиск)

> **Ключевая особенность:** Все данные о фильмах приходят из одного внутреннего источника.
> *   Если запрошен **один фильм** (`/movies/{id}`): Gateway передает его ID в списке `ids` и ставит флаг `enrich_playlists=true`.
> *   Если запрошен **список/поиск** (`/movies`): Gateway передает фильтры (поиск, жанр, плейлист) и опционально флаг `enrich_playlists=false` (по умолчанию), либо `true`, если клиент явно запросил детализацию для списка (редкий кейс, обычно для списка enrichment не нужен для производительности, но для единичного объекта — обязателен).
> *   **Фильтрация по плейлисту:** Если передан `playlist_id`, сервис возвращает **только** фильмы, входящие в этот плейлист текущего пользователя.

#### Получить список фильмов (с поиском и фильтрами)
`GET /api/v1/movies`
Query params:
*   `limit` (int, default 20)
*   `offset` (int, default 0)
*   `search` (string, optional) — поиск по названию (partial match).
*   `genre` (string, optional)
*   `playlist_id` (uuid, optional) — жесткий фильтр: вернуть только фильмы из этого плейлиста.

```json
// Request Example: GET /api/v1/movies?search=inception&limit=10
// Headers: Authorization: Bearer <token>

// Response 200 OK
{
  "items": [
    {
      "id": "mov_001",
      "title": "Inception",
      "description": "A thief who steals corporate secrets...",
      "release_year": 2010,
      "genres": ["Sci-Fi", "Action"],
      "poster_url": "string"
      // Поле playlists отсутствует, так как это список и флаг enrich был false
    },
    {
      "id": "mov_002",
      "title": "Inception: The IMAX Experience",
      "description": "...",
      "release_year": 2010,
      "genres": ["Sci-Fi", "Documentary"],
      "poster_url": "string"
    }
  ],
  "total": 2
}
```

#### Получить один фильм (Детально)
`GET /api/v1/movies/{id}`
*Под капотом Gateway вызывает тот же метод, передавая `ids=["{id}"]` и `enrich_playlists=true`.*

```json
// Response 200 OK
{
  "id": "mov_001",
  "title": "Inception",
  "description": "A thief who steals corporate secrets...",
  "release_year": 2010,
  "genres": ["Sci-Fi", "Action"],
  "poster_url": "string",
  "playlists": [
    {
      "id": "pl_123",
      "name": "My Favorites"
    },
    {
      "id": "pl_456",
      "name": "Watch Later"
    }
  ]
}
```

#### CRUD операции (Только Admin)
*   `POST /api/v1/movies` — Создание.
*   `PUT /api/v1/movies/{id}` — Обновление.
*   `DELETE /api/v1/movies/{id}` — Удаление.

```json
// Request Body для Create/Update
{
  "title": "string",
  "description": "string",
  "release_year": 2024,
  "genres": ["Drama"],
  "poster_url": "string"
}
```

---

### 2.3. Movie List Service (Плейлисты)

#### Управление плейлистами
*   `GET /api/v1/playlists` — Список плейлистов юзера.
*   `POST /api/v1/playlists` — Создать.
*   `PUT /api/v1/playlists/{id}` — Переименовать.
*   `DELETE /api/v1/playlists/{id}` — Удалить.

```json
// Response GET /api/v1/playlists
{
  "items": [
    {
      "id": "pl_123",
      "name": "My Favorites",
      "movies_count": 5
    }
  ]
}
```

#### Управление содержимым
*   `POST /api/v1/playlists/{id}/movies` — Добавить фильм.
    *   Body: `{ "movie_id": "mov_001" }`
*   `DELETE /api/v1/playlists/{id}/movies/{movie_id}` — Удалить фильм.

---

### 2.4. Review Service (Отзывы)

#### Получить отзывы
`GET /api/v1/movies/{movie_id}/reviews`
```json
// Response 200 OK
{
  "items": [
    {
      "id": "rev_001",
      "user_id": "usr_555",
      "username": "john_doe",
      "text": "Great movie!",
      "created_at": "2023-10-10T12:00:00Z"
    }
  ]
}
```

#### Оставить отзыв
`POST /api/v1/movies/{movie_id}/reviews`
```json
// Request
{
  "text": "Amazing plot twist."
}
```

---

## 3. Внутренние контракты (gRPC / Protobuf)

Здесь реализована ключевая логика унификации.

### 3.1. User Service (`user.proto`)

```protobuf
syntax = "proto3";
package user;

service UserService {
  rpc ValidateToken (TokenRequest) returns (UserResponse);
  rpc GetUsersByIds (UserIdsRequest) returns (UsersResponse);
  rpc UserExists (UserIdRequest) returns (ExistsResponse);
}

message TokenRequest { string token = 1; }
message UserIdRequest { string id = 1; }
message UserIdsRequest { repeated string ids = 1; }

message UserResponse {
  string id = 1;
  string username = 2;
  string email = 3;
  string role = 4;
  bool is_valid = 5;
}

message UsersResponse { repeated UserResponse users = 1; }
message ExistsResponse { bool exists = 1; }
```

### 3.2. Movie List Service (`movielist.proto`)

```protobuf
syntax = "proto3";
package movielist;

service MovieListService {
  // Фильтрация списка кандидатов по принадлежности к плейлисту
  rpc FilterMoviesByPlaylist (FilterRequest) returns (MovieIdsResponse);

  // Получение списка плейлистов для конкретного фильма
  rpc GetPlaylistsForMovie (MoviePlaylistRequest) returns (PlaylistsResponse);
  
  // CRUD для плейлистов
  rpc CreatePlaylist (CreatePlaylistRequest) returns (PlaylistResponse);
  rpc AddMovieToPlaylist (AddMovieRequest) returns (Empty);
  rpc RemoveMovieFromPlaylist (RemoveMovieRequest) returns (Empty);
}

message FilterRequest {
  string user_id = 1;
  string playlist_id = 2;
  repeated string candidate_movie_ids = 3; 
}

message MovieIdsResponse {
  repeated string movie_ids = 1; 
}

message MoviePlaylistRequest {
  string user_id = 1;
  string movie_id = 2;
}

message PlaylistInfo {
  string id = 1;
  string name = 2;
}

message PlaylistsResponse {
  repeated PlaylistInfo playlists = 1;
}

message CreatePlaylistRequest {
  string user_id = 1;
  string name = 2;
}

message PlaylistResponse {
  string id = 1;
  string name = 2;
}

message AddMovieRequest {
  string playlist_id = 1;
  string movie_id = 2;
}

message RemoveMovieRequest {
  string playlist_id = 1;
  string movie_id = 2;
}

message Empty {}
```

### 3.3. Movie Service (`movie.proto`)

```protobuf
syntax = "proto3";
package movie;

service MovieService {
  // ЕДИНСТВЕННЫЙ метод для получения фильмов (Batch)
  rpc GetMovies (GetMoviesRequest) returns (MoviesResponse);
  
  // CRUD
  rpc CreateMovie (CreateMovieRequest) returns (MovieResponse);
  rpc UpdateMovie (UpdateMovieRequest) returns (MovieResponse);
  rpc DeleteMovie (DeleteMovieRequest) returns (Empty);
  
  // Валидация существования (для Review Service)
  rpc MovieExists (MovieIdRequest) returns (ExistsResponse);
}

message GetMoviesRequest {
  // Пагинация
  int32 limit = 1;
  int32 offset = 2;
  
  // Фильтры
  optional string search_query = 3;   // Поиск по названию
  optional string genre = 4;          // Фильтр по жанру
  repeated string ids = 5;            // Если заполнено, игнорируем search/genre, берем только эти ID (для получения по одному или нескольким ID)
  optional string playlist_id = 6;    // Если есть, применяем жесткую фильтрацию через MovieListService
  
  // Контекст
  string user_id = 7;                 // Обязателен для работы с playlist_id и enrich_playlists
  
  // Флаги поведения
  bool enrich_playlists = 8;          // Если true, для каждого возвращаемого фильма делаем запрос в MovieListService для получения списка плейлистов
}

message MoviesResponse {
  repeated MovieDetailedResponse items = 1;
  int32 total = 2;
}

message MovieDetailedResponse {
  string id = 1;
  string title = 2;
  string description = 3;
  int32 release_year = 4;
  repeated string genres = 5;
  
  // Заполняется только если enrich_playlists = true
  repeated PlaylistInfo playlists = 6; 
}

message PlaylistInfo {
  string id = 1;
  string name = 2;
}

// CRUD сообщения
message CreateMovieRequest {
  string title = 1;
  string description = 2;
  int32 release_year = 3;
  repeated string genres = 4;
}

message UpdateMovieRequest {
  string id = 1;
  optional string title = 2;
  optional string description = 3;
  optional int32 release_year = 4;
  optional repeated string genres = 5;
}

message DeleteMovieRequest { string id = 1; }
message MovieResponse {
  string id = 1;
  string title = 2;
  string description = 3;
  int32 release_year = 4;
  repeated string genres = 5;
}
message MovieIdRequest { string id = 1; }
message ExistsResponse { bool exists = 1; }
message Empty {}
```

### 3.4. Review Service (`review.proto`)

```protobuf
syntax = "proto3";
package review;

service ReviewService {
  rpc GetReviews (GetReviewsRequest) returns (ReviewsResponse);
  rpc AddReview (AddReviewRequest) returns (ReviewResponse);
}

message GetReviewsRequest { string movie_id = 1; }

message AddReviewRequest {
  string movie_id = 1;
  string user_id = 2;
  string text = 3;
}

message ReviewResponse {
  string id = 1;
  string movie_id = 2;
  string user_id = 3;
  string username = 4;
  string text = 5;
  string created_at = 6;
}

message ReviewsResponse { repeated ReviewResponse items = 1; }
```

---

## 4. Детальные сценарии взаимодействия (Flow)

### Сценарий А: Поиск фильмов по названию
1.  **Client**: `GET /api/v1/movies?search=matrix`
2.  **Gateway**:
    *   Валидирует токен -> получает `user_id`.
    *   Вызывает `MovieService.GetMovies`:
        *   `search_query = "matrix"`
        *   `ids = []` (пусто)
        *   `playlist_id = null`
        *   `enrich_playlists = false` (для списков обычно не нужно грузить лишнее)
3.  **MovieService**:
    *   Делает SQL запрос `SELECT * FROM movies WHERE title ILIKE '%matrix%'`.
    *   Возвращает список базовых объектов.
4.  **Gateway** -> JSON клиенту.

### Сценарий Б: Получение одного фильма по ID (с плейлистами)
1.  **Client**: `GET /api/v1/movies/mov_123`
2.  **Gateway**:
    *   Валидирует токен -> `user_id`.
    *   Вызывает `MovieService.GetMovies`:
        *   `ids = ["mov_123"]`
        *   `search_query = null`
        *   `playlist_id = null`
        *   `enrich_playlists = true` **(Критично)**
3.  **MovieService**:
    *   Видит, что `ids` не пуст. Делает выборку `SELECT * FROM movies WHERE id IN ('mov_123')`.
    *   Видит флаг `enrich_playlists = true`.
    *   **Вызов в MovieListService**: `GetPlaylistsForMovie(user_id, "mov_123")`.
    *   Получает список плейлистов `[ {id: "pl_1", name: "Fav"} ]`.
    *   Инъектирует данные в ответ.
4.  **Gateway** -> JSON клиенту с полем `playlists`.

### Сценарий В: Получение фильмов из конкретного плейлиста
1.  **Client**: `GET /api/v1/movies?playlist_id=pl_999&search=action` (Поиск внутри плейлиста)
2.  **Gateway**:
    *   Валидирует токен -> `user_id`.
    *   Вызывает `MovieService.GetMovies`:
        *   `playlist_id = "pl_999"`
        *   `search_query = "action"`
        *   `enrich_playlists = false`
3.  **MovieService**:
    *   **Шаг 1 (Кандидаты):** Делает поиск по своей базе с учетом поиска: `SELECT id FROM movies WHERE title ILIKE '%action%'`. Получает список кандидатов, например `[ID_A, ID_B, ID_C]`.
    *   **Шаг 2 (Фильтрация):** Вызывает `MovieListService.FilterMoviesByPlaylist(user_id, "pl_999", [ID_A, ID_B, ID_C])`.
    *   **MovieListService** проверяет таблицу связей и возвращает только те ID, которые есть в плейлисте, например `[ID_A, ID_C]`.
    *   **Шаг 3 (Финал):** MovieService делает `SELECT * FROM movies WHERE id IN (ID_A, ID_C)` и формирует ответ.
4.  **Gateway** -> JSON клиенту (только фильмы из плейлиста, соответствующие поиску).

### Сценарий Г: Добавление отзыва
1.  **Client**: `POST /api/v1/movies/mov_123/reviews`
2.  **Gateway**:
    *   Валидирует токен -> `user_id`.
    *   Вызывает `ReviewService.AddReview(movie_id="mov_123", user_id="usr_X", text="...")`.
3.  **ReviewService**:
    *   `UserService.UserExists("usr_X")` -> OK.
    *   `MovieService.MovieExists("mov_123")` -> OK.
    *   Сохраняет отзыв.
    *   `UserService.GetUsersByIds(["usr_X"])` -> получает username.
    *   Возвращает готовый объект отзыва.
4.  **Gateway** -> JSON клиенту.

---

## 5. Структура Баз Данных (Логическая модель)

*   **User Service DB**:
    *   `users`: `id (PK), username, password_hash, email, role, created_at`

*   **Movie Service DB**:
    *   `movies`: `id (PK), title, description, release_year, genres (array/jsonb)`
    *   *Индексы:* `title` (для поиска), `genres`.

*   **Movie List Service DB**:
    *   `playlists`: `id (PK), user_id (FK), name, created_at`
    *   `playlist_movies`: `playlist_id (FK), movie_id (FK), added_at`
    *   *Индексы:* `(playlist_id, movie_id)` для быстрой проверки наличия.

*   **Review Service DB**:
    *   `reviews`: `id (PK), movie_id, user_id, text, created_at`
    *   *Индексы:* `movie_id` (для получения списка отзывов).

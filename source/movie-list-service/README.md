# Movie List Service

Movie List Service is a gRPC microservice that manages user playlists (movie lists) for the Letterboxd replica application.

## Features

- **Playlist Management**: Create, rename, delete playlists
- **Movie Management**: Add/remove movies to/from playlists
- **Playlist Queries**: Get all playlists for a user, get movies in a playlist
- **Filtering**: Filter movies by playlist membership (used by Movie Service)
- **Enrichment**: Get all playlists containing a specific movie

## Architecture

The service follows a clean architecture pattern:

```
cmd/server/          # Application entry point
internal/
  app/               # Application initialization
  config/            # Configuration loading
  handler/           # gRPC request handlers
  repository/        # Data access layer
  service/           # Business logic layer
  transport/         # External adapters (postgres, migrations)
gen/go/              # Generated protobuf code
migrations/          # Database migrations
```

## gRPC Methods

### Playlist CRUD
- `CreatePlaylist` - Create a new playlist
- `RenamePlaylist` - Rename an existing playlist
- `DeletePlaylist` - Delete a playlist
- `GetPlaylistsForUser` - Get all playlists for a user

### Movie Management
- `AddMovieToPlaylist` - Add a movie to a playlist
- `RemoveMovieFromPlaylist` - Remove a movie from a playlist
- `GetPlaylistMovies` - Get all movies in a playlist

### Integration Methods
- `FilterMoviesByPlaylist` - Filter candidate movies by playlist membership
- `GetPlaylistsForMovie` - Get all playlists containing a specific movie

## Database Schema

### playlists
| Column | Type | Description |
|--------|------|-------------|
| id | TEXT (UUID) | Primary key |
| user_id | BIGINT | Owner user ID |
| name | TEXT | Playlist name |
| created_at | TIMESTAMPTZ | Creation timestamp |
| updated_at | TIMESTAMPTZ | Last update timestamp |

### playlist_movies
| Column | Type | Description |
|--------|------|-------------|
| playlist_id | TEXT (FK) | Reference to playlists |
| movie_id | TEXT | Movie ID (reference to movie service) |
| added_at | TIMESTAMPTZ | When the movie was added |

## Configuration

Configuration is loaded from YAML file and environment variables:

```yaml
grpc_server:
  port: ":50052"
  max_connection_idle: 300
  keepalive_time: 7200
  keepalive_timeout: 20

migrate:
  need_to_migrate: true

shutdown:
  shutdown_timeout: 5
```

Environment variables (override YAML):
- `DB_CONNECTION_STRING` - PostgreSQL connection string
- `CONFIG_PATH` - Path to config file

## Running

### Local Development

```bash
# Run migrations and start service
go run ./cmd/server/main.go
```

### Docker

```bash
docker-compose up movie-list-service
```

## Development

### Generate Protobuf Code

```bash
cd source
buf generate api/movie-list-service
```

### Run Tests

```bash
go test ./...
```

## API Contracts

See the proto definition: `api/movie-list-service/movielist/v1/movielist.proto`

## Integration with Other Services

### Movie Service
The Movie Service calls `FilterMoviesByPlaylist` when a `playlist_id` filter is provided to filter movies by playlist membership. It also calls `GetPlaylistsForMovie` when `enrich_playlists=true` to add playlist information to movie responses.

### API Gateway
The API Gateway exposes REST endpoints that translate to gRPC calls:
- `GET /api/v1/playlists` → `GetPlaylistsForUser`
- `POST /api/v1/playlists` → `CreatePlaylist`
- `PUT /api/v1/playlists/{id}` → `RenamePlaylist`
- `DELETE /api/v1/playlists/{id}` → `DeletePlaylist`
- `POST /api/v1/playlists/{id}/movies` → `AddMovieToPlaylist`
- `DELETE /api/v1/playlists/{id}/movies/{movie_id}` → `RemoveMovieFromPlaylist`

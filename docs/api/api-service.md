## API → Service Mapping

| Endpoint | Service |
|--------|--------|
| POST /api/v1/auth/register | Auth Service |
| POST /api/v1/auth/login | Auth Service |
| POST /api/v1/auth/logout | Auth Service |
| POST /api/v1/auth/refresh | Auth Service |
| PUT /api/v1/users/me/password | Auth Service |

| GET /api/v1/users/me | User Service |
| PUT /api/v1/users/me | User Service |
| GET /api/v1/users/{userId} | User Service |
| POST /api/v1/users/{userId}/follow | User Service |
| DELETE /api/v1/users/{userId}/follow | User Service |
| GET /api/v1/users/{userId}/followers | User Service |
| GET /api/v1/users/{userId}/following | User Service |
| POST /api/v1/users/me/avatar | User Service |
| DELETE /api/v1/users/me/avatar | User Service |
| GET /api/v1/users/me/settings | User Service |
| PUT /api/v1/users/me/settings | User Service |

| GET /api/v1/films/popular | Movie Service |
| GET /api/v1/films/search | Movie Service |
| GET /api/v1/films/{filmId} | Movie Service |
| GET /api/v1/films/{filmId}/cast | Movie Service |
| GET /api/v1/films/{filmId}/trailers | Movie Service |
| GET /api/v1/films/{filmId}/similar | Movie Service |

| PUT /api/v1/films/{filmId}/rating | Review Service |
| GET /api/v1/films/{filmId}/rating | Review Service |
| DELETE /api/v1/films/{filmId}/rating | Review Service |

| POST /api/v1/films/{filmId}/reviews | Review Service |
| GET /api/v1/films/{filmId}/reviews | Review Service |
| GET /api/v1/reviews/{reviewId} | Review Service |
| PUT /api/v1/reviews/{reviewId} | Review Service |
| DELETE /api/v1/reviews/{reviewId} | Review Service |
| GET /api/v1/users/{userId}/reviews | Review Service |
| POST /api/v1/reviews/{reviewId}/like | Review Service |
| DELETE /api/v1/reviews/{reviewId}/like | Review Service |
| POST /api/v1/reviews/{reviewId}/comments | Review Service |
| GET /api/v1/reviews/{reviewId}/comments | Review Service |
| PUT /api/v1/comments/{commentId} | Review Service |
| DELETE /api/v1/comments/{commentId} | Review Service |

| GET /api/v1/users/me/watchlist | Movie Lists Service |
| POST /api/v1/users/me/watchlist/{filmId} | Movie Lists Service |
| DELETE /api/v1/users/me/watchlist/{filmId} | Movie Lists Service |

| GET /api/v1/users/me/diary | Movie Lists Service |
| GET /api/v1/diary/{entryId} | Movie Lists Service |
| POST /api/v1/users/me/diary | Movie Lists Service |
| PUT /api/v1/diary/{entryId} | Movie Lists Service |
| DELETE /api/v1/diary/{entryId} | Movie Lists Service |

| GET /api/v1/users/me/lists | Movie Lists Service |
| GET /api/v1/users/{userId}/lists | Movie Lists Service |
| POST /api/v1/lists | Movie Lists Service |
| GET /api/v1/lists/{listId} | Movie Lists Service |
| PUT /api/v1/lists/{listId} | Movie Lists Service |
| DELETE /api/v1/lists/{listId} | Movie Lists Service |
| POST /api/v1/lists/{listId}/films/{filmId} | Movie Lists Service |
| DELETE /api/v1/lists/{listId}/films/{filmId} | Movie Lists Service |
| GET /api/v1/lists/{listId}/films | Movie Lists Service |
| POST /api/v1/lists/{listId}/like | Movie Lists Service |
| DELETE /api/v1/lists/{listId}/like | Movie Lists Service |
| GET /api/v1/lists/{listId}/comments | Movie Lists Service |
| POST /api/v1/lists/{listId}/comments | Movie Lists Service |
| PUT /api/v1/list-comments/{commentId} | Movie Lists Service |
| DELETE /api/v1/list-comments/{commentId} | Movie Lists Service |

| GET /api/v1/feed/following | Notification Service |
| GET /api/v1/notifications | Notification Service |
| GET /api/v1/notifications/unread-count | Notification Service |
| PUT /api/v1/notifications/{notificationId}/read | Notification Service |
| PUT /api/v1/notifications/read-all | Notification Service |
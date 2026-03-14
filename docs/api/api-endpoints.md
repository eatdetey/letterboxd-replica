# API Specification

## Клиентский API

---

### 2.1. Аутентификация

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/v1/auth/register` | Регистрация нового пользователя |
| `POST` | `/api/v1/auth/login` | Вход в систему |
| `POST` | `/api/v1/auth/logout` | Выход из системы |
| `POST` | `/api/v1/auth/refresh` | Обновление токена доступа |
| `PUT` | `/api/v1/users/me/password` | Смена пароля |

---

### 2.2. Каталог фильмов

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/api/v1/films/popular` | Список популярных фильмов (с пагинацией) |
| `GET` | `/api/v1/films/search?q={text}` | Поиск фильмов по названию |
| `GET` | `/api/v1/films/{filmId}` | Детальная информация о фильме |
| `GET` | `/api/v1/films/{filmId}/cast` | Актёры и съемочная группа |

---

### 2.3. Оценки и отзывы

#### Оценки

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `PUT` | `/api/v1/films/{filmId}/rating` | Поставить или обновить свою оценку |
| `DELETE` | `/api/v1/films/{filmId}/rating` | Удалить свою оценку |
| `GET` | `/api/v1/films/{filmId}/rating` | Получить свою оценку за фильм |

#### Отзывы

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/v1/films/{filmId}/reviews` | Написать отзыв |
| `GET` | `/api/v1/films/{filmId}/reviews` | Все отзывы к фильму (с пагинацией) |
| `GET` | `/api/v1/reviews/{reviewId}` | Получить конкретный отзыв |
| `PUT` | `/api/v1/reviews/{reviewId}` | Редактировать свой отзыв |
| `DELETE` | `/api/v1/reviews/{reviewId}` | Удалить свой отзыв |
| `GET` | `/api/v1/users/{userId}/reviews` | Все отзывы пользователя |

---

### 2.4. Watchlist (Список для просмотра)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/api/v1/users/me/watchlist` | Получить свой watchlist (с пагинацией) |
| `POST` | `/api/v1/users/me/watchlist/{filmId}` | Добавить фильм в watchlist |
| `DELETE` | `/api/v1/users/me/watchlist/{filmId}` | Удалить фильм из watchlist |

---

### 2.5. Дневник просмотров

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/api/v1/users/me/diary` | История просмотров (с пагинацией) |
| `GET` | `/api/v1/diary/{entryId}` | Получить запись дневника |
| `POST` | `/api/v1/users/me/diary` | Добавить запись о просмотре |
| `PUT` | `/api/v1/diary/{entryId}` | Изменить запись |
| `DELETE` | `/api/v1/diary/{entryId}` | Удалить запись |

---

### 2.6. Профиль пользователя

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/api/v1/users/me` | Получить свой профиль |
| `PUT` | `/api/v1/users/me` | Редактировать свой профиль |
| `GET` | `/api/v1/users/{userId}` | Публичный профиль другого пользователя |

---

### 3.1. Дополнительные данные о фильме

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/api/v1/films/{filmId}/trailers` | Трейлеры фильма |
| `GET` | `/api/v1/films/{filmId}/similar` | Похожие фильмы |

---

### 3.2. Социальные функции

#### Подписки

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/v1/users/{userId}/follow` | Подписаться на пользователя |
| `DELETE` | `/api/v1/users/{userId}/follow` | Отписаться |
| `GET` | `/api/v1/users/{userId}/followers` | Список подписчиков |
| `GET` | `/api/v1/users/{userId}/following` | Список подписок |

#### Лента и уведомления

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/api/v1/feed/following` | Лента активности друзей |
| `GET` | `/api/v1/notifications` | Уведомления пользователя |
| `GET` | `/api/v1/notifications/unread-count` | Количество непрочитанных уведомлений |
| `PUT` | `/api/v1/notifications/{notificationId}/read` | Отметить уведомление как прочитанное |
| `PUT` | `/api/v1/notifications/read-all` | Отметить все уведомления как прочитанные |

#### Комментарии к отзывам

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/v1/reviews/{reviewId}/comments` | Добавить комментарий |
| `GET` | `/api/v1/reviews/{reviewId}/comments` | Получить комментарии |
| `PUT` | `/api/v1/comments/{commentId}` | Редактировать комментарий |
| `DELETE` | `/api/v1/comments/{commentId}` | Удалить комментарий |

---

### 3.3. Пользовательские списки (коллекции)

#### Управление списками

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/api/v1/users/me/lists` | Список своих коллекций |
| `GET` | `/api/v1/users/{userId}/lists` | Список коллекций пользователя |
| `POST` | `/api/v1/lists` | Создать новую коллекцию |
| `GET` | `/api/v1/lists/{listId}` | Получить коллекцию |
| `PUT` | `/api/v1/lists/{listId}` | Редактировать коллекцию |
| `DELETE` | `/api/v1/lists/{listId}` | Удалить коллекцию |

#### Фильмы в коллекциях

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/v1/lists/{listId}/films/{filmId}` | Добавить фильм в коллекцию |
| `DELETE` | `/api/v1/lists/{listId}/films/{filmId}` | Удалить фильм из коллекции |
| `GET` | `/api/v1/lists/{listId}/films` | Получить фильмы коллекции |

#### Взаимодействие с коллекциями

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/v1/lists/{listId}/like` | Лайкнуть коллекцию |
| `DELETE` | `/api/v1/lists/{listId}/like` | Убрать лайк |
| `GET` | `/api/v1/lists/{listId}/comments` | Комментарии к коллекции |
| `POST` | `/api/v1/lists/{listId}/comments` | Добавить комментарий |
| `PUT` | `/api/v1/list-comments/{commentId}` | Редактировать комментарий |
| `DELETE` | `/api/v1/list-comments/{commentId}` | Удалить комментарий |

---

### 3.4. Настройки пользователя

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `PUT` | `/api/v1/users/me/password` | Смена пароля |
| `POST` | `/api/v1/users/me/avatar` | Загрузить аватар |
| `DELETE` | `/api/v1/users/me/avatar` | Удалить аватар |
| `GET` | `/api/v1/users/me/settings` | Получить настройки приватности |
| `PUT` | `/api/v1/users/me/settings` | Обновить настройки приватности |

---

## Внутренние микросервисы

### Auth Service (User Service)

| Метод | Описание |
|-------|----------|
| `Register` | Регистрация нового пользователя |
| `Login` | Вход в систему |
| `Logout` | Выход из системы |
| `RefreshToken` | Обновление токена |
| `ValidateToken` | Проверка токена (используется всеми сервисами) |
| `ChangePassword` | Смена пароля |

---

### User Service

| Метод | Описание |
|-------|----------|
| `GetProfile` | Получить профиль |
| `UpdateProfile` | Обновить профиль |
| `GetPublicProfile` | Публичный профиль |
| `FollowUser` | Подписаться на пользователя |
| `UnfollowUser` | Отписаться |
| `GetFollowers` | Список подписчиков |
| `GetFollowing` | Список подписок |
| `UpdateSettings` | Обновить настройки |
| `UploadAvatar` | Загрузить аватар |
| `DeleteAvatar` | Удалить аватар |

---

### Movie Service

| Метод | Описание |
|-------|----------|
| `GetPopularFilms` | Получить популярные фильмы |
| `SearchFilms` | Поиск фильмов |
| `GetFilmDetail` | Детальная информация о фильме |
| `GetFilmCast` | Актёры и команда |
| `GetTrailers` | Трейлеры фильма |
| `GetSimilarFilms` | Похожие фильмы |
| `GetFilmsByIds` | Получить несколько фильмов по ID |


---

### Review Service

| Метод | Описание |
|-------|----------|
| `CreateReview` | Создать отзыв |
| `UpdateReview` | Обновить отзыв |
| `DeleteReview` | Удалить отзыв |
| `GetReview` | Получить отзыв |
| `GetFilmReviews` | Отзывы к фильму |
| `GetUserReviews` | Отзывы пользователя |
| `RateFilm` | Поставить оценку |
| `UpdateRating` | Изменить оценку |
| `DeleteRating` | Удалить оценку |
| `GetUserRating` | Получить оценку пользователя за фильм |
| `GetFilmRatings` | Все оценки фильма |
| `GetAverageRating` | Средний рейтинг фильма |
| `GetUserRatings` | Все оценки пользователя |
| `LikeReview` | Лайкнуть отзыв |
| `UnlikeReview` | Убрать лайк |
| `AddComment` | Добавить комментарий |
| `UpdateComment` | Обновить комментарий |
| `DeleteComment` | Удалить комментарий |
| `GetComments` | Получить комментарии |

---

### Watchlist & Diary Service (Movie Lists Service)

| Метод | Описание |
|-------|----------|
| `AddToWatchlist` | Добавить в watchlist |
| `RemoveFromWatchlist` | Удалить из watchlist |
| `GetWatchlist` | Получить watchlist |
| `IsInWatchlist` | Проверить наличие в watchlist |
| `AddDiaryEntry` | Добавить запись в дневник |
| `UpdateDiaryEntry` | Обновить запись |
| `DeleteDiaryEntry` | Удалить запись |
| `GetDiary` | Получить дневник |
| `IsFilmWatched` | Проверить, просмотрен ли фильм |

---

### List Service (Movie Lists Service)

| Метод | Описание |
|-------|----------|
| `CreateList` | Создать коллекцию |
| `UpdateList` | Обновить коллекцию |
| `DeleteList` | Удалить коллекцию |
| `GetList` | Получить коллекцию |
| `GetUserLists` | Коллекции пользователя |
| `GetPublicLists` | Публичные коллекции |
| `AddFilmToList` | Добавить фильм в коллекцию |
| `RemoveFilmFromList` | Удалить фильм из коллекции |
| `GetListFilms` | Фильмы в коллекции |
| `LikeList` | Лайкнуть коллекцию |
| `UnlikeList` | Убрать лайк |
| `AddComment` | Добавить комментарий к коллекции |
| `UpdateComment` | Обновить комментарий |
| `DeleteComment` | Удалить комментарий |
| `GetComments` | Комментарии к коллекции |

---

### Notification Service

| Метод | Описание |
|-------|----------|
| `GetNotifications` | Получить уведомления |
| `MarkAsRead` | Отметить как прочитанное |
| `MarkAllAsRead` | Отметить всё как прочитанное |
| `GetUnreadCount` | Количество непрочитанных |
| `NotifyFollow` | Уведомление о подписке |
| `NotifyLike` | Уведомление о лайке |
| `NotifyComment` | Уведомление о комментарии |
| `NotifyMention` | Уведомление об упоминании |

---

## Условные обозначения

| Символ | Значение |
|--------|----------|
| `{filmId}` | ID фильма |
| `{userId}` | ID пользователя |
| `{reviewId}` | ID отзыва |
| `{listId}` | ID коллекции |
| `{commentId}` | ID комментария |
| `{entryId}` | ID записи в дневнике |
| `{notificationId}` | ID уведомления |
| `?q={text}` | Поисковый запрос |
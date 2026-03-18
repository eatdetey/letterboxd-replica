from datetime import date
from typing import Iterable, List, Optional

import grpc
from django.db import transaction

from movies.models import Genre, Movie
from movies.services.movie_service import MovieService
from grpc_layer.protobuf.movie.v1 import movie_pb2, movie_pb2_grpc


def _parse_int(value: str) -> Optional[int]:
    try:
        return int(value)
    except (TypeError, ValueError):
        return None


def _movie_to_detailed(movie: Movie) -> movie_pb2.MovieDetailed:
    release_year = movie.release_date.year if movie.release_date else 0
    genres = [genre.name for genre in movie.genres.all()]
    # playlists are not stored yet; return an empty list to keep contract stable
    return movie_pb2.MovieDetailed(
        id=str(movie.id),
        title=movie.title,
        description=movie.description or "",
        release_year=release_year,
        genres=genres,
        poster=movie.poster or "",
        playlists=[],
    )


def _movie_to_proto(movie: Movie) -> movie_pb2.Movie:
    release_year = movie.release_date.year if movie.release_date else 0
    genres = [genre.name for genre in movie.genres.all()]
    return movie_pb2.Movie(
        id=str(movie.id),
        title=movie.title,
        description=movie.description or "",
        release_year=release_year,
        genres=genres,
        poster=movie.poster or "",
    )


def _get_genres_by_names(names: Iterable[str]) -> List[Genre]:
    genres: List[Genre] = []
    for name in names:
        genre_obj, _ = Genre.objects.get_or_create(name=name)
        genres.append(genre_obj)
    return genres


class MovieServiceHandler(movie_pb2_grpc.MovieServiceServicer):
    DEFAULT_LIMIT = 20

    def GetMovies(self, request, context):
        limit = request.limit or self.DEFAULT_LIMIT
        offset = max(request.offset, 0)

        if limit <= 0:
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "limit must be positive")

        ids: Optional[List[int]] = None
        if request.ids:
            ids = []
            for raw_id in request.ids:
                parsed = _parse_int(raw_id)
                if parsed is None:
                    context.abort(grpc.StatusCode.INVALID_ARGUMENT, "ids must be integers")
                ids.append(parsed)

        search_query = request.search_query if request.HasField("search_query") else None
        genre = request.genre if request.HasField("genre") else None

        queryset, total = MovieService.get_movies(
            limit=limit,
            offset=offset,
            search_query=search_query,
            genre=genre,
            ids=ids,
        )

        items = [_movie_to_detailed(movie) for movie in queryset]
        return movie_pb2.GetMoviesResponse(
            items=items,
            total=total,
        )

    def CreateMovie(self, request, context):
        if request.release_year <= 0:
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "release_year must be positive")

        with transaction.atomic():
            movie = Movie.objects.create(
                title=request.title,
                original_title=request.title,
                description=request.description or "",
                release_date=date(request.release_year, 1, 1),
                poster=request.poster or "",
            )

            if request.genres:
                movie.genres.set(_get_genres_by_names(request.genres))

        return movie_pb2.CreateMovieResponse(movie=_movie_to_proto(movie))

    def UpdateMovie(self, request, context):
        movie = self._get_movie_or_abort(request.id, context)

        if request.HasField("title"):
            movie.title = request.title
        if request.HasField("description"):
            movie.description = request.description or ""
        if request.HasField("release_year"):
            if request.release_year <= 0:
                context.abort(grpc.StatusCode.INVALID_ARGUMENT, "release_year must be positive")
            movie.release_date = date(request.release_year, 1, 1)

        # Proto3 repeated fields are always present; update only when provided content exists
        if request.genres:
            movie.genres.set(_get_genres_by_names(request.genres))

        if request.poster:
            movie.poster = request.poster

        movie.save()

        return movie_pb2.UpdateMovieResponse(movie=_movie_to_proto(movie))

    def DeleteMovie(self, request, context):
        movie = self._get_movie_or_abort(request.id, context)
        movie.delete()
        return movie_pb2.DeleteMovieResponse()

    def _get_movie_or_abort(self, raw_id: str, context) -> Movie:
        movie_id = _parse_int(raw_id)
        if movie_id is None:
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "id must be an integer")

        try:
            return Movie.objects.get(id=movie_id)
        except Movie.DoesNotExist:
            context.abort(grpc.StatusCode.NOT_FOUND, "movie not found")

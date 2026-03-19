from datetime import date
from typing import Dict, Iterable, List, Optional, Sequence, Tuple

import grpc
from django.db import transaction

from grpc_layer.clients import MovieListClient
from grpc_layer.context_utils import (
    AUTHORIZATION_KEY,
    REQUEST_ID_KEY,
    get_auth_header_from_context,
    get_claims_from_context,
    get_request_id_from_context,
)
from grpc_layer.protobuf.movie.v1 import movie_pb2, movie_pb2_grpc
from movies.models import Genre, Movie
from movies.services.movie_service import MovieService

Metadata = Sequence[Tuple[str, str]]


def _parse_int(value: str) -> Optional[int]:
    try:
        return int(value)
    except (TypeError, ValueError):
        return None


def _movie_to_detailed(
    movie: Movie,
    playlists: Optional[List[movie_pb2.PlaylistInfo]] = None,
) -> movie_pb2.MovieDetailed:
    release_year = movie.release_date.year if movie.release_date else 0
    genres = [genre.name for genre in movie.genres.all()]
    return movie_pb2.MovieDetailed(
        id=str(movie.id),
        title=movie.title,
        description=movie.description or "",
        release_year=release_year,
        genres=genres,
        poster=movie.poster or "",
        playlists=playlists or [],
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

    def __init__(self, movie_list_client: Optional[MovieListClient] = None):
        self.movie_list_client = movie_list_client or MovieListClient()

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

        if request.HasField("playlist_id"):
            metadata = self._build_upstream_metadata(context, require_auth=True)
            candidate_movie_ids = [
                str(movie_id)
                for movie_id in MovieService.get_movie_ids(
                    search_query=search_query,
                    genre=genre,
                    ids=ids,
                )
            ]
            try:
                filter_response = self.movie_list_client.filter_movies_by_playlist(
                    playlist_id=request.playlist_id,
                    candidate_movie_ids=candidate_movie_ids,
                    metadata=metadata,
                )
            except grpc.RpcError as err:
                self._abort_from_upstream_error(err, context)
            filtered_ids: List[int] = []
            for raw_id in filter_response.movie_ids:
                parsed_id = _parse_int(raw_id)
                if parsed_id is None:
                    context.abort(grpc.StatusCode.INTERNAL, "movie-list-service returned invalid movie id")
                filtered_ids.append(parsed_id)
            ids = filtered_ids

        queryset, total = MovieService.get_movies(
            limit=limit,
            offset=offset,
            search_query=search_query,
            genre=genre,
            ids=ids,
        )

        movies = list(queryset)
        playlists_by_movie: Dict[str, List[movie_pb2.PlaylistInfo]] = {}
        if request.enrich_playlists and movies:
            metadata = self._build_upstream_metadata(context, require_auth=True)
            playlists_by_movie = self._get_playlists_by_movie(
                movies=movies,
                metadata=metadata,
                context=context,
            )

        items = [
            _movie_to_detailed(
                movie,
                playlists=playlists_by_movie.get(str(movie.id), []),
            )
            for movie in movies
        ]
        return movie_pb2.GetMoviesResponse(
            items=items,
            total=total,
        )

    def GetPlaylistsForUser(self, request, context):
        metadata = self._build_upstream_metadata(context, require_auth=True)
        try:
            response = self.movie_list_client.get_playlists_for_user(metadata=metadata)
        except grpc.RpcError as err:
            self._abort_from_upstream_error(err, context)
        items = [
            movie_pb2.Playlist(
                id=playlist.id,
                name=playlist.name,
                movies_count=playlist.movies_count,
            )
            for playlist in response.playlists
        ]
        return movie_pb2.GetPlaylistsForUserResponse(items=items)

    def CreatePlaylist(self, request, context):
        metadata = self._build_upstream_metadata(context, require_auth=True)
        try:
            response = self.movie_list_client.create_playlist(
                name=request.name,
                metadata=metadata,
            )
        except grpc.RpcError as err:
            self._abort_from_upstream_error(err, context)
        return movie_pb2.CreatePlaylistResponse(
            playlist=movie_pb2.Playlist(
                id=response.id,
                name=response.name,
                movies_count=response.movies_count,
            )
        )

    def RenamePlaylist(self, request, context):
        metadata = self._build_upstream_metadata(context, require_auth=True)
        try:
            response = self.movie_list_client.rename_playlist(
                playlist_id=request.playlist_id,
                new_name=request.new_name,
                metadata=metadata,
            )
        except grpc.RpcError as err:
            self._abort_from_upstream_error(err, context)
        return movie_pb2.RenamePlaylistResponse(
            playlist=movie_pb2.Playlist(
                id=response.id,
                name=response.name,
                movies_count=response.movies_count,
            )
        )

    def DeletePlaylist(self, request, context):
        metadata = self._build_upstream_metadata(context, require_auth=True)
        try:
            self.movie_list_client.delete_playlist(
                playlist_id=request.playlist_id,
                metadata=metadata,
            )
        except grpc.RpcError as err:
            self._abort_from_upstream_error(err, context)
        return movie_pb2.DeletePlaylistResponse()

    def AddMovieToPlaylist(self, request, context):
        metadata = self._build_upstream_metadata(context, require_auth=True)
        try:
            self.movie_list_client.add_movie_to_playlist(
                playlist_id=request.playlist_id,
                movie_id=request.movie_id,
                metadata=metadata,
            )
        except grpc.RpcError as err:
            self._abort_from_upstream_error(err, context)
        return movie_pb2.AddMovieToPlaylistResponse()

    def RemoveMovieFromPlaylist(self, request, context):
        metadata = self._build_upstream_metadata(context, require_auth=True)
        try:
            self.movie_list_client.remove_movie_from_playlist(
                playlist_id=request.playlist_id,
                movie_id=request.movie_id,
                metadata=metadata,
            )
        except grpc.RpcError as err:
            self._abort_from_upstream_error(err, context)
        return movie_pb2.RemoveMovieFromPlaylistResponse()

    def CreateMovie(self, request, context):
        self._require_admin(context)
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
        self._require_admin(context)
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
        self._require_admin(context)
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

    def _build_upstream_metadata(self, context, require_auth: bool) -> List[Tuple[str, str]]:
        metadata: List[Tuple[str, str]] = []
        auth_header = get_auth_header_from_context(context)
        if auth_header:
            metadata.append((AUTHORIZATION_KEY, auth_header))
        elif require_auth:
            context.abort(grpc.StatusCode.UNAUTHENTICATED, "authorization required")

        request_id = get_request_id_from_context(context)
        if request_id:
            metadata.append((REQUEST_ID_KEY, request_id))
        return metadata

    def _abort_from_upstream_error(self, err: grpc.RpcError, context):
        code = err.code() if hasattr(err, "code") else grpc.StatusCode.INTERNAL
        details = err.details() if hasattr(err, "details") else ""
        context.abort(code, details or "movie-list-service request failed")

    def _get_playlists_by_movie(
        self,
        movies: List[Movie],
        metadata: Metadata,
        context,
    ) -> Dict[str, List[movie_pb2.PlaylistInfo]]:
        result: Dict[str, List[movie_pb2.PlaylistInfo]] = {}
        for movie in movies:
            movie_id = str(movie.id)
            try:
                response = self.movie_list_client.get_playlists_for_movie(
                    movie_id=movie_id,
                    metadata=metadata,
                )
            except grpc.RpcError as err:
                self._abort_from_upstream_error(err, context)
            result[movie_id] = [
                movie_pb2.PlaylistInfo(id=playlist.id, name=playlist.name)
                for playlist in response.playlists
            ]
        return result

    def _require_admin(self, context):
        claims = get_claims_from_context(context)
        if not claims:
            context.abort(grpc.StatusCode.UNAUTHENTICATED, "authorization required")
        roles = [role.lower() for role in (claims.roles or [])]
        if "admin" not in roles:
            context.abort(grpc.StatusCode.PERMISSION_DENIED, "admin role required")

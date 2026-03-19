import os
from typing import Optional, Sequence, Tuple

import grpc

from grpc_layer.protobuf.movielist.v1 import movielist_pb2, movielist_pb2_grpc

Metadata = Sequence[Tuple[str, str]]


def _parse_timeout(raw_value: str) -> float:
    try:
        value = float(raw_value)
    except (TypeError, ValueError):
        return 3.0
    if value <= 0:
        return 3.0
    return value


class MovieListClient:
    def __init__(self, address: Optional[str] = None, timeout_sec: Optional[float] = None):
        self.address = address or os.environ.get("MOVIE_LIST_SERVICE_ADDRESS", "movie-list-service:50052")
        default_timeout = _parse_timeout(os.environ.get("MOVIE_LIST_SERVICE_TIMEOUT_SEC", "3"))
        self.timeout_sec = timeout_sec if timeout_sec and timeout_sec > 0 else default_timeout

        self._channel = grpc.insecure_channel(self.address)
        self._stub = movielist_pb2_grpc.MovieListServiceStub(self._channel)

    def filter_movies_by_playlist(
        self,
        playlist_id: str,
        candidate_movie_ids: Sequence[str],
        metadata: Optional[Metadata] = None,
    ) -> movielist_pb2.FilterMoviesByPlaylistResponse:
        return self._stub.FilterMoviesByPlaylist(
            movielist_pb2.FilterMoviesByPlaylistRequest(
                playlist_id=playlist_id,
                candidate_movie_ids=list(candidate_movie_ids),
            ),
            timeout=self.timeout_sec,
            metadata=metadata,
        )

    def get_playlists_for_movie(
        self,
        movie_id: str,
        metadata: Optional[Metadata] = None,
    ) -> movielist_pb2.GetPlaylistsForMovieResponse:
        return self._stub.GetPlaylistsForMovie(
            movielist_pb2.GetPlaylistsForMovieRequest(movie_id=movie_id),
            timeout=self.timeout_sec,
            metadata=metadata,
        )

    def get_playlists_for_user(
        self,
        metadata: Optional[Metadata] = None,
    ) -> movielist_pb2.GetPlaylistsForUserResponse:
        return self._stub.GetPlaylistsForUser(
            movielist_pb2.GetPlaylistsForUserRequest(),
            timeout=self.timeout_sec,
            metadata=metadata,
        )

    def create_playlist(
        self,
        name: str,
        metadata: Optional[Metadata] = None,
    ) -> movielist_pb2.CreatePlaylistResponse:
        return self._stub.CreatePlaylist(
            movielist_pb2.CreatePlaylistRequest(name=name),
            timeout=self.timeout_sec,
            metadata=metadata,
        )

    def rename_playlist(
        self,
        playlist_id: str,
        new_name: str,
        metadata: Optional[Metadata] = None,
    ) -> movielist_pb2.RenamePlaylistResponse:
        return self._stub.RenamePlaylist(
            movielist_pb2.RenamePlaylistRequest(
                playlist_id=playlist_id,
                new_name=new_name,
            ),
            timeout=self.timeout_sec,
            metadata=metadata,
        )

    def delete_playlist(
        self,
        playlist_id: str,
        metadata: Optional[Metadata] = None,
    ) -> movielist_pb2.DeletePlaylistResponse:
        return self._stub.DeletePlaylist(
            movielist_pb2.DeletePlaylistRequest(playlist_id=playlist_id),
            timeout=self.timeout_sec,
            metadata=metadata,
        )

    def add_movie_to_playlist(
        self,
        playlist_id: str,
        movie_id: str,
        metadata: Optional[Metadata] = None,
    ) -> movielist_pb2.AddMovieToPlaylistResponse:
        return self._stub.AddMovieToPlaylist(
            movielist_pb2.AddMovieToPlaylistRequest(
                playlist_id=playlist_id,
                movie_id=movie_id,
            ),
            timeout=self.timeout_sec,
            metadata=metadata,
        )

    def remove_movie_from_playlist(
        self,
        playlist_id: str,
        movie_id: str,
        metadata: Optional[Metadata] = None,
    ) -> movielist_pb2.RemoveMovieFromPlaylistResponse:
        return self._stub.RemoveMovieFromPlaylist(
            movielist_pb2.RemoveMovieFromPlaylistRequest(
                playlist_id=playlist_id,
                movie_id=movie_id,
            ),
            timeout=self.timeout_sec,
            metadata=metadata,
        )

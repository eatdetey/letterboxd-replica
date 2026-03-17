import grpc
from datetime import datetime

from movies.models import Movie, Genre
from movies.services.movie_service import MovieService

from grpc_layer.protobuf import movie_pb2, movie_pb2_grpc


class MovieServiceHandler(movie_pb2_grpc.MovieServiceServicer):

    def GetMovies(self, request, context):
        movies, total = MovieService.get_movies(
            limit=request.limit,
            offset=request.offset,
            search_query=request.search_query if request.HasField("search_query") else None,
            genre=request.genre if request.HasField("genre") else None,
            ids=list(request.ids)
        )

        return movie_pb2.MoviesResponse(
            items=[self._to_movie_detailed(m, request.enrich_playlists, request.user_id) for m in movies],
            total=total
        )


    def CreateMovie(self, request, context):
        movie = Movie.objects.create(
            title=request.title,
            description=request.description,
            release_date=datetime(request.release_year, 1, 1)
        )

        genres = Genre.objects.filter(name__in=request.genres)
        movie.genres.set(genres)

        return self._to_movie_response(movie)


    def UpdateMovie(self, request, context):
        try:
            movie = Movie.objects.get(id=request.id)
        except Movie.DoesNotExist:
            context.abort(grpc.StatusCode.NOT_FOUND, "Movie not found")

        if request.HasField("title"):
            movie.title = request.title

        if request.HasField("description"):
            movie.description = request.description

        if request.HasField("release_year"):
            movie.release_date = datetime(request.release_year, 1, 1)

        if request.genres:
            genres = Genre.objects.filter(name__in=request.genres)
            movie.genres.set(genres)

        movie.save()

        return self._to_movie_response(movie)


    def DeleteMovie(self, request, context):
        deleted, _ = Movie.objects.filter(id=request.id).delete()

        if deleted == 0:
            context.abort(grpc.StatusCode.NOT_FOUND, "Movie not found")

        return movie_pb2.Empty()


    def MovieExists(self, request, context):
        exists = Movie.objects.filter(id=request.id).exists()
        return movie_pb2.ExistsResponse(exists=exists)


    def _to_movie_response(self, movie):
        return movie_pb2.MovieResponse(
            id=str(movie.id),
            title=movie.title,
            description=movie.description or "",
            release_year=movie.release_date.year if movie.release_date else 0,
            genres=[g.name for g in movie.genres.all()]
        )


    def _to_movie_detailed(self, movie, enrich_playlists, user_id):
        playlists = []

        if enrich_playlists:
            playlists = self._get_playlists_stub(movie.id, user_id)

        return movie_pb2.MovieDetailedResponse(
            id=str(movie.id),
            title=movie.title,
            description=movie.description or "",
            release_year=movie.release_date.year if movie.release_date else 0,
            genres=[g.name for g in movie.genres.all()],
            playlists=playlists
        )


    def _get_playlists_stub(self, movie_id, user_id):
        return []

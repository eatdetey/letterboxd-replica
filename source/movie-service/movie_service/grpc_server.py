import grpc
from concurrent import futures
import os
import django
from grpc_reflection.v1alpha import reflection

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "movie_service.settings")
django.setup()

from movies.services.movie_service import MovieService
import movie_pb2
import movie_pb2_grpc


class MovieServiceHandler(movie_pb2_grpc.MovieServiceServicer):

    def GetMovie(self, request, context):
        movie = MovieService.get_movie(request.id)
        return movie_pb2.MovieResponse(
            id=movie.id,
            title=movie.title,
            original_title=movie.original_title or "",
            description=movie.description or "",
            release_date=str(movie.release_date) if movie.release_date else "",
            duration_minutes=movie.duration_minutes or 0,
            country=movie.country or "",
            age_rating=movie.age_rating or "",
            genres=[g.name for g in movie.genres.all()],
            cast=[movie_pb2.CastResponse(
                person_name=mc.person.name,
                role_type=mc.role_type,
                character_name=mc.character_name
            ) for mc in movie.moviecast_set.all()]
        )

    def ListMovies(self, request, context):
        movies = MovieService.list_movies(request.limit or 10, request.offset or 0)
        return movie_pb2.ListMoviesResponse(
            movies=[self._to_proto(m) for m in movies]
        )

    def SearchMovies(self, request, context):
        movies = MovieService.search_movies(request.query, request.limit or 10, request.offset or 0)
        return movie_pb2.ListMoviesResponse(
            movies=[self._to_proto(m) for m in movies]
        )

    def _to_proto(self, movie):
        return movie_pb2.MovieResponse(
            id=movie.id,
            title=movie.title,
            original_title=movie.original_title or "",
            description=movie.description or "",
            release_date=str(movie.release_date) if movie.release_date else "",
            duration_minutes=movie.duration_minutes or 0,
            country=movie.country or "",
            age_rating=movie.age_rating or "",
            genres=[g.name for g in movie.genres.all()],
            cast=[movie_pb2.CastResponse(
                person_name=mc.person.name,
                role_type=mc.role_type,
                character_name=mc.character_name
            ) for mc in movie.moviecast_set.all()]
        )



def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    movie_pb2_grpc.add_MovieServiceServicer_to_server(MovieServiceHandler(), server)

    SERVICE_NAMES = (
        movie_pb2_grpc.MovieServiceServicer.__name__,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)

    server.add_insecure_port('0.0.0.0:50051')
    server.start()
    print("gRPC server running on 0.0.0.0:50051")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()

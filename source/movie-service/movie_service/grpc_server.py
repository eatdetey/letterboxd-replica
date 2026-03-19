import grpc
import os
import django
from concurrent import futures

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "movie_service.settings")
django.setup()

from grpc_layer.servicers.movie_servicer import MovieServiceHandler
from grpc_layer.protobuf.movie.v1 import movie_pb2_grpc
from grpc_layer.interceptors import RequestIdInterceptor, AuthInterceptor


def serve():
    port = os.environ.get("MOVIE_SERVICE_PORT", "50051")
    access_secret = os.environ.get("AUTH_ACCESS_SECRET")

    if not access_secret:
        raise RuntimeError("AUTH_ACCESS_SECRET env variable is required to start movie-service")

    protected_methods = {
        "/movie.v1.MovieService/GetPlaylistsForUser",
        "/movie.v1.MovieService/CreatePlaylist",
        "/movie.v1.MovieService/RenamePlaylist",
        "/movie.v1.MovieService/DeletePlaylist",
        "/movie.v1.MovieService/AddMovieToPlaylist",
        "/movie.v1.MovieService/RemoveMovieFromPlaylist",
        "/movie.v1.MovieService/CreateMovie",
        "/movie.v1.MovieService/UpdateMovie",
        "/movie.v1.MovieService/DeleteMovie",
    }

    interceptors = [
        RequestIdInterceptor(),
        AuthInterceptor(access_secret, protected_methods),
    ]

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10), interceptors=interceptors)

    movie_pb2_grpc.add_MovieServiceServicer_to_server(
        MovieServiceHandler(),
        server
    )

    server.add_insecure_port(f'0.0.0.0:{port}')
    server.start()
    print(f"gRPC server running on 0.0.0.0:{port}")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()

import grpc
import os
import django
from concurrent import futures

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "movie_service.settings")
django.setup()

from grpc_layer.servicers.movie_servicer import MovieServiceHandler
from grpc_layer.protobuf.movie.v1 import movie_pb2_grpc


def serve():
    port = os.environ.get("MOVIE_SERVICE_PORT", "50051")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

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

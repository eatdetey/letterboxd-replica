import grpc
import os
import django
from concurrent import futures

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "movie_service.settings")
django.setup()

from grpc_layer.servicers.movie_servicer import MovieServiceHandler
from grpc_layer.protobuf.movie.v1 import movie_pb2_grpc


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    movie_pb2_grpc.add_MovieServiceServicer_to_server(
        MovieServiceHandler(),
        server
    )

    server.add_insecure_port('0.0.0.0:50051')
    server.start()
    print("gRPC server running on 0.0.0.0:50051")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()

import grpc
import os
import django
from concurrent import futures


os.environ.setdefault("DJANGO_SETTINGS_MODULE", "review_service.settings")
django.setup()

from grpc_layer.servicers.review_servicer import ReviewServiceHandler
from grpc_layer.protobuf import review_pb2_grpc


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    review_pb2_grpc.add_ReviewServiceServicer_to_server(
        ReviewServiceHandler(),
        server
    )

    server.add_insecure_port('0.0.0.0:50051')
    server.start()
    print("Review gRPC server running on 0.0.0.0:50051")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
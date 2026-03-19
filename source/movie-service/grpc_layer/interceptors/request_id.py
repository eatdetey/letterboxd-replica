import grpc
from grpc import ServerInterceptor, RpcMethodHandler
from typing import Callable

from grpc_layer.context_utils import ensure_request_id


class RequestIdInterceptor(ServerInterceptor):
    def intercept_service(self, continuation: Callable, handler_call_details):
        handler: RpcMethodHandler = continuation(handler_call_details)

        if handler is None:
            return None

        def wrap_unary_unary(request, context):
            ensure_request_id(context)
            return handler.unary_unary(request, context)

        if handler.unary_unary:
            return grpc.unary_unary_rpc_method_handler(
                wrap_unary_unary,
                request_deserializer=handler.request_deserializer,
                response_serializer=handler.response_serializer,
            )

        return handler

import os
from typing import Callable, Iterable, Optional, Set, Tuple

import grpc
from grpc import RpcMethodHandler, ServerInterceptor, StatusCode

from grpc_layer.auth import Claims, parse_token
from grpc_layer.context_utils import AUTHORIZATION_KEY


def _find_auth_header(metadata: Optional[Iterable[Tuple[str, str]]]) -> str:
    if not metadata:
        return ""
    for item in metadata:
        key = getattr(item, "key", None) or (item[0] if isinstance(item, (tuple, list)) and len(item) > 1 else None)
        value = getattr(item, "value", None) or (item[1] if isinstance(item, (tuple, list)) and len(item) > 1 else None)
        if key and isinstance(key, str) and key.lower() == AUTHORIZATION_KEY:
            return value or ""
    return ""


class AuthInterceptor(ServerInterceptor):
    def __init__(self, secret: Optional[str], protected_methods: Set[str]):
        self.secret = secret
        self.protected_methods = protected_methods

    def intercept_service(self, continuation: Callable, handler_call_details):
        handler: RpcMethodHandler = continuation(handler_call_details)
        if handler is None:
            return None

        method = handler_call_details.method
        should_protect = method in self.protected_methods

        if not should_protect:
            return handler

        if not self.secret:
            raise RuntimeError("AUTH_ACCESS_SECRET is not set for protected gRPC methods")

        def wrap_unary_unary(request, context):
            header = _find_auth_header(handler_call_details.invocation_metadata)
            token = self._extract_bearer(header)
            if not token:
                context.abort(StatusCode.UNAUTHENTICATED, "unauthorized")

            try:
                claims = parse_token(token, self.secret)
            except ValueError:
                context.abort(StatusCode.UNAUTHENTICATED, "unauthorized")
                return None  # pragma: no cover

            setattr(context, "claims", claims)
            return handler.unary_unary(request, context)

        if handler.unary_unary:
            return grpc.unary_unary_rpc_method_handler(
                wrap_unary_unary,
                request_deserializer=handler.request_deserializer,
                response_serializer=handler.response_serializer,
            )

        return handler

    @staticmethod
    def _extract_bearer(header: str) -> str:
        if not header:
            return ""
        parts = header.split(" ", 1)
        if len(parts) != 2:
            return ""
        scheme, token = parts
        if scheme.lower() != "bearer":
            return ""
        return token.strip()

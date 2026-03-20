import uuid
from typing import Iterable, Optional

import grpc

AUTHORIZATION_KEY = "authorization"
REQUEST_ID_KEY = "x-request-id"


def _find_metadata_value(metadata: Optional[Iterable]) -> str:
    if not metadata:
        return ""
    for item in metadata:
        # item can be a tuple or an object with key/value attrs depending on grpc version
        key = getattr(item, "key", None) or (item[0] if isinstance(item, (tuple, list)) and len(item) > 1 else None)
        value = getattr(item, "value", None) or (item[1] if isinstance(item, (tuple, list)) and len(item) > 1 else None)
        if key and isinstance(key, str) and key.lower() == AUTHORIZATION_KEY:
            return value or ""
    return ""


def get_auth_header_from_context(context: grpc.ServicerContext) -> str:
    # First try to get from context attributes (set by auth interceptor)
    if hasattr(context, "_auth_header"):
        return getattr(context, "_auth_header", "")
    # Fallback to invocation metadata
    return _find_metadata_value(context.invocation_metadata())


def get_request_id_from_context(context: grpc.ServicerContext) -> str:
    if hasattr(context, "request_id"):
        return getattr(context, "request_id")
    if not context:
        return ""
    for item in context.invocation_metadata():
        key = getattr(item, "key", None) or (item[0] if isinstance(item, (tuple, list)) and len(item) > 1 else None)
        value = getattr(item, "value", None) or (item[1] if isinstance(item, (tuple, list)) and len(item) > 1 else None)
        if key and isinstance(key, str) and key.lower() == REQUEST_ID_KEY:
            return value or ""
    return ""


def ensure_request_id(context: grpc.ServicerContext) -> str:
    req_id = get_request_id_from_context(context)
    if not req_id:
        req_id = str(uuid.uuid4())
    setattr(context, "request_id", req_id)
    try:
        context.set_trailing_metadata(((REQUEST_ID_KEY, req_id),))
    except Exception:
        # trailing metadata is best-effort
        pass
    return req_id


def get_claims_from_context(context: grpc.ServicerContext):
    if hasattr(context, "claims"):
        return getattr(context, "claims")
    return None

from dataclasses import dataclass
from typing import List

import jwt
from jwt import InvalidTokenError


@dataclass
class Claims:
    id: int
    username: str
    email: str
    status: str
    roles: List[str]
    is_deleted: bool = False


def parse_token(token: str, secret: str) -> Claims:
    try:
        payload = jwt.decode(
            token,
            secret,
            algorithms=["HS256"],
            options={"verify_aud": False, "verify_exp": False},
        )
    except InvalidTokenError as exc:
        raise ValueError("invalid token") from exc

    roles = payload.get("Roles") or payload.get("roles") or []
    if isinstance(roles, str):
        roles = [roles]
    return Claims(
        id=int(payload.get("Id") or payload.get("id") or 0),
        username=payload.get("Username") or payload.get("username") or "",
        email=payload.get("Email") or payload.get("email") or "",
        status=payload.get("Status") or payload.get("status") or "",
        roles=list(roles),
        is_deleted=bool(payload.get("IsDeleted") or payload.get("is_deleted") or False),
    )

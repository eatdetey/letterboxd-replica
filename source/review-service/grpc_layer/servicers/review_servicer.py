import os
from typing import List, Sequence, Tuple

import grpc

from grpc_layer.context_utils import (
    AUTHORIZATION_KEY,
    REQUEST_ID_KEY,
    get_auth_header_from_context,
    get_claims_from_context,
    get_request_id_from_context,
)
from grpc_layer.protobuf.movie.v1 import movie_pb2, movie_pb2_grpc
from grpc_layer.protobuf.review.v1 import review_pb2, review_pb2_grpc
from reviews.services.review_service import ReviewService

Metadata = Sequence[Tuple[str, str]]


def _parse_timeout(raw_value: str) -> float:
    try:
        value = float(raw_value)
    except (TypeError, ValueError):
        return 3.0
    if value <= 0:
        return 3.0
    return value


class ReviewServiceHandler(review_pb2_grpc.ReviewServiceServicer):
    def __init__(self):
        self.movie_service_address = os.environ.get("MOVIE_SERVICE_ADDRESS", "movie-service:50053")
        self.movie_service_timeout = _parse_timeout(os.environ.get("MOVIE_SERVICE_TIMEOUT_SEC", "3"))
        self.movie_channel = grpc.insecure_channel(self.movie_service_address)
        self.movie_stub = movie_pb2_grpc.MovieServiceStub(self.movie_channel)

    def GetReviews(self, request, context):
        if not request.movie_id:
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "movie_id is required")

        reviews = ReviewService.get_reviews(request.movie_id)
        return review_pb2.GetReviewsResponse(
            items=[self._to_review_response(item) for item in reviews],
        )

    def AddReview(self, request, context):
        if not request.movie_id:
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "movie_id is required")

        text = (request.text or "").strip()
        if not text:
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "text is required")

        claims = get_claims_from_context(context)
        if not claims or claims.id <= 0:
            context.abort(grpc.StatusCode.UNAUTHENTICATED, "authorization required")

        if not self._movie_exists(request.movie_id, context):
            context.abort(grpc.StatusCode.NOT_FOUND, "movie not found")

        review = ReviewService.add_review(
            movie_id=request.movie_id,
            user_id=str(claims.id),
            username=claims.username or "",
            text=text,
        )
        return review_pb2.AddReviewResponse(review=self._to_review_response(review))

    def _movie_exists(self, movie_id: str, context) -> bool:
        try:
            response = self.movie_stub.GetMovies(
                movie_pb2.GetMoviesRequest(
                    limit=1,
                    offset=0,
                    ids=[movie_id],
                ),
                timeout=self.movie_service_timeout,
                metadata=self._build_movie_service_metadata(context),
            )
            return response.total > 0
        except grpc.RpcError as err:
            code = err.code() if hasattr(err, "code") else grpc.StatusCode.INTERNAL
            details = err.details() if hasattr(err, "details") else ""
            context.abort(code, details or "movie-service request failed")

    def _build_movie_service_metadata(self, context) -> List[Tuple[str, str]]:
        metadata: List[Tuple[str, str]] = []
        auth_header = get_auth_header_from_context(context)
        if auth_header:
            metadata.append((AUTHORIZATION_KEY, auth_header))

        request_id = get_request_id_from_context(context)
        if request_id:
            metadata.append((REQUEST_ID_KEY, request_id))
        return metadata

    @staticmethod
    def _to_review_response(review) -> review_pb2.Review:
        return review_pb2.Review(
            id=str(review.id),
            movie_id=review.movie_id,
            user_id=review.user_id,
            username=review.username or "",
            text=review.text,
            created_at=review.created_at.isoformat() if review.created_at else "",
        )

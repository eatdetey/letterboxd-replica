import grpc
import logging

from reviews.models import Review
from reviews.services.review_service import ReviewService

from grpc_layer.protobuf import review_pb2, review_pb2_grpc
from grpc_layer.protobuf import user_pb2_grpc
from grpc_layer.protobuf.user_pb2 import GetUsersRequest
from grpc_layer.protobuf import movie_pb2_grpc
from grpc_layer.protobuf.movie_pb2 import MovieIdRequest


logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

USER_SERVICE_ADDRESS = "user-service:50051"
MOVIE_SERVICE_ADDRESS = "movie-service:50051"


class ReviewServiceHandler(review_pb2_grpc.ReviewServiceServicer):

    def GetReviews(self, request, context):
        reviews = ReviewService.get_reviews(request.movie_id)

        user_ids = [r.user_id for r in reviews]
        users_map = self._get_users_by_ids(user_ids)

        return review_pb2.ReviewsResponse(
            items=[self._to_review_response(r, users_map) for r in reviews]
        )

    def AddReview(self, request, context):
        # Проверка существования пользователя
        users_response = self._get_users_by_ids([request.user_id])

        if not users_response or request.user_id not in users_response:
            context.abort(
                grpc.StatusCode.NOT_FOUND,
                f"User with id {request.user_id} not found"
            )

        user = users_response[request.user_id]

        # Проверка существования фильма
        if not self._movie_exists(request.movie_id):
            context.abort(
                grpc.StatusCode.NOT_FOUND,
                f"Movie with id {request.movie_id} not found"
            )

        review = ReviewService.add_review(
            movie_id=request.movie_id,
            user_id=request.user_id,
            text=request.text
        )

        return self._to_review_response(review, {request.user_id: user})

    def _to_review_response(self, review, users_map):
        user = users_map.get(review.user_id)
        username = user.username if user else ""

        return review_pb2.ReviewResponse(
            id=str(review.id),
            movie_id=review.movie_id,
            user_id=review.user_id,
            username=username,
            text=review.text,
            created_at=review.created_at.isoformat() if review.created_at else ""
        )

    def _get_users_by_ids(self, user_ids):
        """Вызов user-service для получения пользователей по ID"""
        if not user_ids:
            return {}

        try:
            channel = grpc.insecure_channel(USER_SERVICE_ADDRESS)
            stub = user_pb2_grpc.UserServiceStub(channel)

            # Преобразуем строковые ID в int64 для user-service
            valid_ids = []
            for uid in user_ids:
                try:
                    valid_ids.append(int(uid))
                except (ValueError, TypeError):
                    continue

            if not valid_ids:
                return {}

            response = stub.GetUsers(
                GetUsersRequest(
                    ids=valid_ids,
                    limit=len(valid_ids),
                    offset=0
                )
            )

            # Маппинг: строковый ID -> User
            users_map = {str(u.id): u for u in response.users}
            return users_map

        except grpc.RpcError as e:
            print(f"gRPC error calling user-service: {e}")
            return {}

    def _movie_exists(self, movie_id: str) -> bool:
        """Проверка существования фильма через movie-service"""
        try:
            channel = grpc.insecure_channel(MOVIE_SERVICE_ADDRESS)
            stub = movie_pb2_grpc.MovieServiceStub(channel)

            response = stub.MovieExists(MovieIdRequest(id=movie_id))
            return response.exists

        except grpc.RpcError as e:
            print(f"gRPC error calling movie-service: {e}")
            return False

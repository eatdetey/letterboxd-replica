from reviews.models import Review


class ReviewService:

    @staticmethod
    def get_reviews(movie_id: str):
        """Получить все отзывы для фильма"""
        return Review.objects.filter(movie_id=movie_id).order_by('-created_at')

    @staticmethod
    def add_review(movie_id: str, user_id: str, text: str):
        """Создать новый отзыв"""
        return Review.objects.create(
            movie_id=movie_id,
            user_id=user_id,
            text=text
        )

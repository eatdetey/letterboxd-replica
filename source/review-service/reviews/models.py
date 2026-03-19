from django.db import models


class Review(models.Model):
    movie_id = models.CharField(
        max_length=255,
        verbose_name='ID фильма'
    )
    user_id = models.CharField(
        max_length=255,
        verbose_name='ID пользователя'
    )
    text = models.CharField(
        max_length=1000,
        verbose_name='Текст отзыва'
    )
    created_at = models.DateTimeField(
        auto_now_add=True,
        verbose_name='Дата создания'
    )

    class Meta:
        verbose_name = 'Отзыв'
        verbose_name_plural = 'Отзывы'
        ordering = ['-created_at']
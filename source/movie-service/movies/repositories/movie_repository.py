from movies.models import Movie


class MovieRepository:

    @staticmethod
    def get_movie(movie_id):
        return Movie.objects.prefetch_related('genres', 'moviecast_set__person').get(id=movie_id)

    @staticmethod
    def list_movies(limit=10, offset=0):
        return Movie.objects.prefetch_related('genres', 'moviecast_set__person').all()[offset:offset+limit]

    @staticmethod
    def search_movies(query, limit=10, offset=0):
        return Movie.objects.prefetch_related('genres', 'moviecast_set__person') \
                    .filter(title__icontains=query)[offset:offset+limit]

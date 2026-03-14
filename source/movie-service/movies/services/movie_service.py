from movies.repositories.movie_repository import MovieRepository


class MovieService:

    @staticmethod
    def get_movie(movie_id):
        return MovieRepository.get_movie(movie_id)

    @staticmethod
    def list_movies(limit=10, offset=0):
        return MovieRepository.list_movies(limit, offset)

    @staticmethod
    def search_movies(query, limit=10, offset=0):
        return MovieRepository.search_movies(query, limit, offset)

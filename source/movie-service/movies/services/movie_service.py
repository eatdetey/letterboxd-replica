from movies.models import Movie


class MovieService:
    @staticmethod
    def build_queryset(search_query=None, genre=None, ids=None):
        queryset = Movie.objects.prefetch_related("genres")

        if search_query:
            queryset = queryset.filter(title__icontains=search_query)

        if genre:
            queryset = queryset.filter(genres__name=genre)

        # Preserve semantic difference between "no ids filter" (None)
        # and "ids filter resolved to empty set" ([]).
        if ids is not None:
            queryset = queryset.filter(id__in=ids)

        return queryset.distinct().order_by("id")

    @staticmethod
    def get_movies(limit, offset, search_query=None, genre=None, ids=None):
        queryset = MovieService.build_queryset(
            search_query=search_query,
            genre=genre,
            ids=ids,
        )
        total = queryset.count()
        queryset = queryset[offset:offset + limit]

        return queryset, total

    @staticmethod
    def get_movie_ids(search_query=None, genre=None, ids=None):
        queryset = MovieService.build_queryset(
            search_query=search_query,
            genre=genre,
            ids=ids,
        )
        return list(queryset.values_list("id", flat=True))

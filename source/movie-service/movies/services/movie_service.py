from movies.models import Movie

class MovieService:

    @staticmethod
    def get_movies(limit, offset, search_query=None, genre=None, ids=None):
        queryset = Movie.objects.prefetch_related("genres")

        if ids:
            queryset = queryset.filter(id__in=ids)
        else:
            if search_query:
                queryset = queryset.filter(title__icontains=search_query)

            if genre:
                queryset = queryset.filter(genres__name=genre)

        total = queryset.count()
        queryset = queryset[offset:offset + limit]

        return queryset, total

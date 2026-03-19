import { httpClient } from '@shared/api/httpClient'
import type { MovieReviewsResponseDto } from '../model/types'

export const reviewsApi = {
  getMovieReviews(movieId: string): Promise<MovieReviewsResponseDto> {
    return httpClient.get<MovieReviewsResponseDto>(`/v1/movies/${movieId}/reviews`)
  },
}

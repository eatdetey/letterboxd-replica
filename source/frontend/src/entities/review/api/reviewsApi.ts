import { httpClient } from '@shared/api/httpClient'
import type {
  CreateMovieReviewPayloadDto,
  CreateMovieReviewResponseDto,
  MovieReviewsResponseDto,
} from '../model/types'

function isNotFound(error: unknown) {
  return error instanceof Error && error.message.includes('HTTP 404')
}

export const reviewsApi = {
  async getMovieReviews(movieId: string): Promise<MovieReviewsResponseDto> {
    try {
      return await httpClient.post<MovieReviewsResponseDto>(`/v1/movies/${movieId}/reviews/search`, {})
    } catch (error) {
      if (!isNotFound(error)) {
        throw error
      }

      return httpClient.get<MovieReviewsResponseDto>(`/v1/movies/${movieId}/reviews`)
    }
  },

  createMovieReview(
    movieId: string,
    payload: CreateMovieReviewPayloadDto,
  ): Promise<CreateMovieReviewResponseDto> {
    return httpClient.post<CreateMovieReviewResponseDto>(`/v1/movies/${movieId}/reviews`, payload)
  },
}

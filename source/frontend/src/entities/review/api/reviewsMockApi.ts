import { getMock } from '@shared/api/mockClient'
import { authTokenStorage } from '@shared/api/authTokenStorage'
import type {
  CreateMovieReviewPayloadDto,
  CreateMovieReviewResponseDto,
  MovieReviewsResponseDto,
} from '../model/types'

export const reviewsMockApi = {
  async getMovieReviews(movieId: string): Promise<MovieReviewsResponseDto> {
    try {
      return await getMock<MovieReviewsResponseDto>(`v1/movies/${movieId}/reviews/response.json`)
    } catch {
      return { items: [] }
    }
  },

  async createMovieReview(
    _movieId: string,
    payload: CreateMovieReviewPayloadDto,
  ): Promise<CreateMovieReviewResponseDto> {
    const user = authTokenStorage.getUser<{ id?: string; username?: string }>()

    return {
      review: {
        id: `mock-review-${Date.now()}`,
        user_id: String(user?.id ?? 'current-user'),
        username: user?.username ?? 'you',
        text: payload.text,
        created_at: new Date().toISOString(),
      },
    }
  },
}

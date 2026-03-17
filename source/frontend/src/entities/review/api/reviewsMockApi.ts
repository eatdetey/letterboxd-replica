import { getMock } from '@shared/api/mockClient'
import type { MovieReviewsResponseDto } from '../model/types'

export const reviewsMockApi = {
  async getMovieReviews(movieId: string): Promise<MovieReviewsResponseDto> {
    try {
      return await getMock<MovieReviewsResponseDto>(`v1/movies/${movieId}/reviews/response.json`)
    } catch {
      return { items: [] }
    }
  },
}

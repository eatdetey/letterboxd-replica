import { getMock } from '@shared/api/mockClient'
import type { MoviesListResponseDto } from '../model/types'

export const moviesMockApi = {
  getMovies(): Promise<MoviesListResponseDto> {
    return getMock<MoviesListResponseDto>('v1/movies/response.json')
  },
}

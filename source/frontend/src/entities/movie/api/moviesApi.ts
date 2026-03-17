import { httpClient } from '@shared/api/httpClient'
import type { MoviesListResponseDto } from '../model/types'

export const moviesApi = {
  getMovies(): Promise<MoviesListResponseDto> {
    return httpClient.get<MoviesListResponseDto>('/v1/movies')
  },
}

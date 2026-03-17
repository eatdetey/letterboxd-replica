import { httpClient } from '@shared/api/httpClient'
import type { MovieDetailsDto, MoviesListResponseDto } from '../model/types'

export const moviesApi = {
  getMovies(): Promise<MoviesListResponseDto> {
    return httpClient.get<MoviesListResponseDto>('/v1/movies')
  },
  getMovie(id: string): Promise<MovieDetailsDto> {
    return httpClient.get<MovieDetailsDto>(`/v1/movies/${id}`)
  },
}

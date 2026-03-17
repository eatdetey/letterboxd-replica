// TODO
export interface PlaylistDto {
  id: string
  name: string
}

// TODO
export interface MovieDto {
  id: string
  title: string
  description: string
  release_year: number
  genres: string[]
  playlists?: PlaylistDto[]
  poster_url?: string
  backdrop_url?: string
  rating?: number
  is_featured?: boolean
  is_popular?: boolean
}

// TODO
export interface MoviesListResponseDto {
  items: MovieDto[]
  total: number
}

// TODO
export interface MovieDetailsDto extends MovieDto {
  playlists: PlaylistDto[]
  director?: string
  cast?: string[]
  runtime_minutes?: number
  country?: string
}

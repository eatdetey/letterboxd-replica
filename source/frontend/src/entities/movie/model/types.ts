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
}

// TODO
export interface MoviesListResponseDto {
  items: MovieDto[]
  total: number
}

// TODO
export interface MovieDetailsDto extends MovieDto {
  playlists: PlaylistDto[]
}

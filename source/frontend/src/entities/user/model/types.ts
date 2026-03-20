export interface MoviePreviewDto {
  id: string
  title: string
  release_year: number
  poster_url?: string | null
}

export interface UserStatsDto {
  filmsWatched: number
  reviewsCount: number
  listsCount: number
  averageRating: number | null
  followersCount: number
  followingCount: number
  likesCount: number
}

export interface UserProfileDto {
  id: string
  username: string
  email?: string
  role?: 'user' | 'admin'
  avatar_url?: string | null
  cover_url?: string | null
  bio?: string | null
  tags?: string[]
  stats?: Partial<UserStatsDto>
}

export interface UserListDto {
  id: string
  name: string
  description?: string | null
  is_public?: boolean
  movies_count?: number
  movies?: MoviePreviewDto[]
  preview_movies?: MoviePreviewDto[]
}

export interface UserReviewDto {
  id: string
  movie_id: string
  username: string
  text: string
  created_at: string
  movie?: MoviePreviewDto
}

export interface UserProfileResponseDto {
  user: UserProfileDto
  lists: UserListDto[]
  reviews: UserReviewDto[]
  recentRatings: MoviePreviewDto[]
  watchlist: MoviePreviewDto[]
}
// TODO
export interface MovieReviewDto {
  id: string
  user_id: string
  username: string
  text: string
  created_at: string
}

// TODO
export interface MovieReviewsResponseDto {
  items: MovieReviewDto[]
}

// TODO
export interface CreateMovieReviewPayloadDto {
  text: string
}

// TODO
export interface CreateMovieReviewResponseDto {
  review: MovieReviewDto
}

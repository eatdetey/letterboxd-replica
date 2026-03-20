export interface UserProfileDto {
  id: string
  username: string
  email?: string
  avatar_url?: string
  role?: 'user' | 'admin'
}

export interface UserProfileResponseDto {
  user: UserProfileDto
}

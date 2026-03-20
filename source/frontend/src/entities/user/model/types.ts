export interface UserProfileDto {
  id: string
  username: string
  email?: string
  role?: 'user' | 'admin'
}

export interface UserProfileResponseDto {
  user: UserProfileDto
}

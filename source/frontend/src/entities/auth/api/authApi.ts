import { httpClient } from '@shared/api/httpClient'

export type AuthUserDto = {
  id: string | number
  username: string
  email?: string
  bio?: string
  avatar_url?: string
  status?: string
  role: string
}

export type LoginRequestDto = {
  username: string
  password: string
}

export type RegisterRequestDto = {
  username: string
  password: string
  email: string
}

export type RefreshRequestDto = {
  refresh_token?: string
}

export type LoginResponseDto = {
  user: AuthUserDto
  access_token: string
}

export type RegisterResponseDto = {
  user: AuthUserDto
  access_token: string
}

export type RefreshResponseDto = {
  access_token: string
}

export const authApi = {
  login(payload: LoginRequestDto): Promise<LoginResponseDto> {
    return httpClient.post<LoginResponseDto>('/v1/auth/login', payload)
  },

  register(payload: RegisterRequestDto): Promise<RegisterResponseDto> {
    return httpClient.post<RegisterResponseDto>('/v1/auth/register', payload)
  },

  refresh(payload?: RefreshRequestDto): Promise<RefreshResponseDto> {
    return httpClient.post<RefreshResponseDto>('/v1/auth/refresh', payload)
  },
}

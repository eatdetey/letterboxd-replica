import { httpClient } from '@shared/api/httpClient'

export type AuthUserDto = {
  id: string
  username: string
  role: 'user' | 'admin'
}

export type AuthTokensDto = {
  access: string
  refresh: string
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
  refresh_token: string
}

export type LoginResponseDto = AuthUserDto & {
  tokens: AuthTokensDto
}

export type RegisterResponseDto = AuthUserDto & {
  email: string
  tokens: AuthTokensDto
}

export type RefreshResponseDto = AuthTokensDto

export const authApi = {
  login(payload: LoginRequestDto): Promise<LoginResponseDto> {
    return httpClient.post<LoginResponseDto>('/v1/auth/login', payload)
  },

  register(payload: RegisterRequestDto): Promise<RegisterResponseDto> {
    return httpClient.post<RegisterResponseDto>('/v1/auth/register', payload)
  },

  refresh(payload: RefreshRequestDto): Promise<RefreshResponseDto> {
    return httpClient.post<RefreshResponseDto>('/v1/auth/refresh', payload)
  },
}
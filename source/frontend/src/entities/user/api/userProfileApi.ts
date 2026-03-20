import { httpClient } from '@shared/api/httpClient'
import type {
  UserProfileResponseDto,
} from '../model/types'

type GetProfileParams = {
  idOrUsername: string
}

export const userProfileApi = {
  getProfile(idOrUsername: string): Promise<UserProfileResponseDto> {
    const encoded = encodeURIComponent(idOrUsername)
    return httpClient.get<UserProfileResponseDto>(`/v1/users/${encoded}/profile`)
  },
}
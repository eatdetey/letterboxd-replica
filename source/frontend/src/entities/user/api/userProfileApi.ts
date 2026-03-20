import type {
  UserProfileResponseDto,
} from '../model/types'


export const userProfileApi = {
  getProfile(idOrUsername: string): Promise<UserProfileResponseDto> {
    return Promise.resolve({
      user: {
        id: idOrUsername,
        username: idOrUsername,
      },
    })
  },
}

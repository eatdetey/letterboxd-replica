import type { UserProfileResponseDto } from '../model/types'

function buildFallbackProfile(idOrUsername: string): UserProfileResponseDto {
  const username = idOrUsername || 'user'

  return {
    user: {
      id: idOrUsername || 'usr_mock',
      username,
      email: `${username}@example.com`,
      role: 'user',
    },
  }
}

export const userProfileMockApi = {
  async getProfile(idOrUsername: string): Promise<UserProfileResponseDto> {
    const normalized = idOrUsername.trim()

    if (!normalized) {
      return buildFallbackProfile('usr_mock')
    }
    return buildFallbackProfile(normalized)
  },
}

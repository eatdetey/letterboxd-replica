import { getMock } from '@shared/api/mockClient'
import type { PlaylistsResponseDto } from '../model/types'

export const playlistsMockApi = {
  getPlaylists(): Promise<PlaylistsResponseDto> {
    return getMock<PlaylistsResponseDto>('v1/playlists/response.json')
  },
}

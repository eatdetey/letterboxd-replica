import { httpClient } from '@shared/api/httpClient'
import type { PlaylistsResponseDto } from '../model/types'

export const playlistsApi = {
  getPlaylists(): Promise<PlaylistsResponseDto> {
    return httpClient.get<PlaylistsResponseDto>('/v1/playlists')
  },
}

import { getMock } from '@shared/api/mockClient'
import type { PlaylistDto, PlaylistsResponseDto } from '../model/types'

export const playlistsMockApi = {
  getPlaylists(): Promise<PlaylistsResponseDto> {
    return getMock<PlaylistsResponseDto>('v1/playlists/response.json')
  },

  async getPlaylist(id: string): Promise<PlaylistDto> {
    const response = await getMock<PlaylistsResponseDto>('v1/playlists/response.json')
    const playlist = response.items.find((item) => item.id === id)

    if (!playlist) {
      throw new Error(`Playlist with id "${id}" not found in mocks`)
    }

    return playlist
  },
}

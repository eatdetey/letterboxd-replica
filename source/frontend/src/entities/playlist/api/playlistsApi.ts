import { httpClient } from '@shared/api/httpClient'
import type { PlaylistDto, PlaylistsResponseDto } from '../model/types'

type PlaylistResponseDto = {
  playlist: PlaylistDto
}

type PlaylistPayloadDto = {
  name: string
}

type PlaylistMoviePayloadDto = {
  movie_id: string
}

function isNotFound(error: unknown) {
  return error instanceof Error && error.message.includes('HTTP 404')
}

export const playlistsApi = {
  async getPlaylists(): Promise<PlaylistsResponseDto> {
    try {
      return await httpClient.post<PlaylistsResponseDto>('/v1/playlists/search', {})
    } catch (error) {
      if (!isNotFound(error)) {
        throw error
      }

      return httpClient.get<PlaylistsResponseDto>('/v1/playlists')
    }
  },

  async getPlaylist(id: string): Promise<PlaylistDto> {
    try {
      const response = await httpClient.get<PlaylistResponseDto>(`/v1/playlists/${id}`)
      return response.playlist
    } catch (error) {
      if (!isNotFound(error)) {
        throw error
      }

      const response = await this.getPlaylists()
      const playlist = response.items.find((item) => item.id === id)

      if (!playlist) {
        throw new Error(`HTTP 404: {"error":"playlist not found"}`)
      }

      return playlist
    }
  },

  async createPlaylist(name: string): Promise<PlaylistDto> {
    const response = await httpClient.post<PlaylistResponseDto>('/v1/playlists', {
      name,
    } satisfies PlaylistPayloadDto)

    return response.playlist
  },

  async renamePlaylist(id: string, name: string): Promise<PlaylistDto> {
    const response = await httpClient.put<PlaylistResponseDto>(`/v1/playlists/${id}`, {
      name,
    } satisfies PlaylistPayloadDto)

    return response.playlist
  },

  deletePlaylist(id: string): Promise<void> {
    return httpClient.delete<void>(`/v1/playlists/${id}`)
  },

  addMovieToPlaylist(id: string, movieId: string): Promise<void> {
    return httpClient.post<void>(`/v1/playlists/${id}/movies`, {
      movie_id: movieId,
    } satisfies PlaylistMoviePayloadDto)
  },

  removeMovieFromPlaylist(id: string, movieId: string): Promise<void> {
    return httpClient.delete<void>(`/v1/playlists/${id}/movies/${movieId}`)
  },
}

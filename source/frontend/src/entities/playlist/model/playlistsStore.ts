import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { env } from '@shared/config/env'
import { playlistsApi } from '../api/playlistsApi'
import { playlistsMockApi } from '../api/playlistsMockApi'
import type { PlaylistDto } from './types'

export const usePlaylistsStore = defineStore('playlists', () => {
  const items = ref<PlaylistDto[]>([])
  const isLoading = ref(false)
  const isMutating = ref(false)
  const error = ref<string | null>(null)

  const totalMovies = computed(() => {
    return items.value.reduce((sum, playlist) => sum + playlist.movies_count, 0)
  })

  async function loadPlaylists() {
    isLoading.value = true
    error.value = null

    try {
      const api = env.usePlaylistMocks ? playlistsMockApi : playlistsApi
      const response = await api.getPlaylists()
      items.value = response.items
    } catch (loadError) {
      items.value = []
      error.value = loadError instanceof Error ? loadError.message : 'Unknown error'
    } finally {
      isLoading.value = false
    }
  }

  async function loadPlaylist(playlistId: string) {
    const existingPlaylist = items.value.find((playlist) => playlist.id === playlistId)
    if (existingPlaylist) {
      return existingPlaylist
    }

    isLoading.value = true
    error.value = null

    try {
      const api = env.usePlaylistMocks ? playlistsMockApi : playlistsApi
      const playlist = await api.getPlaylist(playlistId)
      upsertPlaylist(playlist)
      return playlist
    } catch (loadError) {
      error.value = loadError instanceof Error ? loadError.message : 'Unknown error'
      return null
    } finally {
      isLoading.value = false
    }
  }

  async function createPlaylist(name: string) {
    isMutating.value = true
    error.value = null

    try {
      const playlist = await playlistsApi.createPlaylist(name)
      upsertPlaylist(playlist, { prepend: true })
      return playlist
    } catch (mutationError) {
      error.value = mutationError instanceof Error ? mutationError.message : 'Unknown error'
      throw mutationError
    } finally {
      isMutating.value = false
    }
  }

  async function renamePlaylist(playlistId: string, name: string) {
    isMutating.value = true
    error.value = null

    try {
      const playlist = await playlistsApi.renamePlaylist(playlistId, name)
      upsertPlaylist(playlist)
      return playlist
    } catch (mutationError) {
      error.value = mutationError instanceof Error ? mutationError.message : 'Unknown error'
      throw mutationError
    } finally {
      isMutating.value = false
    }
  }

  async function deletePlaylist(playlistId: string) {
    isMutating.value = true
    error.value = null

    try {
      await playlistsApi.deletePlaylist(playlistId)
      items.value = items.value.filter((playlist) => playlist.id !== playlistId)
    } catch (mutationError) {
      error.value = mutationError instanceof Error ? mutationError.message : 'Unknown error'
      throw mutationError
    } finally {
      isMutating.value = false
    }
  }

  function updateMoviesCount(playlistId: string, delta: number) {
    const playlist = items.value.find((item) => item.id === playlistId)
    if (!playlist) {
      return
    }

    playlist.movies_count = Math.max(0, playlist.movies_count + delta)
  }

  function upsertPlaylist(playlist: PlaylistDto, options: { prepend?: boolean } = {}) {
    const index = items.value.findIndex((item) => item.id === playlist.id)

    if (index >= 0) {
      items.value[index] = playlist
      return
    }

    if (options.prepend) {
      items.value.unshift(playlist)
      return
    }

    items.value.push(playlist)
  }

  return {
    items,
    isLoading,
    isMutating,
    error,
    totalMovies,
    loadPlaylists,
    loadPlaylist,
    createPlaylist,
    renamePlaylist,
    deletePlaylist,
    updateMoviesCount,
  }
})

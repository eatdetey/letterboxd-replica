<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import MovieCard from '~/entities/movie/ui/MovieCard.vue'
import { usePlaylistMoviesStore } from '~/entities/movie/model/playlistMoviesStore'
import { playlistsApi } from '~/entities/playlist/api/playlistsApi'
import { usePlaylistsStore } from '~/entities/playlist/model/playlistsStore'
import { useInfiniteScroll } from '@shared/lib/useInfiniteScroll'

const route = useRoute()
const router = useRouter()
const playlistsStore = usePlaylistsStore()
const playlistMoviesStore = usePlaylistMoviesStore()
const isRenameDialogOpen = ref(false)
const isDeleteDialogOpen = ref(false)
const playlistName = ref('')
const actionError = ref<string | null>(null)
const activeMovieId = ref<string | null>(null)
const loadMoreSentinel = ref<HTMLElement | null>(null)

const playlistId = computed(() => String(route.params.id ?? ''))
const playlist = computed(() => {
  return playlistsStore.items.find((item) => item.id === playlistId.value) || null
})
const showLoadMoreAnchor = computed(() => playlistMoviesStore.hasMore && !playlistMoviesStore.error)
const canLoadMoreMovies = computed(() => {
  return showLoadMoreAnchor.value && !playlistMoviesStore.isLoading && !playlistMoviesStore.isLoadingMore
})

async function loadPlaylistPage(id: string) {
  if (!id) {
    return
  }

  const currentPlaylist = await playlistsStore.loadPlaylist(id)

  if (!currentPlaylist) {
    playlistMoviesStore.reset()
    return
  }

  if (currentPlaylist.movies_count === 0) {
    playlistMoviesStore.reset()
    return
  }

  await playlistMoviesStore.loadPlaylistMovies(id)
}

async function renamePlaylist() {
  if (!playlist.value) {
    return
  }

  const name = playlistName.value.trim()
  if (!name) {
    actionError.value = 'Playlist name is required'
    return
  }

  actionError.value = null

  try {
    await playlistsStore.renamePlaylist(playlist.value.id, name)
    isRenameDialogOpen.value = false
  } catch (error) {
    actionError.value = error instanceof Error ? error.message : 'Failed to rename playlist'
  }
}

async function deletePlaylist() {
  if (!playlist.value) {
    return
  }

  actionError.value = null

  try {
    await playlistsStore.deletePlaylist(playlist.value.id)
    isDeleteDialogOpen.value = false
    await router.replace({ name: 'playlists' })
  } catch (error) {
    actionError.value = error instanceof Error ? error.message : 'Failed to delete playlist'
  }
}

async function removeMovie(movieId: string) {
  if (!playlist.value) {
    return
  }

  activeMovieId.value = movieId
  actionError.value = null

  try {
    await playlistsApi.removeMovieFromPlaylist(playlist.value.id, movieId)
    playlistsStore.updateMoviesCount(playlist.value.id, -1)
    await playlistMoviesStore.loadPlaylistMovies(playlist.value.id)
  } catch (error) {
    actionError.value = error instanceof Error ? error.message : 'Failed to remove movie from playlist'
  } finally {
    activeMovieId.value = null
  }
}

onMounted(() => {
  void loadPlaylistPage(playlistId.value)
})

watch(playlistId, (id) => {
  void loadPlaylistPage(id)
})

watch(playlist, (value) => {
  playlistName.value = value?.name ?? ''
})

useInfiniteScroll({
  target: loadMoreSentinel,
  enabled: canLoadMoreMovies,
  onLoadMore: () => playlistMoviesStore.loadMorePlaylistMovies(),
})
</script>

<template>
  <div class="playlist-page">
    <section class="playlist-hero">
      <v-container class="playlist-hero__container">
        <router-link class="playlist-hero__back-link" :to="{ name: 'playlists' }">
          <v-icon icon="mdi-arrow-left" size="18" />
          <span>Back to lists</span>
        </router-link>
        <div class="playlist-hero__badge">Playlist</div>
        <h1 class="playlist-hero__title">{{ playlist?.name || 'Playlist' }}</h1>
        <div class="playlist-hero__meta" v-if="playlist">
          <span>{{ playlist.movies_count }} movies</span>
        </div>
        <div v-if="playlist" class="playlist-hero__actions">
          <v-btn variant="outlined" color="white" @click="isRenameDialogOpen = true">Rename list</v-btn>
          <v-btn variant="outlined" color="white" @click="isDeleteDialogOpen = true">Delete list</v-btn>
        </div>
      </v-container>
    </section>

    <section class="playlist-content">
      <v-container>
        <v-alert
          v-if="actionError"
          type="error"
          variant="tonal"
          class="playlist-page__alert"
        >
          {{ actionError }}
        </v-alert>

        <div v-if="playlistMoviesStore.isLoading" class="playlist-page__state">
          <v-progress-circular indeterminate color="white" />
        </div>

        <v-alert
          v-else-if="playlistsStore.error || playlistMoviesStore.error"
          type="error"
          variant="tonal"
          class="playlist-page__alert"
        >
          {{ playlistsStore.error || playlistMoviesStore.error }}
        </v-alert>

        <template v-else-if="playlist">
          <div class="playlist-content__header">
            <h2 class="playlist-content__title">Movies in this list</h2>
            <div class="playlist-content__count">{{ playlistMoviesStore.total }} items</div>
          </div>

          <v-row v-if="playlistMoviesStore.items.length">
            <v-col
              v-for="movie in playlistMoviesStore.items"
              :key="movie.id"
              cols="12"
              sm="6"
              md="4"
              lg="3"
            >
              <div class="playlist-movie">
                <MovieCard :movie="movie" />
                <v-btn
                  variant="outlined"
                  color="white"
                  block
                  class="playlist-movie__remove"
                  :loading="activeMovieId === movie.id"
                  @click="removeMovie(movie.id)"
                >
                  Remove from list
                </v-btn>
              </div>
            </v-col>
          </v-row>

          <div
            v-if="playlistMoviesStore.items.length && showLoadMoreAnchor"
            ref="loadMoreSentinel"
            class="playlist-content__load-more"
          >
            <v-progress-circular
              v-if="playlistMoviesStore.isLoadingMore"
              indeterminate
              size="24"
              width="2"
              color="white"
            />
          </div>

          <div v-if="!playlistMoviesStore.items.length" class="playlist-page__empty">
            В этом списке пока нет фильмов.
          </div>
        </template>

        <div v-else class="playlist-page__empty">
          Список не найден.
        </div>
      </v-container>
    </section>

    <v-dialog v-model="isRenameDialogOpen" max-width="480">
      <v-card class="playlist-dialog" elevation="0">
        <h2 class="playlist-dialog__title">Rename playlist</h2>
        <v-text-field
          v-model="playlistName"
          label="Playlist name"
          variant="outlined"
          density="comfortable"
          hide-details
          autofocus
        />
        <div class="playlist-dialog__actions">
          <v-btn variant="text" @click="isRenameDialogOpen = false">Cancel</v-btn>
          <v-btn
            color="white"
            variant="flat"
            class="playlist-dialog__submit"
            :loading="playlistsStore.isMutating"
            @click="renamePlaylist"
          >
            Save
          </v-btn>
        </div>
      </v-card>
    </v-dialog>

    <v-dialog v-model="isDeleteDialogOpen" max-width="460">
      <v-card class="playlist-dialog" elevation="0">
        <h2 class="playlist-dialog__title">Delete playlist</h2>
        <p class="playlist-dialog__description">
          This action will remove the playlist and its movie links for the current user.
        </p>
        <div class="playlist-dialog__actions">
          <v-btn variant="text" @click="isDeleteDialogOpen = false">Cancel</v-btn>
          <v-btn
            color="white"
            variant="flat"
            class="playlist-dialog__submit"
            :loading="playlistsStore.isMutating"
            @click="deletePlaylist"
          >
            Delete
          </v-btn>
        </div>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.playlist-page {
  display: grid;
  gap: 40px;
  padding-bottom: 72px;
}

.playlist-hero {
  margin: 24px;
  border-radius: 28px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background:
    radial-gradient(circle at top left, rgba(240, 78, 62, 0.18), transparent 38%),
    radial-gradient(circle at right, rgba(244, 211, 94, 0.12), transparent 32%),
    linear-gradient(135deg, #121418, #0d1014 62%, #141920);
}

.playlist-hero__container {
  padding: 64px 32px;
  display: grid;
  gap: 14px;
}

.playlist-hero__badge {
  width: fit-content;
  padding: 4px 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);
  color: var(--text-primary);
  font-size: 0.82rem;
  text-transform: uppercase;
  letter-spacing: 0.14em;
}

.playlist-hero__back-link {
  width: fit-content;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 1rem;
  font-weight: 600;
}

.playlist-hero__back-link:hover {
  color: var(--text-primary);
}

.playlist-hero__title {
  margin: 0;
  color: var(--text-primary);
  font-size: clamp(2.1rem, 5vw, 4rem);
  line-height: 1;
}

.playlist-hero__meta {
  color: var(--text-secondary);
}

.playlist-hero__actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.playlist-content {
  padding: 0 12px;
}

.playlist-content__header {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.playlist-content__title {
  margin: 0;
  color: var(--text-text);
  opacity: 1;
  font-size: 1.6rem;
}

.playlist-content__count {
  color: var(--text-muted);
}

.playlist-content__load-more {
  min-height: 80px;
  display: grid;
  place-items: center;
}

.playlist-page__state,
.playlist-page__empty {
  min-height: 260px;
  display: grid;
  place-items: center;
  color: var(--text-muted);
}

.playlist-page__alert {
  border-radius: 20px;
}

.playlist-movie {
  display: grid;
  gap: 12px;
  height: 100%;
}

.playlist-movie__remove,
.playlist-dialog__submit {
  color: #0d0f12 !important;
}

.playlist-dialog {
  background: var(--surface);
  color: var(--text-primary);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  padding: 24px;
  display: grid;
  gap: 16px;
}

.playlist-dialog__title {
  margin: 0;
  font-size: 1.35rem;
}

.playlist-dialog__description {
  margin: 0;
  color: var(--text-secondary);
  line-height: 1.6;
}

.playlist-dialog__actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 960px) {
  .playlist-hero {
    margin: 16px;
  }

  .playlist-hero__container {
    padding: 52px 24px;
  }
}

@media (max-width: 600px) {
  .playlist-hero__container {
    padding: 40px 16px;
  }
}
</style>

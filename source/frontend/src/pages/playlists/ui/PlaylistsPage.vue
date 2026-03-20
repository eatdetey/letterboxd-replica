<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { usePlaylistsStore } from '~/entities/playlist/model/playlistsStore'
import PlaylistCard from '~/entities/playlist/ui/PlaylistCard.vue'
import { useInfiniteScroll } from '@shared/lib/useInfiniteScroll'

const PLAYLISTS_BATCH_SIZE = 12

const playlistsBackdropUrl =
  'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?auto=format&fit=crop&w=1600&q=80'

const router = useRouter()
const playlistsStore = usePlaylistsStore()
const isCreateDialogOpen = ref(false)
const playlistName = ref('')
const createError = ref<string | null>(null)
const visiblePlaylistsCount = ref(PLAYLISTS_BATCH_SIZE)
const loadMoreSentinel = ref<HTMLElement | null>(null)

onMounted(() => {
  visiblePlaylistsCount.value = PLAYLISTS_BATCH_SIZE
  void playlistsStore.loadPlaylists()
})

const playlists = computed(() => playlistsStore.items)
const visiblePlaylists = computed(() => playlists.value.slice(0, visiblePlaylistsCount.value))
const totalMovies = computed(() => playlistsStore.totalMovies)
const hasMorePlaylists = computed(() => playlists.value.length > visiblePlaylistsCount.value)
const canLoadMorePlaylists = computed(() => hasMorePlaylists.value && !playlistsStore.isLoading)

function loadMorePlaylists() {
  if (!canLoadMorePlaylists.value) {
    return
  }

  visiblePlaylistsCount.value += PLAYLISTS_BATCH_SIZE
}

watch(() => playlists.value.length, (length, previousLength) => {
  if (length === 0) {
    visiblePlaylistsCount.value = PLAYLISTS_BATCH_SIZE
    return
  }

  if (previousLength === 0) {
    visiblePlaylistsCount.value = Math.min(PLAYLISTS_BATCH_SIZE, length)
    return
  }

  if (visiblePlaylistsCount.value > length) {
    visiblePlaylistsCount.value = length
  }
})

useInfiniteScroll({
  target: loadMoreSentinel,
  enabled: canLoadMorePlaylists,
  onLoadMore: loadMorePlaylists,
})

async function createPlaylist() {
  const name = playlistName.value.trim()
  if (!name) {
    createError.value = 'Playlist name is required'
    return
  }

  createError.value = null

  try {
    const playlist = await playlistsStore.createPlaylist(name)
    playlistName.value = ''
    isCreateDialogOpen.value = false
    await router.push({ name: 'playlist', params: { id: playlist.id } })
  } catch (error) {
    createError.value = error instanceof Error ? error.message : 'Failed to create playlist'
  }
}
</script>

<template>
  <div class="playlists-page">
    <section class="playlists-hero">
      <div class="playlists-hero__backdrop" :style="{ backgroundImage: `url(${playlistsBackdropUrl})` }"></div>
      <v-container class="playlists-hero__container">
        <div class="playlists-hero__badge">Playlists</div>
        <h1 class="playlists-hero__title">Your film lists, all in one place.</h1>
        <p class="playlists-hero__description">
          Browse saved lists and keep track of the films you grouped together.
        </p>
        <div class="playlists-hero__meta">
          <span>{{ playlists.length }} lists</span>
          <span class="dot">•</span>
          <span>{{ totalMovies }} movies in total</span>
        </div>
      </v-container>
    </section>

    <section class="playlists-section">
      <v-container>
        <div class="playlists-section__header">
          <h2 class="playlists-section__title">All lists</h2>
          <v-btn color="white" variant="flat" class="playlists-section__create" @click="isCreateDialogOpen = true">
            Create list
          </v-btn>
        </div>

        <div v-if="playlistsStore.isLoading" class="playlists-page__state">
          <v-progress-circular indeterminate color="white" />
        </div>

        <v-alert
          v-else-if="playlistsStore.error"
          type="error"
          variant="tonal"
          class="playlists-page__alert"
        >
          {{ playlistsStore.error }}
        </v-alert>

        <v-row v-else-if="playlists.length">
          <v-col
            v-for="playlist in visiblePlaylists"
            :key="playlist.id"
            cols="12"
            sm="6"
            lg="4"
          >
            <PlaylistCard :playlist="playlist" />
          </v-col>
        </v-row>

        <div
          v-if="playlists.length && hasMorePlaylists"
          ref="loadMoreSentinel"
          class="playlists-page__load-more"
        >
          <v-progress-circular indeterminate size="24" width="2" color="white" />
        </div>

        <div v-if="!playlistsStore.isLoading && !playlistsStore.error && !playlists.length" class="playlists-page__empty">
          Пока нет пользовательских списков.
        </div>
      </v-container>
    </section>

    <v-dialog v-model="isCreateDialogOpen" max-width="480">
      <v-card class="playlists-dialog" elevation="0">
        <h2 class="playlists-dialog__title">Create playlist</h2>

        <v-alert
          v-if="createError"
          type="error"
          variant="tonal"
          class="playlists-dialog__alert"
        >
          {{ createError }}
        </v-alert>

        <v-text-field
          v-model="playlistName"
          label="Playlist name"
          variant="outlined"
          density="comfortable"
          hide-details
          autofocus
        />

        <div class="playlists-dialog__actions">
          <v-btn variant="text" @click="isCreateDialogOpen = false">Cancel</v-btn>
          <v-btn
            color="white"
            variant="flat"
            class="playlists-dialog__submit"
            :loading="playlistsStore.isMutating"
            @click="createPlaylist"
          >
            Create
          </v-btn>
        </div>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.playlists-page {
  display: grid;
  gap: 40px;
  padding-bottom: 72px;
}

.playlists-hero {
  position: relative;
  margin: 24px;
  border-radius: 28px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background:
    radial-gradient(circle at top left, rgba(240, 78, 62, 0.18), transparent 38%),
    radial-gradient(circle at right, rgba(244, 211, 94, 0.12), transparent 32%),
    rgba(13, 16, 20, 0.58);
}

.playlists-hero__backdrop {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center;
  opacity: 1;
}

.playlists-hero__backdrop::after {
  content: '';
  position: absolute;
  inset: 0;
  background:
    linear-gradient(90deg, rgba(11, 13, 16, 0.82) 18%, rgba(11, 13, 16, 0.52) 58%, rgba(11, 13, 16, 0.8)),
    linear-gradient(180deg, rgba(11, 13, 16, 0.12), rgba(11, 13, 16, 0.64));
}

.playlists-hero__container {
  position: relative;
  z-index: 1;
  padding: 72px 32px;
  display: grid;
  gap: 16px;
}

.playlists-hero__badge {
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

.playlists-hero__title {
  margin: 0;
  max-width: 720px;
  color: var(--text-primary);
  font-size: clamp(2.1rem, 5vw, 4rem);
  line-height: 1;
}

.playlists-hero__description {
  margin: 0;
  max-width: 620px;
  color: var(--text-secondary);
  line-height: 1.7;
}

.playlists-hero__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  color: var(--text-muted);
}

.playlists-section {
  padding: 0 12px;
}

.playlists-section__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.playlists-section__title {
  margin: 0;
  color: var(--text-text);
  opacity: 1;
  font-size: 1.6rem;
}

.playlists-page__state,
.playlists-page__empty {
  min-height: 260px;
  display: grid;
  place-items: center;
  color: var(--text-muted);
}

.playlists-page__load-more {
  min-height: 80px;
  display: grid;
  place-items: center;
}

.playlists-page__alert {
  border-radius: 20px;
}

.playlists-section__create,
.playlists-dialog__submit {
  color: #0d0f12 !important;
}

.playlists-dialog {
  background: var(--surface);
  color: var(--text-primary);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  padding: 24px;
  display: grid;
  gap: 16px;
}

.playlists-dialog__title {
  margin: 0;
  font-size: 1.4rem;
}

.playlists-dialog__alert {
  border-radius: 16px;
}

.playlists-dialog__actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 960px) {
  .playlists-hero {
    margin: 16px;
  }

  .playlists-hero__container {
    padding: 56px 24px;
  }
}

@media (max-width: 600px) {
  .playlists-hero__container {
    padding: 42px 16px;
  }
}
</style>

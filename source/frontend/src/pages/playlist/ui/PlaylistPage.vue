<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import MovieCard from '~/entities/movie/ui/MovieCard.vue'
import { usePlaylistMoviesStore } from '~/entities/movie/model/playlistMoviesStore'
import { usePlaylistsStore } from '~/entities/playlist/model/playlistsStore'

const route = useRoute()
const playlistsStore = usePlaylistsStore()
const playlistMoviesStore = usePlaylistMoviesStore()

const playlistId = computed(() => String(route.params.id ?? ''))
const playlist = computed(() => {
  return playlistsStore.items.find((item) => item.id === playlistId.value) || null
})

async function loadPlaylistPage(id: string) {
  if (!id) {
    return
  }

  if (!playlistsStore.items.length) {
    await playlistsStore.loadPlaylists()
  }

  await playlistMoviesStore.loadPlaylistMovies(id)
}

onMounted(() => {
  void loadPlaylistPage(playlistId.value)
})

watch(playlistId, (id) => {
  void loadPlaylistPage(id)
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
      </v-container>
    </section>

    <section class="playlist-content">
      <v-container>
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
              <MovieCard :movie="movie" />
            </v-col>
          </v-row>

          <div v-else class="playlist-page__empty">
            В этом списке пока нет фильмов.
          </div>
        </template>

        <div v-else class="playlist-page__empty">
          Список не найден.
        </div>
      </v-container>
    </section>
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

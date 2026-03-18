<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { usePlaylistsStore } from '~/entities/playlist/model/playlistsStore'
import PlaylistCard from '~/entities/playlist/ui/PlaylistCard.vue'

const playlistsBackdropUrl =
  'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?auto=format&fit=crop&w=1600&q=80'

const playlistsStore = usePlaylistsStore()

onMounted(() => {
  void playlistsStore.loadPlaylists()
})

const playlists = computed(() => playlistsStore.items)
const totalMovies = computed(() => playlistsStore.totalMovies)
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
            v-for="playlist in playlists"
            :key="playlist.id"
            cols="12"
            sm="6"
            lg="4"
          >
            <PlaylistCard :playlist="playlist" />
          </v-col>
        </v-row>

        <div v-else class="playlists-page__empty">
          Пока нет пользовательских списков.
        </div>
      </v-container>
    </section>
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
  opacity: 0.9;
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
  margin-bottom: 24px;
}

.playlists-section__title {
  margin: 0;
  color: var(--text-primary);
  font-size: 1.6rem;
}

.playlists-page__state,
.playlists-page__empty {
  min-height: 260px;
  display: grid;
  place-items: center;
  color: var(--text-muted);
}

.playlists-page__alert {
  border-radius: 20px;
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

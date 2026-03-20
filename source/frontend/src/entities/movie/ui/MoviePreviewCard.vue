<script setup lang="ts">
import type { MoviePreviewDto } from '../model/types'

defineProps<{ movie: MoviePreviewDto }>()
</script>

<template>
  <router-link class="movie-preview-card__link" :to="{ name: 'movie', params: { id: movie.id } }">
    <v-card class="movie-preview-card" elevation="0">
      <div
        class="movie-preview-card__poster"
        :style="{ backgroundImage: movie.poster_url ? `url(${movie.poster_url})` : undefined }"
      ></div>

      <div class="movie-preview-card__body">
        <div class="movie-preview-card__title">{{ movie.title }}</div>
        <div class="movie-preview-card__meta">
          <span>{{ movie.release_year }}</span>
        </div>
      </div>
    </v-card>
  </router-link>
</template>

<style scoped>
.movie-preview-card__link {
  display: block;
  color: inherit;
  text-decoration: none;
}

.movie-preview-card {
  background: var(--surface);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 18px;
  overflow: hidden;
  color: var(--text-primary);
  display: grid;
  grid-template-columns: 72px minmax(0, 1fr);
  gap: 12px;
  align-items: center;
  transition:
    transform 0.2s ease,
    border-color 0.2s ease,
    box-shadow 0.2s ease;
}

.movie-preview-card__link:hover .movie-preview-card {
  transform: translateY(-4px);
  border-color: rgba(244, 211, 94, 0.28);
  box-shadow: 0 18px 42px rgba(0, 0, 0, 0.2);
}

.movie-preview-card__poster {
  aspect-ratio: 2 / 3;
  background-color: #1d1f23;
  background-size: cover;
  background-position: center;
}

.movie-preview-card__body {
  padding-right: 14px;
  display: grid;
  gap: 6px;
}

.movie-preview-card__title {
  font-size: 0.98rem;
  font-weight: 600;
  line-height: 1.2;
}

.movie-preview-card__meta {
  color: var(--text-muted);
  font-size: 0.85rem;
}
</style>
<script setup lang="ts">
import type { MovieDto } from '../model/types'

defineProps<{ movie: MovieDto; variant?: 'compact' | 'default' }>()
</script>

<template>
  <v-card class="movie-card" :class="variant ? `movie-card--${variant}` : 'movie-card--default'" elevation="0">
    <div class="movie-card__poster" :style="{ backgroundImage: movie.poster_url ? `url(${movie.poster_url})` : undefined }">
      <div class="movie-card__overlay">
        <div class="movie-card__rating" v-if="movie.rating">
          {{ movie.rating.toFixed(1) }}
        </div>
      </div>
    </div>
    <div class="movie-card__body">
      <div class="movie-card__title">{{ movie.title }}</div>
      <div class="movie-card__meta">
        <span>{{ movie.release_year }}</span>
        <span class="dot">•</span>
        <span>{{ movie.genres.join(', ') }}</span>
      </div>
      <div class="movie-card__description" v-if="variant !== 'compact'">
        {{ movie.description }}
      </div>
    </div>
  </v-card>
</template>

<style scoped>
.movie-card {
  background: var(--surface);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 18px;
  overflow: hidden;
  color: var(--text-primary);
  display: grid;
  grid-template-rows: auto 1fr;
  min-height: 100%;
}

.movie-card__poster {
  position: relative;
  padding-top: 140%;
  background-color: #1d1f23;
  background-size: cover;
  background-position: center;
}

.movie-card__overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(9, 10, 12, 0) 45%, rgba(9, 10, 12, 0.8));
  display: flex;
  align-items: flex-end;
  justify-content: flex-end;
  padding: 12px;
}

.movie-card__rating {
  background: rgba(12, 13, 15, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: #f4d35e;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 0.85rem;
}

.movie-card__body {
  padding: 16px;
  display: grid;
  gap: 8px;
}

.movie-card__title {
  font-size: 1.05rem;
  font-weight: 600;
}

.movie-card__meta {
  color: var(--text-muted);
  font-size: 0.85rem;
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.movie-card__description {
  font-size: 0.9rem;
  color: var(--text-secondary);
  line-height: 1.5;
}

.movie-card--compact .movie-card__description {
  display: none;
}
</style>

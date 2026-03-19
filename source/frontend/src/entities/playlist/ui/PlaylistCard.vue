<script setup lang="ts">
import { computed } from 'vue'
import type { PlaylistDto } from '../model/types'

const props = defineProps<{ playlist: PlaylistDto }>()

const movieLabel = computed(() => {
  return props.playlist.movies_count === 1 ? 'movie' : 'movies'
})
</script>

<template>
  <router-link class="playlist-card__link" :to="{ name: 'playlist', params: { id: playlist.id } }">
    <v-card class="playlist-card" elevation="0">
      <div class="playlist-card__eyebrow">User list</div>
      <h3 class="playlist-card__title">{{ playlist.name }}</h3>
      <p class="playlist-card__meta">
        {{ playlist.movies_count }} {{ movieLabel }}
      </p>
      <v-btn variant="text" class="playlist-card__action">Open list</v-btn>
    </v-card>
  </router-link>
</template>

<style scoped>
.playlist-card__link {
  display: block;
  color: inherit;
  text-decoration: none;
}

.playlist-card {
  background:
    linear-gradient(180deg, rgba(18, 20, 24, 0.96), rgba(18, 20, 24, 0.9)),
    linear-gradient(135deg, rgba(240, 78, 62, 0.12), rgba(244, 211, 94, 0.06));
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 22px;
  color: var(--text-primary);
  padding: 24px;
  display: grid;
  gap: 12px;
  min-height: 100%;
  transition:
    transform 0.2s ease,
    border-color 0.2s ease,
    box-shadow 0.2s ease;
}

.playlist-card__link:hover .playlist-card {
  transform: translateY(-6px);
  border-color: rgba(244, 211, 94, 0.22);
  box-shadow: 0 18px 42px rgba(0, 0, 0, 0.18);
}

.playlist-card__eyebrow {
  color: var(--accent-soft);
  font-size: 0.78rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.playlist-card__title {
  margin: 0;
  font-size: 1.4rem;
  line-height: 1.1;
}

.playlist-card__meta {
  margin: 0;
  color: var(--text-secondary);
}

.playlist-card__action {
  justify-self: start;
  padding-left: 0;
  color: var(--text-primary) !important;
}
</style>

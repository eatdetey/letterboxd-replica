<script setup lang="ts">
import { computed } from 'vue'
import type { PlaylistDto } from '../model/types'

const props = defineProps<{ playlist: PlaylistDto }>()

const movieLabel = computed(() => {
  return props.playlist.movies_count === 1 ? 'movie' : 'movies'
})
</script>

<template>
  <v-card class="playlist-card" elevation="0">
    <div class="playlist-card__eyebrow">User list</div>
    <h3 class="playlist-card__title">{{ playlist.name }}</h3>
    <p class="playlist-card__meta">
      {{ playlist.movies_count }} {{ movieLabel }}
    </p>
    <v-btn variant="text" class="playlist-card__action">Open list</v-btn>
  </v-card>
</template>

<style scoped>
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

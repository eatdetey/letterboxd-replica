<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useMoviesStore } from '~/entities/movie/model/moviesStore'
import MovieCard from '~/entities/movie/ui/MovieCard.vue'

const moviesStore = useMoviesStore()

onMounted(() => {
  void moviesStore.loadMovies()
})

const featured = computed(() => moviesStore.featured)
const popular = computed(() => moviesStore.popular)
const allMovies = computed(() => moviesStore.all)
</script>

<template>
  <div class="home-page">
    <section class="hero" v-if="featured">
      <div
        class="hero__backdrop"
        :style="{ backgroundImage: featured.backdrop_url ? `url(${featured.backdrop_url})` : undefined }"
      ></div>
      <v-container class="hero__content">
        <div class="hero__badge">Featured</div>
        <h1 class="hero__title">{{ featured.title }}</h1>
        <p class="hero__description">{{ featured.description }}</p>
        <div class="hero__meta">
          <span>{{ featured.release_year }}</span>
          <span class="dot">•</span>
          <span>{{ featured.genres.join(', ') }}</span>
          <span class="dot">•</span>
          <span v-if="featured.rating !== undefined">{{ featured.rating.toFixed(1) }} rating</span>
        </div>
        <div class="hero__actions">
          <v-btn
            color="white"
            variant="flat"
            class="hero__primary"
            :to="{ name: 'movie', params: { id: featured.id } }"
          >
            View details
          </v-btn>
          <v-btn color="white" variant="outlined" class="hero__secondary">Add to watchlist</v-btn>
        </div>
      </v-container>
    </section>

    <section class="section" v-if="popular.length">
      <v-container>
        <div class="section__header">
          <h2 class="section__title">Popular right now</h2>
          <v-btn variant="text" class="section__action">See all</v-btn>
        </div>
        <v-row>
          <v-col
            v-for="movie in popular"
            :key="movie.id"
            cols="12"
            sm="6"
            md="4"
          >
            <MovieCard :movie="movie" variant="compact" />
          </v-col>
        </v-row>
      </v-container>
    </section>

    <section class="section">
      <v-container>
        <div class="section__header">
          <h2 class="section__title">All films</h2>
          <div class="section__filters">
            <v-btn size="small" variant="outlined" class="filter-btn">Drama</v-btn>
            <v-btn size="small" variant="outlined" class="filter-btn">Sci-Fi</v-btn>
            <v-btn size="small" variant="outlined" class="filter-btn">Thriller</v-btn>
          </div>
        </div>

        <v-row>
          <v-col
            v-for="movie in allMovies"
            :key="movie.id"
            cols="12"
            sm="6"
            md="4"
            lg="3"
          >
            <MovieCard :movie="movie" />
          </v-col>
        </v-row>

        <div class="section__empty" v-if="!moviesStore.isLoading && !allMovies.length">
          Пока нет фильмов для отображения.
        </div>
      </v-container>
    </section>
  </div>
</template>

<style scoped>
.home-page {
  display: grid;
  gap: 64px;
  padding-bottom: 64px;
}

.hero {
  position: relative;
  min-height: 520px;
  border-radius: 28px;
  overflow: hidden;
  margin: 24px;
  background: var(--surface);
}

.hero__backdrop {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center;
  filter: saturate(1.1);
}

.hero__backdrop::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(120deg, rgba(8, 10, 13, 0.92) 20%, rgba(8, 10, 13, 0.4) 70%);
}

.hero__content {
  position: relative;
  z-index: 1;
  color: white;
  padding: 80px 32px;
  display: grid;
  gap: 18px;
}

.hero__badge {
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.2);
  width: fit-content;
  padding: 4px 12px;
  border-radius: 999px;
  font-size: 0.85rem;
  text-transform: uppercase;
  letter-spacing: 0.12em;
}

.hero__title {
  font-size: clamp(2rem, 4vw, 3.4rem);
  line-height: 1.1;
  margin: 0;
}

.hero__description {
  max-width: 520px;
  color: rgba(255, 255, 255, 0.82);
  line-height: 1.6;
}

.hero__meta {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  color: rgba(255, 255, 255, 0.75);
  font-size: 0.95rem;
}

.hero__actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.hero__primary {
  color: #0d0f12 !important;
}

.section {
  padding: 0 12px;
}

.section__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.section__title {
  margin: 0;
  font-size: 1.6rem;
}

.section__action {
  color: var(--text-muted);
}

.section__filters {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-btn {
  border-color: rgba(255, 255, 255, 0.15) !important;
  color: var(--text-muted) !important;
}

.section__empty {
  text-align: center;
  color: var(--text-muted);
  padding: 48px 0;
}

@media (max-width: 960px) {
  .hero {
    margin: 16px;
  }

  .hero__content {
    padding: 56px 24px;
  }
}

@media (max-width: 600px) {
  .hero__content {
    padding: 42px 16px;
  }

  .hero__actions {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>

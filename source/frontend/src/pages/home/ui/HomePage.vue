<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useMoviesStore } from '~/entities/movie/model/moviesStore'
import MovieCard from '~/entities/movie/ui/MovieCard.vue'
import { useInfiniteScroll } from '@shared/lib/useInfiniteScroll'

const heroBackdropUrl =
  'https://images.unsplash.com/photo-1440404653325-ab127d49abc1?auto=format&fit=crop&w=1600&q=80'

const moviesStore = useMoviesStore()
const loadMoreSentinel = ref<HTMLElement | null>(null)

onMounted(() => {
  void moviesStore.loadMovies()
})

const featured = computed(() => moviesStore.featured)
const allMovies = computed(() => moviesStore.all)
const showLoadMoreAnchor = computed(() => moviesStore.hasMore && !moviesStore.error)
const canLoadMoreMovies = computed(() => {
  return showLoadMoreAnchor.value && !moviesStore.isLoading && !moviesStore.isLoadingMore
})

useInfiniteScroll({
  target: loadMoreSentinel,
  enabled: canLoadMoreMovies,
  onLoadMore: () => moviesStore.loadMoreMovies(),
})
</script>

<template>
  <div class="home-page">
    <section class="hero" v-if="featured">
      <div class="hero__backdrop" :style="{ backgroundImage: `url(${heroBackdropUrl})` }"></div>
      <v-container class="hero__content">
        <div class="hero__grid">
          <div class="hero__copy">
            <div class="hero__badge">Now showing</div>
            <h1 class="hero__title">{{ featured.title }}</h1>
            <p class="hero__description">{{ featured.description }}</p>
            <div class="hero__meta">
              <span>{{ featured.release_year }}</span>
              <span class="dot">•</span>
              <span>{{ featured.genres.join(', ') }}</span>
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
            </div>
          </div>

          <div
            class="hero__poster"
            :style="{ backgroundImage: featured.poster_url ? `url(${featured.poster_url})` : undefined }"
          ></div>
        </div>
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

        <div v-if="moviesStore.isLoading && !allMovies.length" class="section__state">
          <v-progress-circular indeterminate color="white" />
        </div>

        <v-row v-else-if="allMovies.length">
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

        <div v-if="showLoadMoreAnchor" ref="loadMoreSentinel" class="section__load-more">
          <v-progress-circular
            v-if="moviesStore.isLoadingMore"
            indeterminate
            size="24"
            width="2"
            color="white"
          />
        </div>

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
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.hero__backdrop {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center;
}

.hero__backdrop::after {
  content: '';
  position: absolute;
  inset: 0;
  background:
    linear-gradient(90deg, rgba(11, 13, 16, 0.92) 22%, rgba(11, 13, 16, 0.7) 58%, rgba(11, 13, 16, 0.92)),
    linear-gradient(180deg, rgba(11, 13, 16, 0.22), rgba(11, 13, 16, 0.82));
}

.hero__content {
  position: relative;
  z-index: 1;
  color: white;
  padding: 80px 32px;
}

.hero__grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(220px, 320px);
  gap: 32px;
  align-items: center;
}

.hero__copy {
  display: grid;
  gap: 18px;
}

.hero__poster {
  min-height: 440px;
  border-radius: 24px;
  background-color: rgba(255, 255, 255, 0.08);
  background-size: cover;
  background-position: center;
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.32);
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

.section__state,
.section__load-more {
  display: grid;
  place-items: center;
}

.section__state {
  min-height: 260px;
}

.section__load-more {
  min-height: 80px;
}

@media (max-width: 960px) {
  .hero {
    margin: 16px;
  }

  .hero__content {
    padding: 56px 24px;
  }

  .hero__grid {
    grid-template-columns: 1fr;
  }

  .hero__poster {
    min-height: 360px;
    max-width: 280px;
  }
}

@media (max-width: 600px) {
  .hero__content {
    padding: 42px 16px;
  }

  .hero__poster {
    min-height: 320px;
    max-width: none;
  }

  .hero__actions {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>

<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useMovieDetailsStore } from '~/entities/movie/model/movieDetailsStore'
import { useMovieReviewsStore } from '~/entities/review/model/movieReviewsStore'
import ReviewCard from '~/entities/review/ui/ReviewCard.vue'

const route = useRoute()
const movieDetailsStore = useMovieDetailsStore()
const movieReviewsStore = useMovieReviewsStore()

const movieId = computed(() => String(route.params.id ?? ''))
const movie = computed(() => movieDetailsStore.item)
const playlists = computed(() => movieDetailsStore.playlists)
const reviews = computed(() => movieReviewsStore.items)
const isLoading = computed(() => movieDetailsStore.isLoading || movieReviewsStore.isLoading)
const error = computed(() => movieDetailsStore.error || movieReviewsStore.error)

async function loadPageData(id: string) {
  if (!id) {
    return
  }

  await Promise.all([
    movieDetailsStore.loadMovie(id),
    movieReviewsStore.loadReviews(id),
  ])
}

onMounted(() => {
  void loadPageData(movieId.value)
})

watch(movieId, (id) => {
  void loadPageData(id)
})
</script>

<template>
  <div class="movie-page">
    <v-container class="movie-page__container">
      <div v-if="isLoading" class="movie-page__state">
        <v-progress-circular indeterminate color="white" />
      </div>

      <v-alert
        v-else-if="error"
        type="error"
        variant="tonal"
        class="movie-page__alert"
      >
        {{ error }}
      </v-alert>

      <template v-else-if="movie">
        <section class="movie-hero">
          <div class="movie-hero__grid">
            <div
              class="movie-poster"
              :style="{ backgroundImage: movie.poster_url ? `url(${movie.poster_url})` : undefined }"
            ></div>

            <div class="movie-hero__content">
              <div class="movie-hero__eyebrow">Film page</div>
              <h1 class="movie-hero__title">{{ movie.title }}</h1>

              <div class="movie-hero__meta">
                <span>{{ movie.release_year }}</span>
                <span class="dot">•</span>
                <span>{{ movie.genres.join(', ') }}</span>
              </div>

              <p class="movie-hero__description">{{ movie.description }}</p>

              <div class="movie-hero__actions">
                <v-btn color="white" variant="flat" class="movie-hero__primary">Add to watchlist</v-btn>
                <v-btn color="white" variant="outlined">Write a review</v-btn>
                <v-btn color="white" variant="outlined">Mark as watched</v-btn>
              </div>

              <div v-if="playlists.length" class="movie-playlists">
                <div class="movie-playlists__label">In your lists</div>
                <div class="movie-playlists__items">
                  <v-chip
                    v-for="playlist in playlists"
                    :key="playlist.id"
                    variant="outlined"
                    class="movie-playlists__chip"
                  >
                    {{ playlist.name }}
                  </v-chip>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section class="reviews-section">
          <div class="reviews-section__header">
            <h2 class="reviews-section__title">Reviews</h2>
            <div class="reviews-section__count">{{ reviews.length }} total</div>
          </div>

          <div v-if="reviews.length" class="reviews-section__list">
            <ReviewCard
              v-for="review in reviews"
              :key="review.id"
              :review="review"
            />
          </div>

          <div v-else class="movie-page__empty">
            Пока нет отзывов для этого фильма.
          </div>
        </section>
      </template>

      <div v-else class="movie-page__empty">
        Фильм не найден.
      </div>
    </v-container>
  </div>
</template>

<style scoped>
.movie-page {
  padding: 24px 0 72px;
}

.movie-page__container {
  display: grid;
  gap: 32px;
}

.movie-page__state,
.movie-page__empty {
  min-height: 320px;
  display: grid;
  place-items: center;
  color: var(--text-muted);
}

.movie-page__alert {
  border-radius: 20px;
}

.movie-hero {
  border-radius: 32px;
  overflow: hidden;
  background: var(--surface);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.movie-hero__grid {
  display: grid;
  grid-template-columns: minmax(240px, 300px) 1fr;
  gap: 32px;
  padding: 40px;
}

.movie-poster {
  min-height: 460px;
  border-radius: 24px;
  background-color: #161a20;
  background-size: cover;
  background-position: center;
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.32);
}

.movie-hero__content {
  display: grid;
  gap: 20px;
  align-content: center;
}

.movie-hero__eyebrow {
  color: var(--accent-soft);
  text-transform: uppercase;
  letter-spacing: 0.16em;
  font-size: 0.8rem;
}

.movie-hero__title {
  margin: 0;
  color: var(--text-primary);
  font-size: clamp(2.4rem, 5vw, 4.5rem);
  line-height: 0.98;
}

.movie-hero__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  color: rgba(242, 245, 247, 0.74);
}

.movie-hero__description {
  margin: 0;
  max-width: 760px;
  color: var(--text-secondary);
  line-height: 1.75;
  font-size: 1.02rem;
}

.movie-hero__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.movie-hero__primary {
  color: #0d0f12 !important;
}

.movie-playlists {
  display: grid;
  gap: 12px;
}

.movie-playlists__label {
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.12em;
  font-size: 0.72rem;
}

.movie-playlists__items {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.movie-playlists__chip {
  color: var(--text-primary) !important;
  border-color: rgba(255, 255, 255, 0.16) !important;
}

.reviews-section {
  display: grid;
  gap: 20px;
}

.reviews-section__header {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
  flex-wrap: wrap;
}

.reviews-section__title {
  margin: 0;
  font-size: 1.8rem;
}

.reviews-section__count {
  color: var(--text-muted);
}

.reviews-section__list {
  display: grid;
  gap: 16px;
}

@media (max-width: 960px) {
  .movie-hero__grid {
    grid-template-columns: 1fr;
    padding: 24px;
  }

  .movie-poster {
    min-height: 400px;
    max-width: 320px;
  }
}

@media (max-width: 600px) {
  .movie-page {
    padding-top: 16px;
  }

  .movie-hero__grid {
    padding: 20px;
  }

  .movie-poster {
    min-height: 320px;
    max-width: none;
  }

  .movie-hero__actions {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>

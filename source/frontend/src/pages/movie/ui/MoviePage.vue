<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '~/entities/auth/model/authStore'
import { useMovieDetailsStore } from '~/entities/movie/model/movieDetailsStore'
import { useMovieReviewsStore } from '~/entities/review/model/movieReviewsStore'
import { playlistsApi } from '~/entities/playlist/api/playlistsApi'
import { usePlaylistsStore } from '~/entities/playlist/model/playlistsStore'
import ReviewCard from '~/entities/review/ui/ReviewCard.vue'
import { useInfiniteScroll } from '@shared/lib/useInfiniteScroll'

const REVIEWS_BATCH_SIZE = 10

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const movieDetailsStore = useMovieDetailsStore()
const movieReviewsStore = useMovieReviewsStore()
const playlistsStore = usePlaylistsStore()
const isPlaylistDialogOpen = ref(false)
const selectedPlaylistId = ref('')
const newPlaylistName = ref('')
const playlistActionError = ref<string | null>(null)
const isSubmittingPlaylist = ref(false)
const isReviewDialogOpen = ref(false)
const reviewText = ref('')
const reviewActionError = ref<string | null>(null)
const visibleReviewsCount = ref(REVIEWS_BATCH_SIZE)
const reviewsLoadMoreSentinel = ref<HTMLElement | null>(null)

const movieId = computed(() => String(route.params.id ?? ''))
const movie = computed(() => movieDetailsStore.item)
const playlists = computed(() => movieDetailsStore.playlists)
const reviews = computed(() => movieReviewsStore.items)
const visibleReviews = computed(() => reviews.value.slice(0, visibleReviewsCount.value))
const isLoading = computed(() => movieDetailsStore.isLoading || movieReviewsStore.isLoading)
const error = computed(() => movieDetailsStore.error || movieReviewsStore.error)
const hasMoreReviews = computed(() => reviews.value.length > visibleReviewsCount.value)
const canLoadMoreReviews = computed(() => hasMoreReviews.value && !movieReviewsStore.isLoading)

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
  visibleReviewsCount.value = REVIEWS_BATCH_SIZE
  void loadPageData(id)
})

watch(() => reviews.value.length, (length, previousLength) => {
  if (length === 0) {
    visibleReviewsCount.value = REVIEWS_BATCH_SIZE
    return
  }

  if (previousLength === 0) {
    visibleReviewsCount.value = Math.min(REVIEWS_BATCH_SIZE, length)
    return
  }

  if (visibleReviewsCount.value > length) {
    visibleReviewsCount.value = length
  }
})

function loadMoreReviews() {
  if (!canLoadMoreReviews.value) {
    return
  }

  visibleReviewsCount.value += REVIEWS_BATCH_SIZE
}

useInfiniteScroll({
  target: reviewsLoadMoreSentinel,
  enabled: canLoadMoreReviews,
  onLoadMore: loadMoreReviews,
})

async function openReviewDialog() {
  if (!authStore.isAuthenticated) {
    await router.push({
      name: 'login',
      query: { redirect: route.fullPath },
    })
    return
  }

  reviewActionError.value = null
  reviewText.value = ''
  isReviewDialogOpen.value = true
}

async function openPlaylistDialog() {
  if (!authStore.isAuthenticated) {
    await router.push({
      name: 'login',
      query: { redirect: route.fullPath },
    })
    return
  }

  playlistActionError.value = null
  selectedPlaylistId.value = ''
  newPlaylistName.value = ''

  if (!playlistsStore.items.length) {
    await playlistsStore.loadPlaylists()
  }

  isPlaylistDialogOpen.value = true
}

async function addMovieToPlaylist() {
  if (!movie.value) {
    return
  }

  if (!selectedPlaylistId.value) {
    playlistActionError.value = 'Select a playlist first'
    return
  }

  isSubmittingPlaylist.value = true
  playlistActionError.value = null

  try {
    await playlistsApi.addMovieToPlaylist(selectedPlaylistId.value, movie.value.id)
    playlistsStore.updateMoviesCount(selectedPlaylistId.value, 1)
    await movieDetailsStore.loadMovie(movie.value.id)
    isPlaylistDialogOpen.value = false
  } catch (error) {
    playlistActionError.value = error instanceof Error ? error.message : 'Failed to add movie to playlist'
  } finally {
    isSubmittingPlaylist.value = false
  }
}

async function createPlaylistAndAddMovie() {
  if (!movie.value) {
    return
  }

  const name = newPlaylistName.value.trim()
  if (!name) {
    playlistActionError.value = 'Playlist name is required'
    return
  }

  isSubmittingPlaylist.value = true
  playlistActionError.value = null

  try {
    const playlist = await playlistsStore.createPlaylist(name)
    await playlistsApi.addMovieToPlaylist(playlist.id, movie.value.id)
    playlistsStore.updateMoviesCount(playlist.id, 1)
    await movieDetailsStore.loadMovie(movie.value.id)
    isPlaylistDialogOpen.value = false
  } catch (error) {
    playlistActionError.value = error instanceof Error ? error.message : 'Failed to create playlist'
  } finally {
    isSubmittingPlaylist.value = false
  }
}

async function submitReview() {
  if (!movie.value) {
    return
  }

  const text = reviewText.value.trim()
  if (!text) {
    reviewActionError.value = 'Review text is required'
    return
  }

  reviewActionError.value = null

  try {
    await movieReviewsStore.createReview(movie.value.id, { text })
    reviewText.value = ''
    isReviewDialogOpen.value = false
  } catch (submitError) {
    reviewActionError.value = submitError instanceof Error ? submitError.message : 'Failed to create review'
  }
}
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
                <v-btn color="white" variant="flat" class="movie-hero__primary" @click="openPlaylistDialog">
                  Add to list
                </v-btn>
                <v-btn color="white" variant="outlined" @click="openReviewDialog">Write a review</v-btn>
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
              v-for="review in visibleReviews"
              :key="review.id"
              :review="review"
            />
          </div>

          <div
            v-if="reviews.length && hasMoreReviews"
            ref="reviewsLoadMoreSentinel"
            class="reviews-section__load-more"
          >
            <v-progress-circular indeterminate size="24" width="2" color="white" />
          </div>

          <div v-if="!reviews.length" class="reviews-section__empty">
            No reviews yet for this film.
          </div>
        </section>
      </template>

      <div v-else class="movie-page__empty">
        Фильм не найден.
      </div>
    </v-container>

    <v-dialog v-model="isPlaylistDialogOpen" max-width="520">
      <v-card class="playlist-dialog" elevation="0">
        <h2 class="playlist-dialog__title">Add to playlist</h2>

        <v-alert
          v-if="playlistActionError"
          type="error"
          variant="tonal"
          class="playlist-dialog__alert"
        >
          {{ playlistActionError }}
        </v-alert>

        <v-select
          v-model="selectedPlaylistId"
          label="Choose playlist"
          variant="outlined"
          density="comfortable"
          :items="playlistsStore.items"
          item-title="name"
          item-value="id"
          hide-details
        />

        <div class="playlist-dialog__actions">
          <v-btn variant="text" @click="isPlaylistDialogOpen = false">Cancel</v-btn>
          <v-btn
            color="white"
            variant="flat"
            class="playlist-dialog__submit"
            :loading="isSubmittingPlaylist"
            @click="addMovieToPlaylist"
          >
            Add
          </v-btn>
        </div>

        <div class="playlist-dialog__divider"></div>

        <div class="playlist-dialog__subtitle">Or create a new playlist</div>

        <v-text-field
          v-model="newPlaylistName"
          label="New playlist name"
          variant="outlined"
          density="comfortable"
          hide-details
        />

        <div class="playlist-dialog__actions">
          <v-btn
            color="white"
            variant="flat"
            class="playlist-dialog__submit"
            :loading="isSubmittingPlaylist"
            @click="createPlaylistAndAddMovie"
          >
            Create and add
          </v-btn>
        </div>
      </v-card>
    </v-dialog>

    <v-dialog v-model="isReviewDialogOpen" max-width="640">
      <v-card class="playlist-dialog" elevation="0">
        <h2 class="playlist-dialog__title">Write a review</h2>

        <v-alert
          v-if="reviewActionError"
          type="error"
          variant="tonal"
          class="playlist-dialog__alert"
        >
          {{ reviewActionError }}
        </v-alert>

        <v-textarea
          v-model="reviewText"
          label="Your review"
          variant="outlined"
          rows="6"
          auto-grow
          hide-details
        />

        <div class="playlist-dialog__actions">
          <v-btn variant="text" @click="isReviewDialogOpen = false">Cancel</v-btn>
          <v-btn
            color="white"
            variant="flat"
            class="playlist-dialog__submit"
            :loading="movieReviewsStore.isSubmitting"
            @click="submitReview"
          >
            Publish
          </v-btn>
        </div>
      </v-card>
    </v-dialog>
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

.playlist-dialog {
  background: var(--surface);
  color: var(--text-primary);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  padding: 24px;
  display: grid;
  gap: 16px;
}

.playlist-dialog__title {
  margin: 0;
  font-size: 1.35rem;
}

.playlist-dialog__subtitle {
  color: var(--text-secondary);
  font-size: 0.95rem;
}

.playlist-dialog__alert {
  border-radius: 16px;
}

.playlist-dialog__divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.08);
}

.playlist-dialog__actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.playlist-dialog__submit {
  color: #0d0f12 !important;
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
  color: #000;
}

.reviews-section__list {
  display: grid;
  gap: 16px;
}

.reviews-section__load-more {
  min-height: 80px;
  display: grid;
  place-items: center;
}

.reviews-section__empty {
  min-height: 120px;
  display: grid;
  place-items: center;
  text-align: center;
  color: #000;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 16px;
  padding: 20px;
  font-weight: 600;
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

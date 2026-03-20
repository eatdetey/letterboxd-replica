<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useUserProfileStore } from '~/entities/user/model/userProfileStore'
import ReviewCard from '~/entities/review/ui/ReviewCard.vue'
import MoviePreviewCard from '~/entities/movie/ui/MoviePreviewCard.vue'

const route = useRoute()
const profileStore = useUserProfileStore()

const profileId = computed(() => String(route.params.username ?? route.params.id ?? ''))
const user = computed(() => profileStore.item)
const userLists = computed(() => profileStore.lists)
const reviews = computed(() => profileStore.reviews)
const recentRatings = computed(() => profileStore.recentRatings)
const watchlist = computed(() => profileStore.watchlist)
const isLoading = computed(() => profileStore.isLoading)
const error = computed(() => profileStore.error)

async function loadPageData(id: string) {
  if (!id) {
    return
  }

  await profileStore.loadProfile(id)
}

onMounted(() => {
  void loadPageData(profileId.value)
})

watch(profileId, (id) => {
  void loadPageData(id)
})
</script>

<template>
  <div class="profile-page">
    <v-container class="profile-page__container">
      <div v-if="isLoading" class="profile-page__state">
        <v-progress-circular indeterminate color="white" />
      </div>

      <v-alert
        v-else-if="error"
        type="error"
        variant="tonal"
        class="profile-page__alert"
      >
        {{ error }}
      </v-alert>

      <template v-else-if="user">
        <section class="profile-hero">
          <div
            class="profile-hero__backdrop"
            :style="{ backgroundImage: user.cover_url ? `url(${user.cover_url})` : undefined }"
          ></div>

          <div class="profile-hero__overlay"></div>

          <div class="profile-hero__content">
            <div class="profile-hero__grid">
              <div class="profile-avatar-wrap">
                <div
                  class="profile-avatar"
                  :style="{ backgroundImage: user.avatar_url ? `url(${user.avatar_url})` : undefined }"
                ></div>
              </div>

              <div class="profile-hero__copy">
                <div class="profile-hero__eyebrow">User profile</div>
                <h1 class="profile-hero__title">{{ user.username }}</h1>

                <div class="profile-hero__meta">
                  <span>{{ user.stats?.filmsWatched ?? 0 }} films watched</span>
                  <span class="dot">•</span>
                  <span>{{ user.stats?.reviewsCount ?? reviews.length }} reviews</span>
                  <span class="dot">•</span>
                  <span>{{ user.stats?.listsCount ?? userLists.length }} lists</span>
                </div>

                <p class="profile-hero__description">
                  {{ user.bio || 'Пользователь пока не добавил описание профиля.' }}
                </p>

                <div class="profile-hero__actions">
                  <v-btn color="white" variant="flat" class="profile-hero__primary">
                    Follow
                  </v-btn>
                  <v-btn color="white" variant="outlined">
                    Message
                  </v-btn>
                  <v-btn color="white" variant="outlined">
                    Share profile
                  </v-btn>
                </div>

                <div v-if="user.tags?.length" class="profile-tags">
                  <div class="profile-tags__label">Interests</div>
                  <div class="profile-tags__items">
                    <v-chip
                      v-for="tag in user.tags"
                      :key="tag"
                      variant="outlined"
                      class="profile-tags__chip"
                    >
                      {{ tag }}
                    </v-chip>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <div class="profile-layout">
          <main class="profile-main">
            <section class="profile-section">
              <div class="profile-section__header">
                <h2 class="profile-section__title">User lists</h2>
                <div class="profile-section__count">{{ userLists.length }} total</div>
              </div>

              <div v-if="userLists.length" class="lists-grid">
                <v-card
                  v-for="list in userLists"
                  :key="list.id"
                  class="list-card"
                  variant="flat"
                >
                  <v-card-text class="list-card__content">
                    <div class="list-card__top">
                      <div>
                        <div class="list-card__title">{{ list.name }}</div>
                        <div class="list-card__meta">
                          {{ list.movies_count ?? list.movies?.length ?? 0 }} films
                        </div>
                      </div>

                      <v-chip size="small" variant="outlined" class="list-card__chip">
                        {{ list.is_public ? 'Public' : 'Private' }}
                      </v-chip>
                    </div>

                    <p class="list-card__description">
                      {{ list.description || 'Описание списка пока не заполнено.' }}
                    </p>

                    <div v-if="list.preview_movies?.length" class="list-card__preview">
                      <div
                        v-for="movie in list.preview_movies.slice(0, 3)"
                        :key="movie.id"
                        class="list-card__preview-poster"
                        :style="{ backgroundImage: movie.poster_url ? `url(${movie.poster_url})` : undefined }"
                      ></div>
                    </div>
                  </v-card-text>
                </v-card>
              </div>

              <div v-else class="profile-page__empty">
                У пользователя пока нет списков.
              </div>
            </section>

            <section class="profile-section">
              <div class="profile-section__header">
                <h2 class="profile-section__title">Reviews</h2>
                <div class="profile-section__count">{{ reviews.length }} total</div>
              </div>

              <div v-if="reviews.length" class="profile-section__list">
                <ReviewCard
                  v-for="review in reviews"
                  :key="review.id"
                  :review="review"
                />
              </div>

              <div v-else class="profile-page__empty">
                Пока нет обзоров.
              </div>
            </section>

            <section class="profile-section">
              <div class="profile-section__header">
                <h2 class="profile-section__title">Recently rated</h2>
                <div class="profile-section__count">{{ recentRatings.length }} films</div>
              </div>

              <v-row v-if="recentRatings.length" dense>
                <v-col
                  v-for="movie in recentRatings"
                  :key="movie.id"
                  cols="12"
                  sm="6"
                  lg="4"
                >
                  <MoviePreviewCard :movie="movie" />
                </v-col>
              </v-row>

              <div v-else class="profile-page__empty">
                Недавно оценённых фильмов пока нет.
              </div>
            </section>
          </main>

          <aside class="profile-sidebar">
            <section class="sidebar-panel">
              <div class="sidebar-panel__header">
                <h2 class="sidebar-panel__title">Watchlist</h2>
                <div class="sidebar-panel__count">{{ watchlist.length }}</div>
              </div>

              <div v-if="watchlist.length" class="watchlist-stack">
                <v-card
                  v-for="movie in watchlist"
                  :key="movie.id"
                  class="watchlist-card"
                  variant="flat"
                  :to="{ name: 'movie', params: { id: movie.id } }"
                >
                  <div
                    class="watchlist-card__poster"
                    :style="{ backgroundImage: movie.poster_url ? `url(${movie.poster_url})` : undefined }"
                  ></div>

                  <div class="watchlist-card__content">
                    <div class="watchlist-card__title">{{ movie.title }}</div>
                    <div class="watchlist-card__meta">
                      {{ movie.release_year }}
                    </div>
                  </div>
                </v-card>
              </div>

              <div v-else class="profile-page__empty profile-page__empty--compact">
                Watchlist пуст.
              </div>
            </section>

            <section class="sidebar-panel">
              <div class="sidebar-panel__header">
                <h2 class="sidebar-panel__title">Profile stats</h2>
              </div>

              <div class="stats-grid">
                <div class="stat-card">
                  <div class="stat-card__value">{{ user.stats?.averageRating ?? '—' }}</div>
                  <div class="stat-card__label">Average rating</div>
                </div>

                <div class="stat-card">
                  <div class="stat-card__value">{{ user.stats?.followersCount ?? 0 }}</div>
                  <div class="stat-card__label">Followers</div>
                </div>

                <div class="stat-card">
                  <div class="stat-card__value">{{ user.stats?.followingCount ?? 0 }}</div>
                  <div class="stat-card__label">Following</div>
                </div>

                <div class="stat-card">
                  <div class="stat-card__value">{{ user.stats?.likesCount ?? 0 }}</div>
                  <div class="stat-card__label">Likes</div>
                </div>
              </div>
            </section>
          </aside>
        </div>
      </template>

      <div v-else class="profile-page__empty profile-page__empty--big">
        Пользователь не найден.
      </div>
    </v-container>
  </div>
</template>

<style scoped>
.profile-page {
  padding: 24px 0 72px;
  color: var(--text-primary);
}

.profile-page__container {
  display: grid;
  gap: 32px;
}

.profile-page__state,
.profile-page__empty {
  min-height: 320px;
  display: grid;
  place-items: center;
  color: var(--text-muted);
}

.profile-page__empty--compact {
  min-height: 160px;
}

.profile-page__empty--big {
  min-height: 420px;
}

.profile-page__alert {
  border-radius: 20px;
}

.profile-hero {
  position: relative;
  border-radius: 32px;
  overflow: hidden;
  background: var(--surface);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.profile-hero__backdrop {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center;
  filter: saturate(0.9);
}

.profile-hero__overlay {
  position: absolute;
  inset: 0;
  background:
    linear-gradient(90deg, rgba(11, 13, 16, 0.94) 20%, rgba(11, 13, 16, 0.66) 58%, rgba(11, 13, 16, 0.92)),
    linear-gradient(180deg, rgba(11, 13, 16, 0.18), rgba(11, 13, 16, 0.82));
}

.profile-hero__content {
  position: relative;
  z-index: 1;
  color: white;
  padding: 40px;
}

.profile-hero__grid {
  display: grid;
  grid-template-columns: 180px minmax(0, 1fr);
  gap: 28px;
  align-items: center;
}

.profile-avatar-wrap {
  display: flex;
  justify-content: flex-start;
}

.profile-avatar {
  width: 180px;
  height: 180px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.08);
  background-size: cover;
  background-position: center;
  border: 4px solid rgba(255, 255, 255, 0.18);
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.32);
}

.profile-hero__copy {
  display: grid;
  gap: 18px;
  align-content: center;
}

.profile-hero__eyebrow {
  color: var(--accent-soft);
  text-transform: uppercase;
  letter-spacing: 0.16em;
  font-size: 0.8rem;
}

.profile-hero__title {
  margin: 0;
  color: var(--text-primary);
  font-size: clamp(2.4rem, 5vw, 4.5rem);
  line-height: 0.98;
}

.profile-hero__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  color: rgba(242, 245, 247, 0.74);
}

.profile-hero__description {
  margin: 0;
  max-width: 760px;
  color: var(--text-secondary);
  line-height: 1.75;
  font-size: 1.02rem;
}

.profile-hero__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.profile-hero__primary {
  color: #0d0f12 !important;
}

.profile-tags {
  display: grid;
  gap: 12px;
}

.profile-tags__label {
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.12em;
  font-size: 0.72rem;
}

.profile-tags__items {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.profile-tags__chip {
  color: var(--text-primary) !important;
  border-color: rgba(255, 255, 255, 0.16) !important;
}

.profile-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 320px;
  gap: 28px;
  align-items: start;
}

.profile-main {
  display: grid;
  gap: 28px;
}

.profile-section {
  display: grid;
  gap: 20px;
}

.profile-section__header,
.sidebar-panel__header {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
  flex-wrap: wrap;
}

.profile-section__title {
  margin: 0;
  font-size: 1.8rem;
  color: var(--text-primary);
}

.sidebar-panel__title {
  margin: 0;
  font-size: 1.2rem;
  color: var(--text-secondary);
}

.profile-section__count,
.sidebar-panel__count {
  color: var(--text-muted);
}

.profile-section__list {
  display: grid;
  gap: 16px;
}

.lists-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.list-card {
  border-radius: 24px;
  background: var(--surface);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.list-card__content {
  display: grid;
  gap: 14px;
}

.list-card__top {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: flex-start;
}

.list-card__title {
  font-size: 1.08rem;
  font-weight: 600;
  color: var(--text-primary);
}

.list-card__meta {
  color: var(--text-muted);
  font-size: 0.92rem;
  margin-top: 4px;
}

.list-card__chip {
  color: var(--text-primary) !important;
  border-color: rgba(255, 255, 255, 0.16) !important;
}

.list-card__description {
  margin: 0;
  color: var(--text-secondary);
  line-height: 1.6;
}

.list-card__preview {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.list-card__preview-poster {
  aspect-ratio: 2 / 3;
  border-radius: 14px;
  background-color: #161a20;
  background-size: cover;
  background-position: center;
}

.profile-sidebar {
  display: grid;
  gap: 20px;
  position: sticky;
  top: 24px;
}

.sidebar-panel {
  background: var(--surface);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 28px;
  padding: 20px;
  display: grid;
  gap: 18px;
  color: var(--text-primary);
}

.watchlist-stack {
  display: grid;
  gap: 12px;
}

.watchlist-card {
  display: grid;
  grid-template-columns: 64px minmax(0, 1fr);
  gap: 12px;
  align-items: center;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.03);
  text-decoration: none;
  overflow: hidden;
}

.watchlist-card__poster {
  width: 64px;
  height: 96px;
  background-color: #161a20;
  background-size: cover;
  background-position: center;
}

.watchlist-card__content {
  display: grid;
  gap: 6px;
  padding-right: 12px;
}

.watchlist-card__title {
  color: var(--text-primary);
  font-weight: 600;
  line-height: 1.2;
}

.watchlist-card__meta {
  color: var(--text-muted);
  font-size: 0.9rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.stat-card {
  border-radius: 18px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.04);
  display: grid;
  gap: 6px;
}

.stat-card__value {
  font-size: 1.6rem;
  line-height: 1;
  color: var(--text-primary);
  font-weight: 700;
}

.stat-card__label {
  color: var(--text-muted);
  font-size: 0.88rem;
}

@media (max-width: 1080px) {
  .profile-layout {
    grid-template-columns: 1fr;
  }

  .profile-sidebar {
    position: static;
  }
}

@media (max-width: 960px) {
  .profile-hero__content {
    padding: 24px;
  }

  .profile-hero__grid {
    grid-template-columns: 1fr;
    justify-items: start;
  }

  .profile-avatar {
    width: 140px;
    height: 140px;
  }

  .lists-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 600px) {
  .profile-page {
    padding-top: 16px;
  }

  .profile-hero__content {
    padding: 20px;
  }

  .profile-hero__actions {
    flex-direction: column;
    align-items: stretch;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
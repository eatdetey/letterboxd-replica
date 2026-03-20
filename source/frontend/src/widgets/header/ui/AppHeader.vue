<template>
  <v-app-bar flat class="app-header">
    <v-container class="app-header__container">
      <router-link class="app-header__brand" :to="{ name: 'home' }">Betterboxd</router-link>

      <div class="app-header__nav">
        <v-btn variant="text" class="app-header__link" :to="{ name: 'home' }">Movies</v-btn>
        <v-btn variant="text" class="app-header__link" :to="{ name: 'playlists' }">Lists</v-btn>
        <v-btn variant="text" class="app-header__link" @click="goToProfile">Profile</v-btn>
      </div>

      <div class="app-header__actions">
        <v-menu
          v-model="isSearchOpen"
          :close-on-content-click="false"
          location="bottom"
          offset="10"
          max-width="520"
          min-width="100%"
        >
          <template #activator="{ props }">
            <v-text-field
              v-bind="props"
              v-model="searchQuery"
              density="compact"
              variant="solo"
              hide-details
              placeholder="Search for a film"
              class="app-header__search"
              @focus="handleSearchFocus"
              @keydown.enter.prevent="goToFirstResult"
              @keydown.esc.prevent="closeSearch"
            />
          </template>

          <v-card class="search-dropdown" elevation="8">
            <v-list v-if="searchLoading">
              <v-list-item title="Searching..." />
            </v-list>

            <v-list v-else-if="searchQuery.trim().length >= 2 && searchResults.length">
              <v-list-item
                v-for="movie in searchResults"
                :key="movie.id"
                :title="movie.title"
                :subtitle="`${movie.release_year} • ${movie.genres.join(', ')}`"
                :prepend-avatar="movie.poster_url || undefined"
                @click="openMovie(movie.id)"
              />
            </v-list>

            <v-list v-else>
              <v-list-item
                :title="searchQuery.trim().length < 2 ? 'Type at least 2 characters' : 'No films found'"
              />
            </v-list>
          </v-card>
        </v-menu>

        <v-btn
          variant="outlined"
          class="app-header__logout"
          @click="handleLogout"
        >
          {{ authStore.isAuthenticated ? 'Log out' : 'Log in' }}
        </v-btn>
      </div>
    </v-container>
  </v-app-bar>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '~/entities/auth/model/authStore'
import { moviesApi } from '~/entities/movie/api/moviesApi'
import { moviesMockApi } from '~/entities/movie/api/moviesMockApi'
import { env } from '@shared/config/env'
import type { MovieDto } from '~/entities/movie/model/types'

const router = useRouter()
const authStore = useAuthStore()

const searchQuery = ref('')
const searchResults = ref<MovieDto[]>([])
const searchLoading = ref(false)
const isSearchOpen = ref(false)

let searchTimer: number | undefined

const moviesSearchApi = env.useMocks ? moviesMockApi : moviesApi

const canSearch = computed(() => searchQuery.value.trim().length >= 2)

function clearSearchTimer() {
  if (searchTimer !== undefined) {
    window.clearTimeout(searchTimer)
    searchTimer = undefined
  }
}

async function runSearch(query: string) {
  const normalized = query.trim()

  if (normalized.length < 2) {
    searchResults.value = []
    isSearchOpen.value = false
    return
  }

  searchLoading.value = true

  try {
    const response = await moviesSearchApi.getMovies({
      search: normalized,
      limit: 6,
      offset: 0,
    })

    searchResults.value = response.items
    isSearchOpen.value = true
  } catch {
    searchResults.value = []
    isSearchOpen.value = true
  } finally {
    searchLoading.value = false
  }
}

watch(searchQuery, () => {
  clearSearchTimer()

  const query = searchQuery.value.trim()
  if (!query) {
    searchResults.value = []
    isSearchOpen.value = false
    return
  }

  searchTimer = window.setTimeout(() => {
    void runSearch(query)
  }, 250)
})

function handleSearchFocus() {
  if (canSearch.value) {
    isSearchOpen.value = true
    void runSearch(searchQuery.value)
  }
}

function closeSearch() {
  isSearchOpen.value = false
}

async function openMovie(movieId: string) {
  closeSearch()
  searchQuery.value = ''
  await router.push({ name: 'movie', params: { id: movieId } })
}

async function goToFirstResult() {
  if (searchResults.value.length) {
    const first = searchResults.value[0]
    if (first) {
      await openMovie(first.id)
    }
    return
  }

  if (canSearch.value) {
    await runSearch(searchQuery.value)
    if (searchResults.value.length) {
      const first = searchResults.value[0]
      if (first) {
        await openMovie(first.id)
      }
    }
  }
}

async function goToProfile() {
  if (!authStore.isAuthenticated) {
    await router.push({
      name: 'login',
      query: { redirect: '/' },
    })
    return
  }

  if (!authStore.username) {
    await router.push({ name: 'home' })
    return
  }

  await router.push({
    name: 'user',
    params: { username: authStore.username },
  })
}

async function handleLogout() {
  if (authStore.isAuthenticated) {
    authStore.logout()
    await router.push({ name: 'home' })
    return
  }

  await router.push({ name: 'login' })
}

onBeforeUnmount(() => {
  clearSearchTimer()
})
</script>

<style scoped>
.app-header {
  background: rgba(9, 10, 12, 0.9);
  backdrop-filter: blur(18px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.app-header__container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.app-header__brand {
  color: var(--text-primary);
  text-decoration: none;
  font-family: var(--display-font);
  font-size: 1.4rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.app-header__nav {
  display: flex;
  gap: 12px;
}

.app-header__link {
  color: var(--text-primary) !important;
  font-weight: 500;
}

.app-header__actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1 1 auto;
  justify-content: flex-end;
}

.app-header__search {
  flex: 0 1 420px;
  min-width: 320px;
}

.app-header__search :deep(.v-field) {
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.08);
}

.app-header__search :deep(input) {
  color: var(--text-primary);
}

.search-dropdown {
  width: 100%;
  border-radius: 18px;
  overflow: hidden;
  background: #11151b;
  color: var(--text-primary);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.app-header__logout {
  border-color: rgba(255, 255, 255, 0.2) !important;
  color: var(--text-primary) !important;
}

@media (max-width: 960px) {
  .app-header__nav {
    display: none;
  }

  .app-header__search {
    min-width: 260px;
    flex: 1 1 320px;
  }
}

@media (max-width: 600px) {
  .app-header__container {
    flex-direction: column;
    align-items: stretch;
  }

  .app-header__actions {
    width: 100%;
    justify-content: space-between;
  }

  .app-header__search {
    flex: 1;
    min-width: 0;
  }
}
</style>
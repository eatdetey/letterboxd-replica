<template>
  <v-navigation-drawer
    v-model="drawer"
    location="right"
    temporary
    class="mobile-drawer"
  >
    <v-list nav>
      <v-list-item>
        <v-text-field
          v-model="searchQuery"
          density="compact"
          variant="solo-filled"
          hide-details
          placeholder="Search..."
          prepend-inner-icon="mdi-magnify"
          @keydown.enter.prevent="goToFirstResult"
        />
      </v-list-item>

      <v-divider class="my-2"></v-divider>

      <v-list-item prepend-icon="mdi-movie-open" @click="router.push({ name: 'home' })">Movies</v-list-item>
      <v-list-item prepend-icon="mdi-format-list-bulleted" @click="router.push({ name: 'playlists' })">Lists</v-list-item>
      <v-list-item prepend-icon="mdi-account" @click="goToProfile">Profile</v-list-item>
      
      <v-divider class="my-2"></v-divider>
      
      <v-list-item 
        :prepend-icon="authStore.isAuthenticated ? 'mdi-logout' : 'mdi-login'" 
        @click="handleLogout"
      >
        {{ authStore.isAuthenticated ? 'Log out' : 'Log in' }}
      </v-list-item>
    </v-list>
  </v-navigation-drawer>

  <v-app-bar flat class="app-header" height="70">
    <v-container class="app-header__container d-flex align-center">
      <router-link class="app-header__brand" :to="{ name: 'home' }">Betterboxd</router-link>

      <v-spacer></v-spacer>

      <div class="app-header__nav d-none d-md-flex">
        <v-btn variant="text" class="app-header__link" :to="{ name: 'home' }">Movies</v-btn>
        <v-btn variant="text" class="app-header__link" :to="{ name: 'playlists' }">Lists</v-btn>
        <v-btn variant="text" class="app-header__link" @click="goToProfile">Profile</v-btn>
      </div>

      <div class="app-header__actions d-none d-md-flex ml-4">
        <v-menu
          v-model="isSearchOpen"
          :close-on-content-click="false"
          location="bottom"
          offset="10"
          max-width="520"
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
            <v-list v-else-if="canSearch && searchResults.length">
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
              <v-list-item :title="!canSearch ? 'Type at least 2 characters' : 'No films found'" />
            </v-list>
          </v-card>
        </v-menu>

        <v-btn
          variant="outlined"
          class="app-header__logout ml-4"
          @click="handleLogout"
        >
          {{ authStore.isAuthenticated ? 'Log out' : 'Log in' }}
        </v-btn>
      </div>

      <v-app-bar-nav-icon 
        variant="text" 
        class="d-flex d-md-none" 
        @click.stop="drawer = !drawer"
      ></v-app-bar-nav-icon>
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

const drawer = ref(false)
const searchQuery = ref('')
const searchResults = ref<MovieDto[]>([])
const searchLoading = ref(false)
const isSearchOpen = ref(false)

let searchTimer: number | undefined
const moviesSearchApi = env.useMovieMocks ? moviesMockApi : moviesApi
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
  drawer.value = false
  searchQuery.value = ''
  await router.push({ name: 'movie', params: { id: movieId } })
}

async function goToFirstResult() {
  if (searchResults.value.length) {
    const first = searchResults.value[0]
    if (first) {
      await openMovie(first.id)
    }
  }
}

async function goToProfile() {
  drawer.value = false
  if (!authStore.isAuthenticated) {
    await router.push({ name: 'login', query: { redirect: '/' } })
    return
  }
  await router.push({ name: 'user', params: { username: authStore.username } })
}

async function handleLogout() {
  drawer.value = false
  if (authStore.isAuthenticated) {
    authStore.logout()
    await router.push({ name: 'home' })
  } else {
    await router.push({ name: 'login' })
  }
}

onBeforeUnmount(() => {
  clearSearchTimer()
})
</script>

<style scoped>
.app-header {
  background: rgba(9, 10, 12, 0.9) !important;
  backdrop-filter: blur(18px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  color: #fff;
}

.app-header__brand {
  color: #fff;
  text-decoration: none;
  font-size: 1.4rem;
  font-weight: bold;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  white-space: nowrap;
}

.app-header__search {
  width: 300px;
  transition: width 0.3s ease;
}

.app-header__search :deep(.v-field) {
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}

.mobile-drawer {
  background: #11151b !important;
  color: white !important;
}

@media (max-width: 600px) {
  .app-header__brand {
    font-size: 1.1rem;
  }
}
</style>
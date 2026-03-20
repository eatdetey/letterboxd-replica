<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '~/entities/auth/model/authStore'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const username = ref('')
const password = ref('')
const showPassword = ref(false)

const isLoading = computed(() => authStore.isLoading)
const error = computed(() => authStore.error)

const redirectTo = computed(() => {
  const value = route.query.redirect
  return typeof value === 'string' && value.length ? value : '/'
})

async function submit() {
  await authStore.login({
    username: username.value.trim(),
    password: password.value,
  })

  await router.replace(redirectTo.value)
}

onMounted(async () => {
  if (!authStore.isReady) {
    await authStore.init()
  }

  if (authStore.isAuthenticated) {
    await router.replace(redirectTo.value)
  }
})
</script>

<template>
  <div class="login-page">
    <v-container class="login-page__container">
      <v-card class="login-card" elevation="0">
        <div class="login-card__copy">
          <div class="login-card__eyebrow">Welcome back</div>
          <h1 class="login-card__title">Sign in to Betterboxd</h1>
          <p class="login-card__description">
            Enter your username and password to continue.
          </p>
        </div>

        <v-alert
          v-if="error"
          type="error"
          variant="tonal"
          class="login-card__alert"
        >
          {{ error }}
        </v-alert>

        <v-form class="login-form" @submit.prevent="submit">
          <v-text-field
            v-model="username"
            label="Username"
            variant="outlined"
            density="comfortable"
            autocomplete="username"
            required
          />

          <v-text-field
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            label="Password"
            variant="outlined"
            density="comfortable"
            autocomplete="current-password"
            required
          />

          <v-checkbox
            v-model="showPassword"
            label="Show password"
            density="comfortable"
            hide-details
          />

          <v-btn
            type="submit"
            color="white"
            variant="flat"
            class="login-form__submit"
            :loading="isLoading"
            block
          >
            Sign in
          </v-btn>
        </v-form>
      </v-card>
    </v-container>
  </div>
</template>

<style scoped>
.login-page {
  min-height: calc(100vh - 80px);
  display: grid;
  place-items: center;
  padding: 32px 0;
}

.login-page__container {
  width: 100%;
  max-width: 560px;
}

.login-card {
  background: var(--surface);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 28px;
  padding: 32px;
  color: var(--text-primary);
  display: grid;
  gap: 24px;
}

.login-card__copy {
  display: grid;
  gap: 10px;
}

.login-card__eyebrow {
  color: var(--accent-soft);
  text-transform: uppercase;
  letter-spacing: 0.14em;
  font-size: 0.75rem;
}

.login-card__title {
  margin: 0;
  font-size: clamp(2rem, 4vw, 3rem);
  line-height: 1.05;
}

.login-card__description {
  margin: 0;
  color: var(--text-secondary);
  line-height: 1.6;
}

.login-card__alert {
  border-radius: 16px;
}

.login-form {
  display: grid;
  gap: 12px;
}

.login-form__submit {
  color: #0d0f12 !important;
}

@media (max-width: 600px) {
  .login-card {
    padding: 20px;
  }
}
</style>
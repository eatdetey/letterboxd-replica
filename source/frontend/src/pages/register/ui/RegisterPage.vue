<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '~/entities/auth/model/authStore'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const username = ref('')
const email = ref('')
const password = ref('')
const showPassword = ref(false)

const isLoading = computed(() => authStore.isLoading)
const error = computed(() => authStore.error)

const redirectTo = computed(() => {
  const value = route.query.redirect
  return typeof value === 'string' && value.length ? value : '/'
})

async function submit() {
  await authStore.register({
    username: username.value.trim(),
    email: email.value.trim(),
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
  <div class="register-page">
    <v-container class="register-page__container">
      <v-card class="register-card" elevation="0">
        <div class="register-card__copy">
          <div class="register-card__eyebrow">Create account</div>
          <h1 class="register-card__title">Join Betterboxd</h1>
          <p class="register-card__description">
            Create an account to save playlists and continue exploring films.
          </p>
        </div>

        <v-alert
          v-if="error"
          type="error"
          variant="tonal"
          class="register-card__alert"
        >
          {{ error }}
        </v-alert>

        <v-form class="register-form" @submit.prevent="submit">
          <v-text-field
            v-model="username"
            label="Username"
            variant="outlined"
            density="comfortable"
            autocomplete="username"
            required
          />

          <v-text-field
            v-model="email"
            label="Email"
            variant="outlined"
            density="comfortable"
            autocomplete="email"
            required
          />

          <v-text-field
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            label="Password"
            variant="outlined"
            density="comfortable"
            autocomplete="new-password"
            required
            hint="At least 8 characters with letters and digits"
            persistent-hint
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
            class="register-form__submit"
            :loading="isLoading"
            block
          >
            Create account
          </v-btn>

          <div class="register-form__footer">
            <span>Already have an account?</span>
            <router-link class="register-form__link" :to="{ name: 'login', query: route.query }">
              Sign in
            </router-link>
          </div>
        </v-form>
      </v-card>
    </v-container>
  </div>
</template>

<style scoped>
.register-page {
  min-height: calc(100vh - 80px);
  display: grid;
  place-items: center;
  padding: 32px 0;
}

.register-page__container {
  width: 100%;
  max-width: 560px;
}

.register-card {
  background: var(--surface);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 28px;
  padding: 32px;
  color: var(--text-primary);
  display: grid;
  gap: 24px;
}

.register-card__copy {
  display: grid;
  gap: 10px;
}

.register-card__eyebrow {
  color: var(--accent-soft);
  text-transform: uppercase;
  letter-spacing: 0.14em;
  font-size: 0.75rem;
}

.register-card__title {
  margin: 0;
  font-size: clamp(2rem, 4vw, 3rem);
  line-height: 1.05;
}

.register-card__description {
  margin: 0;
  color: var(--text-secondary);
  line-height: 1.6;
}

.register-card__alert {
  border-radius: 16px;
}

.register-form {
  display: grid;
  gap: 12px;
}

.register-form__submit {
  color: #0d0f12 !important;
}

.register-form__footer {
  display: flex;
  justify-content: center;
  gap: 6px;
  flex-wrap: wrap;
  color: var(--text-muted);
  font-size: 0.95rem;
}

.register-form__link {
  color: var(--text-primary);
  text-decoration: none;
  font-weight: 600;
}

@media (max-width: 600px) {
  .register-card {
    padding: 20px;
  }
}
</style>

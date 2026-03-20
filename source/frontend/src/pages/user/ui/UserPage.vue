<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useUserProfileStore } from '~/entities/user/model/userProfileStore'

const route = useRoute()
const profileStore = useUserProfileStore()

const profileId = computed(() => String(route.params.username ?? route.params.id ?? ''))
const user = computed(() => profileStore.item)
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
        <section class="profile-card">
          <div class="profile-card__avatar" aria-hidden="true">
            {{ user.username.charAt(0).toUpperCase() }}
          </div>

          <div class="profile-card__content">
            <div class="profile-card__eyebrow">Profile</div>
            <h1 class="profile-card__title">{{ user.username }}</h1>

            <div class="profile-card__grid">
              <div class="profile-field">
                <div class="profile-field__label">ID</div>
                <div class="profile-field__value">{{ user.id }}</div>
              </div>

              <div class="profile-field">
                <div class="profile-field__label">Username</div>
                <div class="profile-field__value">{{ user.username }}</div>
              </div>

              <div class="profile-field">
                <div class="profile-field__label">Email</div>
                <div class="profile-field__value">{{ user.email || 'Not provided' }}</div>
              </div>

              <div class="profile-field">
                <div class="profile-field__label">Role</div>
                <div class="profile-field__value">{{ user.role || 'user' }}</div>
              </div>
            </div>
          </div>
        </section>
      </template>

      <div v-else class="profile-page__empty">
        Пользователь не найден.
      </div>
    </v-container>
  </div>
</template>

<style scoped>
.profile-page {
  padding: 24px 0 72px;
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

.profile-page__alert {
  border-radius: 20px;
}

.profile-card {
  display: grid;
  grid-template-columns: 160px minmax(0, 1fr);
  gap: 28px;
  align-items: center;
  padding: 36px;
  border-radius: 32px;
  background: var(--surface);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.profile-card__avatar {
  width: 160px;
  height: 160px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  background:
    radial-gradient(circle at 30% 30%, rgba(255, 255, 255, 0.18), rgba(255, 255, 255, 0.04)),
    #171b22;
  color: var(--text-primary);
  font-size: 3.5rem;
  font-weight: 700;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.profile-card__content {
  display: grid;
  gap: 20px;
}

.profile-card__eyebrow {
  color: var(--accent-soft);
  text-transform: uppercase;
  letter-spacing: 0.16em;
  font-size: 0.8rem;
}

.profile-card__title {
  margin: 0;
  color: var(--text-primary);
  font-size: clamp(2.2rem, 4vw, 4rem);
  line-height: 1;
}

.profile-card__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.profile-field {
  display: grid;
  gap: 6px;
  padding: 16px 18px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.profile-field__label {
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.1em;
  font-size: 0.72rem;
}

.profile-field__value {
  color: var(--text-primary);
  font-size: 1rem;
  word-break: break-word;
}

@media (max-width: 900px) {
  .profile-card {
    grid-template-columns: 1fr;
    justify-items: center;
    text-align: center;
    padding: 28px;
  }

  .profile-card__grid {
    grid-template-columns: 1fr;
    width: 100%;
  }
}

@media (max-width: 600px) {
  .profile-page {
    padding-top: 16px;
  }

  .profile-card {
    padding: 22px;
  }

  .profile-card__avatar {
    width: 120px;
    height: 120px;
    font-size: 2.8rem;
  }
}
</style>

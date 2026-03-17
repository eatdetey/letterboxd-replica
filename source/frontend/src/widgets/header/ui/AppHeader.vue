<template>
  <v-app-bar flat class="app-header">
    <v-container class="app-header__container">
      <div class="app-header__brand">CineVault</div>
      <div class="app-header__nav">
        <v-btn variant="text" class="app-header__link">Movies</v-btn>
        <v-btn variant="text" class="app-header__link">Lists</v-btn>
        <v-btn variant="text" class="app-header__link">Profile</v-btn>
      </div>
      <div class="app-header__actions">
        <v-text-field
          density="compact"
          variant="solo"
          hide-details
          placeholder="Search for a film"
          class="app-header__search"
        />
        <v-btn variant="outlined" class="app-header__logout">Log out</v-btn>
      </div>
    </v-container>
  </v-app-bar>
</template>

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
}

.app-header__search :deep(.v-field) {
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.08);
}

.app-header__search :deep(input) {
  color: var(--text-primary);
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
    max-width: 180px;
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
  }
}
</style>

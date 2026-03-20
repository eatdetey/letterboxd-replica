<template>
  <v-footer class="app-footer" height="auto">
    <v-container class="app-footer__container">
      <div class="app-footer__brand">Betterboxd</div>
      <div class="app-footer__meta">Легенды ФИ 41/1, 2026</div>
    </v-container>
  </v-footer>
</template>

<style scoped>
.app-footer {
  background: rgba(9, 10, 12, 0.95);
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  color: var(--text-muted);
}

.app-footer__container {
  display: grid;
  gap: 12px;
  padding: 32px 0;
  text-align: center;
}

.app-footer__brand {
  font-family: var(--display-font);
  color: var(--text-primary);
  letter-spacing: 0.2em;
  text-transform: uppercase;
}

.app-footer__links {
  display: flex;
  justify-content: center;
  gap: 16px;
  flex-wrap: wrap;
}

.app-footer__link {
  color: var(--text-muted);
  text-decoration: none;
}

.app-footer__link:hover {
  color: var(--text-primary);
}

.app-footer__meta {
  font-size: 0.85rem;
}
</style>

<script setup lang="ts">
import AppHeader from '~/widgets/header/ui/AppHeader.vue'
import AppFooter from '~/widgets/footer/ui/AppFooter.vue'
</script>

<template>
  <v-app>
    <AppHeader />
    <v-main>
      <router-view></router-view>
    </v-main>
    <AppFooter />
  </v-app>
</template>

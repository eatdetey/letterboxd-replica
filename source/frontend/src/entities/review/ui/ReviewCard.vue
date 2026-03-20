<script setup lang="ts">
import { computed } from 'vue'
import type { MovieReviewDto } from '../model/types'

type ReviewCardDto = Pick<MovieReviewDto, 'id' | 'username' | 'text' | 'created_at'>

const props = defineProps<{ review: ReviewCardDto }>()

const formattedDate = computed(() => {
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  }).format(new Date(props.review.created_at))
})
</script>

<template>
  <v-card class="review-card" elevation="0">
    <div class="review-card__header">
      <div>
        <div class="review-card__author">{{ review.username }}</div>
        <div class="review-card__date">{{ formattedDate }}</div>
      </div>
    </div>
    <p class="review-card__text">{{ review.text }}</p>
  </v-card>
</template>

<style scoped>
.review-card {
  background: var(--surface-soft);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 20px;
  color: var(--text-primary);
  padding: 20px;
}

.review-card__header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.review-card__author {
  font-size: 1rem;
  font-weight: 600;
}

.review-card__date {
  color: var(--text-muted);
  font-size: 0.85rem;
}

.review-card__text {
  margin: 0;
  color: var(--text-secondary);
  line-height: 1.7;
}
</style>

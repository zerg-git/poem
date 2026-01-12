<template>
  <div class="poem-content paper-card">
    <h2 class="poem-title">{{ poem?.title || '无题' }}</h2>
    <div class="poem-author">{{ poem?.author || '佚名' }}</div>
    <div class="poem-body">
      <p v-for="(line, index) in (poem?.paragraphs || [])" :key="index" class="poem-line">
        {{ line }}
      </p>
    </div>
    <div v-if="poem?.rhythmic" class="poem-rhythmic">
      词牌：{{ poem.rhythmic }}
    </div>
  </div>
</template>

<script setup>
defineProps({
  poem: {
    type: Object,
    required: true
  }
})
</script>

<style scoped>
.poem-content {
  max-width: 600px;
  margin: 0 auto;
}

.poem-title {
  font-size: 1.8rem;
  margin-bottom: 0.5rem;
  color: var(--ink-black);
}

.poem-author {
  color: #666;
  margin-bottom: 2rem;
  font-size: 1rem;
}

.poem-body {
  margin: 2rem 0;
}

.poem-line {
  margin: 0.8em 0;
  line-height: 2;
  font-size: 1.25rem;
}

.poem-rhythmic {
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid #eee;
  color: #666;
  font-size: 0.9rem;
}
</style>

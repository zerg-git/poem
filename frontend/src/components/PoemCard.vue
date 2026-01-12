<template>
  <div class="card paper-card" @click="handleClick">
    <div class="card-title">{{ poem?.title || '无题' }}</div>
    <div class="card-subtitle">{{ poem?.author || '佚名' }}</div>
    <div class="poem-preview">
      <p v-for="(line, index) in previewLines" :key="index" class="poem-line">
        {{ line }}
      </p>
    </div>
    <div class="card-meta">
      <span class="tag">{{ dynastyName }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  poem: {
    type: Object,
    required: true
  }
})

const router = useRouter()

const previewLines = computed(() => {
  return props.poem?.paragraphs?.slice(0, 2) || []
})

const dynastyName = computed(() => {
  const dynastyMap = {
    'preqin': '先秦',
    'tang': '唐诗',
    'wudai': '五代',
    'song': '宋词',
    'yuan': '元曲',
    'ming': '明代',
    'qing': '清代',
    'other': '古诗词'
  }
  return dynastyMap[props.poem?.dynastyID] || '古诗词'
})

const handleClick = () => {
  if (props.poem?.id) {
    router.push(`/poem/${props.poem.id}`)
  }
}
</script>

<style scoped>
.poem-preview {
  margin: 1rem 0;
}

.poem-line {
  margin: 0.3em 0;
  font-size: 0.95rem;
  color: #333;
}

.card-meta {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}

.tag {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  background: var(--rice-paper);
  border: 1px solid #ccc;
  border-radius: 12px;
  font-size: 0.8rem;
  color: #666;
}
</style>

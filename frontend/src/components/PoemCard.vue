<template>
  <div class="card paper-card" @click="handleClick">
    <div class="card-title">{{ poem?.title || '无题' }}</div>
    <div class="card-subtitle">{{ authorName }}</div>
    <div class="poem-preview">
      <p v-for="(line, index) in previewLines" :key="index" class="poem-line">
        {{ line }}
      </p>
    </div>
    <div class="card-meta">
      <span class="tag">{{ categoryOrDynasty }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { formatDynasty } from '@/utils/common'

const props = defineProps({
  poem: {
    type: Object,
    required: true
  }
})

const router = useRouter()

const previewLines = computed(() => {
  // New backend returns 'content' as array of strings
  return props.poem?.content?.slice(0, 2) || []
})

const authorName = computed(() => {
  // Check if author is object (new structure) or string
  if (props.poem?.author?.name) {
    return props.poem.author.name
  }
  return props.poem?.author || '佚名'
})

const categoryOrDynasty = computed(() => {
  // Prioritize category display_name if available
  if (props.poem?.category?.display_name) {
    return props.poem.category.display_name
  }
  
  // Fallback to author's dynasty
  if (props.poem?.author?.dynasty) {
    return formatDynasty(props.poem.author.dynasty)
  }

  // Fallback to direct dynasty field (legacy/compatibility)
  if (props.poem?.dynasty) {
    return formatDynasty(props.poem.dynasty)
  }

  return '古诗词'
})

const handleClick = () => {
  // Use original_id if available (for permalinks), otherwise ID
  // Actually, let's use the DB ID for consistency now
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

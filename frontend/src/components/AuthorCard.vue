<template>
  <div class="card paper-card" @click="handleClick">
    <div class="card-title">{{ author.name }}</div>
    <div class="card-meta">
      <span class="tag dynasty-tag">{{ formatDynasty(author.dynasty) }}</span>
    </div>
    <div class="author-bio" v-if="author.biography">
      {{ truncate(author.biography, 50) }}
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { formatDynasty } from '@/utils/common'

const props = defineProps({
  author: {
    type: Object,
    required: true
  }
})

const router = useRouter()

const handleClick = () => {
  if (props.author.name) {
    router.push({ name: 'Author', params: { name: props.author.name } })
  }
}

const truncate = (text, length) => {
  if (!text) return ''
  if (text.length <= length) return text
  return text.substring(0, length) + '...'
}
</script>

<style scoped>
.author-bio {
  margin-top: 1rem;
  font-size: 0.9rem;
  color: #666;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-meta {
  margin-top: 0.5rem;
}

.tag {
  display: inline-block;
  padding: 0.2rem 0.6rem;
  background: var(--rice-paper);
  border: 1px solid #ccc;
  border-radius: 12px;
  font-size: 0.8rem;
  color: #666;
}

.dynasty-tag {
  color: var(--indigo);
  border-color: var(--indigo-light, #ccc);
}
</style>

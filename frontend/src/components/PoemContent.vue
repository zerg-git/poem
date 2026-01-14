<template>
  <div class="poem-content paper-card">
    <h2 class="poem-title">{{ displayTitle }}</h2>
    <div class="poem-author">
      <span v-if="dynasty" class="dynasty">[{{ dynasty }}]</span>
      {{ authorName }}
    </div>
    <div class="poem-body">
      <p v-for="(line, index) in contentLines" :key="index" class="poem-line">
        {{ line }}
      </p>
    </div>
    <div v-if="poem?.rhythmic" class="poem-rhythmic">
      词牌：{{ poem.rhythmic }}
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { formatDynasty } from '@/utils/common'

const props = defineProps({
  poem: {
    type: Object,
    required: true
  }
})

const displayTitle = computed(() => {
  return props.poem?.title || props.poem?.rhythmic || '无题'
})

const authorName = computed(() => {
  if (props.poem?.author?.name) {
    return props.poem.author.name
  }
  return props.poem?.author || '佚名'
})

const dynasty = computed(() => {
  if (props.poem?.author?.dynasty) {
    return formatDynasty(props.poem.author.dynasty)
  }
  if (props.poem?.dynasty) {
    return formatDynasty(props.poem.dynasty)
  }
  return ''
})

const contentLines = computed(() => {
  // New backend uses 'content', old one used 'paragraphs'
  return props.poem?.content || props.poem?.paragraphs || []
})
</script>

<style scoped>
.poem-content {
  max-width: 600px;
  margin: 0 auto;
  text-align: center;
  padding: 3rem 2rem;
}

.poem-title {
  font-size: 1.8rem;
  margin-bottom: 0.5rem;
  color: var(--ink-black);
  font-weight: bold;
}

.poem-author {
  color: #666;
  margin-bottom: 2rem;
  font-size: 1rem;
}

.dynasty {
  margin-right: 0.5rem;
  color: var(--cinnabar);
}

.poem-body {
  margin: 2rem 0;
  text-align: center;
}

.poem-line {
  margin: 0.8em 0;
  line-height: 2;
  font-size: 1.25rem;
  font-family: "Kaiti SC", "STKaiti", "KaiTi", serif;
}

.poem-rhythmic {
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid #eee;
  color: #666;
  font-size: 0.9rem;
}
</style>

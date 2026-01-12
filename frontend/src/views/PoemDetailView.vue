<template>
  <div class="poem-detail">
    <div class="container">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="poem">
        <!-- 返回按钮 -->
        <button class="btn back-btn" @click="goBack">
          ← 返回
        </button>

        <!-- 诗词内容 -->
        <PoemContent :poem="poem" />

        <!-- 诗词信息 -->
        <div class="poem-info paper-card">
          <div class="info-item">
            <span class="info-label">作者：</span>
            <router-link
              v-if="poem.author"
              :to="{ name: 'Author', params: { name: poem.author } }"
              class="info-value author-link"
            >
              {{ poem.author }}
            </router-link>
          </div>
          <div v-if="poem.rhythmic" class="info-item">
            <span class="info-label">词牌：</span>
            <span class="info-value">{{ poem.rhythmic }}</span>
          </div>
          <div v-if="poem.chapter" class="info-item">
            <span class="info-label">章节：</span>
            <span class="info-value">{{ poem.chapter }}</span>
            <span v-if="poem.section" class="info-value">・{{ poem.section }}</span>
          </div>
          <div v-else-if="poem.section" class="info-item">
            <span class="info-label">分类：</span>
            <span class="info-value">{{ poem.section }}</span>
          </div>
          <div v-if="poem.dynasty" class="info-item">
            <span class="info-label">朝代：</span>
            <span class="info-value">{{ poem.dynasty }}</span>
          </div>
        </div>

        <!-- 相关诗词 -->
        <div class="related-section">
          <h3 class="section-title">更多{{ poem.author }}的作品</h3>
          <div v-if="relatedPoems.length > 0" class="grid grid-cols-1 grid-cols-sm-2">
            <PoemCard
              v-for="relatedPoem in relatedPoems"
              :key="relatedPoem.id"
              :poem="relatedPoem"
            />
          </div>
          <div v-else class="empty">
            <p>暂无更多作品</p>
          </div>
        </div>
      </div>
      <div v-else class="empty">
        <p>诗词不存在</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { poetryAPI } from '@/api/poetry-api'
import { usePoetry } from '@/composables/usePoetry'
import PoemContent from '@/components/PoemContent.vue'
import PoemCard from '@/components/PoemCard.vue'

const route = useRoute()
const router = useRouter()
const { fetchPoemById } = usePoetry()

const poem = ref(null)
const relatedPoems = ref([])
const loading = ref(false)

onMounted(async () => {
  await loadPoem()
  if (poem.value && poem.value.author) {
    await loadRelatedPoems()
  }
})

const loadPoem = async () => {
  loading.value = true
  try {
    const data = await fetchPoemById(route.params.id)
    poem.value = data
  } catch (e) {
    console.error('加载诗词失败:', e)
  } finally {
    loading.value = false
  }
}

const loadRelatedPoems = async () => {
  try {
    const response = await poetryAPI.getAuthorPoems(poem.value.author, {
      page: 1,
      page_size: 4
    })
    if (response.data.success) {
      // 过滤掉当前诗词
      relatedPoems.value = response.data.data.poems
        .filter(p => p.id !== poem.value.id)
        .slice(0, 4)
    }
  } catch (e) {
    console.error('加载相关诗词失败:', e)
  }
}

const goBack = () => {
  router.back()
}
</script>

<style scoped>
.poem-detail {
  min-height: 100vh;
  padding: 2rem 0;
}

.back-btn {
  margin-bottom: 2rem;
}

.poem-info {
  max-width: 600px;
  margin: 2rem auto;
  padding: 1.5rem;
}

.info-item {
  margin-bottom: 1rem;
}

.info-item:last-child {
  margin-bottom: 0;
}

.info-label {
  font-weight: bold;
  color: #666;
}

.info-value {
  color: var(--ink-black);
}

.author-link {
  color: var(--indigo);
}

.author-link:hover {
  color: var(--cinnabar);
  text-decoration: underline;
}

.related-section {
  margin-top: 3rem;
  max-width: 900px;
  margin-left: auto;
  margin-right: auto;
}

.section-title {
  text-align: center;
  font-size: 1.5rem;
  margin-bottom: 1.5rem;
}

.empty {
  text-align: center;
  padding: 2rem;
  color: #999;
}
</style>

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
              v-if="authorName"
              :to="{ name: 'Author', params: { name: authorName } }"
              class="info-value author-link"
            >
              {{ authorName }}
            </router-link>
            <span v-else class="info-value">佚名</span>
          </div>
          <div v-if="poem.rhythmic" class="info-item">
            <span class="info-label">词牌：</span>
            <span class="info-value">{{ poem.rhythmic }}</span>
          </div>
          <div v-if="poem.volume" class="info-item">
            <span class="info-label">卷：</span>
            <span class="info-value">{{ poem.volume }}</span>
          </div>
          <div v-if="poem.section" class="info-item">
            <span class="info-label">篇章：</span>
            <span class="info-value">{{ poem.section }}</span>
          </div>
          <div v-if="poem.category?.display_name" class="info-item">
            <span class="info-label">收录：</span>
            <span class="info-value">{{ poem.category.display_name }}</span>
          </div>
          <div v-if="poem.author?.dynasty" class="info-item">
            <span class="info-label">朝代：</span>
            <span class="info-value">{{ formatDynasty(poem.author.dynasty) }}</span>
          </div>
        </div>

        <!-- 评论/注解 (新功能) -->
        <div v-if="poem.comments && poem.comments.length > 0" class="comments-section paper-card">
          <h3 class="section-title">注解与评析</h3>
          <div v-for="comment in poem.comments" :key="comment.id" class="comment-item">
             <div class="comment-meta">
               <span class="comment-type tag">{{ getCommentType(comment.type) }}</span>
               <span v-if="comment.commenter" class="commenter">{{ comment.commenter }}</span>
             </div>
             <div class="comment-content">{{ comment.content }}</div>
          </div>
        </div>

        <!-- 相关诗词 -->
        <div class="related-section">
          <h3 class="section-title">更多{{ authorName }}的作品</h3>
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
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { poetryAPI } from '@/api/poetry-api'
import { formatDynasty } from '@/utils/common'
import PoemContent from '@/components/PoemContent.vue'
import PoemCard from '@/components/PoemCard.vue'

const route = useRoute()
const router = useRouter()

const poem = ref(null)
const relatedPoems = ref([])
const loading = ref(false)

const authorName = computed(() => {
  return poem.value?.author?.name || poem.value?.author
})

const initData = async () => {
  // 重置状态
  poem.value = null
  relatedPoems.value = []
  
  await loadPoem()
  if (authorName.value) {
    await loadRelatedPoems()
  }
}

onMounted(initData)

watch(
  () => route.params.id,
  (newId) => {
    if (newId) {
      initData()
    }
  }
)

const loadPoem = async () => {
  loading.value = true
  try {
    const res = await poetryAPI.getPoemById(route.params.id)
    if (res.data.success) {
      poem.value = res.data.data
    }
  } catch (e) {
    console.error('加载诗词失败:', e)
  } finally {
    loading.value = false
  }
}

const loadRelatedPoems = async () => {
  try {
    const response = await poetryAPI.getAuthorPoems(authorName.value, {
      page: 1,
      pageSize: 4
    })
    if (response.data.success) {
      // 过滤掉当前诗词
      relatedPoems.value = response.data.data.works
        .filter(p => p.id !== poem.value.id)
        .slice(0, 4)
    }
  } catch (e) {
    console.error('加载相关诗词失败:', e)
  }
}

const getCommentType = (type) => {
  const map = {
    'note': '注解',
    'comment': '评析',
    'translation': '译文'
  }
  return map[type] || type
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

.comments-section {
  max-width: 800px;
  margin: 2rem auto;
  padding: 1.5rem;
}

.comment-item {
  margin-bottom: 1.5rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px dashed #eee;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-meta {
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.commenter {
  font-weight: bold;
  color: var(--indigo);
}

.comment-content {
  line-height: 1.6;
  color: #444;
}

.tag {
  display: inline-block;
  padding: 0.15rem 0.5rem;
  background: #eee;
  border-radius: 4px;
  font-size: 0.75rem;
  color: #666;
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

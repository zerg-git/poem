<template>
  <div class="author">
    <div class="container">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else>
        <!-- 作者信息 -->
        <div class="author-header paper-card">
          <h1 class="author-name">{{ name }}</h1>
          <div v-if="author" class="author-meta">
            <p class="author-dynasty">{{ formatDynasty(author.dynasty) }}</p>
            <p v-if="author.biography" class="author-bio">{{ author.biography }}</p>
          </div>
        </div>

        <!-- 诗词列表 -->
        <div class="poems-section">
          <h2 class="section-title">诗词作品</h2>
          <div v-if="poems.length > 0" class="grid grid-cols-1 grid-cols-sm-2 grid-cols-lg-3">
            <PoemCard
              v-for="poem in poems"
              :key="poem.id"
              :poem="poem"
            />
          </div>
          <div v-else class="empty">
            <p>暂无诗词</p>
          </div>
        </div>

        <!-- 分页 -->
        <div v-if="totalPages > 1" class="pagination">
          <button
            class="pagination-btn"
            :disabled="currentPage === 1"
            @click="goToPage(currentPage - 1)"
          >
            上一页
          </button>
          <span v-for="page in displayedPages" :key="page">
            <button
              v-if="page !== '...'"
              class="pagination-btn"
              :class="{ active: page === currentPage }"
              @click="goToPage(page)"
            >
              {{ page }}
            </button>
            <span v-else class="pagination-ellipsis">...</span>
          </span>
          <button
            class="pagination-btn"
            :disabled="currentPage === totalPages"
            @click="goToPage(currentPage + 1)"
          >
            下一页
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { poetryAPI } from '@/api/poetry-api'
import { formatDynasty } from '@/utils/common'
import PoemCard from '@/components/PoemCard.vue'

const route = useRoute()
// use computed for reactive name from route
const name = computed(() => route.params.name)
const author = ref(null)
const poems = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const displayedPages = computed(() => {
  const pages = []
  const showPages = 5
  let start = Math.max(1, currentPage.value - Math.floor(showPages / 2))
  let end = Math.min(totalPages.value, start + showPages - 1)

  if (end - start < showPages - 1) {
    start = Math.max(1, end - showPages + 1)
  }

  if (start > 1) {
    pages.push(1)
    if (start > 2) pages.push('...')
  }

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }

  if (end < totalPages.value) {
    if (end < totalPages.value - 1) pages.push('...')
    pages.push(totalPages.value)
  }

  return pages
})

const loadData = async () => {
  loading.value = true
  try {
    // 加载作者信息
    try {
      const authorRes = await poetryAPI.getAuthorByName(name.value)
      if (authorRes.data.success) {
        author.value = authorRes.data.data
      }
    } catch (e) {
      // 作者信息获取失败不影响诗词列表显示
      console.warn('获取作者信息失败:', e)
    }

    // 加载诗词列表
    const poemsRes = await poetryAPI.getAuthorPoems(name.value, {
      page: currentPage.value,
      pageSize: pageSize.value
    })
    if (poemsRes.data.success) {
      poems.value = poemsRes.data.data.works
      total.value = poemsRes.data.data.total
    }
  } catch (e) {
    console.error('加载数据失败:', e)
  } finally {
    loading.value = false
  }
}

const initData = async () => {
  author.value = null
  poems.value = []
  currentPage.value = 1
  total.value = 0
  await loadData()
}

onMounted(initData)

watch(
  () => route.params.name,
  (newName) => {
    if (newName) {
      initData()
    }
  }
)

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    loadData()
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}
</script>

<style scoped>
.author {
  min-height: 100vh;
  padding: 2rem 0;
}

.author-header {
  max-width: 600px;
  margin: 0 auto 3rem;
  padding: 2rem;
  text-align: center;
}

.author-name {
  font-size: 2rem;
  margin-bottom: 1rem;
  color: var(--ink-black);
}

.author-dynasty {
  font-size: 1.2rem;
  color: var(--indigo);
  margin-bottom: 0.5rem;
}

.author-bio {
  color: #666;
  font-size: 0.95rem;
  line-height: 1.6;
  margin-top: 1rem;
  text-align: left;
}

.poems-section {
  max-width: 1200px;
  margin: 0 auto;
}

.section-title {
  text-align: center;
  font-size: 1.5rem;
  margin-bottom: 2rem;
}

.empty {
  text-align: center;
  padding: 3rem;
  color: #999;
}

.pagination-ellipsis {
  padding: 0.5rem;
}
</style>

<template>
  <div class="search">
    <div class="container">
      <h1 class="page-title">搜索诗词</h1>

      <!-- 搜索框 -->
      <div class="search-box">
        <input
          v-model="query"
          type="text"
          class="input search-input"
          placeholder="搜索诗词标题、作者或内容..."
          @keyup.enter="handleSearch"
        >
        <button class="btn search-btn" @click="handleSearch">
          搜索
        </button>
      </div>

      <!-- 搜索结果 -->
      <div v-if="hasSearched">
        <div v-if="loading" class="loading">搜索中...</div>
        <div v-else-if="results.length > 0">
          <div class="search-info">
            找到 {{ total }} 首相关诗词
          </div>
          <div class="grid grid-cols-1 grid-cols-sm-2 grid-cols-lg-3">
            <PoemCard
              v-for="poem in results"
              :key="poem.id"
              :poem="poem"
            />
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
        <div v-else class="empty">
          <p>未找到相关诗词</p>
          <p class="empty-hint">请尝试其他关键词</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useSearch } from '@/composables/useSearch'
import PoemCard from '@/components/PoemCard.vue'

const query = ref('')
const hasSearched = ref(false)
const currentPage = ref(1)
const pageSize = ref(12)

const { results, loading, total, search } = useSearch()

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

const handleSearch = async () => {
  if (!query.value.trim()) {
    return
  }

  hasSearched.value = true
  currentPage.value = 1

  try {
    await search(query.value, {
      page: currentPage.value,
      page_size: pageSize.value
    })
  } catch (e) {
    console.error('搜索失败:', e)
  }
}

const goToPage = async (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    await search(query.value, {
      page: currentPage.value,
      page_size: pageSize.value
    })
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}
</script>

<style scoped>
.search {
  min-height: 100vh;
  padding: 2rem 0;
}

.page-title {
  text-align: center;
  font-size: 2rem;
  margin-bottom: 2rem;
}

.search-box {
  max-width: 600px;
  margin: 0 auto 3rem;
  display: flex;
  gap: 0.5rem;
}

.search-input {
  flex: 1;
  padding: 0.75rem 1rem;
  font-size: 1rem;
}

.search-info {
  text-align: center;
  margin-bottom: 2rem;
  color: #666;
}

.empty {
  text-align: center;
  padding: 3rem;
  color: #999;
}

.empty-hint {
  margin-top: 0.5rem;
  font-size: 0.9rem;
  color: #bbb;
}

.pagination-ellipsis {
  padding: 0.5rem;
}

@media (max-width: 640px) {
  .search-box {
    flex-direction: column;
  }

  .search-btn {
    width: 100%;
  }
}
</style>

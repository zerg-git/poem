<template>
  <div class="catalog">
    <div class="container">
      <h1 class="page-title">诗词目录</h1>

      <!-- 分类筛选 -->
      <div class="filters">
        <button
          class="filter-btn"
          :class="{ active: !selectedCategory }"
          @click="selectCategory(null)"
        >
          全部
        </button>
        <button
          class="filter-btn"
          :class="{ active: selectedCategory === 'quantangshi' }"
          @click="selectCategory('quantangshi')"
        >
          全唐诗
        </button>
        <button
          class="filter-btn"
          :class="{ active: selectedCategory === 'songci' }"
          @click="selectCategory('songci')"
        >
          宋词
        </button>
        <button
          class="filter-btn"
          :class="{ active: selectedCategory === 'yuanqu' }"
          @click="selectCategory('yuanqu')"
        >
          元曲
        </button>
        <button
          class="filter-btn"
          :class="{ active: selectedCategory === 'sishuwujing' }"
          @click="selectCategory('sishuwujing')"
        >
          四书五经
        </button>
        <button
          class="filter-btn"
          :class="{ active: selectedCategory === 'youmengying' }"
          @click="selectCategory('youmengying')"
        >
          幽梦影
        </button>
         <button
          class="filter-btn"
          :class="{ active: selectedCategory === 'shijing' }"
          @click="selectCategory('shijing')"
        >
          诗经
        </button>
         <button
          class="filter-btn"
          :class="{ active: selectedCategory === 'chuci' }"
          @click="selectCategory('chuci')"
        >
          楚辞
        </button>
         <button
          class="filter-btn"
          :class="{ active: selectedCategory === 'lunyu' }"
          @click="selectCategory('lunyu')"
        >
          论语
        </button>
         <button
          class="filter-btn"
          :class="{ active: selectedCategory === 'shuimotangshi' }"
          @click="selectCategory('shuimotangshi')"
        >
          水墨唐诗
        </button>
      </div>

      <!-- 诗词列表 -->
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="poems.length > 0" class="grid grid-cols-1 grid-cols-sm-2 grid-cols-lg-3">
        <PoemCard
          v-for="poem in poems"
          :key="poem.id"
          :poem="poem"
        />
      </div>
      <div v-else class="empty">
        <p>暂无诗词</p>
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
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { poetryAPI } from '@/api/poetry-api'
import PoemCard from '@/components/PoemCard.vue'

const route = useRoute()
const poems = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)
const selectedCategory = ref(null)

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

const fetchPoems = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value
    }
    
    // 如果选择了分类，传递 category 参数 (现在是 category name)
    if (selectedCategory.value) {
      params.category = selectedCategory.value
    }

    const res = await poetryAPI.getPoems(params)
    if (res.data && res.data.success) {
      poems.value = res.data.data.works || []
      total.value = res.data.data.total || 0
    }
  } catch (error) {
    console.error('Failed to fetch poems:', error)
  } finally {
    loading.value = false
  }
}

const selectCategory = (categoryId) => {
  selectedCategory.value = categoryId
  currentPage.value = 1
  fetchPoems()
}

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    fetchPoems()
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

onMounted(() => {
  // 检查 URL 参数是否有分类
  const category = route.query.category
  if (category) {
    selectedCategory.value = category
  }
  fetchPoems()
})
</script>

<style scoped>
.catalog {
  min-height: 100vh;
  padding: 2rem 0;
}

.page-title {
  text-align: center;
  font-size: 2rem;
  margin-bottom: 2rem;
}

.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 2rem;
  justify-content: center;
}

.filter-btn {
  padding: 0.5rem 1.5rem;
  border: 1px solid #ccc;
  background: white;
  cursor: pointer;
  border-radius: 20px;
  transition: all 0.3s;
}

.filter-btn:hover {
  border-color: var(--ink-black);
}

.filter-btn.active {
  background: var(--ink-black);
  color: white;
  border-color: var(--ink-black);
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

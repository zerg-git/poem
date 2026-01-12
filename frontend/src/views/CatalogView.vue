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
          v-for="cat in categories"
          :key="cat.id"
          class="filter-btn"
          :class="{ active: selectedCategory === cat.id }"
          @click="selectCategory(cat.id)"
        >
          {{ cat.name }}
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { poetryAPI } from '@/api/poetry-api'
import PoemCard from '@/components/PoemCard.vue'

const route = useRoute()
const poems = ref([])
const categories = ref([])
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

onMounted(async () => {
  await loadCategories()

  // 从路由查询参数获取分类
  if (route.query.category) {
    selectedCategory.value = route.query.category
  }

  await loadPoems()
})

watch(() => route.query.category, async (newCategory) => {
  selectedCategory.value = newCategory
  currentPage.value = 1
  await loadPoems()
})

const loadCategories = async () => {
  try {
    const response = await poetryAPI.getCategories()
    if (response.data.success) {
      categories.value = response.data.data
    }
  } catch (e) {
    console.error('加载分类失败:', e)
  }
}

const loadPoems = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }

    // 根据选择的分类映射到正确的category参数
    if (selectedCategory.value) {
      const category = categories.value.find(c => c.id === selectedCategory.value)
      if (category) {
        params.category = category.path
      }
    }

    const response = await poetryAPI.getPoems(params)
    if (response.data.success) {
      poems.value = response.data.data.poems
      total.value = response.data.data.total
    }
  } catch (e) {
    console.error('加载诗词失败:', e)
  } finally {
    loading.value = false
  }
}

const selectCategory = (categoryId) => {
  selectedCategory.value = categoryId
  currentPage.value = 1
  loadPoems()
}

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    loadPoems()
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}
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

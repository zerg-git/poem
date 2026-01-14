<template>
  <div class="authors-catalog">
    <div class="container">
      <h1 class="page-title">诗人名录</h1>

      <!-- 朝代筛选 -->
      <div class="filters">
        <button
          class="filter-btn"
          :class="{ active: !selectedDynasty }"
          @click="selectDynasty(null)"
        >
          全部
        </button>
        <button
          v-for="dyn in dynasties"
          :key="dyn.key"
          class="filter-btn"
          :class="{ active: selectedDynasty === dyn.key }"
          @click="selectDynasty(dyn.key)"
        >
          {{ dyn.label }}
        </button>
      </div>

      <!-- 作者列表 -->
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="authors.length > 0" class="grid grid-cols-1 grid-cols-sm-2 grid-cols-lg-3 grid-cols-xl-4">
        <AuthorCard
          v-for="author in authors"
          :key="author.id"
          :author="author"
        />
      </div>
      <div v-else class="empty">
        <p>暂无诗人信息</p>
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
import { useRoute, useRouter } from 'vue-router'
import { poetryAPI } from '@/api/poetry-api'
import AuthorCard from '@/components/AuthorCard.vue'

const route = useRoute()
const router = useRouter()

const authors = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)
const selectedDynasty = ref(null)

// 常用朝代列表
const dynasties = [
  { key: 'tang', label: '唐代' },
  { key: 'song', label: '宋代' },
  { key: 'yuan', label: '元代' },
  { key: 'ming', label: '明代' },
  { key: 'qing', label: '清代' },
  { key: 'han', label: '汉代' },
  { key: 'wei-jin', label: '魏晋' },
  { key: 'pre-qin', label: '先秦' }
]

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

const fetchAuthors = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value
    }
    
    if (selectedDynasty.value) {
      params.dynasty = selectedDynasty.value
    }

    const res = await poetryAPI.getAuthors(params)
    if (res.data && res.data.success) {
      authors.value = res.data.data.authors || []
      total.value = res.data.data.total || 0
    }
  } catch (error) {
    console.error('Failed to fetch authors:', error)
  } finally {
    loading.value = false
  }
}

const selectDynasty = (dynasty) => {
  selectedDynasty.value = dynasty
  currentPage.value = 1
  fetchAuthors()
  
  // Update URL query
  const query = { ...route.query }
  if (dynasty) {
    query.dynasty = dynasty
  } else {
    delete query.dynasty
  }
  router.replace({ query })
}

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    fetchAuthors()
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

onMounted(() => {
  const dynasty = route.query.dynasty
  if (dynasty) {
    selectedDynasty.value = dynasty
  }
  fetchAuthors()
})
</script>

<style scoped>
.authors-catalog {
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

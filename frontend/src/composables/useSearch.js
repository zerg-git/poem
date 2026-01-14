import { ref } from 'vue'
import { poetryAPI } from '@/api/poetry-api'

export function useSearch() {
  const results = ref([])
  const loading = ref(false)
  const error = ref(null)
  const total = ref(0)

  const search = async (query, params = {}) => {
    if (!query || query.trim() === '') {
      results.value = []
      total.value = 0
      return
    }

    loading.value = true
    error.value = null
    try {
      const response = await poetryAPI.search(query, params)
      if (response.data.success) {
        results.value = response.data.data.works
        total.value = response.data.data.total
        return response.data.data
      }
      throw new Error(response.data.error || '搜索失败')
    } catch (e) {
      error.value = e.message
      results.value = []
      total.value = 0
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    results,
    loading,
    error,
    total,
    search
  }
}

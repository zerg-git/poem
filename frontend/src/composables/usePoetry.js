import { ref } from 'vue'
import { poetryAPI } from '@/api/poetry-api'

export function usePoetry() {
  const poems = ref([])
  const currentPoem = ref(null)
  const loading = ref(false)
  const error = ref(null)

  // 获取诗词列表
  const fetchPoems = async (params = {}) => {
    loading.value = true
    error.value = null
    try {
      const response = await poetryAPI.getPoems(params)
      if (response.data.success) {
        poems.value = response.data.data.poems
        return response.data.data
      }
      throw new Error(response.data.error || '获取诗词列表失败')
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 获取单首诗词
  const fetchPoemById = async (id) => {
    loading.value = true
    error.value = null
    try {
      const response = await poetryAPI.getPoemById(id)
      if (response.data.success) {
        currentPoem.value = response.data.data
        return response.data.data
      }
      throw new Error(response.data.error || '获取诗词详情失败')
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 获取随机诗词
  const fetchRandomPoem = async (count = 1, category = '') => {
    loading.value = true
    error.value = null
    try {
      const response = await poetryAPI.getRandomPoem(count, category)
      if (response.data.success) {
        return response.data.data.poems
      }
      throw new Error(response.data.error || '获取随机诗词失败')
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 获取作者的诗词
  const fetchPoemsByAuthor = async (authorName, params = {}) => {
    loading.value = true
    error.value = null
    try {
      const response = await poetryAPI.getAuthorPoems(authorName, params)
      if (response.data.success) {
        return response.data.data
      }
      throw new Error(response.data.error || '获取作者诗词失败')
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    poems,
    currentPoem,
    loading,
    error,
    fetchPoems,
    fetchPoemById,
    fetchRandomPoem,
    fetchPoemsByAuthor
  }
}

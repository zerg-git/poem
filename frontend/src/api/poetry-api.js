import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response
  },
  error => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

export const poetryAPI = {
  // 获取诗词列表
  getPoems(params = {}) {
    // 映射前端参数到后端
    const backendParams = {
      page: params.page,
      page_size: params.pageSize,
      category: params.category
    }
    return api.get('/poems', { params: backendParams })
  },

  // 获取单首诗词
  getPoemById(id) {
    return api.get(`/poems/${id}`)
  },

  // 获取随机诗词
  getRandomPoem(count = 1, category = '') {
    return api.get('/poems/random', { params: { count, category } })
  },

  // 搜索
  search(query, params = {}) {
    return api.get('/search', { 
      params: { 
        q: query, 
        page: params.page,
        page_size: params.pageSize 
      } 
    })
  },

  // 获取朝代列表 (Legacy)
  getDynasties() {
    return api.get('/dynasties')
  },

  // 获取分类列表
  getCategories() {
    return api.get('/categories')
  },

  // 获取作者列表
  getAuthors(params = {}) {
    return api.get('/authors', { 
      params: {
        page: params.page,
        page_size: params.pageSize,
        dynasty: params.dynasty
      }
    })
  },

  // 获取作者详情
  getAuthorByName(name) {
    return api.get(`/authors/${name}`)
  },

  // 获取作者的诗词
  getAuthorPoems(name, params = {}) {
    return api.get(`/authors/${name}/poems`, { 
      params: {
        page: params.page,
        page_size: params.pageSize
      } 
    })
  }
}

export default api

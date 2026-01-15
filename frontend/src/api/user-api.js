import axios from 'axios'

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

// 创建axios实例
const apiClient = axios.create({
  baseURL: API_BASE,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器：添加token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器：处理错误
apiClient.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      // Token过期或无效，清除本地存储
      localStorage.removeItem('token')
      localStorage.removeItem('currentUser')
      // 可以在这里触发重新登录
      window.location.href = '/login'
    }
    return Promise.reject(error.response?.data || error.message)
  }
)

// 用户API
export const userApi = {
  // 注册
  register(data) {
    return apiClient.post('/api/v2/auth/register', data)
  },

  // 登录
  login(data) {
    return apiClient.post('/api/v2/auth/login', data)
  },

  // 获取个人资料
  getProfile() {
    return apiClient.get('/api/v2/users/profile')
  },

  // 更新个人资料
  updateProfile(data) {
    return apiClient.put('/api/v2/users/profile', data)
  },

  // 获取用户资料（公开）
  getUserProfile(id) {
    return apiClient.get(`/api/v2/users/${id}`)
  },

  // 刷新Token
  refreshToken(data) {
    return apiClient.post('/api/v2/auth/refresh', data)
  },

  // 登出（前端清除token即可）
  logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('currentUser')
    return Promise.resolve()
  }
}

export default apiClient

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref(localStorage.getItem('token') || '')
  const currentUser = ref(JSON.parse(localStorage.getItem('currentUser') || 'null'))
  const isAuthenticated = computed(() => !!token.value)

  // 设置认证信息
  function setAuth(tokenValue, user) {
    token.value = tokenValue
    currentUser.value = user
    localStorage.setItem('token', tokenValue)
    localStorage.setItem('currentUser', JSON.stringify(user))
  }

  // 清除认证信息
  function clearAuth() {
    token.value = ''
    currentUser.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('currentUser')
  }

  // 更新用户信息
  function updateUser(user) {
    currentUser.value = { ...currentUser.value, ...user }
    localStorage.setItem('currentUser', JSON.stringify(currentUser.value))
  }

  return {
    token,
    currentUser,
    isAuthenticated,
    setAuth,
    clearAuth,
    updateUser
  }
})

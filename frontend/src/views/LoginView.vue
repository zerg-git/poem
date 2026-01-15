<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h1>{{ isLogin ? '登录' : '注册' }}</h1>
        <p class="subtitle">中国古诗词平台</p>
      </div>

      <form @submit.prevent="handleSubmit" class="login-form">
        <!-- 注册时显示昵称 -->
        <div v-if="!isLogin" class="form-group">
          <label for="nickname">昵称</label>
          <input
            id="nickname"
            v-model="form.nickname"
            type="text"
            placeholder="请输入昵称（可选）"
            maxlength="50"
          />
        </div>

        <!-- 用户名 -->
        <div class="form-group">
          <label for="username">用户名</label>
          <input
            id="username"
            v-model="form.username"
            type="text"
            placeholder="请输入用户名"
            minlength="3"
            maxlength="50"
            required
          />
        </div>

        <!-- 密码 -->
        <div class="form-group">
          <label for="password">密码</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            minlength="6"
            maxlength="50"
            required
          />
        </div>

        <!-- 错误提示 -->
        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>

        <!-- 提交按钮 -->
        <button type="submit" class="submit-btn" :disabled="loading">
          {{ loading ? '处理中...' : (isLogin ? '登录' : '注册') }}
        </button>

        <!-- 切换登录/注册 -->
        <div class="toggle-mode">
          <span>{{ isLogin ? '还没有账号？' : '已有账号？' }}</span>
          <a href="#" @click.prevent="toggleMode">
            {{ isLogin ? '立即注册' : '立即登录' }}
          </a>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { userApi } from '@/api/user-api'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const isLogin = ref(true)
const loading = ref(false)
const errorMessage = ref('')

const form = ref({
  username: '',
  password: '',
  nickname: ''
})

// 如果用户已登录，跳转到首页或重定向地址
onMounted(() => {
  if (userStore.isAuthenticated) {
    const redirectTarget = route.query.redirect || '/'
    router.replace(redirectTarget)
  }
})

const toggleMode = () => {
  isLogin.value = !isLogin.value
  errorMessage.value = ''
  form.value = {
    username: '',
    password: '',
    nickname: ''
  }
}

const handleSubmit = async () => {
  errorMessage.value = ''
  loading.value = true

  try {
    if (isLogin.value) {
      // 登录
      const response = await userApi.login({
        username: form.value.username,
        password: form.value.password
      })

      if (response.success) {
        userStore.setAuth(response.data.token, response.data.user)
        // 登录成功后跳转到原目标页面或首页
        const redirectTarget = route.query.redirect || '/'
        router.push(redirectTarget)
      }
    } else {
      // 注册
      const response = await userApi.register({
        username: form.value.username,
        password: form.value.password,
        nickname: form.value.nickname || form.value.username
      })

      if (response.success) {
        userStore.setAuth(response.data.token, response.data.user)
        // 注册成功后跳转到原目标页面或首页
        const redirectTarget = route.query.redirect || '/'
        router.push(redirectTarget)
      }
    }
  } catch (error) {
    errorMessage.value = error.error || error.message || '操作失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
  padding: 40px;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h1 {
  font-size: 28px;
  color: #333;
  margin: 0 0 8px 0;
}

.subtitle {
  color: #666;
  font-size: 14px;
  margin: 0;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.form-group input {
  padding: 12px 16px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.2s;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.error-message {
  background: #fee;
  color: #c33;
  padding: 12px;
  border-radius: 8px;
  font-size: 14px;
}

.submit-btn {
  padding: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: opacity 0.2s;
}

.submit-btn:hover:not(:disabled) {
  opacity: 0.9;
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.toggle-mode {
  text-align: center;
  font-size: 14px;
  color: #666;
}

.toggle-mode a {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.toggle-mode a:hover {
  text-decoration: underline;
}
</style>

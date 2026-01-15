<template>
  <div id="app">
    <!-- 导航栏 -->
    <nav class="navbar">
      <div class="container">
        <router-link to="/" class="nav-logo">
          中国古诗词
        </router-link>
        <div class="nav-links">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link to="/catalog" class="nav-link">目录</router-link>
          <router-link to="/authors" class="nav-link">诗人</router-link>
          <router-link to="/search" class="nav-link">搜索</router-link>
        </div>
        <!-- 用户导航区域 -->
        <div class="nav-user">
          <!-- 未登录：显示登录按钮 -->
          <router-link v-if="!isAuthenticated" to="/login" class="nav-login-btn">
            登录/注册
          </router-link>

          <!-- 已登录：显示用户信息和下拉菜单 -->
          <div v-else class="user-dropdown">
            <div class="user-trigger" @click="toggleDropdown">
              <span class="user-name">{{ currentUser?.nickname || currentUser?.username }}</span>
              <span class="dropdown-arrow">▼</span>
            </div>
            <div v-if="showDropdown" class="dropdown-menu">
              <router-link to="/profile" class="dropdown-item" @click="showDropdown = false">个人中心</router-link>
              <div class="dropdown-divider"></div>
              <a @click="handleLogout" class="dropdown-item">退出登录</a>
            </div>
          </div>
        </div>
      </div>
    </nav>

    <!-- 主内容 -->
    <main>
      <router-view />
    </main>

    <!-- 页脚 -->
    <footer class="footer">
      <div class="container">
        <p>中国古诗词 - 收录逾55万首古典诗词</p>
        <p class="footer-note">数据来源：chinese-poetry</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { poetryAPI } from '@/api/poetry-api'
import { usePoetryStore } from '@/stores/poetry'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const poetryStore = usePoetryStore()
const userStore = useUserStore()
const showDropdown = ref(false)

const isAuthenticated = computed(() => userStore.isAuthenticated)
const currentUser = computed(() => userStore.currentUser)

onMounted(async () => {
  // 预加载朝代和分类数据
  try {
    const [dynastiesRes, categoriesRes] = await Promise.all([
      poetryAPI.getDynasties(),
      poetryAPI.getCategories()
    ])

    if (dynastiesRes.data.success) {
      poetryStore.setDynasties(dynastiesRes.data.data)
    }

    if (categoriesRes.data.success) {
      poetryStore.setCategories(categoriesRes.data.data)
    }
  } catch (e) {
    console.error('预加载数据失败:', e)
  }

  // 添加点击外部关闭下拉菜单的事件监听
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

function toggleDropdown() {
  showDropdown.value = !showDropdown.value
}

function handleLogout() {
  userStore.clearAuth()
  showDropdown.value = false
  router.push('/')
}

function handleClickOutside(event) {
  if (!event.target.closest('.user-dropdown')) {
    showDropdown.value = false
  }
}
</script>

<style scoped>
.navbar {
  background: white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.navbar .container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 20px;
}

.nav-logo {
  font-size: 1.3rem;
  font-weight: bold;
  color: var(--ink-black);
}

.nav-links {
  display: flex;
  gap: 2rem;
}

.nav-link {
  color: var(--ink-black);
  text-decoration: none;
  transition: color 0.3s;
}

.nav-link:hover,
.nav-link.router-link-active {
  color: var(--cinnabar);
}

/* 用户导航样式 */
.nav-user {
  display: flex;
  align-items: center;
}

.nav-login-btn {
  padding: 8px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 20px;
  text-decoration: none;
  font-weight: 500;
  transition: opacity 0.2s;
}

.nav-login-btn:hover {
  opacity: 0.9;
}

.user-dropdown {
  position: relative;
}

.user-trigger {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  cursor: pointer;
  border-radius: 20px;
  background: #f5f5f5;
  transition: background 0.2s;
}

.user-trigger:hover {
  background: #eee;
}

.user-name {
  font-size: 14px;
  color: #333;
}

.dropdown-arrow {
  font-size: 10px;
  color: #666;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  min-width: 150px;
  overflow: hidden;
  z-index: 200;
}

.dropdown-item {
  display: block;
  padding: 12px 16px;
  color: #333;
  text-decoration: none;
  font-size: 14px;
  cursor: pointer;
}

.dropdown-item:hover {
  background: #f5f5f5;
}

.dropdown-divider {
  height: 1px;
  background: #eee;
  margin: 0;
}

main {
  min-height: calc(100vh - 140px);
}

.footer {
  background: var(--ink-black);
  color: var(--rice-paper);
  padding: 2rem 0;
  text-align: center;
}

.footer-note {
  margin-top: 0.5rem;
  font-size: 0.85rem;
  opacity: 0.7;
}

@media (max-width: 640px) {
  .navbar .container {
    flex-wrap: wrap;
    gap: 1rem;
  }

  .nav-links {
    gap: 1rem;
    order: 2;
    width: 100%;
    justify-content: center;
  }

  .nav-user {
    order: 3;
    width: 100%;
    justify-content: center;
  }
}
</style>

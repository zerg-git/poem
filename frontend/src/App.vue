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
          <router-link to="/search" class="nav-link">搜索</router-link>
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
import { onMounted } from 'vue'
import { poetryAPI } from '@/api/poetry-api'
import { usePoetryStore } from '@/stores/poetry'

const poetryStore = usePoetryStore()

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
})
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
    flex-direction: column;
    gap: 1rem;
  }

  .nav-links {
    gap: 1rem;
  }
}
</style>

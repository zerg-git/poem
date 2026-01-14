<template>
  <div class="home">
    <!-- 头部 -->
    <header class="header">
      <div class="container">
        <h1 class="title">中国古诗词</h1>
        <p class="subtitle">收录唐诗、宋词、元曲等海量古典文学作品</p>
      </div>
    </header>

    <!-- 快捷入口 -->
    <section class="section">
      <div class="container">
        <h2 class="section-title">浏览诗词</h2>
        <div class="grid grid-cols-1 grid-cols-sm-2 grid-cols-lg-3">
          <div
            v-for="category in categories"
            :key="category.id"
            class="category-card paper-card card"
            @click="goToCatalog(category)"
          >
            <h3 class="category-title">{{ category.name }}</h3>
            <p class="category-desc">{{ category.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 每日一诗 -->
    <section class="section">
      <div class="container">
        <h2 class="section-title">每日一诗</h2>
        <div v-if="randomPoem" class="random-poem">
          <PoemContent :poem="randomPoem" />
          <div class="random-actions">
            <button class="btn" @click="loadRandomPoem">换一首</button>
          </div>
        </div>
        <div v-else-if="loading" class="loading">加载中...</div>
        <div v-else class="error-message">
          <p>暂时无法加载随机诗词</p>
          <button class="btn" @click="loadRandomPoem">重试</button>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { poetryAPI } from '@/api/poetry-api'
import PoemContent from '@/components/PoemContent.vue'

const router = useRouter()
const randomPoem = ref(null)
const loading = ref(false)
const categories = ref([
  {
    id: 'quantangshi',
    name: '全唐诗',
    description: '全唐诗收录唐诗四万八千九百余首'
  },
  {
    id: 'songci',
    name: '宋词',
    description: '全宋词收录宋词二万余首'
  },
  {
    id: 'yuanqu',
    name: '元曲',
    description: '元代文学形式，包括散曲和杂剧'
  },
  {
    id: 'sishuwujing',
    name: '四书五经',
    description: '儒家经典著作'
  },
  {
    id: 'lunyu',
    name: '论语',
    description: '儒家经典'
  },
  {
    id: 'shijing',
    name: '诗经',
    description: '中国古代诗歌开端'
  },
  {
    id: 'chuci',
    name: '楚辞',
    description: '屈原创作的诗歌总集'
  },
  {
    id: 'shuimotangshi',
    name: '水墨唐诗',
    description: '水墨风格唐诗精选'
  },
  {
    id: 'youmengying',
    name: '幽梦影',
    description: '清代张潮著'
  }
])

onMounted(async () => {
  await loadRandomPoem()
})

const loadRandomPoem = async () => {
  loading.value = true
  try {
    const response = await poetryAPI.getRandomPoem(1)
    if (response.data.success) {
      // 兼容两种返回格式：直接返回数组或在 data.poems 中
      const data = response.data.data
      const poems = Array.isArray(data) ? data : (data.poems || [])
      
      if (poems.length > 0) {
        randomPoem.value = poems[0]
      } else {
        console.warn('未获取到随机诗词')
      }
    } else {
      console.warn('API返回失败:', response.data)
    }
  } catch (e) {
    console.error('加载随机诗词失败:', e)
  } finally {
    loading.value = false
  }
}

const goToCatalog = (category) => {
  router.push({ name: 'Catalog', query: { category: category.id } })
}
</script>

<style scoped>
.home {
  min-height: 100vh;
}

.header {
  background: linear-gradient(135deg, var(--ink-black) 0%, #333 100%);
  color: var(--rice-paper);
  padding: 3rem 0;
  text-align: center;
}

.title {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
}

.subtitle {
  font-size: 1.1rem;
  opacity: 0.9;
}

.section {
  padding: 3rem 0;
}

.section-title {
  text-align: center;
  font-size: 1.8rem;
  margin-bottom: 2rem;
}

.category-card {
  padding: 2rem;
  text-align: center;
}

.category-title {
  font-size: 1.5rem;
  margin-bottom: 1rem;
  color: var(--ink-black);
}

.category-desc {
  color: #666;
  font-size: 0.9rem;
}

.random-poem {
  max-width: 600px;
  margin: 0 auto;
}

.random-actions {
  text-align: center;
  margin-top: 2rem;
}

.error-message {
  text-align: center;
  padding: 2rem;
  color: #999;
}

.error-message .btn {
  margin-top: 1rem;
}
</style>

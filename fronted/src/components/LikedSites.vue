<template>
  <div class="liked-sites-container">
    <!-- <header class="header">
      <h1>我的喜欢</h1>
    </header> -->

    <!-- 导航栏 -->
    <nav class="nav-bar">
      <div class="nav-container">
        <div class="nav-center">
          <router-link to="/sites" class="nav-item">
            <el-icon><List /></el-icon>
            <span>景点列表</span>
          </router-link>
          <router-link to="/liked-sites" class="nav-item active">
            <el-icon><Star /></el-icon>
            <span>我的喜欢</span>
          </router-link>
          <router-link to="/viewed-sites" class="nav-item">
            <el-icon><View /></el-icon>
            <span>我的足迹</span>
          </router-link>
        </div>
        <div class="nav-item user-section">
          <span class="username">{{ userInfo.username }}</span>
          <el-button type="primary" size="small" @click="logout">退出登录</el-button>
        </div>
      </div>
    </nav>

    <main class="main-content">
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-container">
        <el-loading-directive></el-loading-directive>
        <p>加载中...</p>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error" class="error-container">
        <el-alert
          title="加载失败"
          :description="error"
          type="error"
          show-icon
        />
      </div>

      <!-- 空状态 -->
      <div v-else-if="likedSites.length === 0" class="empty-container">
        <div class="empty-content">
          <el-icon class="empty-icon"><Star /></el-icon>
          <h3>还没有点赞的景点</h3>
          <p>快去景点列表发现喜欢的景点吧！</p>
          <el-button type="primary" @click="goToSiteList">浏览景点</el-button>
        </div>
      </div>

      <!-- 点赞景点列表 -->
      <div v-else class="sites-wrapper">
        <div class="sites-grid">
          <div
            v-for="site in likedSites"
            :key="site.id"
            class="site-card"
            @click="goToSiteDetail(site.siteIndex)"
          >
            <div class="site-image">
              <img :src="getFirstImage(site.images)" :alt="site.name" />
              <div class="liked-badge">
                <el-icon><StarFilled /></el-icon>
                <span>已点赞</span>
              </div>
            </div>
            <div class="site-info">
              <h3 class="site-name">{{ site.name }}</h3>
              <p class="site-address">{{ site.address }}</p>
              <div class="site-stats">
                <div class="stat-item">
                  <span class="label">评分</span>
                  <span class="value score">{{ site.score }}</span>
                </div>
                <div class="stat-item">
                  <span class="label">热度</span>
                  <span class="value heat">{{ site.hogDegree }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLikeStore } from '@/stores/likeStore'
import { useSiteDetailStore } from '@/stores/siteDetail'
import { ElMessage } from 'element-plus'
import { List, Star, StarFilled } from '@element-plus/icons-vue'
import api from '@/axios'

const router = useRouter()
const likeStore = useLikeStore()
const siteDetailStore = useSiteDetailStore()

// 响应式数据
const likedSites = ref([])
const loading = ref(false)
const error = ref(null)

// 用户信息
const userInfo = ref({
  username: localStorage.getItem('userInfo') ? JSON.parse(localStorage.getItem('userInfo')).username : '游客'
})

// 获取景点第一张图片
const getFirstImage = (images) => {
  if (!images) return '/placeholder.jpg'
  const imageArray = images.split(',')
  return imageArray[0] || '/placeholder.jpg'
}

// 退出登录
const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('userInfo')
  ElMessage.success('退出登录成功')
  router.push('/login')
}

// 前往景点列表
const goToSiteList = () => {
  router.push('/sites')
}

// 前往景点详情
const goToSiteDetail = (siteIndex) => {
  router.push(`/sites/${siteIndex}`)
}

// 获取用户点赞的景点详细信息
const fetchLikedSitesDetails = async () => {
  loading.value = true
  error.value = null

  try {
    // 首先获取用户点赞的景点索引列表
    const likedSiteIndices = await likeStore.fetchUserLikedSites()

    if (likedSiteIndices.length === 0) {
      likedSites.value = []
      return
    }

    // 根据景点索引列表获取景点详细信息
    const response = await api.post('/site/query/siteList', {
      siteIndexList: likedSiteIndices
    })

    if (response.data.success) {
      likedSites.value = response.data.data.data || []
    } else {
      error.value = response.data.error || '获取点赞景点详情失败'
    }
  } catch (err) {
    error.value = '网络错误，请稍后重试'
    console.error('Error fetching liked sites details:', err)
  } finally {
    loading.value = false
  }
}

// 组件挂载时获取数据
onMounted(() => {
  fetchLikedSitesDetails()
})
</script>

<style scoped>
.liked-sites-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.header {
  background: rgba(255, 255, 255, 0.95);
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.header h1 {
  color: #333;
  font-size: 2rem;
  font-weight: 600;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-info span {
  font-weight: 500;
  color: #666;
}

/* 导航栏样式 */
.nav-bar {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.3);
}

.nav-container {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  padding: 0 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.nav-center {
  display: flex;
  gap: 1rem;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem 1.5rem;
  color: #666;
  text-decoration: none;
  font-weight: 500;
  border-radius: 8px;
  transition: all 0.3s ease;
  position: relative;
}

.nav-item:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.nav-item.active {
  color: #667eea;
  background: rgba(102, 126, 234, 0.1);
}

.nav-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 40px;
  height: 3px;
  background: #667eea;
  border-radius: 2px;
}

.user-section {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-left: auto;
}

.username {
  font-weight: 500;
  color: #666;
  font-size: 0.9rem;
}

.main-content {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.loading-container,
.error-container {
  text-align: center;
  padding: 3rem;
}

/* 空状态样式 */
.empty-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.empty-content {
  text-align: center;
  color: white;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.8;
}

.empty-content h3 {
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
  font-weight: 600;
}

.empty-content p {
  font-size: 1rem;
  margin-bottom: 2rem;
  opacity: 0.9;
}

.sites-grid {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.site-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.18);
  display: flex;
  align-items: center;
  cursor: pointer;
}

.site-card:hover {
  transform: translateX(5px);
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.15);
}

.site-image {
  width: 300px;
  height: 180px;
  overflow: hidden;
  flex-shrink: 0;
  position: relative;
}

.site-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.site-card:hover .site-image img {
  transform: scale(1.05);
}

.liked-badge {
  position: absolute;
  top: 1rem;
  right: 1rem;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  background: rgba(255, 107, 107, 0.9);
  color: white;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
  backdrop-filter: blur(10px);
}

.site-info {
  padding: 1.5rem;
  flex: 1;
}

.site-name {
  font-size: 1.25rem;
  font-weight: 600;
  color: #333;
  line-height: 1.4;
  margin-bottom: 0.5rem;
}

.site-address {
  color: #666;
  font-size: 0.9rem;
  margin-bottom: 1rem;
  line-height: 1.4;
}

.site-stats {
  display: flex;
  gap: 1.5rem;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-item .label {
  font-size: 0.8rem;
  color: #999;
  margin-bottom: 0.25rem;
}

.stat-item .value {
  font-size: 1.1rem;
  font-weight: 600;
}

.stat-item .value.score {
  color: #f39c12;
}

.stat-item .value.heat {
  color: #e74c3c;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header {
    padding: 1rem;
    flex-direction: column;
    gap: 1rem;
  }

  .main-content {
    padding: 1rem;
  }

  .nav-container {
    padding: 0 1rem;
    gap: 1rem;
  }

  .nav-item {
    padding: 0.75rem 1rem;
    font-size: 0.9rem;
  }

  .sites-grid {
    gap: 1rem;
  }

  .site-card {
    flex-direction: column;
  }

  .site-image {
    width: 100%;
    height: 200px;
  }

  .empty-icon {
    font-size: 3rem;
  }

  .empty-content h3 {
    font-size: 1.2rem;
  }
}
</style>
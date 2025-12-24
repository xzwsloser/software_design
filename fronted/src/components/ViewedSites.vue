<template>
  <div class="viewed-sites-container">
    <!-- 导航栏 -->
    <nav class="nav-bar">
      <div class="nav-container">
        <div class="nav-center">
          <router-link to="/sites" class="nav-item">
            <el-icon><List /></el-icon>
            <span>景点列表</span>
          </router-link>
          <router-link to="/liked-sites" class="nav-item">
            <el-icon><Star /></el-icon>
            <span>我的喜欢</span>
          </router-link>
          <router-link to="/viewed-sites" class="nav-item active">
            <el-icon><View /></el-icon>
            <span>我的足迹</span>
          </router-link>
          <router-link to="/data-visualization" class="nav-item">
            <el-icon><DataAnalysis /></el-icon>
            <span>数据可视化</span>
          </router-link>
          <router-link to="/user-info" class="nav-item">
            <el-icon><User /></el-icon>
            <span>我的信息</span>
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
      <div v-else-if="viewedSites.length === 0" class="empty-container">
        <div class="empty-content">
          <el-icon class="empty-icon"><View /></el-icon>
          <h3>还没有浏览过任何景点</h3>
          <p>快去景点列表发现精彩的景点吧！</p>
          <el-button type="primary" @click="goToSiteList">浏览景点</el-button>
        </div>
      </div>

      <!-- 浏览记录列表 -->
      <div v-else class="sites-wrapper">
        <div class="sites-grid">
          <div
            v-for="site in viewedSites"
            :key="site.id"
            class="site-card"
            @click="goToSiteDetail(site.siteIndex)"
          >
            <div class="site-image">
              <img :src="getFirstImage(site.images)" :alt="site.name" />
            </div>
            <div class="site-info">
              <div class="site-header">
                <h3 class="site-name">{{ site.name }}</h3>
                <div class="viewed-badge">
                  <el-icon><View /></el-icon>
                  <span>已浏览</span>
                </div>
              </div>
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
import { useViewStore } from '../stores/viewStore'
import { useSiteDetailStore } from '@/stores/siteDetail'
import { ElMessage } from 'element-plus'
import { View, List, Star, User, DataAnalysis } from '@element-plus/icons-vue'
import api from '../axios'

const router = useRouter()
const viewStore = useViewStore()
const siteDetailStore = useSiteDetailStore()

// 响应式数据
const viewedSites = ref([])
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

// 跳转到景点列表
const goToSiteList = () => {
  router.push('/sites')
}

// 跳转到景点详情页
const goToSiteDetail = (siteIndex) => {
  siteDetailStore.setPrevPageIndex(1)
  router.push(`/sites/${siteIndex}`)
}

// 获取浏览记录的景点详细信息
const fetchViewedSitesDetails = async () => {
  try {
    loading.value = true
    error.value = null

    // 首先获取浏览记录中的景点索引列表
    const viewResponse = await api.get('/view/siteList')

    if (!viewResponse.data.success) {
      throw new Error(viewResponse.data.error || '获取浏览记录失败')
    }

    const siteIndexList = viewResponse.data.data.data || []

    if (siteIndexList.length === 0) {
      viewedSites.value = []
      return
    }

    // 根据景点索引列表获取景点详细信息
    const siteResponse = await api.post('/site/query/siteList', {
      siteIndexList: siteIndexList
    })

    if (!siteResponse.data.success) {
      throw new Error(siteResponse.data.error || '获取景点详情失败')
    }

    viewedSites.value = siteResponse.data.data.data || []
  } catch (err) {
    error.value = err.message || '网络错误，获取浏览记录失败'
    console.error('Error fetching viewed sites:', err)
  } finally {
    loading.value = false
  }
}

// 组件挂载时获取数据
onMounted(() => {
  // 同时获取浏览记录状态和详细信息
  viewStore.fetchViewedSites()
  fetchViewedSitesDetails()
})
</script>

<style scoped>
.viewed-sites-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

/* 导航栏样式 */
.nav-bar {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.3);
  position: sticky;
  top: 0;
  z-index: 100;
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
.error-container,
.empty-container {
  text-align: center;
  padding: 3rem;
}

.empty-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.empty-icon {
  font-size: 4rem;
  color: rgba(255, 255, 255, 0.6);
}

.empty-content h3 {
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
  font-size: 1.5rem;
}

.empty-content p {
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
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

.site-info {
  padding: 1.5rem;
  flex: 1;
}

.site-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 0.5rem;
}

.site-name {
  font-size: 1.25rem;
  font-weight: 600;
  color: #333;
  line-height: 1.4;
  flex: 1;
  margin: 0;
}

.viewed-badge {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  background: rgba(76, 175, 80, 0.1);
  color: #4CAF50;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
  border: 1px solid rgba(76, 175, 80, 0.2);
  flex-shrink: 0;
  margin-left: 0.5rem;
}

.viewed-badge .el-icon {
  font-size: 0.8rem;
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

@media (max-width: 768px) {
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

  .site-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }

  .viewed-badge {
    margin-left: 0;
  }

  .nav-container {
    padding: 0 1rem;
    gap: 1rem;
  }

  .nav-item {
    padding: 0.75rem 1rem;
    font-size: 0.9rem;
  }

  .main-content {
    padding: 1rem;
  }
}
</style>
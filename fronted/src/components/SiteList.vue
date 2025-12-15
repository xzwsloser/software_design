<template>
  <div class="site-list-container">
    <header class="header">
      <h1>景点列表</h1>
      <div class="user-info">
        <span>{{ userInfo.username }}</span>
        <el-button type="primary" @click="logout">退出登录</el-button>
      </div>
    </header>

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

      <!-- 景点列表 -->
      <div v-else class="sites-wrapper">
        <div class="sites-grid">
          <div
            v-for="site in sites"
            :key="site.id"
            class="site-card"
            @click="goToSiteDetail(site.siteIndex)"
          >
            <div class="site-image">
              <img :src="getFirstImage(site.images)" :alt="site.name" />
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

        <!-- 分页组件 -->
        <div class="pagination-wrapper">
          <div class="pagination">
            <!-- 上一页按钮 -->
            <button
              class="pagination-nav-btn"
              :disabled="!hasPrevPage"
              @click="prevPage"
            >
              &lt;
            </button>

            <!-- 页码按钮 -->
            <div class="page-numbers">
              <button
                v-for="page in pageNumbers"
                :key="page"
                :class="[
                  'page-btn',
                  { active: page === currentPage, 'ellipsis': page === '...' }
                ]"
                :disabled="page === '...'"
                @click="goToPage(page)"
              >
                {{ page }}
              </button>
            </div>

            <!-- 下一页按钮 -->
            <button
              class="pagination-nav-btn"
              :disabled="!hasNextPage"
              @click="nextPage"
            >
              &gt;
            </button>
          </div>

          <!-- 页面信息 -->
          <div class="page-info">
            <span>共  1000 条记录</span>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useSiteStore } from '../stores/siteStore'
import { ElMessage } from 'element-plus'

const router = useRouter()
const siteStore = useSiteStore()

// 从store获取响应式数据
const sites = computed(() => siteStore.sites)
const currentPage = computed(() => siteStore.currentPage)
const totalPages = computed(() => siteStore.totalPages)
const total = computed(() => siteStore.total)
const loading = computed(() => siteStore.loading)
const error = computed(() => siteStore.error)
const hasPrevPage = computed(() => siteStore.hasPrevPage)
const hasNextPage = computed(() => siteStore.hasNextPage)
const pageNumbers = computed(() => siteStore.pageNumbers)

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

// 分页方法
const nextPage = () => {
  siteStore.nextPage()
}

const prevPage = () => {
  siteStore.prevPage()
}

const goToPage = (page) => {
  if (page !== '...') {
    siteStore.goToPage(page)
  }
}

// 跳转到景点详情页面
const goToSiteDetail = (siteIndex) => {
  router.push(`/sites/${siteIndex}`)
}

// 组件挂载时获取数据
onMounted(() => {
  siteStore.fetchSites()
})
</script>

<style scoped>
.site-list-container {
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

.sites-grid {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  margin-bottom: 3rem;
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

.site-name {
  font-size: 1.25rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 0.5rem;
  line-height: 1.4;
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

/* 分页样式 */
.pagination-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.pagination {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: rgba(255, 255, 255, 0.9);
  padding: 0.5rem 1rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.pagination-nav-btn {
  width: 32px;
  height: 32px;
  border: 1px solid #ddd;
  background: white;
  color: #666;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 1rem;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pagination-nav-btn:hover:not(:disabled) {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.pagination-nav-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-numbers {
  display: flex;
  gap: 0.25rem;
  margin: 0 0.5rem;
}

.page-btn {
  min-width: 32px;
  height: 32px;
  border: 1px solid #ddd;
  background: white;
  color: #333;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 0.9rem;
  padding: 0 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.page-btn:hover:not(.ellipsis):not(:disabled) {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.page-btn.active {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.page-btn.ellipsis {
  border: none;
  background: transparent;
  cursor: default;
  color: #999;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  color: rgba(255, 255, 255, 0.9);
  font-size: 0.9rem;
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

  .header {
    padding: 1rem;
    flex-direction: column;
    gap: 1rem;
  }

  .main-content {
    padding: 1rem;
  }

  .pagination {
    flex-wrap: wrap;
    justify-content: center;
  }

  .page-info {
    flex-direction: column;
    text-align: center;
    gap: 0.5rem;
  }
}
</style>
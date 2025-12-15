<template>
  <div class="site-detail-container">
    <!-- 头部导航 -->
    <header class="header">
      <div class="header-left">
        <el-button type="text" @click="goBack" class="back-btn">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
      </div>
      <div class="header-title">景点详情</div>
      <div class="header-right">
        <span class="username">{{ userInfo.username }}</span>
        <el-button type="primary" @click="logout" size="small">退出登录</el-button>
      </div>
    </header>

    <!-- 主要内容 -->
    <main class="main-content" v-if="!loading && !error && siteDetail">
      <!-- 景点头部信息 -->
      <section class="site-hero">
        <div class="hero-background">
          <img :src="getFirstImage(siteDetail.images)" :alt="siteDetail.name" />
          <div class="hero-overlay"></div>
        </div>
        <div class="hero-content">
          <div class="hero-info">
            <h1 class="site-name">{{ siteDetail.name }}</h1>
            <div class="site-meta">
              <div class="meta-item">
                <el-icon><Star /></el-icon>
                <span class="score">{{ siteDetail.score }}</span>
              </div>
              <div class="meta-item">
                <el-icon><TrendCharts /></el-icon>
                <span class="heat">{{ siteDetail.hotDegree }}</span>
              </div>
            </div>
            <div class="site-address">
              <el-icon><LocationInformation /></el-icon>
              <span>{{ siteDetail.address }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- 景点详细信息 -->
      <section class="site-info-section">
        <div class="info-container">
          <!-- 景点介绍 -->
          <div class="info-card">
            <h2 class="card-title">
              <el-icon><Document /></el-icon>
              景点介绍
            </h2>
            <div class="introduction" v-html="formatIntroduction(siteDetail.introduce)"></div>
          </div>

          <!-- 实用信息 -->
          <div class="info-card">
            <h2 class="card-title">
              <el-icon><InfoFilled /></el-icon>
              实用信息
            </h2>
            <div class="practical-info">
              <div class="info-item">
                <div class="info-label">
                  <el-icon><Clock /></el-icon>
                  开放时间
                </div>
                <div class="info-value">{{ siteDetail.openTime }}</div>
              </div>
              <div class="info-item">
                <div class="info-label">
                  <el-icon><Phone /></el-icon>
                  联系电话
                </div>
                <div class="info-value">{{ siteDetail.phone }}</div>
              </div>
            </div>
          </div>

          <!-- 景点图片 -->
          <div class="info-card" v-if="imageList.length > 0">
            <h2 class="card-title">
              <el-icon><Picture /></el-icon>
              景点图片
            </h2>
            <div class="image-gallery">
              <div
                v-for="(image, index) in imageList"
                :key="index"
                class="image-item"
                @click="previewImage(image)"
              >
                <img :src="image" :alt="`${siteDetail.name} - 图片${index + 1}`" />
                <div class="image-overlay">
                  <el-icon><ZoomIn /></el-icon>
                </div>
              </div>
            </div>
          </div>

          <!-- 游客评价 -->
          <div class="info-card">
            <CommentSection :siteIndex="siteDetail?.siteIndex" />
          </div>
        </div>
      </section>
    </main>

    <!-- 加载状态 -->
    <div v-else-if="loading" class="loading-container">
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

    <!-- 图片预览对话框 -->
    <el-dialog v-model="imagePreviewVisible" width="80%" :show-close="true">
      <img :src="previewImageUrl" style="width: 100%; height: auto;" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useSiteDetailStore } from '@/stores/siteDetail'
import { ElMessage } from 'element-plus'
import api from '@/axios'
import CommentSection from './CommentSection.vue'
import {
  ArrowLeft,
  Star,
  TrendCharts,
  LocationInformation,
  Document,
  InfoFilled,
  Clock,
  Phone,
  Picture,
  ZoomIn
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

const siteDetailStore = useSiteDetailStore()

// 响应式数据
const siteDetail = ref(null)
const loading = ref(false)
const error = ref(null)
const imagePreviewVisible = ref(false)
const previewImageUrl = ref('')

// 用户信息
const userInfo = ref({
  username: localStorage.getItem('userInfo') ? JSON.parse(localStorage.getItem('userInfo')).username : '游客'
})

// 计算属性
// const imageList = computed(() => {
//   if (!siteDetail.value?.images) return []
//   const images = siteDetail.value.images.split(',').filter(img => img.trim())
//   return [...new Set(images)]
// })

// 计算属性
const imageList = computed(() => {
  if (!siteDetail.value?.images) return []
  
  const images = siteDetail.value.images.split(',').filter(img => img.trim())
  
  // 创建一个Map来存储去重后的图片，key为去除尺寸后的URL
  const uniqueImages = new Map()
  
  images.forEach(img => {
    // 提取基础URL（去除尺寸部分）
    const baseImg = img.replace(/_[WD]_\d+_\d+\.(jpg|jpeg|png|gif)$/i, '')
    
    // 如果是大图(_W_)或者还没有这张图片，则保存
    if (img.includes('_W_') || !uniqueImages.has(baseImg)) {
      uniqueImages.set(baseImg, img)
    }
  })
  
  // 返回去重后的图片数组
  return Array.from(uniqueImages.values())
})

// 获取景点第一张图片
const getFirstImage = (images) => {
  if (!images) return '/placeholder.jpg'
  const imageArray = images.split(',')
  return imageArray[0] || '/placeholder.jpg'
}

// 格式化介绍文本，将换行符转换为HTML换行
const formatIntroduction = (text) => {
  if (!text) return ''
  return text.replace(/\n/g, '<br>')
}

// 预览图片
const previewImage = (image) => {
  previewImageUrl.value = image
  imagePreviewVisible.value = true
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 退出登录
const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('userInfo')
  ElMessage.success('退出登录成功')
  router.push('/login')
}

// 获取景点详情
const fetchSiteDetail = async (siteIndex) => {
  loading.value = true
  error.value = null

  try {
    // const token = localStorage.getItem('token')
    // const response = await fetch(`http://localhost:8079/site/query/${siteIndex}`, {
    //   method: 'GET',
    //   headers: {
    //     'Authorization': `Bearer ${token}`,
    //     'Content-Type': 'application/json'
    //   }
    // })
    const response = await api.get(
      `/site/query/${siteIndex}`
    )

    const data = response.data

    if (data.success) {
      siteDetail.value = data.data
    } else {
      error.value = data.error || '获取景点详情失败'
    }
  } catch (err) {
    error.value = '网络错误，请稍后重试'
    console.error('Error fetching site detail:', err)
  } finally {
    loading.value = false
  }
}

// 组件挂载时获取数据
onMounted(() => {
  const siteIndex = route.params.siteIndex
  if (siteIndex) {
    fetchSiteDetail(siteIndex)
  } else {
    error.value = '缺少景点编号参数'
  }
})
</script>

<style scoped>
.site-detail-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

/* 头部样式 */
.header {
  background: rgba(255, 255, 255, 0.95);
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.back-btn {
  color: #667eea;
  font-size: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.back-btn:hover {
  background: rgba(102, 126, 234, 0.1);
}

.header-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #333;
}

.username {
  font-weight: 500;
  color: #666;
}

/* 主要内容 */
.main-content {
  padding-bottom: 2rem;
}

/* 景点头部英雄区域 */
.site-hero {
  position: relative;
  height: 60vh;
  min-height: 400px;
  overflow: hidden;
}

.hero-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.hero-background img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.hero-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    to bottom,
    rgba(0, 0, 0, 0.2),
    rgba(0, 0, 0, 0.6)
  );
}

.hero-content {
  position: relative;
  height: 100%;
  display: flex;
  align-items: flex-end;
  padding: 2rem;
}

.hero-info {
  color: white;
  max-width: 800px;
}

.site-name {
  font-size: 3rem;
  font-weight: 700;
  margin-bottom: 1rem;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
  line-height: 1.2;
}

.site-meta {
  display: flex;
  gap: 2rem;
  margin-bottom: 1rem;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.2rem;
  font-weight: 600;
}

.meta-item .score {
  color: #ffd700;
}

.meta-item .heat {
  color: #ff6b6b;
}

.site-address {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1rem;
  opacity: 0.9;
}

/* 信息区域 */
.site-info-section {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.info-container {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.info-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.18);
}

.card-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 1.5rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.card-title .el-icon {
  color: #667eea;
}

.introduction {
  line-height: 1.8;
  color: #555;
  font-size: 1.1rem;
}

.practical-info {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.info-item {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
}

.info-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 600;
  color: #333;
  min-width: 120px;
  flex-shrink: 0;
}

.info-label .el-icon {
  color: #667eea;
}

.info-value {
  color: #555;
  line-height: 1.6;
  flex: 1;
}

/* 图片画廊 */
.image-gallery {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.image-item {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  aspect-ratio: 16/9;
}

.image-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.image-item:hover img {
  transform: scale(1.05);
}

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.image-item:hover .image-overlay {
  opacity: 1;
}

.image-overlay .el-icon {
  color: white;
  font-size: 2rem;
}

/* 加载和错误状态 */
.loading-container,
.error-container {
  text-align: center;
  padding: 3rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header {
    padding: 1rem;
  }

  .header-title {
    font-size: 1.2rem;
  }

  .site-hero {
    height: 50vh;
    min-height: 300px;
  }

  .hero-content {
    padding: 1rem;
  }

  .site-name {
    font-size: 2rem;
  }

  .site-meta {
    gap: 1rem;
    flex-direction: column;
  }

  .site-info-section {
    padding: 1rem;
  }

  .info-card {
    padding: 1.5rem;
  }

  .card-title {
    font-size: 1.2rem;
  }

  .practical-info {
    gap: 1rem;
  }

  .info-item {
    flex-direction: column;
    gap: 0.5rem;
  }

  .info-label {
    min-width: auto;
  }

  .image-gallery {
    grid-template-columns: 1fr;
  }
}
</style>
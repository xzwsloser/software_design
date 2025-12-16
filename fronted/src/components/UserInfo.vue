<template>
  <div class="user-info-container">
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
          <router-link to="/viewed-sites" class="nav-item">
            <el-icon><View /></el-icon>
            <span>我的足迹</span>
          </router-link>
          <router-link to="/user-info" class="nav-item active">
            <el-icon><User /></el-icon>
            <span>我的信息</span>
          </router-link>
        </div>
        <div class="nav-item user-section">
          <span class="username">{{ basicUserInfo.username }}</span>
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
        <div class="retry-button">
          <el-button type="primary" @click="fetchUserInfo">重试</el-button>
        </div>
      </div>

      <!-- 用户信息展示 -->
      <div v-else class="user-profile">
        <div class="profile-card">
          <!-- 头像部分 -->
          <div class="avatar-section">
            <div class="avatar-container">
              <img
                src="/src/assets/images/default_avatar.jpg"
                alt="用户头像"
                class="avatar"
              />
            </div>
          </div>

          <!-- 基本信息部分 -->
          <div class="info-section">
            <h2 class="profile-title">个人信息</h2>

            <div class="info-grid">
              <div class="info-item">
                <div class="info-label">用户名</div>
                <div class="info-value">{{ userInfo.username }}</div>
              </div>

              <div class="info-item">
                <div class="info-label">性别</div>
                <div class="info-value">{{ genderText }}</div>
              </div>

              <div class="info-item">
                <div class="info-label">城市</div>
                <div class="info-value">{{ userInfo.city || '未设置' }}</div>
              </div>

              <div class="info-item">
                <div class="info-label">用户ID</div>
                <div class="info-value">{{ userInfo.id }}</div>
              </div>
            </div>

            <div class="info-note">
              <el-icon><InfoFilled /></el-icon>
              <span>当前为只读模式，用户信息暂不支持修改</span>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserInfoStore } from '../stores/userInfo'
import { ElMessage } from 'element-plus'
import { List, Star, View, User, InfoFilled } from '@element-plus/icons-vue'

const router = useRouter()
const userInfoStore = useUserInfoStore()

// 从store获取响应式数据
const userInfo = computed(() => userInfoStore.userInfo)
const loading = computed(() => userInfoStore.loading)
const error = computed(() => userInfoStore.error)
const genderText = computed(() => userInfoStore.genderText)

// 基础用户信息（从localStorage获取）
const basicUserInfo = ref({
  username: localStorage.getItem('userInfo') ? JSON.parse(localStorage.getItem('userInfo')).username : '游客'
})

// 获取用户详细信息
const fetchUserInfo = async () => {
  const result = await userInfoStore.fetchUserInfo()
  if (result.success) {
    ElMessage.success('用户信息加载成功')
  }
}

// 退出登录
const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('userInfo')
  userInfoStore.clearUserInfo()
  ElMessage.success('退出登录成功')
  router.push('/login')
}

// 组件挂载时获取数据
onMounted(() => {
  fetchUserInfo()
})
</script>

<style scoped>
.user-info-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

/* 导航栏样式 - 复用现有样式 */
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
  max-width: 800px;
  margin: 0 auto;
}

.loading-container,
.error-container {
  text-align: center;
  padding: 3rem;
}

.retry-button {
  margin-top: 1rem;
}

.user-profile {
  display: flex;
  justify-content: center;
}

.profile-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.18);
  width: 100%;
  max-width: 500px;
}

.avatar-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 3rem 2rem;
  display: flex;
  justify-content: center;
  align-items: center;
}

.avatar-container {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  overflow: hidden;
  border: 4px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
}

.avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.info-section {
  padding: 2.5rem 2rem;
}

.profile-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 2rem;
  text-align: center;
}

.info-grid {
  display: grid;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-label {
  font-size: 0.9rem;
  color: #666;
  font-weight: 500;
}

.info-value {
  font-size: 1.1rem;
  color: #333;
  font-weight: 600;
  padding: 0.75rem 1rem;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 8px;
  border: 1px solid rgba(102, 126, 234, 0.1);
}

.info-note {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem;
  background: rgba(230, 162, 60, 0.1);
  border: 1px solid rgba(230, 162, 60, 0.2);
  border-radius: 8px;
  color: #e6a23c;
  font-size: 0.9rem;
}

.info-note .el-icon {
  font-size: 1rem;
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .nav-container {
    padding: 0 1rem;
    gap: 0.5rem;
  }

  .nav-item {
    padding: 0.75rem 1rem;
    font-size: 0.9rem;
  }

  .nav-item span {
    display: none;
  }

  .main-content {
    padding: 1rem;
  }

  .profile-card {
    max-width: 100%;
  }

  .info-section {
    padding: 2rem 1.5rem;
  }

  .avatar-section {
    padding: 2rem 1.5rem;
  }

  .avatar-container {
    width: 100px;
    height: 100px;
  }
}
</style>
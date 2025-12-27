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
          <router-link to="/recommended-sites" class="nav-item">
            <el-icon><MagicStick /></el-icon>
            <span>推荐景点</span>
          </router-link>
          <router-link to="/liked-sites" class="nav-item">
            <el-icon><Star /></el-icon>
            <span>我的喜欢</span>
          </router-link>
          <router-link to="/viewed-sites" class="nav-item">
            <el-icon><View /></el-icon>
            <span>我的足迹</span>
          </router-link>
          <router-link to="/data-visualization" class="nav-item">
            <el-icon><DataAnalysis /></el-icon>
            <span>数据可视化</span>
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
            <div class="profile-header">
              <h2 class="profile-title">个人信息</h2>
              <el-button
                v-if="!isEditing"
                type="primary"
                size="small"
                @click="startEdit"
              >
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <div v-else class="edit-actions">
                <el-button size="small" @click="cancelEdit">取消</el-button>
                <el-button type="primary" size="small" @click="saveEdit">保存</el-button>
              </div>
            </div>

            <!-- 只读模式 -->
            <div v-if="!isEditing" class="info-display">
              <div class="info-grid">
                <div class="info-item">
                  <div class="info-label">用户名</div>
                  <div class="info-value">{{ userInfo?.username || '-' }}</div>
                </div>

                <div class="info-item">
                  <div class="info-label">用户ID</div>
                  <div class="info-value">{{ userInfo?.id || '-' }}</div>
                </div>

                <div class="info-item">
                  <div class="info-label">性别</div>
                  <div class="info-value">{{ genderText || '-' }}</div>
                </div>

                <div class="info-item">
                  <div class="info-label">所在省份</div>
                  <div class="info-value">{{ provinceText ? provinceText : '-' }}</div>
                </div>

                <div class="info-item">
                  <div class="info-label">出游类型</div>
                  <div class="info-value">{{ touristTypeText ? touristTypeText : '-' }}</div>
                </div>

                <div class="info-item">
                  <div class="info-label">喜欢的景点类型</div>
                  <div class="info-value">{{ likeTypeText ? likeTypeText : '-' }}</div>
                </div>

                <div class="info-item">
                  <div class="info-label">出游动机</div>
                  <div class="info-value">{{ targetsText ? targetsText : '-' }}</div>
                </div>

                <div class="info-item">
                  <div class="info-label">价格敏感度</div>
                  <div class="info-value">{{ priceSensitiveText ? priceSensitiveText : '-' }}</div>
                </div>

                <div class="info-item">
                  <div class="info-label">体验关注点</div>
                  <div class="info-value">{{ attentionText ? attentionText : '-' }}</div>
                </div>
              </div>
            </div>

            <!-- 编辑模式 -->
            <div v-else class="info-edit">
              <el-form :model="editForm" label-width="100px" class="edit-form">
                <div class="readonly-info">
                  <div class="readonly-item">
                    <span class="readonly-label">用户名：</span>
                    <span class="readonly-value">{{ userInfo?.username || '-' }}</span>
                  </div>
                  <div class="readonly-item">
                    <span class="readonly-label">用户ID：</span>
                    <span class="readonly-value">{{ userInfo?.id || '-' }}</span>
                  </div>
                  <div class="readonly-item">
                    <span class="readonly-label">性别：</span>
                    <span class="readonly-value">{{ genderText || '-' }}</span>
                  </div>
                </div>

                <el-form-item label="所在省份">
                  <el-select v-model="editForm.addressId" placeholder="请选择省份" style="width: 100%">
                    <el-option
                      v-for="(name, id) in dictionaries.provinceDict"
                      :key="id"
                      :label="name"
                      :value="id"
                    />
                  </el-select>
                </el-form-item>

                <el-form-item label="出游类型">
                  <el-select v-model="editForm.touristType" placeholder="请选择出游类型" style="width: 100%">
                    <el-option
                      v-for="(name, id) in dictionaries.touristTypeDict"
                      :key="id"
                      :label="name"
                      :value="id"
                    />
                  </el-select>
                </el-form-item>

                <el-form-item label="喜欢的景点类型">
                  <el-select
                    v-model="selectedLikeTypes"
                    placeholder="请选择喜欢的景点类型(最多6项)"
                    style="width: 100%"
                    multiple
                    collapse-tags
                    collapse-tags-tooltip
                  >
                    <el-option
                      v-for="(name, id) in dictionaries.likeTypeDict"
                      :key="id"
                      :label="name"
                      :value="id"
                    />
                  </el-select>
                </el-form-item>

                <el-form-item label="出游动机">
                  <el-select
                    v-model="selectedTargets"
                    placeholder="请选择出游动机(最多6项)"
                    style="width: 100%"
                    multiple
                    collapse-tags
                    collapse-tags-tooltip
                  >
                    <el-option
                      v-for="(name, id) in dictionaries.targetsDict"
                      :key="id"
                      :label="name"
                      :value="id"
                    />
                  </el-select>
                </el-form-item>

                <el-form-item label="价格敏感度">
                  <el-select v-model="editForm.priceSensitive" placeholder="请选择价格敏感度" style="width: 100%">
                    <el-option label="价格敏感型" :value="0" />
                    <el-option label="价格不敏感" :value="1" />
                  </el-select>
                </el-form-item>

                <el-form-item label="体验关注点">
                  <el-select
                    v-model="selectedAttentions"
                    placeholder="请选择体验关注点(最多4项)"
                    style="width: 100%"
                    multiple
                    collapse-tags
                    collapse-tags-tooltip
                  >
                    <el-option
                      v-for="(name, id) in dictionaries.attentionDict"
                      :key="id"
                      :label="name"
                      :value="id"
                    />
                  </el-select>
                </el-form-item>
              </el-form>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useUserInfoStore } from '../stores/userInfo'
import { ElMessage } from 'element-plus'
import { List, Star, View, User, Edit, DataAnalysis, MagicStick } from '@element-plus/icons-vue'

const router = useRouter()
const userInfoStore = useUserInfoStore()

// 使用 storeToRefs 获取响应式数据
const {
  userInfo,
  loading,
  error,
  genderText,
  provinceText,
  touristTypeText,
  likeTypeText,
  targetsText,
  priceSensitiveText,
  attentionText
} = storeToRefs(userInfoStore)

// 获取字典数据（不是响应式的，直接调用即可）
const dictionaries = computed(() => userInfoStore.getDictionaries())

// 基础用户信息（从localStorage获取）
const basicUserInfo = ref({
  username: localStorage.getItem('userInfo') ? JSON.parse(localStorage.getItem('userInfo')).username : '游客'
})

// 编辑状态
const isEditing = ref(false)

// 编辑表单
const editForm = reactive({
  addressId: '',
  touristType: '',
  likeType: '',
  targets: '',
  priceSensitive: '',
  attention: ''
})

// 多选项的选中值
const selectedLikeTypes = ref([])
const selectedTargets = ref([])
const selectedAttentions = ref([])

// 获取用户详细信息
const fetchUserInfo = async () => {
  const result = await userInfoStore.fetchUserInfo()
  if (result.success) {
    ElMessage.success('用户信息加载成功')
  }
}

// 开始编辑
const startEdit = () => {
  // 确保转换为数字类型
  editForm.addressId = Number(userInfo.value.addressId)
  editForm.touristType = Number(userInfo.value.touristType)
  editForm.priceSensitive = Number(userInfo.value.priceSensitive)

  // 解析多选字段（保持数字类型以确保与字典key匹配）
  selectedLikeTypes.value = userInfo.value.likeType
    ? userInfo.value.likeType.split(',').map(id => Number(id.trim())).filter(id => !isNaN(id))
    : []

  selectedTargets.value = userInfo.value.targets
    ? userInfo.value.targets.split(',').map(id => Number(id.trim())).filter(id => !isNaN(id))
    : []

  selectedAttentions.value = userInfo.value.attention
    ? userInfo.value.attention.split(',').map(id => Number(id.trim())).filter(id => !isNaN(id))
    : []

  isEditing.value = true
}

// 取消编辑
const cancelEdit = () => {
  isEditing.value = false
}

// 保存编辑
const saveEdit = async () => {
  // 验证
  if (selectedLikeTypes.value.length === 0) {
    ElMessage.error('请选择喜欢的景点类型')
    return
  }

  if (selectedLikeTypes.value.length > 6) {
    ElMessage.error('喜欢的景点类型最多选择6项')
    return
  }

  if (selectedTargets.value.length === 0) {
    ElMessage.error('请选择出游动机')
    return
  }

  if (selectedTargets.value.length > 6) {
    ElMessage.error('出游动机最多选择6项')
    return
  }

  if (selectedAttentions.value.length === 0) {
    ElMessage.error('请选择体验关注点')
    return
  }

  if (selectedAttentions.value.length > 4) {
    ElMessage.error('体验关注点最多选择4项')
    return
  }

  // 将多选项转换为逗号分隔的字符串
  editForm.likeType = selectedLikeTypes.value.sort((a, b) => a - b).join(',')
  editForm.targets = selectedTargets.value.sort((a, b) => a - b).join(',')
  editForm.attention = selectedAttentions.value.sort((a, b) => a - b).join(',')

  const result = await userInfoStore.updateUserInfo(editForm)
  if (result.success) {
    ElMessage.success('用户信息更新成功')
    isEditing.value = false
  } else {
    ElMessage.error(result.error || '用户信息更新失败')
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
  max-width: 900px;
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

.profile-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.profile-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.edit-actions {
  display: flex;
  gap: 0.5rem;
}

/* 只读模式样式 */
.info-display {
  width: 100%;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
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
  font-size: 1rem;
  color: #333;
  font-weight: 500;
  padding: 0.75rem 1rem;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 8px;
  border: 1px solid rgba(102, 126, 234, 0.1);
  line-height: 1.5;
  word-break: break-word;
}

/* 编辑模式样式 */
.info-edit {
  width: 100%;
}

.readonly-info {
  padding: 1rem;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 8px;
  margin-bottom: 1.5rem;
}

.readonly-item {
  display: flex;
  margin-bottom: 0.75rem;
}

.readonly-item:last-child {
  margin-bottom: 0;
}

.readonly-label {
  font-size: 0.9rem;
  color: #666;
  font-weight: 500;
  min-width: 80px;
}

.readonly-value {
  font-size: 0.95rem;
  color: #333;
  font-weight: 500;
}

.edit-form {
  max-width: 600px;
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

  .info-grid {
    grid-template-columns: 1fr;
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

  .profile-header {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }
}
</style>

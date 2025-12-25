<template>
  <div class="login-register-container">
    <div class="form-card">
      <!-- 标题 -->
      <div class="title">
        <h2>可视化景点推荐系统</h2>
      </div>

      <!-- Tab切换 -->
      <el-tabs v-model="activeTab" class="tabs" @tab-change="handleTabChange">
        <el-tab-pane label="登录" name="login">
          <el-form
            :model="loginForm"
            label-width="0"
            class="form-content"
          >
            <el-form-item>
              <el-input
                v-model="loginForm.username"
                placeholder="用户名"
                size="large"
              >
                <template #prefix>
                  <el-icon><User /></el-icon>
                </template>
              </el-input>

            </el-form-item>
            <el-form-item>
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="密码"
                size="large"
                show-password
              >
                <template #prefix>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                size="large"
                class="submit-btn"
                @click="handleLogin"
              >
                登录
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="注册" name="register">
          <el-form
            :model="registerForm"
            label-width="0"
            class="form-content"
          >
            <el-form-item>
              <el-input
                v-model="registerForm.username"
                placeholder="用户名"
                size="large"
              >
                <template #prefix>
                  <el-icon><User /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item>
              <el-select
                v-model="registerForm.gender"
                placeholder="请选择性别"
                size="large"
                style="width: 100%"
              >
                <el-option
                  label="男"
                  :value="0"
                />
                <el-option
                  label="女"
                  :value="1"
                />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-select
                v-model="registerForm.addressId"
                placeholder="请选择省份"
                size="large"
                style="width: 100%"
              >
                <el-option
                  v-for="(name, id) in provinceDict"
                  :key="id"
                  :label="name"
                  :value="id"
                />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-select
                v-model="registerForm.touristType"
                placeholder="请选择出游类型"
                size="large"
                style="width: 100%"
              >
                <el-option
                  v-for="(name, id) in touristTypeDict"
                  :key="id"
                  :label="name"
                  :value="id"
                />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-select
                v-model="selectedLikeTypes"
                placeholder="请选择喜欢的景点类型(最多6项)"
                size="large"
                style="width: 100%"
                multiple
                collapse-tags
                collapse-tags-tooltip
              >
                <el-option
                  v-for="(name, id) in likeTypeDict"
                  :key="id"
                  :label="name"
                  :value="id"
                />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-select
                v-model="selectedTargets"
                placeholder="请选择出游动机(最多6项)"
                size="large"
                style="width: 100%"
                multiple
                collapse-tags
                collapse-tags-tooltip
              >
                <el-option
                  v-for="(name, id) in targetsDict"
                  :key="id"
                  :label="name"
                  :value="id"
                />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-select
                v-model="registerForm.priceSensitive"
                placeholder="请选择价格敏感度"
                size="large"
                style="width: 100%"
              >
                <el-option
                  label="价格敏感型"
                  :value="0"
                />
                <el-option
                  label="价格不敏感"
                  :value="1"
                />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-select
                v-model="selectedAttentions"
                placeholder="请选择体验关注点(最多4项)"
                size="large"
                style="width: 100%"
                multiple
                collapse-tags
                collapse-tags-tooltip
              >
                <el-option
                  v-for="(name, id) in attentionDict"
                  :key="id"
                  :label="name"
                  :value="id"
                />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-input
                v-model="registerForm.password"
                type="password"
                placeholder="密码"
                size="large"
                show-password
              >
                <template #prefix>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item>
              <el-input
                v-model="registerForm.confirmPassword"
                type="password"
                placeholder="确认密码"
                size="large"
                show-password
              >
                <template #prefix>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                size="large"
                class="submit-btn"
                @click="handleRegister"
              >
                注册
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'
import { User, Lock } from '@element-plus/icons-vue'

// Tab切换
const activeTab = ref('login')
const isLogin = computed(() => activeTab.value === 'login')

// 使用路由和用户store
const router = useRouter()
const userStore = useUserStore()

// 登录表单
const loginForm = reactive({
  username: '',
  password: '',
})

// 字典数据
const provinceDict = {
  0: '未知',
  1: '河北省',
  2: '山西省',
  3: '辽宁省',
  4: '吉林省',
  5: '黑龙江省',
  6: '江苏省',
  7: '浙江省',
  8: '安徽省',
  9: '福建省',
  10: '江西省',
  11: '山东省',
  12: '河南省',
  13: '湖北省',
  14: '湖南省',
  15: '广东省',
  16: '海南省',
  17: '四川省',
  18: '贵州省',
  19: '云南省',
  20: '陕西省',
  21: '甘肃省',
  22: '青海省',
  23: '台湾省',
  24: '内蒙古自治区',
  25: '广西壮族自治区',
  26: '西藏自治区',
  27: '宁夏回族自治区',
  28: '新疆维吾尔自治区',
  29: '北京市',
  30: '天津市',
  31: '上海市',
  32: '重庆市',
  33: '香港特别行政区',
  34: '澳门特别行政区'
}

const touristTypeDict = {
  0: '其他出游',
  1: '单独旅行',
  2: '商务出差',
  3: '家庭亲子',
  4: '情侣夫妻',
  5: '朋友出游',
  6: '陪同父母'
}

const likeTypeDict = {
  0: '亲子同乐',
  1: '观光游览',
  2: '夜游观景',
  3: '自然风光',
  4: '名胜古迹',
  5: '户外活动',
  6: '展馆展览',
  7: '动植物园',
  8: '冬季滑雪',
  9: '主题乐园',
  10: '休闲娱乐',
  11: '温泉泡汤',
  12: '水上活动',
  13: '空中体验'
}

const targetsDict = {
  0: '其他',
  1: '历史文化溯源',
  2: '自然景观观赏',
  3: '亲子遛娃互动',
  4: '主题乐园狂欢',
  5: '城市地标打卡',
  6: '休闲度假放松',
  7: '网红地标打卡',
  8: '文化艺术体验',
  9: '户外探险猎奇',
  10: '家庭团聚出游',
  11: '治愈系散心',
  12: '节庆主题体验'
}

const attentionDict = {
  0: '排队效率敏感',
  1: '设备完善度敏感',
  2: '服务质量敏感',
  3: '行程规划偏好',
  4: '舒适度敏感',
  5: '导览体验敏感',
  6: '消费透明敏感',
  7: '无障碍设施敏感'
}

// 注册表单
const registerForm = reactive({
  username: '',
  gender: '',
  addressId: '',
  touristType: '',
  likeType: '',
  targets: '',
  priceSensitive: '',
  attention: '',
  password: '',
  confirmPassword: ''
})

// 多选项的选中值
const selectedLikeTypes = ref([])
const selectedTargets = ref([])
const selectedAttentions = ref([])

// Tab切换处理
const handleTabChange = (tabName) => {
  activeTab.value = tabName
}

// 登录处理
const handleLogin = async () => {
  // 表单验证
  if (!loginForm.username || !loginForm.password) {
    ElMessage.error('请填写用户名和密码')
    return
  }

  try {
    const result = await userStore.login(loginForm)
    if (result.success) {
      ElMessage.success('登录成功！')
      // 跳转到景点列表页面
      router.push('/sites')
    } else {
      ElMessage.error(result.error || '登录失败')
    }
  } catch (error) {
    ElMessage.error('登录失败，请稍后重试')
  }
}

// 注册处理
const handleRegister = async () => {
  // 表单验证
  if (!registerForm.username || !registerForm.password) {
    ElMessage.error('请填写用户名和密码')
    return
  }

  if (registerForm.password !== registerForm.confirmPassword) {
    ElMessage.error('两次输入的密码不一致')
    return
  }

  if (registerForm.gender === '') {
    ElMessage.error('请选择性别')
    return
  }

  if (registerForm.addressId === '') {
    ElMessage.error('请选择省份')
    return
  }

  if (registerForm.touristType === '') {
    ElMessage.error('请选择出游类型')
    return
  }

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

  if (registerForm.priceSensitive === '') {
    ElMessage.error('请选择价格敏感度')
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
  registerForm.likeType = selectedLikeTypes.value.sort((a, b) => a - b).join(',')
  registerForm.targets = selectedTargets.value.sort((a, b) => a - b).join(',')
  registerForm.attention = selectedAttentions.value.sort((a, b) => a - b).join(',')

  try {
    const result = await userStore.register(registerForm)
    if (result.success) {
      ElMessage.success('注册成功！')
      // 跳转到景点列表页面
      router.push('/sites')
    } else {
      ElMessage.error(result.error || '注册失败')
    }
  } catch (error) {
    ElMessage.error('注册失败，请稍后重试')
  }
}
</script>

<style scoped>
.login-register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.form-card {
  width: 100%;
  max-width: 450px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  overflow: hidden;
}

.title {
  text-align: center;
  padding: 30px 0 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.title h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 500;
}

.tabs {
  padding: 0 30px 30px;
}

.tabs :deep(.el-tabs__header) {
  margin-bottom: 20px;
}

.tabs :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.tabs :deep(.el-tabs__item) {
  padding: 0 30px;
  font-size: 16px;
  font-weight: 500;
}

.tabs :deep(.el-tabs__item.is-active) {
  color: #667eea;
}

.tabs :deep(.el-tabs__active-bar) {
  background-color: #667eea;
}

.form-content {
  width: 100%;
}

.form-content :deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
  transition: all 0.3s;
}

.form-content :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

.form-content :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2) inset;
}

.submit-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  margin-top: 10px;
}

.submit-btn:hover {
  background: linear-gradient(135deg, #5a72d1 0%, #6a4190 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.submit-btn:active {
  transform: translateY(0);
}

/* 滚动条样式 */
.tabs::-webkit-scrollbar {
  width: 6px;
}

.tabs::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.tabs::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.tabs::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>

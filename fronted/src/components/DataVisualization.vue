<template>
  <div class="data-visualization-container">
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
          <router-link to="/data-visualization" class="nav-item active">
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
      <!-- Tab切换 -->
      <div class="tab-container">
        <div class="tab-buttons">
          <button
            :class="['tab-btn', { active: activeTab === 'heatmap' }]"
            @click="activeTab = 'heatmap'"
          >
            <el-icon><Location /></el-icon>
            <span>景点热力图</span>
          </button>
          <button
            :class="['tab-btn', { active: activeTab === 'piechart' }]"
            @click="activeTab = 'piechart'"
          >
            <el-icon><PieChart /></el-icon>
            <span>旅游动机分布</span>
          </button>
        </div>

        <!-- 热力图内容 -->
        <div v-show="activeTab === 'heatmap'" class="tab-content">
          <div class="chart-card">
            <h2 class="chart-title">中国各省份景点数量分布热力图</h2>
            <div class="heatmap-container" id="heatmap"></div>
            <div class="info-panel">
              <h3>数据概览</h3>
              <p>共有 <strong>34</strong> 个省份/直辖市数据，总景点数量：<strong>987</strong> 个</p>
              <p>景点数量前三的省份：</p>
              <div class="top-provinces">
                <div class="province-tag">云南 <span>72</span></div>
                <div class="province-tag">四川 <span>61</span></div>
                <div class="province-tag">江苏 <span>52</span></div>
              </div>
            </div>
          </div>
        </div>

        <!-- 饼图内容 -->
        <div v-show="activeTab === 'piechart'" class="tab-content">
          <div class="chart-card">
            <div v-if="loading" class="loading-container">
              <el-loading-directive></el-loading-directive>
              <p>加载中...</p>
            </div>
            <div v-else-if="error" class="error-container">
              <el-alert
                title="加载失败"
                :description="error"
                type="error"
                show-icon
              />
            </div>
            <div v-else class="piechart-container">
              <h2 class="chart-title">所有游客旅游动机分布</h2>
              <div class="piechart-image-wrapper">
                <img
                  :src="touristTypeImageUrl"
                  alt="游客旅游动机分布饼图"
                  class="piechart-image"
                />
              </div>
              <div class="info-panel">
                <h3>说明</h3>
                <p>该图表展示了系统中所有游客的旅游动机分布情况，包括不同类型游客的比例统计。</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useVisualizationStore } from '../stores/visualizationStore'
import { ElMessage } from 'element-plus'
import { View, List, Star, User, DataAnalysis, Location, PieChart, MagicStick } from '@element-plus/icons-vue'
import * as echarts from 'echarts'

const router = useRouter()
const visualizationStore = useVisualizationStore()

// 当前激活的Tab
const activeTab = ref('heatmap')

// 用户信息
const userInfo = ref({
  username: localStorage.getItem('userInfo') ? JSON.parse(localStorage.getItem('userInfo')).username : '游客'
})

// 从store获取响应式数据
const touristTypeImageUrl = computed(() => visualizationStore.touristTypeImageUrl)
const loading = computed(() => visualizationStore.loading)
const error = computed(() => visualizationStore.error)

// 退出登录
const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('userInfo')
  ElMessage.success('退出登录成功')
  router.push('/login')
}

// 初始化热力图
const initHeatmap = () => {
  nextTick(() => {
    const container = document.getElementById('heatmap')
    if (!container) return

    const chart = echarts.init(container)

    // 热力图数据
    const rawData = [
      { name: '云南', value: 72 },
      { name: '四川', value: 61 },
      { name: '江苏', value: 52 },
      { name: '浙江', value: 51 },
      { name: '山东', value: 51 },
      { name: '广西', value: 47 },
      { name: '福建', value: 39 },
      { name: '山西', value: 38 },
      { name: '海南', value: 38 },
      { name: '陕西', value: 38 },
      { name: '河南', value: 37 },
      { name: '北京', value: 36 },
      { name: '湖南', value: 34 },
      { name: '新疆', value: 30 },
      { name: '广东', value: 29 },
      { name: '上海', value: 28 },
      { name: '湖北', value: 28 },
      { name: '重庆', value: 26 },
      { name: '江西', value: 25 },
      { name: '贵州', value: 25 },
      { name: '内蒙古', value: 25 },
      { name: '黑龙江', value: 21 },
      { name: '甘肃', value: 21 },
      { name: '河北', value: 19 },
      { name: '辽宁', value: 17 },
      { name: '安徽', value: 17 },
      { name: '香港', value: 16 },
      { name: '西藏', value: 11 },
      { name: '青海', value: 10 },
      { name: '澳门', value: 10 },
      { name: '吉林', value: 9 },
      { name: '天津', value: 7 },
      { name: '宁夏', value: 7 }
    ]

    // 省份简称 -> geojson 标准名映射
    const toGeoName = {
      '北京': '北京市',
      '天津': '天津市',
      '上海': '上海市',
      '重庆': '重庆市',
      '河北': '河北省',
      '山西': '山西省',
      '辽宁': '辽宁省',
      '吉林': '吉林省',
      '黑龙江': '黑龙江省',
      '江苏': '江苏省',
      '浙江': '浙江省',
      '安徽': '安徽省',
      '福建': '福建省',
      '江西': '江西省',
      '山东': '山东省',
      '河南': '河南省',
      '湖北': '湖北省',
      '湖南': '湖南省',
      '广东': '广东省',
      '海南': '海南省',
      '四川': '四川省',
      '贵州': '贵州省',
      '云南': '云南省',
      '陕西': '陕西省',
      '甘肃': '甘肃省',
      '青海': '青海省',
      '内蒙古': '内蒙古自治区',
      '广西': '广西壮族自治区',
      '西藏': '西藏自治区',
      '宁夏': '宁夏回族自治区',
      '新疆': '新疆维吾尔自治区',
      '香港': '香港特别行政区',
      '澳门': '澳门特别行政区'
    }

    // 转换数据
    const data = rawData.map(d => ({
      name: toGeoName[d.name] || d.name,
      value: d.value
    }))

    const option = {
      backgroundColor: '#fff',
      tooltip: {
        trigger: 'item',
        formatter: (params) => {
          const value = (params.value === undefined || params.value === null || Number.isNaN(params.value))
            ? '数据未提供'
            : params.value
          return `<div style="padding:5px;">
            <strong>${params.name}</strong><br>
            景点数量: <span style="color:#e63946; font-size:16px;">${value}</span>个
          </div>`
        }
      },
      visualMap: {
        type: 'continuous',
        min: 0,
        max: 75,
        left: 'left',
        bottom: 'bottom',
        text: ['高', '低'],
        calculable: true,
        inRange: {
          color: ['#e6f7ff', '#bae7ff', '#91d5ff', '#69c0ff', '#40a9ff', '#1890ff', '#096dd9', '#0050b3']
        },
        textStyle: { color: '#333' }
      },
      series: [{
        name: '景点数量',
        type: 'map',
        map: 'china',
        roam: true,
        zoom: 1.2,
        itemStyle: {
          areaColor: '#f5f5f5',
          borderColor: '#999',
          borderWidth: 0.5
        },
        emphasis: {
          itemStyle: { areaColor: '#ffcc00', borderWidth: 2 },
          label: { show: true, color: '#000', fontWeight: 'bold' }
        },
        label: { show: true, fontSize: 12, color: '#666' },
        data
      }]
    }

    // 加载地图数据
    fetch('https://geo.datav.aliyun.com/areas_v3/bound/100000_full.json')
      .then(r => {
        if (!r.ok) throw new Error('地图数据加载失败: ' + r.status)
        return r.json()
      })
      .then(chinaJson => {
        echarts.registerMap('china', chinaJson)
        chart.setOption(option)
      })
      .catch(err => {
        console.error('获取地图数据失败:', err)
        ElMessage.error('地图加载失败，请刷新页面重试')
      })

    // 窗口大小变化时重新渲染
    const resizeHandler = () => chart.resize()
    window.addEventListener('resize', resizeHandler)

    // 保存resizeHandler用于清理
    container._resizeHandler = resizeHandler
  })
}

// 清理热力图
const cleanupHeatmap = () => {
  const container = document.getElementById('heatmap')
  if (container && container._resizeHandler) {
    window.removeEventListener('resize', container._resizeHandler)
    delete container._resizeHandler
  }
}

// 监听Tab切换
watch(activeTab, (newTab) => {
  if (newTab === 'heatmap') {
    initHeatmap()
  } else {
    cleanupHeatmap()
  }
})

// 组件挂载时
onMounted(() => {
  // 初始化热力图
  initHeatmap()
  // 获取游客类型图片
  visualizationStore.fetchTouristTypeImage()
})

// 组件卸载时
onUnmounted(() => {
  cleanupHeatmap()
})
</script>

<style scoped>
.data-visualization-container {
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

/* Tab容器样式 */
.tab-container {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.18);
}

.tab-buttons {
  display: flex;
  gap: 0.5rem;
  padding: 1rem 1rem 0;
  background: rgba(255, 255, 255, 0.5);
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border: none;
  background: transparent;
  color: #666;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  border-radius: 8px 8px 0 0;
  transition: all 0.3s ease;
  position: relative;
}

.tab-btn:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.tab-btn.active {
  background: rgba(255, 255, 255, 0.95);
  color: #667eea;
}

.tab-btn.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 2px;
  background: #667eea;
}

.tab-content {
  padding: 1.5rem;
}

.chart-card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
}

.chart-title {
  text-align: center;
  color: #333;
  font-size: 1.5rem;
  margin: 0 0 1.5rem;
  font-weight: 600;
}

/* 热力图样式 */
.heatmap-container {
  width: 100%;
  height: 700px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

/* 饼图样式 */
.piechart-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.piechart-image-wrapper {
  width: 100%;
  max-width: 800px;
  margin: 0 auto 1.5rem;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  background: #fff;
}

.piechart-image {
  width: 100%;
  height: auto;
  display: block;
}

/* 信息面板样式 */
.info-panel {
  margin-top: 1.5rem;
  padding: 1.5rem;
  background: #f8f9fa;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.info-panel h3 {
  margin: 0 0 1rem;
  color: #1890ff;
  font-size: 1.1rem;
}

.info-panel p {
  margin: 0.5rem 0;
  color: #666;
  line-height: 1.6;
}

.top-provinces {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  margin-top: 1rem;
}

.province-tag {
  background: #e6f4ff;
  color: #1890ff;
  padding: 0.5rem 1rem;
  border-radius: 15px;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  font-weight: 500;
}

.province-tag span {
  font-weight: 700;
  margin-left: 0.5rem;
}

.loading-container,
.error-container {
  text-align: center;
  padding: 3rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .nav-container {
    padding: 0 1rem;
    gap: 0.5rem;
  }

  .nav-item {
    padding: 0.75rem 0.75rem;
    font-size: 0.85rem;
  }

  .nav-item span {
    display: none;
  }

  .main-content {
    padding: 1rem;
  }

  .tab-buttons {
    flex-direction: column;
    padding: 0.5rem;
  }

  .tab-btn {
    width: 100%;
    justify-content: center;
    border-radius: 8px;
  }

  .heatmap-container {
    height: 400px;
  }

  .chart-title {
    font-size: 1.2rem;
  }

  .info-panel {
    padding: 1rem;
  }

  .top-provinces {
    justify-content: center;
  }
}
</style>

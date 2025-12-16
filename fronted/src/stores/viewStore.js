import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/axios'

export const useViewStore = defineStore('view', () => {
  // 状态
  const viewedSites = ref(new Set()) // 使用Set存储已浏览的景点siteIndex
  const loading = ref(false)
  const error = ref(null)

  // 计算属性
  const isViewed = (siteIndex) => {
    return viewedSites.value.has(siteIndex)
  }

  const viewedSitesList = computed(() => {
    return Array.from(viewedSites.value)
  })

  // 行为
  // 记录浏览
  const recordView = async (siteIndex) => {
    if (!siteIndex) return false

    try {
      loading.value = true
      error.value = null

      const response = await api.get(`/view/view/${siteIndex}`)

      if (response.data.success) {
        // 如果返回的data为true，表示是第一次浏览；false表示已经浏览过了
        // 无论哪种情况，都将该景点标记为已浏览
        viewedSites.value.add(siteIndex)
        return response.data.data // 返回是否为首次浏览
      } else {
        error.value = response.data.error || '记录浏览失败'
        return false
      }
    } catch (err) {
      error.value = '网络错误，记录浏览失败'
      console.error('Error recording view:', err)
      return false
    } finally {
      loading.value = false
    }
  }

  // 获取用户浏览过的景点列表
  const fetchViewedSites = async () => {
    try {
      loading.value = true
      error.value = null

      const response = await api.get('/view/siteList')

      if (response.data.success) {
        const siteIndexList = response.data.data.data || []
        viewedSites.value = new Set(siteIndexList)
        return siteIndexList
      } else {
        error.value = response.data.error || '获取浏览记录失败'
        return []
      }
    } catch (err) {
      error.value = '网络错误，获取浏览记录失败'
      console.error('Error fetching viewed sites:', err)
      return []
    } finally {
      loading.value = false
    }
  }

  // 检查景点是否已浏览（用于初始化时检查特定景点）
  const checkSiteViewedStatus = async (siteIndex) => {
    if (!siteIndex) return false

    // 如果已经在本地缓存中，直接返回
    if (viewedSites.value.has(siteIndex)) {
      return true
    }

    // 否则获取完整的浏览列表并检查
    await fetchViewedSites()
    return viewedSites.value.has(siteIndex)
  }

  // 清除浏览记录（本地状态）
  const clearViewedSites = () => {
    viewedSites.value.clear()
  }

  return {
    // 状态
    viewedSites,
    loading,
    error,

    // 计算属性
    isViewed,
    viewedSitesList,

    // 行为
    recordView,
    fetchViewedSites,
    checkSiteViewedStatus,
    clearViewedSites
  }
})
import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/axios'

export const useRecommendationStore = defineStore('recommendation', () => {
  // 状态
  const recommendedSites = ref([])
  const loading = ref(false)
  const error = ref(null)

  // 获取用户推荐景点列表
  const fetchRecommendedSites = async () => {
    if (loading.value) return []

    loading.value = true
    error.value = null

    try {
      const response = await api.get('/rec/siteIdxList')

      if (response.data.success) {
        const siteIndices = response.data.data.data || []
        recommendedSites.value = siteIndices
        return siteIndices
      } else {
        error.value = response.data.error || '获取推荐景点失败'
        return []
      }
    } catch (err) {
      error.value = '网络错误，请稍后重试'
      console.error('获取推荐景点失败:', err)
      return []
    } finally {
      loading.value = false
    }
  }

  // 重置状态
  const resetState = () => {
    recommendedSites.value = []
    loading.value = false
    error.value = null
  }

  return {
    // 状态
    recommendedSites,
    loading,
    error,

    // 方法
    fetchRecommendedSites,
    resetState
  }
})

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/axios'

export const useLikeStore = defineStore('like', () => {
  // 状态
  const likedSites = ref(new Set()) // 存储已点赞的景点siteIndex
  const likeLoading = ref(false)
  const error = ref(null)

  // 计算属性
  const isSiteLiked = (siteIndex) => {
    return likedSites.value.has(siteIndex)
  }

  // 点赞景点
  const likeSite = async (siteIndex) => {
    if (likeLoading.value) return false

    likeLoading.value = true
    error.value = null

    try {
      const response = await api.get(`/like/like/${siteIndex}`)

      if (response.data.success) {
        likedSites.value.add(siteIndex)
        return true
      } else {
        error.value = response.data.error || '点赞失败'
        return false
      }
    } catch (err) {
      error.value = '网络错误，请稍后重试'
      console.error('点赞失败:', err)
      return false
    } finally {
      likeLoading.value = false
    }
  }

  // 取消点赞
  const cancelLikeSite = async (siteIndex) => {
    if (likeLoading.value) return false

    likeLoading.value = true
    error.value = null

    try {
      const response = await api.get(`/like/cancel/${siteIndex}`)

      if (response.data.success) {
        likedSites.value.delete(siteIndex)
        return true
      } else {
        error.value = response.data.error || '取消点赞失败'
        return false
      }
    } catch (err) {
      error.value = '网络错误，请稍后重试'
      console.error('取消点赞失败:', err)
      return false
    } finally {
      likeLoading.value = false
    }
  }

  // 切换点赞状态（点赞/取消点赞）
  const toggleLike = async (siteIndex) => {
    if (isSiteLiked(siteIndex)) {
      return await cancelLikeSite(siteIndex)
    } else {
      return await likeSite(siteIndex)
    }
  }

  // 检查景点是否已点赞
  const checkSiteLikeStatus = async (siteIndex) => {
    try {
      const response = await api.get(`/like/isLike/${siteIndex}`)

      if (response.data.success) {
        if (response.data.data) {
          likedSites.value.add(siteIndex)
        } else {
          likedSites.value.delete(siteIndex)
        }
        return response.data.data
      }
    } catch (err) {
      console.error('检查点赞状态失败:', err)
    }
    return false
  }

  // 获取用户点赞的景点列表
  const fetchUserLikedSites = async () => {
    try {
      const response = await api.get('/like/siteList')

      if (response.data.success) {
        const siteIndices = response.data.data.data
        likedSites.value = new Set(siteIndices)
        return siteIndices
      }
    } catch (err) {
      console.error('获取点赞列表失败:', err)
    }
    return []
  }

  // 重置状态
  const resetState = () => {
    likedSites.value.clear()
    likeLoading.value = false
    error.value = null
  }

  return {
    // 状态
    likedSites,
    likeLoading,
    error,

    // 计算属性
    isSiteLiked,

    // 方法
    likeSite,
    cancelLikeSite,
    toggleLike,
    checkSiteLikeStatus,
    fetchUserLikedSites,
    resetState
  }
})
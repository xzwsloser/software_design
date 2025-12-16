import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../axios'

export const useUserInfoStore = defineStore('userInfo', () => {
  const userInfo = ref({
    id: 0,
    username: '',
    gender: 0, // 0 -> 男, 1 -> 女
    city: ''
  })

  const loading = ref(false)
  const error = ref('')

  // 获取用户详细信息
  const fetchUserInfo = async () => {
    try {
      loading.value = true
      error.value = ''

      const response = await api.get('/userInfo/user')

      if (response.data.success) {
        userInfo.value = response.data.data
        return { success: true }
      } else {
        error.value = response.data.error
        return { success: false, error: response.data.error }
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
      const errorMsg = error.response?.data?.error || '获取用户信息失败，请稍后重试'
      error.value = errorMsg
      return { success: false, error: errorMsg }
    } finally {
      loading.value = false
    }
  }

  // 清空用户信息
  const clearUserInfo = () => {
    userInfo.value = {
      id: 0,
      username: '',
      gender: 0,
      city: ''
    }
    error.value = ''
  }

  // 格式化性别显示
  const genderText = computed(() => {
    return userInfo.value.gender === 0 ? '男' : '女'
  })

  return {
    userInfo,
    loading,
    error,
    genderText,
    fetchUserInfo,
    clearUserInfo
  }
})
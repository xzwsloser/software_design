import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../axios'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))

  const isLoggedIn = computed(() => !!token.value)

  const login = async (loginForm) => {
    try {
      const response = await api.post('/user/login', {
        id: 0,
        username: loginForm.username,
        password: loginForm.password,
        gender: 0,
        city: ''
      })

      if (response.data.success) {
        token.value = response.data.data
        localStorage.setItem('token', token.value)

        const decodedToken = decodeJWT(token.value)
        userInfo.value = {
          id: decodedToken.id,
          username: decodedToken.username
        }
        localStorage.setItem('userInfo', JSON.stringify(userInfo.value))

        return { success: true }
      } else {
        return { success: false, error: response.data.error }
      }
    } catch (error) {
      console.error('登录失败:', error)
      return {
        success: false,
        error: error.response?.data?.error || '登录失败，请稍后重试'
      }
    }
  }

  const register = async (registerForm) => {
    try {
      const city = registerForm.province && registerForm.city
        ? `${registerForm.province}${registerForm.city}`
        : registerForm.city || ''

      const response = await api.post('/user/register', {
        id: 0,
        username: registerForm.username,
        password: registerForm.password,
        gender: registerForm.gender,
        city: city
      })

      if (response.data.success) {
        token.value = response.data.data
        localStorage.setItem('token', token.value)

        const decodedToken = decodeJWT(token.value)
        userInfo.value = {
          id: decodedToken.id,
          username: decodedToken.username
        }
        localStorage.setItem('userInfo', JSON.stringify(userInfo.value))

        return { success: true }
      } else {
        return { success: false, error: response.data.error }
      }
    } catch (error) {
      console.error('注册失败:', error)
      return {
        success: false,
        error: error.response?.data?.error || '注册失败，请稍后重试'
      }
    }
  }

  const logout = () => {
    token.value = ''
    userInfo.value = {}
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  function decodeJWT(token) {
    try {
      const base64Url = token.split('.')[1]
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split('')
          .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
          .join('')
      )
      return JSON.parse(jsonPayload)
    } catch (error) {
      console.error('JWT解码失败:', error)
      return {}
    }
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    login,
    register,
    logout
  }
})
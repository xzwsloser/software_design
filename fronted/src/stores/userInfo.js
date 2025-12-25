import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../axios'

// 省份字典
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

// 出游类型字典
const touristTypeDict = {
  0: '其他出游',
  1: '单独旅行',
  2: '商务出差',
  3: '家庭亲子',
  4: '情侣夫妻',
  5: '朋友出游',
  6: '陪同父母'
}

// 景点喜好类型字典
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

// 出游动机字典
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

// 价格敏感字典 (0: 价格敏感型, 1: 价格不敏感)
const priceSensitiveDict = {
  0: '价格敏感型',
  1: '价格不敏感'
}

// 体验关注细节字典
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

export const useUserInfoStore = defineStore('userInfo', () => {
  const userInfo = ref({
    id: 0,
    username: '',
    password: '',
    gender: 0,
    addressId: 0,
    touristType: 0,
    likeType: '',
    targets: '',
    priceSensitive: 0,
    attention: ''
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

  // 更新用户信息
  const updateUserInfo = async (updateData) => {
    try {
      loading.value = true
      error.value = ''

      const response = await api.post('/userInfo/update', {
        id: userInfo.value.id,
        username: userInfo.value.username,
        password: userInfo.value.password,
        gender: Number(userInfo.value.gender),
        addressId: Number(updateData.addressId),
        touristType: Number(updateData.touristType),
        likeType: updateData.likeType,
        targets: updateData.targets,
        priceSensitive: Number(updateData.priceSensitive),
        attention: updateData.attention
      })

      if (response.data.success) {
        // 更新本地userInfo
        userInfo.value = {
          ...userInfo.value,
          addressId: updateData.addressId,
          touristType: updateData.touristType,
          likeType: updateData.likeType,
          targets: updateData.targets,
          priceSensitive: updateData.priceSensitive,
          attention: updateData.attention
        }
        return { success: true }
      } else {
        error.value = response.data.error
        return { success: false, error: response.data.error }
      }
    } catch (error) {
      console.error('更新用户信息失败:', error)
      const errorMsg = error.response?.data?.error || '更新用户信息失败，请稍后重试'
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
      password: '',
      gender: 0,
      addressId: 0,
      touristType: 0,
      likeType: '',
      targets: '',
      priceSensitive: 0,
      attention: ''
    }
    error.value = ''
  }

  // 格式化性别显示
  const genderText = computed(() => {
    return userInfo.value.gender === 0 ? '男' : '女'
  })

  // 格式化省份显示
  const provinceText = computed(() => {
    return provinceDict[Number(userInfo.value.addressId)] || '未知'
  })

  // 格式化出游类型显示
  const touristTypeText = computed(() => {
    return touristTypeDict[Number(userInfo.value.touristType)] || '未知'
  })

  // 格式化景点喜好类型显示
  const likeTypeText = computed(() => {
    if (!userInfo.value.likeType) return '未设置'
    const ids = userInfo.value.likeType.split(',').map(id => parseInt(id.trim()))
    return ids.map(id => likeTypeDict[id] || '未知').join('、')
  })

  // 格式化出游动机显示
  const targetsText = computed(() => {
    if (!userInfo.value.targets) return '未设置'
    const ids = userInfo.value.targets.split(',').map(id => parseInt(id.trim()))
    return ids.map(id => targetsDict[id] || '未知').join('、')
  })

  // 格式化价格敏感显示
  const priceSensitiveText = computed(() => {
    return priceSensitiveDict[Number(userInfo.value.priceSensitive)] || '未知'
  })

  // 格式化体验关注细节显示
  const attentionText = computed(() => {
    if (!userInfo.value.attention) return '未设置'
    const ids = userInfo.value.attention.split(',').map(id => parseInt(id.trim()))
    return ids.map(id => attentionDict[id] || '未知').join('、')
  })

  // 获取字典数据供组件使用
  const getDictionaries = () => {
    return {
      provinceDict,
      touristTypeDict,
      likeTypeDict,
      targetsDict,
      priceSensitiveDict,
      attentionDict
    }
  }

  return {
    userInfo,
    loading,
    error,
    genderText,
    provinceText,
    touristTypeText,
    likeTypeText,
    targetsText,
    priceSensitiveText,
    attentionText,
    getDictionaries,
    fetchUserInfo,
    updateUserInfo,
    clearUserInfo
  }
})
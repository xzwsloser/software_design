import { defineStore } from 'pinia'
import api from '../axios'

export const useCommentStore = defineStore('comment', {
  state: () => ({
    positiveComments: {
      data: [],
      currentPage: 1,
      pageSize: 10,
      total: 100,
      loading: false,
      error: null
    },
    negativeComments: {
      data: [],
      currentPage: 1,
      pageSize: 10,
      total: 50,
      loading: false,
      error: null
    },
    currentSiteIndex: null,
    defaultPositiveCommentCount: 100,
    defaultNegativeCommentCount: 50
  }),

  getters: {
    getPositiveCommentsByPage: (state) => {
      return state.positiveComments.data
    },
    getNegativeCommentsByPage: (state) => {
      return state.negativeComments.data
    },
    positiveTotalPages: (state) => {
      return Math.ceil(state.positiveComments.total / state.positiveComments.pageSize)
    },
    negativeTotalPages: (state) => {
      return Math.ceil(state.negativeComments.total / state.negativeComments.pageSize)
    }
  },

  actions: {
    // 获取好评数量
    async fetchCountPositiveComments(siteIndex) {
      try {
        const response = await api.get(`/comment/count/positive/${siteIndex}`)

        if (response.data.success) {
          this.positiveComments.total = response.data.data
          this.defaultPositiveCommentCount = response.data.data
        } else {
          this.positiveComments.error = response.data.error || '获取好评数量失败'
        }
      } catch(error) {
        this.positiveComments.error = '网络错误, 请稍后重试'
        console.error('获取好评数失败', error)
      }
    },
    // 获取差评数量
    async fetchCountNegativeComments(siteIndex) {
      try {
        const response = await api.get(`/comment/count/negative/${siteIndex}`)

        if (response.data.success) {
          this.negativeComments.total = response.data.data
          this.defaultNegativeCommentCount = response.data.data
        } else {
          this.negativeComments.error = response.data.error || '获取差评数量失败'
        }
      } catch(error) {
        this.positiveComments.error = '网络错误, 请稍后重试'
        console.error('获取差评数失败', error)
      }
    },
    // 获取好评
    async fetchPositiveComments(siteIndex, pageIndex = 1) {
      if (this.currentSiteIndex !== siteIndex) {
        this.currentSiteIndex = siteIndex
        // 重置分页状态
        this.positiveComments.currentPage = 1
        this.positiveComments.data = []
      }

      this.positiveComments.loading = true
      this.positiveComments.error = null

      try {
        const response = await api.post(`/comment/positive/${siteIndex}`, {
          pageIndex,
          pageSize: this.positiveComments.pageSize
        })

        if (response.data.success) {
          this.positiveComments.data = response.data.data.data
          // this.positiveComments.total = response.data.data.total
          this.positiveComments.currentPage = pageIndex
        } else {
          this.positiveComments.error = response.data.error || '获取好评失败'
        }
      } catch (error) {
        this.positiveComments.error = '网络错误，请稍后重试'
        console.error('获取好评失败:', error)
      } finally {
        this.positiveComments.loading = false
      }
    },

    // 获取差评
    async fetchNegativeComments(siteIndex, pageIndex = 1) {
      if (this.currentSiteIndex !== siteIndex) {
        this.currentSiteIndex = siteIndex
        // 重置分页状态
        this.negativeComments.currentPage = 1
        this.negativeComments.data = []
      }

      this.negativeComments.loading = true
      this.negativeComments.error = null

      try {
        const response = await api.post(`/comment/negative/${siteIndex}`, {
          pageIndex,
          pageSize: this.negativeComments.pageSize
        })

        if (response.data.success) {
          this.negativeComments.data = response.data.data.data
          // this.negativeComments.total = response.data.data.total
          this.negativeComments.currentPage = pageIndex
        } else {
          this.negativeComments.error = response.data.error || '获取差评失败'
        }
      } catch (error) {
        this.negativeComments.error = '网络错误，请稍后重试'
        console.error('获取差评失败:', error)
      } finally {
        this.negativeComments.loading = false
      }
    },

    // 切换好评页码
    async changePositivePage(pageIndex) {
      if (this.currentSiteIndex && pageIndex !== this.positiveComments.currentPage) {
        await this.fetchPositiveComments(this.currentSiteIndex, pageIndex)
      }
    },

    // 切换差评页码
    async changeNegativePage(pageIndex) {
      if (this.currentSiteIndex && pageIndex !== this.negativeComments.currentPage) {
        await this.fetchNegativeComments(this.currentSiteIndex, pageIndex)
      }
    },

    // 重置评论状态
    resetComments() {
      this.positiveComments = {
        data: [],
        currentPage: 1,
        pageSize: 10,
        total: this.defaultPositiveCommentCount,
        loading: false,
        error: null
      }
      this.negativeComments = {
        data: [],
        currentPage: 1,
        pageSize: 10,
        total: this.defaultNegativeCommentCount,
        loading: false,
        error: null
      }
      this.currentSiteIndex = null
    }
  }
})
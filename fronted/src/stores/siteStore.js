import { defineStore } from 'pinia'
import api from '../axios'

export const useSiteStore = defineStore('site', {
  state: () => ({
    sites: [],
    currentPage: 1,
    pageSize: 10,
    total: 1000,
    loading: false,
    error: null
  }),

  getters: {
    totalPages: (state) => Math.ceil(1000 / state.pageSize),
    hasNextPage: (state) => state.currentPage < Math.ceil(1000 / state.pageSize),
    hasPrevPage: (state) => state.currentPage > 1,
    // 生成页码数组，用于显示页码按钮
    pageNumbers: (state) => {
      const totalPages = Math.ceil(1000 / state.pageSize)
      const current = state.currentPage
      const delta = 2 // 当前页前后显示的页码数量

      let range = []
      let rangeWithDots = []
      let l

      for (let i = 1; i <= totalPages; i++) {
        if (i == 1 || i == totalPages || (i >= current - delta && i <= current + delta)) {
          range.push(i)
        }
      }

      range.forEach((i) => {
        if (l) {
          if (i - l === 2) {
            rangeWithDots.push(l + 1)
          } else if (i - l !== 1) {
            rangeWithDots.push('...')
          }
        }
        rangeWithDots.push(i)
        l = i
      })

      return rangeWithDots
    }
  },

  actions: {
    async fetchSites(pageIndex = 1, pageSize = 10) {
      this.loading = true
      this.error = null

      try {
        const response = await api.post('/site/query/list', {
          pageIndex,
          pageSize
        })

        if (response.data.success) {
          this.sites = response.data.data.data
          this.total = response.data.data.total
          this.currentPage = pageIndex
          this.pageSize = pageSize
        } else {
          this.error = response.data.error || 'Failed to fetch sites'
        }
      } catch (error) {
        this.error = error.message || 'Network error'
        console.error('Error fetching sites:', error)
      } finally {
        this.loading = false
      }
    },

    async nextPage() {
      if (this.hasNextPage) {
        await this.fetchSites(this.currentPage + 1, this.pageSize)
      }
    },

    async prevPage() {
      if (this.hasPrevPage) {
        await this.fetchSites(this.currentPage - 1, this.pageSize)
      }
    },

    async goToPage(page) {
      const totalPages = Math.ceil(1000 / this.pageSize)
      if (page >= 1 && page <= totalPages && page !== '...') {
        await this.fetchSites(page, this.pageSize)
      }
    }
  }
})
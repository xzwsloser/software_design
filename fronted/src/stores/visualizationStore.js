import { defineStore } from 'pinia'
import api from '../axios'

export const useVisualizationStore = defineStore('visualization', {
  state: () => ({
    touristTypeImageUrl: '',
    loading: false,
    error: null
  }),

  getters: {
    hasImage: (state) => !!state.touristTypeImageUrl
  },

  actions: {
    async fetchTouristTypeImage() {
      this.loading = true
      this.error = null

      try {
        const response = await api.get('/oss/touristType')

        if (response.data.success) {
          this.touristTypeImageUrl = response.data.data
        } else {
          this.error = response.data.error || 'Failed to fetch tourist type image'
        }
      } catch (error) {
        this.error = error.message || 'Network error'
        console.error('Error fetching tourist type image:', error)
      } finally {
        this.loading = false
      }
    }
  }
})

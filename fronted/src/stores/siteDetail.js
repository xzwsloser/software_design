import { defineStore } from "pinia"
import { ref } from 'vue'

/**
 *  @Description: 用于详细页与景点列表页之间通信
 */

export const useSiteDetailStore = defineStore('siteDetail', () => {
    // 景点列表中对应页
    const prevPageIdxInList = ref(1)

    function setPrevPageIndex(pageIndex) {
        prevPageIdxInList.value = pageIndex
    }

      return {
        prevPageIdxInList,
        setPrevPageIndex
    }
})

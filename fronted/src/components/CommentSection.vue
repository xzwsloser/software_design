<template>
  <div class="comment-section">
    <!-- 评论区标题 -->
    <div class="section-header">
      <h2 class="section-title">
        <el-icon><ChatDotRound /></el-icon>
        游客评价
      </h2>
    </div>

    <!-- 评价类型选择器 -->
    <div class="comment-tabs">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="好评" name="positive">
          <div class="comment-count">共 {{ commentStore.positiveComments.total }} 条好评</div>
        </el-tab-pane>
        <el-tab-pane label="差评" name="negative">
          <div class="comment-count">共 {{ commentStore.negativeComments.total }} 条差评</div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 评论内容区域 -->
    <div class="comment-content">
      <!-- 好评内容 -->
      <div v-if="activeTab === 'positive'">
        <div v-if="commentStore.positiveComments.loading" class="loading-state">
          <el-skeleton :rows="3" animated />
          <el-skeleton :rows="3" animated style="margin-top: 1rem;" />
          <el-skeleton :rows="3" animated style="margin-top: 1rem;" />
        </div>

        <div v-else-if="commentStore.positiveComments.error" class="error-state">
          <el-alert
            title="加载失败"
            :description="commentStore.positiveComments.error"
            type="error"
            show-icon
          />
        </div>

        <div v-else-if="positiveComments.length === 0" class="empty-state">
          <el-empty description="暂无好评" />
        </div>

        <div v-else class="comment-list">
          <div
            v-for="comment in positiveComments"
            :key="comment.id"
            class="comment-item positive-comment"
          >
            <div class="comment-wrapper">
              <img
                src="/src/assets/images/default_avatar.jpg"
                alt="用户头像"
                class="user-avatar"
              />
              <div class="comment-content">
                <div class="comment-text">{{ comment.content }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 差评内容 -->
      <div v-if="activeTab === 'negative'">
        <div v-if="commentStore.negativeComments.loading" class="loading-state">
          <el-skeleton :rows="3" animated />
          <el-skeleton :rows="3" animated style="margin-top: 1rem;" />
          <el-skeleton :rows="3" animated style="margin-top: 1rem;" />
        </div>

        <div v-else-if="commentStore.negativeComments.error" class="error-state">
          <el-alert
            title="加载失败"
            :description="commentStore.negativeComments.error"
            type="error"
            show-icon
          />
        </div>

        <div v-else-if="negativeComments.length === 0" class="empty-state">
          <el-empty description="暂无差评" />
        </div>

        <div v-else class="comment-list">
          <div
            v-for="comment in negativeComments"
            :key="comment.id"
            class="comment-item negative-comment"
          >
            <div class="comment-wrapper">
              <img
                src="/src/assets/images/default_avatar.jpg"
                alt="用户头像"
                class="user-avatar"
              />
              <div class="comment-content">
                <div class="comment-text">{{ comment.content }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页组件 -->
    <div class="pagination-container" v-if="showPagination">
      <div class="pagination-info">
        第 {{ currentPage }} 页，共 {{ totalPages }} 页
      </div>

      <div class="pagination-controls">
        <el-button-group>
          <el-button
            :disabled="currentPage <= 1"
            @click="goToPrevPage"
            icon="ArrowLeft"
            size="small"
          >
            上一页
          </el-button>

          <el-button
            :disabled="currentPage >= totalPages"
            @click="goToNextPage"
            icon="ArrowRight"
            size="small"
            icon-position="right"
          >
            下一页
          </el-button>
        </el-button-group>

        <!-- 页码选择器 -->
        <div class="page-selector">
          <span>跳转到</span>
          <el-select
            v-model="selectedPage"
            @change="goToSelectedPage"
            size="small"
            style="width: 80px; margin-left: 8px;"
          >
            <el-option
              v-for="page in totalPages"
              :key="page"
              :label="page"
              :value="page"
            />
          </el-select>
          <span style="margin-left: 8px;">页</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useCommentStore } from '../stores/commentStore'
import { ChatDotRound, ArrowLeft, ArrowRight } from '@element-plus/icons-vue'

// Props
const props = defineProps({
  siteIndex: {
    type: [String, Number],
    required: true
  }
})

// Store
const commentStore = useCommentStore()

// 响应式数据
const activeTab = ref('positive')
const selectedPage = ref(1)

// 计算属性
const positiveComments = computed(() => commentStore.getPositiveCommentsByPage)
const negativeComments = computed(() => commentStore.getNegativeCommentsByPage)

const currentPage = computed(() => {
  return activeTab.value === 'positive'
    ? commentStore.positiveComments.currentPage
    : commentStore.negativeComments.currentPage
})

const totalPages = computed(() => {
  return activeTab.value === 'positive'
    ? commentStore.positiveTotalPages
    : commentStore.negativeTotalPages
})

const showPagination = computed(() => {
  return totalPages.value > 1
})

// 方法
const handleTabChange = (tabName) => {
  activeTab.value = tabName
  selectedPage.value = 1

  // 切换标签时重新加载数据
  if (tabName === 'positive') {
    commentStore.fetchPositiveComments(props.siteIndex, 1)
  } else {
    commentStore.fetchNegativeComments(props.siteIndex, 1)
  }
}

const goToPrevPage = () => {
  if (currentPage.value > 1) {
    const newPage = currentPage.value - 1
    selectedPage.value = newPage

    if (activeTab.value === 'positive') {
      commentStore.changePositivePage(newPage)
    } else {
      commentStore.changeNegativePage(newPage)
    }
  }
}

const goToNextPage = () => {
  if (currentPage.value < totalPages.value) {
    const newPage = currentPage.value + 1
    selectedPage.value = newPage

    if (activeTab.value === 'positive') {
      commentStore.changePositivePage(newPage)
    } else {
      commentStore.changeNegativePage(newPage)
    }
  }
}

const goToSelectedPage = (page) => {
  if (activeTab.value === 'positive') {
    commentStore.changePositivePage(page)
  } else {
    commentStore.changeNegativePage(page)
  }
}

// 监听当前页面变化，更新选择器
watch(currentPage, (newPage) => {
  selectedPage.value = newPage
})

// 初始化数据
onMounted(() => {
  commentStore.fetchCountPositiveComments(props.siteIndex)
  commentStore.fetchCountNegativeComments(props.siteIndex)
  commentStore.fetchPositiveComments(props.siteIndex, 1)
})
</script>

<style scoped>
.comment-section {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.18);
}

/* 头部样式 */
.section-header {
  margin-bottom: 1.5rem;
}

.section-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #333;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0;
}

.section-title .el-icon {
  color: #667eea;
}

/* 标签页样式 */
.comment-tabs {
  margin-bottom: 1.5rem;
}

.comment-count {
  font-size: 0.9rem;
  color: #666;
  margin-top: 0.5rem;
  margin-bottom: 1rem;
}

/* 评论列表样式 */
.comment-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.comment-item {
  border-radius: 8px;
  padding: 1rem;
  transition: all 0.3s ease;
}

.positive-comment {
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  border-left: 4px solid #10b981;
}

.positive-comment:hover {
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.1);
}

.negative-comment {
  background: linear-gradient(135deg, #fef2f2 0%, #fee2e2 100%);
  border-left: 4px solid #ef4444;
}

.negative-comment:hover {
  box-shadow: 0 2px 8px rgba(239, 68, 68, 0.1);
}

/* 评论包装器样式 */
.comment-wrapper {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

/* 用户头像样式 */
.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
  border: 2px solid #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.comment-content {
  flex: 1;
  min-width: 0;
}

.comment-text {
  line-height: 1.6;
  color: #333;
  font-size: 1rem;
  word-wrap: break-word;
}

/* 状态样式 */
.loading-state,
.error-state,
.empty-state {
  margin: 2rem 0;
}

/* 分页样式 */
.pagination-container {
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
}

.pagination-info {
  color: #666;
  font-size: 0.9rem;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.page-selector {
  display: flex;
  align-items: center;
  color: #666;
  font-size: 0.9rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .comment-section {
    padding: 1.5rem;
  }

  .section-title {
    font-size: 1.2rem;
  }

  .pagination-container {
    flex-direction: column;
    align-items: stretch;
    gap: 1rem;
  }

  .pagination-controls {
    justify-content: center;
    flex-wrap: wrap;
  }

  .page-selector {
    order: -1;
    justify-content: center;
  }

  .comment-item {
    padding: 0.75rem;
  }

  .comment-wrapper {
    gap: 10px;
  }

  .user-avatar {
    width: 36px;
    height: 36px;
  }

  .comment-text {
    font-size: 0.9rem;
  }
}
</style>
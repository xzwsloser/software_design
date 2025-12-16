<template>
  <el-button
    :type="isLiked ? 'danger' : 'default'"
    :loading="loading"
    @click="handleLikeClick"
    class="like-button"
    :class="{ 'is-liked': isLiked }"
    size="large"
  >
    <el-icon class="like-icon">
      <StarFilled v-if="isLiked" />
      <Star v-else />
    </el-icon>
    <span class="like-text">{{ isLiked ? '已点赞' : '点赞' }}</span>
  </el-button>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useLikeStore } from '@/stores/likeStore'
import { ElMessage } from 'element-plus'
import { Star, StarFilled } from '@element-plus/icons-vue'

const props = defineProps({
  siteIndex: {
    type: [String, Number],
    required: true
  }
})

const likeStore = useLikeStore()

// 计算属性
const isLiked = computed(() => likeStore.isSiteLiked(Number(props.siteIndex)))
const loading = computed(() => likeStore.likeLoading)

// 处理点赞点击
const handleLikeClick = async () => {
  const success = await likeStore.toggleLike(Number(props.siteIndex))

  if (success) {
    if (isLiked.value) {
      ElMessage.success('点赞成功')
    } else {
      ElMessage.info('已取消点赞')
    }
  } else {
    ElMessage.error(likeStore.error || '操作失败，请稍后重试')
  }
}

onMounted(() => {
  likeStore.checkSiteLikeStatus(Number(props.siteIndex))
})
</script>

<style scoped>
.like-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border-radius: 50px;
  font-size: 1rem;
  font-weight: 500;
  transition: all 0.3s ease;
  border: 2px solid transparent;
  background: rgba(255, 255, 255, 0.9);
  color: #666;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.like-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.like-button.is-liked {
  background: linear-gradient(135deg, #ff6b6b 0%, #ff5252 100%);
  border-color: #ff5252;
  color: white;
  box-shadow: 0 4px 16px rgba(255, 82, 82, 0.3);
}

.like-button.is-liked:hover {
  background: linear-gradient(135deg, #ff5252 0%, #ff3838 100%);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(255, 82, 82, 0.4);
}

.like-icon {
  font-size: 1.2rem;
  transition: transform 0.3s ease;
}

.like-button:hover .like-icon {
  transform: scale(1.2);
}

.like-text {
  font-weight: 600;
}

/* 加载状态样式 */
.like-button.is-loading {
  pointer-events: none;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .like-button {
    padding: 0.6rem 1.2rem;
    font-size: 0.9rem;
  }

  .like-icon {
    font-size: 1rem;
  }
}
</style>
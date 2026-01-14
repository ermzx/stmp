<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h2>仪表板</h2>
      <p class="dashboard-subtitle">欢迎使用SMTP邮件管理系统</p>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :md="8" :lg="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon-wrapper stat-icon-total">
              <el-icon class="stat-icon"><Message /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-title">总发送数</div>
              <div class="stat-value">{{ statistics.total || 0 }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="8" :lg="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon-wrapper stat-icon-success">
              <el-icon class="stat-icon"><SuccessFilled /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-title">成功数</div>
              <div class="stat-value">{{ statistics.success || 0 }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="8" :lg="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon-wrapper stat-icon-failed">
              <el-icon class="stat-icon"><CircleCloseFilled /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-title">失败数</div>
              <div class="stat-value">{{ statistics.failed || 0 }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="8" :lg="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon-wrapper stat-icon-rate">
              <el-icon class="stat-icon"><TrendCharts /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-title">成功率</div>
              <div class="stat-value">{{ successRate }}%</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快速操作 -->
    <el-card class="quick-actions-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>快速操作</span>
        </div>
      </template>
      <div class="quick-actions">
        <el-button type="primary" size="large" @click="goToCompose">
          <el-icon><Edit /></el-icon>
          撰写邮件
        </el-button>
        <el-button type="success" size="large" @click="goToSmtp">
          <el-icon><Setting /></el-icon>
          SMTP配置
        </el-button>
        <el-button type="warning" size="large" @click="goToTemplates">
          <el-icon><Document /></el-icon>
          邮件模板
        </el-button>
        <el-button type="info" size="large" @click="goToHistory">
          <el-icon><Clock /></el-icon>
          发送历史
        </el-button>
      </div>
    </el-card>

    <!-- 最近发送记录 -->
    <el-card class="recent-history-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>最近发送记录</span>
          <el-button type="primary" link @click="goToHistory">查看全部</el-button>
        </div>
      </template>
      <div v-loading="loading" class="history-table">
        <el-table :data="recentHistory" stripe style="width: 100%">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="to_email" label="收件人" min-width="150" />
          <el-table-column prop="subject" label="主题" min-width="200" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="发送时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
        </el-table>
        <el-empty v-if="!loading && recentHistory.length === 0" description="暂无发送记录" />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Message,
  SuccessFilled,
  CircleCloseFilled,
  TrendCharts,
  Edit,
  Setting,
  Document,
  Clock
} from '@element-plus/icons-vue'
import { getHistoryStatistics, getHistory } from '../api'

const router = useRouter()
const loading = ref(false)
const statistics = ref({
  total: 0,
  success: 0,
  failed: 0
})
const recentHistory = ref([])

// 计算成功率
const successRate = computed(() => {
  if (statistics.value.total === 0) return 0
  return ((statistics.value.success / statistics.value.total) * 100).toFixed(1)
})

// 获取状态类型
const getStatusType = (status) => {
  const statusMap = {
    success: 'success',
    failed: 'danger',
    pending: 'warning'
  }
  return statusMap[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    success: '成功',
    failed: '失败',
    pending: '发送中'
  }
  return statusMap[status] || '未知'
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 加载统计数据
const loadStatistics = async () => {
  try {
    const response = await getHistoryStatistics()
    if (response && response.data) {
      statistics.value = response.data
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
    ElMessage.error('加载统计数据失败')
  }
}

// 加载最近发送记录
const loadRecentHistory = async () => {
  loading.value = true
  try {
    const response = await getHistory({ page: 1, pageSize: 10 })
    if (response && response.data) {
      recentHistory.value = response.data.list || []
    }
  } catch (error) {
    console.error('加载最近记录失败:', error)
    ElMessage.error('加载最近记录失败')
  } finally {
    loading.value = false
  }
}

// 快速操作导航
const goToCompose = () => {
  router.push('/compose')
}

const goToSmtp = () => {
  router.push('/smtp')
}

const goToTemplates = () => {
  router.push('/templates')
}

const goToHistory = () => {
  router.push('/history')
}

// 刷新数据
const refreshData = () => {
  loadStatistics()
  loadRecentHistory()
}

onMounted(() => {
  refreshData()
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.dashboard-header {
  margin-bottom: 30px;
}

.dashboard-header h2 {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 10px 0;
}

.dashboard-subtitle {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  height: 100%;
  transition: transform 0.3s, box-shadow 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.stat-icon-wrapper {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon-total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon-success {
  background: linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%);
}

.stat-icon-failed {
  background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%);
}

.stat-icon-rate {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon {
  font-size: 32px;
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
}

.quick-actions-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
}

.quick-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
}

.quick-actions .el-button {
  flex: 1;
  min-width: 150px;
}

.recent-history-card {
  margin-bottom: 20px;
}

.history-table {
  min-height: 200px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .dashboard {
    padding: 10px;
  }

  .dashboard-header h2 {
    font-size: 24px;
  }

  .stat-value {
    font-size: 24px;
  }

  .stat-icon-wrapper {
    width: 50px;
    height: 50px;
  }

  .stat-icon {
    font-size: 24px;
  }

  .quick-actions {
    flex-direction: column;
  }

  .quick-actions .el-button {
    width: 100%;
  }
}
</style>
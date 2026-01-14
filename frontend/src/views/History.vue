<template>
  <div class="history">
    <div class="history-header">
      <h2>发送历史</h2>
      <p class="history-subtitle">查看和管理邮件发送记录</p>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :md="6">
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
      <el-col :xs="24" :sm="12" :md="6">
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
      <el-col :xs="24" :sm="12" :md="6">
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
      <el-col :xs="24" :sm="12" :md="6">
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

    <!-- 历史记录列表 -->
    <el-card class="history-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>发送记录</span>
          <el-button type="primary" link @click="loadData">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-select v-model="filterStatus" placeholder="选择状态" @change="handleFilterChange">
          <el-option label="全部" value="all" />
          <el-option label="成功" value="success" />
          <el-option label="失败" value="failed" />
        </el-select>
      </div>

      <!-- 表格 -->
      <div v-loading="loading" class="history-table">
        <el-table :data="historyList" stripe style="width: 100%">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="to_email" label="收件人" min-width="150" show-overflow-tooltip />
          <el-table-column label="抄送" min-width="120" show-overflow-tooltip>
            <template #default="{ row }">
              {{ row.cc_email && row.cc_email.length > 0 ? row.cc_email.join(', ') : '-' }}
            </template>
          </el-table-column>
          <el-table-column label="密送" min-width="120" show-overflow-tooltip>
            <template #default="{ row }">
              {{ row.bcc_email && row.bcc_email.length > 0 ? row.bcc_email.join(', ') : '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="subject" label="邮件主题" min-width="200" show-overflow-tooltip />
          <el-table-column prop="status" label="发送状态" width="100">
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
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="viewDetail(row)">
                查看详情
              </el-button>
              <el-button type="danger" link size="small" @click="handleDelete(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-empty v-if="!loading && historyList.length === 0" description="暂无发送记录" />
      </div>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="邮件详情"
      width="800px"
      :close-on-click-modal="false"
    >
      <div v-if="currentDetail" class="detail-content">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID">{{ currentDetail.id }}</el-descriptions-item>
          <el-descriptions-item label="发送状态">
            <el-tag :type="getStatusType(currentDetail.status)">
              {{ getStatusText(currentDetail.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="收件人" :span="2">{{ currentDetail.to_email }}</el-descriptions-item>
          <el-descriptions-item label="抄送" :span="2">
            {{ currentDetail.cc_email && currentDetail.cc_email.length > 0 ? currentDetail.cc_email.join(', ') : '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="密送" :span="2">
            {{ currentDetail.bcc_email && currentDetail.bcc_email.length > 0 ? currentDetail.bcc_email.join(', ') : '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="邮件主题" :span="2">{{ currentDetail.subject }}</el-descriptions-item>
          <el-descriptions-item label="发送时间" :span="2">{{ formatDate(currentDetail.created_at) }}</el-descriptions-item>
        </el-descriptions>

        <!-- 邮件正文 -->
        <div class="detail-section">
          <div class="detail-section-title">邮件正文</div>
          <div class="detail-body" v-html="currentDetail.body"></div>
        </div>

        <!-- 附件列表 -->
        <div v-if="currentDetail.attachments && currentDetail.attachments.length > 0" class="detail-section">
          <div class="detail-section-title">附件列表</div>
          <el-table :data="currentDetail.attachments" size="small">
            <el-table-column prop="filename" label="文件名" />
            <el-table-column prop="size" label="大小">
              <template #default="{ row }">
                {{ formatFileSize(row.size) }}
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 错误信息 -->
        <div v-if="currentDetail.status === 'failed' && currentDetail.error_message" class="detail-section">
          <div class="detail-section-title">错误信息</div>
          <el-alert type="error" :closable="false">
            {{ currentDetail.error_message }}
          </el-alert>
        </div>
      </div>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Message,
  SuccessFilled,
  CircleCloseFilled,
  TrendCharts,
  Refresh
} from '@element-plus/icons-vue'
import { getHistory, getHistoryStatistics, deleteHistory } from '../api'

const loading = ref(false)
const historyList = ref([])
const statistics = ref({
  total: 0,
  success: 0,
  failed: 0
})
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const filterStatus = ref('all')
const detailDialogVisible = ref(false)
const currentDetail = ref(null)

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

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
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

// 加载历史记录列表
const loadHistoryList = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      status: filterStatus.value
    }
    const response = await getHistory(params)
    if (response && response.data) {
      historyList.value = response.data.list || []
      total.value = response.data.total || 0
    }
  } catch (error) {
    console.error('加载历史记录失败:', error)
    ElMessage.error('加载历史记录失败')
  } finally {
    loading.value = false
  }
}

// 加载所有数据
const loadData = () => {
  loadStatistics()
  loadHistoryList()
}

// 筛选状态改变
const handleFilterChange = () => {
  currentPage.value = 1
  loadHistoryList()
}

// 每页数量改变
const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadHistoryList()
}

// 页码改变
const handlePageChange = (page) => {
  currentPage.value = page
  loadHistoryList()
}

// 查看详情
const viewDetail = async (row) => {
  try {
    // 这里可以直接使用row中的数据，因为列表已经包含了大部分信息
    // 如果需要获取更完整的信息，可以调用API
    currentDetail.value = row
    detailDialogVisible.value = true
  } catch (error) {
    console.error('获取详情失败:', error)
    ElMessage.error('获取详情失败')
  }
}

// 删除历史记录
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除这条发送记录吗？`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await deleteHistory(row.id)
      ElMessage.success('删除成功')
      loadData()
    } catch (error) {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {
    // 用户取消删除
  })
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.history {
  padding: 20px;
}

.history-header {
  margin-bottom: 30px;
}

.history-header h2 {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 10px 0;
}

.history-subtitle {
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

.history-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
}

.filter-bar {
  margin-bottom: 20px;
}

.history-table {
  min-height: 400px;
  margin-bottom: 20px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #ebeef5;
}

.detail-content {
  max-height: 600px;
  overflow-y: auto;
}

.detail-section {
  margin-top: 20px;
}

.detail-section-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 10px;
}

.detail-body {
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
  min-height: 100px;
  max-height: 300px;
  overflow-y: auto;
  word-wrap: break-word;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .history {
    padding: 10px;
  }

  .history-header h2 {
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

  .pagination-wrapper {
    justify-content: center;
  }
}
</style>
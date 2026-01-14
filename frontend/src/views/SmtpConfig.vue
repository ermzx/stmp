<template>
  <div class="smtp-config-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">SMTP配置管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加配置
          </el-button>
        </div>
      </template>

      <!-- 配置列表表格 -->
      <el-table
        v-loading="loading"
        :data="configList"
        stripe
        border
        style="width: 100%"
        empty-text="暂无SMTP配置"
      >
        <el-table-column prop="name" label="配置名称" min-width="120" />
        <el-table-column prop="host" label="服务器地址" min-width="150" />
        <el-table-column prop="port" label="端口" width="80" align="center" />
        <el-table-column prop="from_email" label="发件人邮箱" min-width="180" />
        <el-table-column prop="encryption" label="加密方式" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getEncryptionType(row.encryption)">
              {{ getEncryptionLabel(row.encryption) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="默认状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.is_default" type="success">默认</el-tag>
            <el-tag v-else type="info">普通</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" align="center" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              link
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              type="warning"
              size="small"
              link
              @click="handleTest(row)"
              :loading="testingId === row.id"
            >
              测试
            </el-button>
            <el-button
              v-if="!row.is_default"
              type="success"
              size="small"
              link
              @click="handleSetDefault(row)"
            >
              设为默认
            </el-button>
            <el-button
              type="danger"
              size="small"
              link
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      :close-on-click-modal="false"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="120px"
      >
        <el-form-item label="配置名称" prop="name">
          <el-input
            v-model="formData.name"
            placeholder="请输入配置名称"
            clearable
          />
        </el-form-item>
        <el-form-item label="SMTP服务器" prop="host">
          <el-input
            v-model="formData.host"
            placeholder="请输入SMTP服务器地址"
            clearable
          />
        </el-form-item>
        <el-form-item label="端口号" prop="port">
          <el-input-number
            v-model="formData.port"
            :min="1"
            :max="65535"
            controls-position="right"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="formData.username"
            placeholder="请输入用户名（可选）"
            clearable
          />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="formData.password"
            type="password"
            placeholder="请输入密码（可选）"
            show-password
            clearable
          />
        </el-form-item>
        <el-form-item label="发件人邮箱" prop="from_email">
          <el-input
            v-model="formData.from_email"
            placeholder="请输入发件人邮箱"
            clearable
          />
        </el-form-item>
        <el-form-item label="发件人名称" prop="from_name">
          <el-input
            v-model="formData.from_name"
            placeholder="请输入发件人名称（可选）"
            clearable
          />
        </el-form-item>
        <el-form-item label="加密方式" prop="encryption">
          <el-select
            v-model="formData.encryption"
            placeholder="请选择加密方式"
            style="width: 100%"
          >
            <el-option label="无" value="none" />
            <el-option label="TLS" value="tls" />
            <el-option label="StartTLS" value="starttls" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getSmtpConfigs,
  createSmtpConfig,
  updateSmtpConfig,
  deleteSmtpConfig,
  testSmtpConnection,
  setDefaultSmtpConfig
} from '@/api'

// 数据状态
const loading = ref(false)
const configList = ref([])
const testingId = ref(null)
const submitting = ref(false)

// 对话框状态
const dialogVisible = ref(false)
const dialogTitle = ref('添加SMTP配置')
const isEdit = ref(false)
const currentId = ref(null)

// 表单引用
const formRef = ref(null)

// 表单数据
const formData = reactive({
  name: '',
  host: '',
  port: 25,
  username: '',
  password: '',
  from_email: '',
  from_name: '',
  encryption: 'none'
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入配置名称', trigger: 'blur' }
  ],
  host: [
    { required: true, message: '请输入SMTP服务器地址', trigger: 'blur' }
  ],
  port: [
    { required: true, message: '请输入端口号', trigger: 'blur' }
  ],
  from_email: [
    { required: true, message: '请输入发件人邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  encryption: [
    { required: true, message: '请选择加密方式', trigger: 'change' }
  ]
}

// 获取加密方式标签
const getEncryptionLabel = (encryption) => {
  const map = {
    none: '无',
    tls: 'TLS',
    starttls: 'StartTLS'
  }
  return map[encryption] || encryption
}

// 获取加密方式标签类型
const getEncryptionType = (encryption) => {
  const map = {
    none: 'info',
    tls: 'success',
    starttls: 'warning'
  }
  return map[encryption] || 'info'
}

// 获取配置列表
const fetchConfigList = async () => {
  loading.value = true
  try {
    const res = await getSmtpConfigs()
    if (res.code === 200) {
      configList.value = res.data || []
    } else {
      ElMessage.error(res.message || '获取配置列表失败')
    }
  } catch (error) {
    console.error('获取配置列表失败:', error)
    ElMessage.error('获取配置列表失败')
  } finally {
    loading.value = false
  }
}

// 添加配置
const handleAdd = () => {
  isEdit.value = false
  dialogTitle.value = '添加SMTP配置'
  resetForm()
  dialogVisible.value = true
}

// 编辑配置
const handleEdit = (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑SMTP配置'
  currentId.value = row.id
  
  // 填充表单数据
  Object.assign(formData, {
    name: row.name,
    host: row.host,
    port: row.port,
    username: row.username || '',
    password: '', // 密码不回显
    from_email: row.from_email,
    from_name: row.from_name || '',
    encryption: row.encryption || 'none'
  })
  
  dialogVisible.value = true
}

// 删除配置
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除SMTP配置"${row.name}"吗？此操作不可恢复。`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const res = await deleteSmtpConfig(row.id)
      if (res.code === 200) {
        ElMessage.success('删除成功')
        await fetchConfigList()
      } else {
        ElMessage.error(res.message || '删除失败')
      }
    } catch (error) {
      console.error('删除配置失败:', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {
    // 用户取消删除
  })
}

// 测试连接
const handleTest = async (row) => {
  testingId.value = row.id
  try {
    const res = await testSmtpConnection(row.id)
    if (res.code === 200) {
      ElMessage.success('SMTP连接测试成功')
    } else {
      ElMessage.error(res.message || 'SMTP连接测试失败')
    }
  } catch (error) {
    console.error('测试连接失败:', error)
    ElMessage.error('SMTP连接测试失败')
  } finally {
    testingId.value = null
  }
}

// 设为默认
const handleSetDefault = async (row) => {
  try {
    const res = await setDefaultSmtpConfig(row.id)
    if (res.code === 200) {
      ElMessage.success('已设置为默认配置')
      await fetchConfigList()
    } else {
      ElMessage.error(res.message || '设置默认配置失败')
    }
  } catch (error) {
    console.error('设置默认配置失败:', error)
    ElMessage.error('设置默认配置失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    submitting.value = true
    try {
      let res
      if (isEdit.value) {
        // 编辑模式
        res = await updateSmtpConfig(currentId.value, formData)
      } else {
        // 添加模式
        res = await createSmtpConfig(formData)
      }
      
      if (res.code === 200) {
        ElMessage.success(isEdit.value ? '更新成功' : '添加成功')
        dialogVisible.value = false
        await fetchConfigList()
      } else {
        ElMessage.error(res.message || (isEdit.value ? '更新失败' : '添加失败'))
      }
    } catch (error) {
      console.error('提交失败:', error)
      ElMessage.error(isEdit.value ? '更新失败' : '添加失败')
    } finally {
      submitting.value = false
    }
  })
}

// 对话框关闭
const handleDialogClose = () => {
  resetForm()
}

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  Object.assign(formData, {
    name: '',
    host: '',
    port: 25,
    username: '',
    password: '',
    from_email: '',
    from_name: '',
    encryption: 'none'
  })
  currentId.value = null
}

// 页面加载时获取配置列表
onMounted(() => {
  fetchConfigList()
})
</script>

<style scoped>
.smtp-config-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.box-card {
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

:deep(.el-table) {
  border-radius: 4px;
}

:deep(.el-dialog__body) {
  padding: 20px 30px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}
</style>
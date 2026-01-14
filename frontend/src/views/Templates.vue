<template>
  <div class="templates-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">邮件模板管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加模板
          </el-button>
        </div>
      </template>

      <!-- 加载状态 -->
      <div v-loading="loading" class="templates-content">
        <!-- 空状态 -->
        <el-empty
          v-if="!loading && templates.length === 0"
          description="暂无邮件模板"
        >
          <el-button type="primary" @click="handleAdd">创建第一个模板</el-button>
        </el-empty>

        <!-- 模板卡片列表 -->
        <div v-else class="templates-grid">
          <el-card
            v-for="template in templates"
            :key="template.id"
            class="template-card"
            shadow="hover"
          >
            <template #header>
              <div class="template-card-header">
                <span class="template-name">{{ template.name }}</span>
              </div>
            </template>

            <div class="template-info">
              <div class="info-item">
                <el-icon><Message /></el-icon>
                <span class="info-label">主题：</span>
                <span class="info-value">{{ template.subject }}</span>
              </div>
              <div class="info-item">
                <el-icon><Clock /></el-icon>
                <span class="info-label">创建时间：</span>
                <span class="info-value">{{ formatDate(template.created_at) }}</span>
              </div>
              <div class="info-item">
                <el-icon><Refresh /></el-icon>
                <span class="info-label">更新时间：</span>
                <span class="info-value">{{ formatDate(template.updated_at) }}</span>
              </div>
            </div>

            <div class="template-actions">
              <el-button-group>
                <el-button size="small" @click="handlePreview(template)">
                  <el-icon><View /></el-icon>
                  预览
                </el-button>
                <el-button size="small" @click="handleEdit(template)">
                  <el-icon><Edit /></el-icon>
                  编辑
                </el-button>
                <el-button size="small" @click="handleUse(template)">
                  <el-icon><Promotion /></el-icon>
                  使用
                </el-button>
                <el-button size="small" type="danger" @click="handleDelete(template)">
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </el-button-group>
            </div>
          </el-card>
        </div>
      </div>
    </el-card>

    <!-- 添加/编辑模板对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑模板' : '添加模板'"
      width="800px"
      :close-on-click-modal="false"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="模板名称" prop="name">
          <el-input
            v-model="formData.name"
            placeholder="请输入模板名称"
            clearable
          />
        </el-form-item>

        <el-form-item label="邮件主题" prop="subject">
          <el-input
            v-model="formData.subject"
            placeholder="请输入邮件主题"
            clearable
          />
        </el-form-item>

        <el-form-item label="邮件正文" prop="body">
          <div class="editor-container">
            <div class="editor-toolbar">
              <el-button-group>
                <el-button size="small" @click="execCommand('bold')" title="粗体">
                  <span style="font-weight: bold;">B</span>
                </el-button>
                <el-button size="small" @click="execCommand('italic')" title="斜体">
                  <span style="font-style: italic;">I</span>
                </el-button>
                <el-button size="small" @click="execCommand('underline')" title="下划线">
                  <span style="text-decoration: underline;">U</span>
                </el-button>
              </el-button-group>
              <el-divider direction="vertical" />
              <el-button-group>
                <el-button size="small" @click="execCommand('justifyLeft')" title="左对齐">
                  <el-icon><DArrowLeft /></el-icon>
                </el-button>
                <el-button size="small" @click="execCommand('justifyCenter')" title="居中">
                  <el-icon><MoreFilled /></el-icon>
                </el-button>
                <el-button size="small" @click="execCommand('justifyRight')" title="右对齐">
                  <el-icon><DArrowRight /></el-icon>
                </el-button>
              </el-button-group>
              <el-divider direction="vertical" />
              <el-button-group>
                <el-button size="small" @click="insertLink" title="插入链接">
                  <el-icon><Link /></el-icon>
                </el-button>
                <el-button size="small" @click="insertImage" title="插入图片">
                  <el-icon><Picture /></el-icon>
                </el-button>
              </el-button-group>
              <el-divider direction="vertical" />
              <el-button size="small" @click="execCommand('removeFormat')" title="清除格式">
                <el-icon><Remove /></el-icon>
              </el-button>
            </div>
            <div
              ref="editorRef"
              class="editor-content"
              contenteditable="true"
              @input="handleEditorInput"
            ></div>
          </div>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEdit ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 预览对话框 -->
    <el-dialog
      v-model="previewVisible"
      title="模板预览"
      width="800px"
    >
      <div class="preview-content">
        <div class="preview-item">
          <span class="preview-label">主题：</span>
          <span class="preview-value">{{ previewData.subject }}</span>
        </div>
        <el-divider />
        <div class="preview-body">
          <div class="preview-label">正文：</div>
          <div class="preview-html" v-html="previewData.body"></div>
        </div>
      </div>
      <template #footer>
        <el-button @click="previewVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 插入链接对话框 -->
    <el-dialog
      v-model="linkDialogVisible"
      title="插入链接"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form label-width="80px">
        <el-form-item label="链接地址">
          <el-input v-model="linkUrl" placeholder="请输入链接地址" />
        </el-form-item>
        <el-form-item label="链接文本">
          <el-input v-model="linkText" placeholder="请输入链接文本" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="linkDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmInsertLink">确定</el-button>
      </template>
    </el-dialog>

    <!-- 插入图片对话框 -->
    <el-dialog
      v-model="imageDialogVisible"
      title="插入图片"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form label-width="80px">
        <el-form-item label="图片地址">
          <el-input v-model="imageUrl" placeholder="请输入图片地址" />
        </el-form-item>
        <el-form-item label="或上传图片">
          <el-upload
            :auto-upload="false"
            :on-change="handleImageChange"
            :show-file-list="false"
            accept="image/*"
          >
            <el-button type="primary">选择图片</el-button>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="imageDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmInsertImage">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Edit,
  Delete,
  View,
  Promotion,
  Message,
  Clock,
  Refresh,
  DArrowLeft,
  DArrowRight,
  MoreFilled,
  Link,
  Picture,
  Remove
} from '@element-plus/icons-vue'
import { getTemplates, createTemplate, updateTemplate, deleteTemplate } from '@/api'

const router = useRouter()

// 数据状态
const loading = ref(false)
const submitting = ref(false)
const templates = ref([])

// 对话框状态
const dialogVisible = ref(false)
const previewVisible = ref(false)
const linkDialogVisible = ref(false)
const imageDialogVisible = ref(false)

// 表单引用
const formRef = ref(null)
const editorRef = ref(null)

// 编辑状态
const isEdit = ref(false)
const currentTemplateId = ref(null)

// 表单数据
const formData = reactive({
  name: '',
  subject: '',
  body: ''
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入模板名称', trigger: 'blur' }
  ],
  subject: [
    { required: true, message: '请输入邮件主题', trigger: 'blur' }
  ],
  body: [
    { required: true, message: '请输入邮件正文', trigger: 'blur' }
  ]
}

// 预览数据
const previewData = reactive({
  subject: '',
  body: ''
})

// 链接对话框数据
const linkUrl = ref('')
const linkText = ref('')

// 图片对话框数据
const imageUrl = ref('')
const imageBase64 = ref('')

// 获取模板列表
const fetchTemplates = async () => {
  loading.value = true
  try {
    const res = await getTemplates()
    if (res.code === 200) {
      templates.value = res.data || []
    } else {
      ElMessage.error(res.message || '获取模板列表失败')
    }
  } catch (error) {
    console.error('获取模板列表失败:', error)
    ElMessage.error('获取模板列表失败')
  } finally {
    loading.value = false
  }
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
    minute: '2-digit'
  })
}

// 添加模板
const handleAdd = () => {
  isEdit.value = false
  currentTemplateId.value = null
  Object.assign(formData, {
    name: '',
    subject: '',
    body: ''
  })
  if (editorRef.value) {
    editorRef.value.innerHTML = ''
  }
  dialogVisible.value = true
}

// 编辑模板
const handleEdit = (template) => {
  isEdit.value = true
  currentTemplateId.value = template.id
  Object.assign(formData, {
    name: template.name,
    subject: template.subject,
    body: template.body
  })
  if (editorRef.value) {
    editorRef.value.innerHTML = template.body || ''
  }
  dialogVisible.value = true
}

// 删除模板
const handleDelete = (template) => {
  ElMessageBox.confirm(
    `确定要删除模板"${template.name}"吗？此操作不可恢复。`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const res = await deleteTemplate(template.id)
      if (res.code === 200) {
        ElMessage.success('删除成功')
        await fetchTemplates()
      } else {
        ElMessage.error(res.message || '删除失败')
      }
    } catch (error) {
      console.error('删除模板失败:', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {
    // 用户取消删除
  })
}

// 预览模板
const handlePreview = (template) => {
  previewData.subject = template.subject
  previewData.body = template.body
  previewVisible.value = true
}

// 使用模板
const handleUse = (template) => {
  // 将模板信息存储到sessionStorage，以便在撰写页面使用
  sessionStorage.setItem('templateData', JSON.stringify({
    subject: template.subject,
    body: template.body
  }))
  // 跳转到撰写页面
  router.push('/compose')
}

// 编辑器输入事件
const handleEditorInput = () => {
  if (editorRef.value) {
    formData.body = editorRef.value.innerHTML
  }
}

// 执行编辑器命令
const execCommand = (command, value = null) => {
  document.execCommand(command, false, value)
  handleEditorInput()
}

// 插入链接
const insertLink = () => {
  linkUrl.value = ''
  linkText.value = ''
  linkDialogVisible.value = true
}

// 确认插入链接
const confirmInsertLink = () => {
  if (!linkUrl.value) {
    ElMessage.warning('请输入链接地址')
    return
  }
  const text = linkText.value || linkUrl.value
  const html = `<a href="${linkUrl.value}" target="_blank">${text}</a>`
  execCommand('insertHTML', html)
  linkDialogVisible.value = false
}

// 插入图片
const insertImage = () => {
  imageUrl.value = ''
  imageBase64.value = ''
  imageDialogVisible.value = true
}

// 图片文件变化
const handleImageChange = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    imageBase64.value = e.target.result
  }
  reader.readAsDataURL(file.raw)
}

// 确认插入图片
const confirmInsertImage = () => {
  if (!imageUrl.value && !imageBase64.value) {
    ElMessage.warning('请输入图片地址或上传图片')
    return
  }
  const src = imageBase64.value || imageUrl.value
  const html = `<img src="${src}" style="max-width: 100%;" />`
  execCommand('insertHTML', html)
  imageDialogVisible.value = false
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    // 验证邮件正文
    if (!formData.body || formData.body.trim() === '') {
      ElMessage.warning('请输入邮件正文')
      return
    }

    submitting.value = true
    try {
      const submitData = {
        name: formData.name,
        subject: formData.subject,
        body: formData.body
      }

      let res
      if (isEdit.value) {
        res = await updateTemplate(currentTemplateId.value, submitData)
      } else {
        res = await createTemplate(submitData)
      }

      if (res.code === 200) {
        ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
        dialogVisible.value = false
        await fetchTemplates()
      } else {
        ElMessage.error(res.message || (isEdit.value ? '更新失败' : '创建失败'))
      }
    } catch (error) {
      console.error('提交失败:', error)
      ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
    } finally {
      submitting.value = false
    }
  })
}

// 对话框关闭事件
const handleDialogClose = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  Object.assign(formData, {
    name: '',
    subject: '',
    body: ''
  })
  if (editorRef.value) {
    editorRef.value.innerHTML = ''
  }
}

// 页面加载时获取数据
onMounted(() => {
  fetchTemplates()
})
</script>

<style scoped>
.templates-container {
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

.templates-content {
  min-height: 400px;
}

.templates-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.template-card {
  border-radius: 8px;
  transition: all 0.3s;
}

.template-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.15);
}

.template-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.template-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.template-info {
  margin-bottom: 16px;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
}

.info-item .el-icon {
  margin-right: 6px;
  color: #909399;
}

.info-label {
  color: #909399;
  margin-right: 4px;
}

.info-value {
  color: #606266;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.template-actions {
  display: flex;
  justify-content: flex-end;
}

.editor-container {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
}

.editor-toolbar {
  background-color: #f5f7fa;
  padding: 8px;
  border-bottom: 1px solid #dcdfe6;
  display: flex;
  align-items: center;
  gap: 8px;
}

.editor-content {
  min-height: 300px;
  padding: 12px;
  outline: none;
  overflow-y: auto;
}

.editor-content:focus {
  background-color: #fafafa;
}

.editor-content:empty:before {
  content: attr(placeholder);
  color: #909399;
  pointer-events: none;
}

.preview-content {
  padding: 10px;
}

.preview-item {
  display: flex;
  align-items: center;
  font-size: 14px;
}

.preview-label {
  font-weight: 600;
  color: #303133;
  margin-right: 8px;
  white-space: nowrap;
}

.preview-value {
  color: #606266;
  flex: 1;
}

.preview-body {
  margin-top: 16px;
}

.preview-html {
  padding: 16px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background-color: #fafafa;
  min-height: 200px;
  max-height: 400px;
  overflow-y: auto;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .templates-container {
    padding: 10px;
  }

  .templates-grid {
    grid-template-columns: 1fr;
  }

  .template-actions {
    flex-wrap: wrap;
  }

  .template-actions .el-button-group {
    width: 100%;
  }

  .template-actions .el-button {
    flex: 1;
  }

  .editor-toolbar {
    flex-wrap: wrap;
  }

  .editor-content {
    min-height: 200px;
  }
}
</style>
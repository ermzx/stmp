<template>
  <div class="compose-email-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">撰写邮件</span>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        v-loading="loading"
      >
        <!-- SMTP配置选择 -->
        <el-form-item label="SMTP配置" prop="smtp_config_id">
          <el-select
            v-model="formData.smtp_config_id"
            placeholder="请选择SMTP配置"
            style="width: 100%"
            filterable
          >
            <el-option
              v-for="config in smtpConfigs"
              :key="config.id"
              :label="`${config.name} ${config.is_default ? '(默认)' : ''}`"
              :value="config.id"
            />
          </el-select>
        </el-form-item>

        <!-- 邮件模板选择 -->
        <el-form-item label="邮件模板">
          <el-select
            v-model="selectedTemplateId"
            placeholder="请选择邮件模板（可选）"
            style="width: 100%"
            clearable
            @change="handleTemplateChange"
          >
            <el-option
              v-for="template in templates"
              :key="template.id"
              :label="template.name"
              :value="template.id"
            />
          </el-select>
        </el-form-item>

        <!-- 收件人 -->
        <el-form-item label="收件人" prop="to">
          <el-input
            v-model="formData.to"
            type="textarea"
            :rows="2"
            placeholder="请输入收件人邮箱，多个邮箱用逗号分隔"
            clearable
          />
        </el-form-item>

        <!-- 抄送 -->
        <el-form-item label="抄送">
          <el-input
            v-model="formData.cc"
            type="textarea"
            :rows="2"
            placeholder="请输入抄送邮箱，多个邮箱用逗号分隔（可选）"
            clearable
          />
        </el-form-item>

        <!-- 密送 -->
        <el-form-item label="密送">
          <el-input
            v-model="formData.bcc"
            type="textarea"
            :rows="2"
            placeholder="请输入密送邮箱，多个邮箱用逗号分隔（可选）"
            clearable
          />
        </el-form-item>

        <!-- 主题 -->
        <el-form-item label="主题" prop="subject">
          <el-input
            v-model="formData.subject"
            placeholder="请输入邮件主题"
            clearable
          />
        </el-form-item>

        <!-- HTML编辑器工具栏 -->
        <el-form-item label="邮件正文" prop="html_content">
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

        <!-- 附件上传 -->
        <el-form-item label="附件">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :on-change="handleFileChange"
            :on-remove="handleFileRemove"
            :file-list="fileList"
            :limit="10"
            multiple
            accept=".pdf,.doc,.docx,.jpg,.jpeg,.png,.txt"
          >
            <el-button type="primary">
              <el-icon><Upload /></el-icon>
              选择文件
            </el-button>
            <template #tip>
              <div class="el-upload__tip">
                支持上传 .pdf, .doc, .docx, .jpg, .png, .txt 格式文件，单个文件不超过10MB
              </div>
            </template>
          </el-upload>
        </el-form-item>

        <!-- 操作按钮 -->
        <el-form-item>
          <el-button type="primary" @click="handleSend" :loading="sending">
            <el-icon><Promotion /></el-icon>
            发送邮件
          </el-button>
          <el-button @click="handleClear">
            <el-icon><Delete /></el-icon>
            清空表单
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

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
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  DArrowLeft,
  DArrowRight,
  MoreFilled,
  Link,
  Picture,
  Remove,
  Upload,
  Promotion,
  Delete
} from '@element-plus/icons-vue'
import { getSmtpConfigs, getTemplates, sendEmail } from '@/api'

// 数据状态
const loading = ref(false)
const sending = ref(false)
const smtpConfigs = ref([])
const templates = ref([])
const selectedTemplateId = ref(null)
const fileList = ref([])

// 表单引用
const formRef = ref(null)
const editorRef = ref(null)
const uploadRef = ref(null)

// 表单数据
const formData = reactive({
  smtp_config_id: null,
  to: '',
  cc: '',
  bcc: '',
  subject: '',
  html_content: '',
  attachments: []
})

// 表单验证规则
const formRules = {
  smtp_config_id: [
    { required: true, message: '请选择SMTP配置', trigger: 'change' }
  ],
  to: [
    { required: true, message: '请输入收件人邮箱', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (!value) {
          callback()
          return
        }
        const emails = value.split(',').map(e => e.trim()).filter(e => e)
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
        for (const email of emails) {
          if (!emailRegex.test(email)) {
            callback(new Error('请输入正确的邮箱格式'))
            return
          }
        }
        callback()
      },
      trigger: 'blur'
    }
  ],
  cc: [
    {
      validator: (rule, value, callback) => {
        if (!value) {
          callback()
          return
        }
        const emails = value.split(',').map(e => e.trim()).filter(e => e)
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
        for (const email of emails) {
          if (!emailRegex.test(email)) {
            callback(new Error('请输入正确的邮箱格式'))
            return
          }
        }
        callback()
      },
      trigger: 'blur'
    }
  ],
  bcc: [
    {
      validator: (rule, value, callback) => {
        if (!value) {
          callback()
          return
        }
        const emails = value.split(',').map(e => e.trim()).filter(e => e)
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
        for (const email of emails) {
          if (!emailRegex.test(email)) {
            callback(new Error('请输入正确的邮箱格式'))
            return
          }
        }
        callback()
      },
      trigger: 'blur'
    }
  ],
  subject: [
    { required: true, message: '请输入邮件主题', trigger: 'blur' }
  ],
  html_content: [
    { required: true, message: '请输入邮件正文', trigger: 'blur' }
  ]
}

// 链接对话框
const linkDialogVisible = ref(false)
const linkUrl = ref('')
const linkText = ref('')

// 图片对话框
const imageDialogVisible = ref(false)
const imageUrl = ref('')
const imageBase64 = ref('')

// 获取SMTP配置列表
const fetchSmtpConfigs = async () => {
  try {
    const res = await getSmtpConfigs()
    if (res.code === 200) {
      smtpConfigs.value = res.data || []
      // 自动选择默认配置
      const defaultConfig = smtpConfigs.value.find(c => c.is_default)
      if (defaultConfig) {
        formData.smtp_config_id = defaultConfig.id
      }
    } else {
      ElMessage.error(res.message || '获取SMTP配置失败')
    }
  } catch (error) {
    console.error('获取SMTP配置失败:', error)
    ElMessage.error('获取SMTP配置失败')
  }
}

// 获取邮件模板列表
const fetchTemplates = async () => {
  try {
    const res = await getTemplates()
    if (res.code === 200) {
      templates.value = res.data || []
    } else {
      ElMessage.error(res.message || '获取邮件模板失败')
    }
  } catch (error) {
    console.error('获取邮件模板失败:', error)
    ElMessage.error('获取邮件模板失败')
  }
}

// 模板选择变化
const handleTemplateChange = (templateId) => {
  if (!templateId) {
    return
  }
  const template = templates.value.find(t => t.id === templateId)
  if (template) {
    formData.subject = template.subject || ''
    // 使用nextTick确保DOM更新后再设置内容
    setTimeout(() => {
      if (editorRef.value) {
        editorRef.value.innerHTML = template.html_content || ''
        formData.html_content = template.html_content || ''
      }
    }, 100)
  }
}

// 编辑器输入事件
const handleEditorInput = () => {
  if (editorRef.value) {
    formData.html_content = editorRef.value.innerHTML
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

// 文件变化
const handleFileChange = (file) => {
  // 验证文件大小（10MB）
  const maxSize = 10 * 1024 * 1024
  if (file.size > maxSize) {
    ElMessage.error('文件大小不能超过10MB')
    return false
  }

  // 验证文件类型
  const allowedTypes = ['.pdf', '.doc', '.docx', '.jpg', '.jpeg', '.png', '.txt']
  const fileName = file.name.toLowerCase()
  const isValidType = allowedTypes.some(type => fileName.endsWith(type))
  if (!isValidType) {
    ElMessage.error('不支持的文件类型')
    return false
  }

  // 读取文件并转换为base64
  const reader = new FileReader()
  reader.onload = (e) => {
    const base64 = e.target.result
    formData.attachments.push({
      filename: file.name,
      content: base64
    })
  }
  reader.readAsDataURL(file.raw)
}

// 文件移除
const handleFileRemove = (file) => {
  const index = formData.attachments.findIndex(a => a.filename === file.name)
  if (index > -1) {
    formData.attachments.splice(index, 1)
  }
}

// 发送邮件
const handleSend = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    // 验证邮件正文
    if (!formData.html_content || formData.html_content.trim() === '') {
      ElMessage.warning('请输入邮件正文')
      return
    }

    sending.value = true
    try {
      // 准备发送数据
      const sendData = {
        smtp_config_id: formData.smtp_config_id,
        to: formData.to.split(',').map(e => e.trim()).filter(e => e),
        cc: formData.cc ? formData.cc.split(',').map(e => e.trim()).filter(e => e) : [],
        bcc: formData.bcc ? formData.bcc.split(',').map(e => e.trim()).filter(e => e) : [],
        subject: formData.subject,
        body: formData.html_content,
        attachments: formData.attachments
      }

      const res = await sendEmail(sendData)
      if (res.code === 200) {
        ElMessage.success('邮件发送成功')
        handleClear()
      } else {
        ElMessage.error(res.message || '邮件发送失败')
      }
    } catch (error) {
      console.error('发送邮件失败:', error)
      ElMessage.error('邮件发送失败')
    } finally {
      sending.value = false
    }
  })
}

// 清空表单
const handleClear = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  
  // 重置表单数据
  Object.assign(formData, {
    smtp_config_id: null,
    to: '',
    cc: '',
    bcc: '',
    subject: '',
    html_content: '',
    attachments: []
  })

  // 清空编辑器
  if (editorRef.value) {
    editorRef.value.innerHTML = ''
  }

  // 清空文件列表
  fileList.value = []
  if (uploadRef.value) {
    uploadRef.value.clearFiles()
  }

  // 清空模板选择
  selectedTemplateId.value = null

  // 重新选择默认配置
  const defaultConfig = smtpConfigs.value.find(c => c.is_default)
  if (defaultConfig) {
    formData.smtp_config_id = defaultConfig.id
  }
}

// 页面加载时获取数据
onMounted(async () => {
  loading.value = true
  try {
    await Promise.all([
      fetchSmtpConfigs(),
      fetchTemplates()
    ])

    // 使用nextTick确保DOM完全渲染后再设置编辑器内容
    setTimeout(() => {
      // 检查是否有从模板页面传递过来的数据
      const templateData = sessionStorage.getItem('templateData')
      if (templateData) {
        try {
          const data = JSON.parse(templateData)
          if (data.subject) {
            formData.subject = data.subject
          }
          if (data.body && editorRef.value) {
            editorRef.value.innerHTML = data.body
            formData.html_content = data.body
          }
          // 清除sessionStorage中的数据
          sessionStorage.removeItem('templateData')
        } catch (error) {
          console.error('解析模板数据失败:', error)
        }
      }
    }, 200)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.compose-email-container {
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

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-upload__tip) {
  color: #909399;
  font-size: 12px;
  margin-top: 8px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .compose-email-container {
    padding: 10px;
  }

  .editor-toolbar {
    flex-wrap: wrap;
  }

  .editor-content {
    min-height: 200px;
  }
}
</style>
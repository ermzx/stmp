<template>
  <el-container class="app-container">
    <el-header class="app-header">
      <div class="header-content">
        <h1 class="app-title">SMTP邮件管理系统</h1>
        <div class="header-actions">
          <el-button type="primary" size="small" @click="handleRefresh">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-header>
    
    <el-container>
      <el-aside width="200px" class="app-aside">
        <el-menu
          :default-active="activeMenu"
          router
          class="app-menu"
          @select="handleMenuSelect"
        >
          <el-menu-item index="/">
            <el-icon><DataBoard /></el-icon>
            <span>仪表板</span>
          </el-menu-item>
          <el-menu-item index="/smtp">
            <el-icon><Setting /></el-icon>
            <span>SMTP配置</span>
          </el-menu-item>
          <el-menu-item index="/compose">
            <el-icon><Edit /></el-icon>
            <span>撰写邮件</span>
          </el-menu-item>
          <el-menu-item index="/templates">
            <el-icon><Document /></el-icon>
            <span>邮件模板</span>
          </el-menu-item>
          <el-menu-item index="/history">
            <el-icon><Clock /></el-icon>
            <span>发送历史</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      
      <el-main class="app-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Refresh, DataBoard, Setting, Edit, Document, Clock } from '@element-plus/icons-vue'

const route = useRoute()
const activeMenu = ref('/')

onMounted(() => {
  activeMenu.value = route.path
})

const handleMenuSelect = (index) => {
  activeMenu.value = index
}

const handleRefresh = () => {
  window.location.reload()
}
</script>

<style scoped>
.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background-color: #409eff;
  color: white;
  padding: 0 20px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.app-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.app-aside {
  background-color: #f5f7fa;
  border-right: 1px solid #e4e7ed;
  overflow-y: auto;
}

.app-menu {
  border-right: none;
  height: 100%;
}

.app-main {
  background-color: #ffffff;
  padding: 20px;
  overflow-y: auto;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-aside {
    width: 64px !important;
  }
  
  .app-menu :deep(.el-menu-item span) {
    display: none;
  }
  
  .app-menu :deep(.el-menu-item) {
    justify-content: center;
  }
  
  .app-title {
    font-size: 16px;
  }
}
</style>
import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import './assets/styles/main.css'

// 配置Axios默认baseURL
axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:50022'

const app = createApp(App)

// 注册Element Plus
app.use(ElementPlus)
app.use(router)

app.mount('#app')
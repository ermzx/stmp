import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建axios实例
// 开发环境使用代理（相对路径），生产环境使用相对路径（同域名）
const service = axios.create({
  baseURL: import.meta.env.DEV ? '/api' : '.',
  timeout: 30000
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    // 从localStorage获取token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data
    
    // 检查业务状态码（后端返回的code字段）
    if (res && res.code !== 200) {
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    
    return res
  },
  error => {
    console.error('响应错误:', error)
    
    if (error.response) {
      switch (error.response.status) {
        case 400:
          ElMessage.error(error.response.data.message || '请求参数错误')
          break
        case 401:
          ElMessage.error('未授权，请重新登录')
          localStorage.removeItem('token')
          window.location.href = '/login'
          break
        case 403:
          ElMessage.error('拒绝访问')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          ElMessage.error(error.response.data.message || '网络错误')
      }
    } else if (error.request) {
      ElMessage.error('网络连接失败，请检查网络设置')
    } else {
      ElMessage.error('请求配置错误')
    }
    
    return Promise.reject(error)
  }
)

// SMTP配置相关API
export const getSmtpConfigs = () => {
  return service({
    url: '/smtp/configs',
    method: 'get'
  })
}

export const createSmtpConfig = (data) => {
  return service({
    url: '/smtp/configs',
    method: 'post',
    data
  })
}

export const updateSmtpConfig = (id, data) => {
  return service({
    url: `/smtp/configs/${id}`,
    method: 'put',
    data
  })
}

export const deleteSmtpConfig = (id) => {
  return service({
    url: `/smtp/configs/${id}`,
    method: 'delete'
  })
}

export const testSmtpConnection = (id) => {
  return service({
    url: `/smtp/configs/${id}/test`,
    method: 'post'
  })
}

export const setDefaultSmtpConfig = (id) => {
  return service({
    url: `/smtp/configs/${id}/default`,
    method: 'post'
  })
}

// 邮件发送相关API
export const sendEmail = (data) => {
  return service({
    url: '/email/send',
    method: 'post',
    data
  })
}

// 邮件模板相关API
export const getTemplates = () => {
  return service({
    url: '/templates',
    method: 'get'
  })
}

export const createTemplate = (data) => {
  return service({
    url: '/templates',
    method: 'post',
    data
  })
}

export const updateTemplate = (id, data) => {
  return service({
    url: `/templates/${id}`,
    method: 'put',
    data
  })
}

export const deleteTemplate = (id) => {
  return service({
    url: `/templates/${id}`,
    method: 'delete'
  })
}

// 发送历史相关API
export const getHistory = (params) => {
  return service({
    url: '/history',
    method: 'get',
    params
  })
}

export const getHistoryStatistics = () => {
  return service({
    url: '/history/statistics',
    method: 'get'
  })
}

export const deleteHistory = (id) => {
  return service({
    url: `/history/${id}`,
    method: 'delete'
  })
}

export default service
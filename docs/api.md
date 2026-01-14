# API文档

## 基础信息

- **Base URL**: `http://localhost:7700/api`
- **Content-Type**: `application/json`

## SMTP配置API

### 获取所有SMTP配置

```http
GET /api/smtp/configs
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "Gmail",
      "host": "smtp.gmail.com",
      "port": 587,
      "username": "user@gmail.com",
      "from_email": "user@gmail.com",
      "from_name": "User",
      "encryption": "tls",
      "is_default": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 创建SMTP配置

```http
POST /api/smtp/configs
Content-Type: application/json

{
  "name": "Gmail",
  "host": "smtp.gmail.com",
  "port": 587,
  "username": "user@gmail.com",
  "password": "password",
  "from_email": "user@gmail.com",
  "from_name": "User",
  "encryption": "tls",
  "is_default": false
}
```

### 获取单个SMTP配置

```http
GET /api/smtp/configs/:id
```

### 更新SMTP配置

```http
PUT /api/smtp/configs/:id
Content-Type: application/json

{
  "name": "Gmail",
  "host": "smtp.gmail.com",
  "port": 587,
  "username": "user@gmail.com",
  "password": "password",
  "from_email": "user@gmail.com",
  "from_name": "User",
  "encryption": "tls",
  "is_default": true
}
```

### 删除SMTP配置

```http
DELETE /api/smtp/configs/:id
```

### 测试SMTP连接

```http
POST /api/smtp/configs/:id/test
```

## 邮件发送API

### 发送邮件

```http
POST /api/email/send
Content-Type: application/json

{
  "smtp_config_id": 1,
  "to": "recipient@example.com",
  "cc": ["cc1@example.com", "cc2@example.com"],
  "bcc": ["bcc@example.com"],
  "subject": "邮件主题",
  "body": "<p>邮件正文</p>",
  "attachments": [
    {
      "filename": "file.pdf",
      "content": "base64编码的文件内容"
    }
  ]
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "邮件发送成功",
  "data": {
    "history_id": 1,
    "status": "success"
  }
}
```

## 邮件模板API

### 获取所有模板

```http
GET /api/templates
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "欢迎邮件",
      "subject": "欢迎加入我们",
      "body": "<p>欢迎内容...</p>",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 创建模板

```http
POST /api/templates
Content-Type: application/json

{
  "name": "欢迎邮件",
  "subject": "欢迎加入我们",
  "body": "<p>欢迎内容...</p>"
}
```

### 获取单个模板

```http
GET /api/templates/:id
```

### 更新模板

```http
PUT /api/templates/:id
Content-Type: application/json

{
  "name": "欢迎邮件",
  "subject": "欢迎加入我们",
  "body": "<p>更新后的内容...</p>"
}
```

### 删除模板

```http
DELETE /api/templates/:id
```

## 发送历史API

### 获取发送历史

```http
GET /api/history?status=success&page=1&page_size=10
```

**查询参数**:
- `status`: 可选，筛选状态（success/failed）
- `page`: 页码，默认1
- `page_size`: 每页数量，默认10

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "page_size": 10,
    "items": [
      {
        "id": 1,
        "smtp_config_id": 1,
        "to_email": "recipient@example.com",
        "cc_email": "[\"cc1@example.com\"]",
        "bcc_email": "[\"bcc@example.com\"]",
        "subject": "邮件主题",
        "body": "<p>邮件正文</p>",
        "attachments": "[{\"filename\":\"file.pdf\"}]",
        "status": "success",
        "error_message": "",
        "sent_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 获取单条历史记录

```http
GET /api/history/:id
```

### 删除历史记录

```http
DELETE /api/history/:id
```

## 健康检查

```http
GET /health
```

**响应示例**:
```json
{
  "status": "ok",
  "message": "服务运行正常"
}
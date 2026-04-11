# FinalAI API 文档

## 1. 基本信息

- 基础路径: `/api/v1`
- 健康检查: `GET /ping`
- 鉴权方式: `Authorization: Bearer <token>`

## 2. 统一响应格式

### 2.1 普通 JSON 响应

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

- `code=0` 表示成功。
- 失败时 `data` 可能不存在。

### 2.2 常见错误码

- `10001` 请求参数错误
- `10002` 未登录或无权限
- `10003` 用户名或密码错误
- `10004` 用户不存在
- `10005` 用户已存在
- `10006` 缺少 Authorization 头
- `10007` Authorization 格式错误
- `10008` token 无效或过期
- `20000` 服务器内部错误
- `20001` 服务器繁忙
- `20002` AI 模型调用失败
- `20003` 会话不存在

## 3. 用户模块

### 3.1 用户注册

- 方法: `POST /api/v1/user/register`
- 鉴权: 否
- 请求体:

```json
{
  "username": "alice",
  "password": "123456",
  "confirm_password": "123456"
}
```

- 成功响应示例:

```json
{
  "code": 0,
  "msg": "用户注册成功",
  "data": {
    "username": "alice"
  }
}
```

### 3.2 用户登录

- 方法: `POST /api/v1/user/login`
- 鉴权: 否
- 请求体:

```json
{
  "username": "alice",
  "password": "123456"
}
```

- 成功响应示例:

```json
{
  "code": 0,
  "msg": "用户登录成功",
  "data": {
    "username": "alice",
    "token": "<jwt-token>"
  }
}
```

### 3.3 获取用户信息

- 方法: `GET /api/v1/user/profile`
- 鉴权: 是
- 请求体: 无
- 成功响应示例:

```json
{
  "code": 0,
  "msg": "获取用户信息成功",
  "data": {
    "id": 1,
    "email": "",
    "username": "alice",
    "created_at": "2026-04-11T00:00:00Z",
    "updated_at": "2026-04-11T00:00:00Z"
  }
}
```

## 4. 会话与聊天模块

### 4.1 获取会话列表

- 方法: `GET /api/v1/chat/sessions`
- 鉴权: 是
- 请求体: 无
- 成功响应示例:

```json
{
  "code": 0,
  "msg": "获取会话列表成功",
  "data": {
    "sessions": [
      {
        "session_id": "f4c4210d-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
        "title": "f4c4210d-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
      }
    ]
  }
}
```

### 4.2 创建会话并发送消息

- 方法: `POST /api/v1/chat/create`
- 鉴权: 是
- 请求体:

```json
{
  "question": "你好，请介绍一下你自己",
  "modelType": "1"
}
```

- 成功响应示例:

```json
{
  "code": 0,
  "msg": "创建会话并发送成功",
  "data": {
    "sessionId": "f4c4210d-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "information": "你好，我是..."
  }
}
```

### 4.3 创建会话并流式发送消息

- 方法: `POST /api/v1/chat/create/stream`
- 鉴权: 是
- Content-Type: `application/json`
- 返回类型: `text/event-stream`
- 请求体:

```json
{
  "question": "你好，请流式回答",
  "modelType": "1"
}
```

- SSE 返回示例:

```text
data: {"sessionId": "f4c4210d-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}

data: 你好

data: ，我是...

data: [DONE]

```

### 4.4 在已有会话发送消息

- 方法: `POST /api/v1/chat/send`
- 鉴权: 是
- 请求体:

```json
{
  "sessionId": "f4c4210d-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "question": "继续上一个问题",
  "modelType": "1"
}
```

- 成功响应示例:

```json
{
  "code": 0,
  "msg": "发送消息成功",
  "data": {
    "information": "继续回答..."
  }
}
```

### 4.5 在已有会话流式发送消息

- 方法: `POST /api/v1/chat/send/stream`
- 鉴权: 是
- Content-Type: `application/json`
- 返回类型: `text/event-stream`
- 请求体:

```json
{
  "sessionId": "f4c4210d-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "question": "继续流式回答",
  "modelType": "1"
}
```

- SSE 返回示例:

```text
data: 分片1

data: 分片2

data: [DONE]

```

### 4.6 获取聊天历史

- 方法: `POST /api/v1/chat/history`
- 鉴权: 是
- 请求体:

```json
{
  "sessionId": "f4c4210d-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

- 成功响应示例:

```json
{
  "code": 0,
  "msg": "获取聊天历史成功",
  "data": {
    "history": [
      {
        "is_user": true,
        "content": "你好"
      },
      {
        "is_user": false,
        "content": "你好，我是AI助手"
      }
    ]
  }
}
```

## 5. 图片识别模块

### 5.1 图片识别

- 方法: `POST /api/v1/image/recognize`
- 鉴权: 是
- Content-Type: `multipart/form-data`
- 表单字段:
  - `image`: 文件，必填

- 成功响应示例:

```json
{
  "code": 0,
  "msg": "图片识别成功",
  "data": {
    "class_name": "golden retriever"
  }
}
```

## 6. 联调建议

- SSE 接口不要用普通 `fetch().json()` 解析，建议使用 `EventSource` 或流式读取器。
- 除注册/登录外，其余接口都需要 `Authorization` 请求头。
- 前端按 `code` 判断业务成功，HTTP 状态码用于传输层补充判断。

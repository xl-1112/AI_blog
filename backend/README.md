# Liang Blog CMS Backend

这是给当前个人网站准备的 Go 后台 API。它负责持久化站点资料、Logo、首页介绍、精选文章、标签、文章、关于我、工作栈、工作经历、邮箱和 GitHub 地址。

## 启动

```powershell
cd backend
$env:ADMIN_TOKEN="change-me"
go run ./cmd/server
```

默认地址：

- API: `http://127.0.0.1:8080`
- 数据文件: `backend/data/site.json`
- 上传目录: `backend/uploads`

如果没有设置 `ADMIN_TOKEN`，开发环境会使用 `dev-admin-token`。正式部署前一定要设置自己的 token。

## 环境变量

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `ADDR` | `:8080` | 服务监听地址 |
| `ADMIN_TOKEN` | `dev-admin-token` | 后台写接口鉴权 token |
| `DATA_PATH` | `data/site.json` | JSON 数据文件 |
| `UPLOAD_DIR` | `uploads` | Logo 等上传文件目录 |
| `CORS_ORIGIN` | `http://127.0.0.1:5173,http://localhost:5173` | 允许访问 API 的前端地址，多个用逗号分隔 |

## 公开接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `GET` | `/healthz` | 健康检查 |
| `GET` | `/api/site` | 获取完整公开站点内容，不返回草稿文章 |
| `GET` | `/api/tags` | 获取文章标签 |
| `GET` | `/api/articles` | 获取公开文章列表 |
| `GET` | `/api/articles/{id}` | 获取公开文章详情 |
| `GET` | `/uploads/{filename}` | 访问上传后的 Logo 文件 |

## 后台接口

后台接口需要请求头：

```http
Authorization: Bearer change-me
```

也可以用：

```http
X-Admin-Token: change-me
```

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `GET` | `/api/admin/content` | 获取完整后台内容，包含草稿 |
| `PUT` | `/api/admin/content` | 整体替换所有内容 |
| `PUT` | `/api/admin/site` | 修改网站名称、描述、Logo URL、定位、位置、邮箱、GitHub |
| `PUT` | `/api/admin/home` | 修改首页介绍和精选文章 ID |
| `PUT` | `/api/admin/about` | 修改关于我描述、工作栈、工作经历 |
| `PUT` | `/api/admin/tags` | 修改标签列表 |
| `POST` | `/api/admin/logo` | 上传 Logo，表单字段名为 `logo` |
| `GET` | `/api/admin/articles` | 获取所有文章列表，包含草稿 |
| `POST` | `/api/admin/articles` | 新建文章 |
| `GET` | `/api/admin/articles/{id}` | 获取文章详情 |
| `PUT` | `/api/admin/articles/{id}` | 更新文章 |
| `DELETE` | `/api/admin/articles/{id}` | 删除文章 |

## 示例

获取完整内容：

```powershell
Invoke-RestMethod http://127.0.0.1:8080/api/site
```

修改站点基础信息：

```powershell
$headers = @{ Authorization = "Bearer change-me" }
$body = @{
  name = "Liang"
  description = "一个产品经理的个人网站"
  logoUrl = "/uploads/blog_logo_header.png"
  role = "产品经理"
  location = "广州"
  contact = @{
    email = "xl1258763@gmail.com"
    github = "https://github.com/xl-1112"
  }
} | ConvertTo-Json -Depth 5

Invoke-RestMethod `
  -Method Put `
  -Uri http://127.0.0.1:8080/api/admin/site `
  -Headers $headers `
  -ContentType "application/json" `
  -Body $body
```

上传 Logo：

```powershell
$headers = @{ Authorization = "Bearer change-me" }
$form = @{
  logo = Get-Item "..\public\blog_logo_header.png"
}

Invoke-RestMethod `
  -Method Post `
  -Uri http://127.0.0.1:8080/api/admin/logo `
  -Headers $headers `
  -Form $form
```

新建文章：

```powershell
$headers = @{ Authorization = "Bearer change-me" }
$article = @{
  id = "new-product-note"
  title = "一次需求优先级判断记录"
  date = "2026-05-15"
  summary = "记录如何从用户价值、业务收益和实现成本判断需求优先级。"
  tags = @("产品复盘", "需求分析")
  featured = $false
  draft = $false
  content = "## 背景`n`n这里写 Markdown 正文。"
} | ConvertTo-Json -Depth 5

Invoke-RestMethod `
  -Method Post `
  -Uri http://127.0.0.1:8080/api/admin/articles `
  -Headers $headers `
  -ContentType "application/json" `
  -Body $article
```

## 数据结构

`site.json` 是单文件 JSON 存储，适合当前个人网站首版后台。后续如果需要多人协作、版本历史或大量文章，可以把 `Store` 换成 SQLite 或 PostgreSQL，路由层不用大改。

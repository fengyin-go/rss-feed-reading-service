# RSS 订阅系统

纯 Go 标准库实现的 RSS 订阅管理服务。

## 运行

```bash
cd origin
go run ./cmd/server
```

默认监听 `:8080`，通过 HTTP API 提供 RSS 订阅与阅读数据。

环境变量：
- `PORT` / `ADDR`：监听地址
- `MAX_PAGE_SIZE`：最大分页大小（默认 100）
- `ADMIN_TOKEN`：admin 接口鉴权令牌（默认 `rss-admin-secret`）
- `LOG_LEVEL`：日志级别（debug/info/warn/error）

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/feeds | 创建订阅源 |
| GET | /api/feeds | 订阅源列表（支持 category/status/keyword 筛选） |
| GET | /api/feeds/{id} | 获取订阅源 |
| PUT | /api/feeds/{id} | 更新订阅源 |
| DELETE | /api/feeds/{id} | 删除订阅源 |
| POST | /api/feeds/{id}/pause | 暂停订阅源 |
| POST | /api/feeds/{id}/resume | 恢复订阅源 |
| POST | /api/articles | 创建文章 |
| GET | /api/articles | 文章列表（支持 feed_id/author/is_read/is_starred/keyword 筛选） |
| GET | /api/articles/{id} | 获取文章 |
| PUT | /api/articles/{id} | 更新文章 |
| DELETE | /api/articles/{id} | 删除文章 |
| POST | /api/articles/{id}/read | 标记已读 |
| POST | /api/articles/{id}/star | 收藏 |
| POST | /api/articles/{id}/unstar | 取消收藏 |
| GET | /api/feeds/{id}/export | 导出某订阅源的全部文章 |
| POST | /api/categories | 创建分类 |
| GET | /api/categories | 分类列表 |
| GET | /api/categories/{id} | 获取分类 |
| PUT | /api/categories/{id} | 更新分类 |
| DELETE | /api/categories/{id} | 删除分类 |
| POST | /api/fetch-tasks | 创建抓取任务 |
| GET | /api/fetch-tasks | 任务列表 |
| GET | /api/fetch-tasks/{id} | 获取任务 |
| POST | /api/fetch-tasks/{id}/start | 启动任务 |
| POST | /api/fetch-tasks/{id}/complete | 完成任务 |
| DELETE | /api/fetch-tasks/{id} | 删除任务 |
| POST | /api/fetch-logs | 创建日志 |
| GET | /api/fetch-logs | 日志列表 |
| GET | /api/fetch-logs/{id} | 获取日志 |
| DELETE | /api/fetch-logs/{id} | 删除日志 |
| POST | /api/users | 创建用户 |
| GET | /api/users | 用户列表 |
| GET | /api/users/{id} | 获取用户 |
| PUT | /api/users/{id} | 更新用户 |
| DELETE | /api/users/{id} | 删除用户 |
| POST | /api/subscriptions | 创建订阅关系 |
| GET | /api/subscriptions | 订阅列表 |
| GET | /api/subscriptions/{id} | 获取订阅 |
| DELETE | /api/subscriptions/{id} | 删除订阅 |
| POST | /api/bookmarks | 创建收藏 |
| GET | /api/bookmarks | 收藏列表 |
| GET | /api/bookmarks/{id} | 获取收藏 |
| DELETE | /api/bookmarks/{id} | 删除收藏 |
| POST | /api/notifications | 创建通知 |
| GET | /api/notifications | 通知列表 |
| GET | /api/notifications/{id} | 获取通知 |
| POST | /api/notifications/{id}/read | 标记通知已读 |
| DELETE | /api/notifications/{id} | 删除通知 |
| POST | /api/tags | 创建标签 |
| GET | /api/tags | 标签列表 |
| GET | /api/tags/{id} | 获取标签 |
| PUT | /api/tags/{id} | 更新标签 |
| DELETE | /api/tags/{id} | 删除标签 |
| GET | /api/tags/{id}/articles | 标签下文章 |
| POST | /api/article-tags | 创建文章标签关联 |
| GET | /api/article-tags | 关联列表 |
| GET | /api/article-tags/{id} | 获取关联 |
| DELETE | /api/article-tags/{id} | 删除关联 |
| GET | /api/articles/{id}/tags | 文章标签 |
| POST | /api/batch/articles | 批量创建文章 |
| POST | /api/batch/articles/read | 批量标记已读 |
| POST | /api/batch/articles/star | 批量收藏 |
| POST | /api/batch/articles/unstar | 批量取消收藏 |
| POST | /api/batch/articles/delete | 批量删除文章 |
| POST | /api/batch/feeds | 批量创建订阅源 |
| GET | /api/search | 全局搜索 |
| GET | /api/articles/recent | 最近文章 |
| GET | /api/articles/unread | 未读文章 |
| GET | /api/articles/starred | 收藏文章 |
| POST | /api/import/feeds | 导入订阅源 |
| POST | /api/import/feeds/{id}/articles | 导入文章 |
| GET | /api/export/feeds | 导出订阅源 |
| GET | /api/export/all | 导出全部数据 |
| GET | /api/stats/global | 全局统计 |
| GET | /api/stats/articles-by-category | 按分类文章数 |
| GET | /api/stats/articles-by-feed | 按订阅源文章数 TOP10 |
| GET | /api/stats/fetch-rate-by-feed | 每订阅源抓取成功率 |
| GET | /api/reports/articles-by-date | 按日期文章统计 |
| GET | /api/reports/user-activity | 用户活跃度 |
| GET | /api/reports/feed-health | 订阅源健康度 |
| GET | /api/reports/unread-by-feed | 每订阅源未读数 |
| GET | /api/reports/task-duration | 任务耗时统计 |

## 结构

```
origin/
├── cmd/server/main.go
├── internal/
│   ├── app/app.go
│   ├── config/config.go
│   ├── model/        # 8 个实体模型
│   ├── store/        # Store 接口 + 内存实现 + 测试
│   ├── service/      # 业务逻辑 + 统计 + 测试
│   └── handler/      # HTTP 处理器 + 中间件 + 测试
└── pkg/
    ├── httpx/httpx.go
    ├── idgen/idgen.go
    └── logger/logger.go
```

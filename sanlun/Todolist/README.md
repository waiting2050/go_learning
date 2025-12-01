# West2 Online Todo List Backend (备忘录后端)

这是一个基于 Go 语言 (Gin 框架) 开发的 RESTful API 备忘录后端服务。项目采用了标准的企业级**三层架构**，实现了用户认证、待办事项管理、**Redis 缓存加速**以及**优先级排序**等高级功能。

## 🛠 技术栈

- **编程语言**: Go (Golang)
- **Web 框架**: Gin
- **数据库 ORM**: GORM
- **数据库**: MySQL 8.0+
- **缓存**: Redis
- **鉴权**: JWT (JSON Web Token)
- **密码加密**: Bcrypt

## Mw 项目结构说明 (三层架构)

本项目遵循 `Controller` -> `Service` -> `DAO` 的分层设计模式，职责清晰，易于维护。

```text
todo-list/
├── main.go                # 程序入口：负责初始化资源(MySQL, Redis)并启动路由
├── go.mod                 # Go 模块依赖管理文件
├── controllers/           # 【控制层】负责处理 HTTP 请求与响应
│   ├── user_controller.go # 处理用户注册、登录请求
│   └── todo_controller.go # 处理待办事项的增删改查、参数校验、结果封装
├── service/               # 【业务层】负责核心业务逻辑与缓存策略
│   └── todo_service.go    # 封装待办业务：Read-Through 缓存、优先级排序、批量逻辑
├── dao/                   # 【数据访问层】负责数据库连接与配置
│   ├── mysql.go           # MySQL 连接初始化与表结构自动迁移
│   └── redis.go           # Redis 连接初始化
├── models/                # 【模型层】定义数据结构
│   ├── user.go            # 用户表结构定义
│   ├── todo.go            # 待办事项表结构定义
│   └── common.go          # 通用响应结构体 (Code, Msg, Data)
├── middleware/            # 【中间件】
│   └── jwt.go             # JWT 鉴权中间件：拦截请求并解析 Token
├── pkg/                   # 【公共包】
│   └── utils/
│       └── jwt.go         # JWT 工具包：Token 的签发与校验
└── routers/               # 【路由层】
    └── routers.go         # 定义 API 路由规则，绑定 Controller 与 Middleware
```
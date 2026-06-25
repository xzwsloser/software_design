# 基于携程旅游的可视化推荐系统设计与实现

> 华中科技大学电信学院软件课程设计项目

## 项目简介

本项目是一个基于携程旅游数据的可视化推荐系统，采用 **前后端分离 + 微服务** 架构，分为前端、后端、推荐系统三个核心模块，并包含数据爬取工具。系统通过双塔模型实现个性化景点推荐，支持用户登录注册、景点浏览、点赞收藏、评论查看、数据可视化等功能。

## 项目架构

```
┌─────────────────────────────────────────────────────────┐
│                      用户浏览器                           │
│              (Vue 3 + Element Plus + ECharts)            │
└─────────────────┬───────────────────────────────────────┘
                  │ HTTP (Nginx 反向代理, 端口 8888)
┌─────────────────▼───────────────────────────────────────┐
│                   后端服务 (Gin, 端口 9999)                │
│  ┌──────────┐  ┌──────────┐  ┌────────────────────┐    │
│  │ Handler  │  │ Service  │  │ Middleware (JWT/CORS)│    │
│  └────┬─────┘  └────┬─────┘  └────────────────────┘    │
│       │              │                                    │
│  ┌────▼──────────────▼─────┐  ┌──────────┐              │
│  │      Model (GORM)       │  │  Cache   │              │
│  └───────────┬─────────────┘  │ (Redis)  │              │
│              │                 └──────────┘              │
└──────────────┼──────────────────────────────────────────┘
               │                          │
      ┌────────▼────────┐    ┌───────────▼──────────────┐
      │   MySQL 数据库    │    │   gRPC Client (端口 7777)  │
      └─────────────────┘    └───────────┬──────────────┘
                                         │
┌────────────────────────────────────────▼──────────────┐
│              推荐系统服务 (Python, gRPC Server)          │
│  ┌──────────────────┐  ┌──────────────────────────┐   │
│  │  双塔模型 (PyTorch) │  │  Milvus 向量数据库         │   │
│  │  UserTower        │  │  (向量存储与相似度检索)      │   │
│  │  SiteTower        │  └──────────────────────────┘   │
│  └──────────────────┘                                   │
└────────────────────────────────────────────────────────┘
```

## 目录结构

```
software_design/
├── fronted/                    # 前端项目 (Vue 3)
│   ├── src/
│   │   ├── components/         # Vue 组件
│   │   │   ├── AuthView.vue        # 登录/注册页
│   │   │   ├── SiteList.vue        # 景点列表页
│   │   │   ├── SiteDetail.vue      # 景点详情页
│   │   │   ├── RecommendedSites.vue # 推荐景点页
│   │   │   ├── LikedSites.vue      # 点赞景点页
│   │   │   ├── ViewedSites.vue     # 浏览记录页
│   │   │   ├── DataVisualization.vue # 数据可视化页
│   │   │   ├── UserInfo.vue        # 用户信息页
│   │   │   ├── CommentSection.vue  # 评论展示组件
│   │   │   ├── LikeButton.vue      # 点赞按钮组件
│   │   │   └── LoginRegister.vue   # 登录注册表单组件
│   │   ├── router/index.js     # 路由配置
│   │   ├── stores/             # Pinia 状态管理
│   │   ├── axios.js            # HTTP 请求封装
│   │   ├── App.vue             # 根组件
│   │   └── main.js             # 入口文件
│   ├── run/
│   │   ├── dist/               # 构建产物
│   │   └── nginx.conf          # Nginx 配置
│   ├── vite.config.js          # Vite 构建配置
│   └── package.json
├── backend/                    # 后端项目 (Go)
│   ├── handler/                # 路由处理器 (Controller 层)
│   │   ├── router.go           # 路由注册
│   │   ├── user.go             # 用户登录/注册
│   │   ├── user_info.go        # 用户信息管理
│   │   ├── site.go             # 景点查询
│   │   ├── comment.go          # 评论查询
│   │   ├── like.go             # 点赞管理
│   │   ├── view.go             # 浏览记录
│   │   ├── oss.go              # OSS 图片获取
│   │   └── rec_sys.go          # 推荐结果查询
│   ├── service/                # 业务逻辑层
│   │   ├── user_service.go     # 用户服务
│   │   ├── site_service.go     # 景点服务
│   │   ├── comment_service.go  # 评论服务
│   │   ├── like_service.go     # 点赞服务
│   │   ├── view_service.go     # 浏览记录服务
│   │   └── rec_sys_service.go  # 推荐系统服务 (Pipeline 模式)
│   ├── model/                  # 数据模型层 (GORM)
│   │   ├── mysql.go            # MySQL 连接初始化
│   │   ├── user.go             # 用户模型
│   │   ├── site.go             # 景点模型
│   │   ├── comment.go          # 评论模型
│   │   ├── like.go             # 点赞模型
│   │   └── view.go             # 浏览记录模型
│   ├── cache/                  # Redis 缓存层
│   │   ├── redis.go            # Redis 连接初始化
│   │   ├── like.go             # 点赞缓存
│   │   ├── view.go             # 浏览记录缓存
│   │   └── rec_sys.go          # 推荐结果缓存
│   ├── middleware/             # 中间件
│   │   ├── jwt.go              # JWT 鉴权中间件
│   │   └── cors.go             # 跨域中间件
│   ├── rpc/                    # gRPC 客户端
│   │   ├── rec_sys_client.go   # 推荐系统 gRPC 客户端
│   │   ├── proto/rec_sys.proto # Protobuf 定义
│   │   └── pb/                 # 生成的 Go 代码
│   ├── dto/                    # 数据传输对象
│   ├── utils/                  # 工具类
│   │   ├── config.go           # 配置加载
│   │   ├── logrus.go           # 日志初始化
│   │   ├── oss.go              # 腾讯云 COS 客户端
│   │   └── common.go           # 通用工具函数
│   ├── main.go                 # 入口文件
│   ├── Makefile                # 构建脚本
│   └── go.mod
├── rec-sys/                    # 推荐系统 (Python)
│   └── code/
│       ├── rec-sys/
│       │   ├── rec_model.py        # 双塔模型定义与训练
│       │   ├── rec_sys.py          # 推荐系统核心逻辑 (召回+精排+重排)
│       │   ├── rec_sys_server.py   # gRPC 服务端
│       │   ├── application.py      # 应用层 (向量入库/命令行推荐)
│       │   └── rpc/                # gRPC 协议定义与生成代码
│       ├── data_process.py         # 数据预处理
│       ├── feature_extract.py      # 特征提取
│       ├── comment_data_analysis.ipynb  # 评论数据分析
│       ├── comment_feature_extract.ipynb # 评论特征提取
│       ├── site_data_analysis.ipynb     # 景点数据分析
│       ├── milvus_init.ipynb       # Milvus 初始化
│       ├── valid.ipynb             # 模型验证
│       ├── weights/                # 模型权重
│       └── run/                    # 训练结果图表
├── crawler/                    # 数据爬虫 (Python)
│   ├── crawler.py                  # 主爬虫脚本
│   ├── fetch_positive_comment_extra.py  # 好评数据爬取
│   ├── fetch_negative_comment_extra.py  # 差评数据爬取
│   ├── fetch_site_extra.py         # 景点额外信息爬取
│   └── requirements.txt
├── startup.sh                  # 一键启动脚本
├── stop.sh                     # 一键停止脚本
└── README.md
```

## 技术栈

### 前端

| 技术 | 说明 |
|------|------|
| Vue 3 | 前端框架，采用 Composition API |
| Vite 7 | 构建工具 |
| Element Plus | UI 组件库 |
| Pinia | 状态管理 |
| Vue Router 4 | 路由管理 |
| Axios | HTTP 请求库 |
| ECharts 6 | 数据可视化图表 |

### 后端

| 技术 | 说明 |
|------|------|
| Go | 编程语言 |
| Gin | HTTP Web 框架 |
| GORM | ORM 框架，操作 MySQL |
| go-redis | Redis 客户端 |
| gRPC + Protobuf | 与推荐系统通信 |
| golang-jwt | JWT 用户认证 |
| Logrus | 日志框架 |
| 腾讯云 COS SDK | 对象存储（词云图、饼图等） |

### 推荐系统

| 技术 | 说明 |
|------|------|
| Python | 编程语言 |
| PyTorch | 深度学习框架 |
| Milvus | 向量数据库，存储景点/用户 Embedding 并支持相似度检索 |
| gRPC | 提供推荐服务接口 |
| PyMySQL | MySQL 数据库连接 |
| scikit-learn | 模型评估 (AUC) |
| pandas | 数据处理 |

### 数据爬虫

| 技术 | 说明 |
|------|------|
| Python | 编程语言 |
| requests | HTTP 请求 |
| BeautifulSoup4 | HTML 解析 |
| PyMySQL | 数据写入 MySQL |

### 基础设施

| 技术 | 说明 |
|------|------|
| MySQL | 关系型数据库，存储用户、景点、评论等数据 |
| Redis | 缓存数据库，缓存点赞、浏览记录、推荐结果 |
| Milvus | 向量数据库，存储 Embedding 向量 |
| Nginx | 反向代理 + 静态资源服务 |
| 腾讯云 COS | 对象存储，存储词云图等静态资源 |

## 核心功能

### 1. 用户系统
- **注册/登录**：基于 JWT Token 的身份认证，Token 有效期 7 天
- **用户信息管理**：支持修改个人资料（城市、出游类型、偏好等）
- **路由守卫**：未登录用户自动跳转登录页

### 2. 景点浏览
- **景点列表**：分页查询，展示景点名称、评分、热度、地址等信息
- **景点详情**：查看景点详细介绍、开放时间、联系方式等
- **评论查看**：支持好评/差评的分页查询与数量统计

### 3. 用户互动
- **点赞功能**：用户可对景点进行点赞/取消点赞，查看点赞列表
- **浏览记录**：自动记录用户浏览的景点，支持历史查看

### 4. 数据可视化
- 使用 ECharts 展示景点相关统计图表
- 通过腾讯云 COS 加载词云图、饼图等静态可视化资源

### 5. 个性化推荐（核心）

推荐系统采用 **召回 → 精排 → 重排** 三层架构：

#### 召回层（双塔模型）
- **User Tower**：将用户特征（城市、出游类型、偏好类型、价格敏感度、关注点等）编码为 64 维 Embedding 向量
- **Site Tower**：将景点特征（评分、热度、介绍 Embedding、地址、价格、好评率等）编码为 64 维 Embedding 向量
- 使用 Attention Pooling 处理多标签稀疏特征（如偏好类型、出游动机等）
- 训练时使用 BCEWithLogitsLoss 作为损失函数，Adam 优化器

#### 精排层
- 基于 Milvus 向量数据库进行内积（IP）相似度检索，召回 Top-K 候选景点
- 结合景点自身评分、热度、评论数等特征进行加权重算分数

#### 重排层
- 多路召回融合：将双塔模型召回结果与基于用户点赞景点的相似景点召回结果进行融合
- 按设定比例交替采样，保证推荐结果的多样性

#### 推荐 Pipeline
- 后端使用 **Channel + Goroutine** 实现异步推荐 Pipeline
- 用户登录/注册/更新信息时自动触发推荐计算
- 推荐结果缓存至 Redis，减少重复计算

## API 接口

所有接口均以 `/api` 为前缀，除登录/注册外均需 JWT 鉴权（Header 中携带 `Authorization`）。

| 模块 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 用户 | POST | `/api/user/login` | 用户登录 |
| 用户 | POST | `/api/user/register` | 用户注册 |
| 用户信息 | GET | `/api/userInfo/user` | 获取用户信息 |
| 用户信息 | POST | `/api/userInfo/update` | 更新用户信息 |
| 景点 | POST | `/api/site/query/list` | 分页查询景点列表 |
| 景点 | GET | `/api/site/query/:siteIndex` | 按索引查询景点详情 |
| 景点 | POST | `/api/site/query/siteList` | 批量查询景点 |
| 评论 | POST | `/api/comment/positive/:siteIndex` | 分页查询好评 |
| 评论 | POST | `/api/comment/negative/:siteIndex` | 分页查询差评 |
| 评论 | GET | `/api/comment/count/positive/:siteIndex` | 好评数量统计 |
| 评论 | GET | `/api/comment/count/negative/:siteIndex` | 差评数量统计 |
| 点赞 | GET | `/api/like/like/:siteIndex` | 点赞 |
| 点赞 | GET | `/api/like/cancel/:siteIndex` | 取消点赞 |
| 点赞 | GET | `/api/like/isLike/:siteIndex` | 是否已点赞 |
| 点赞 | GET | `/api/like/siteList` | 获取点赞景点列表 |
| 点赞 | GET | `/api/like/userList/:siteIndex` | 获取点赞用户列表 |
| 浏览 | GET | `/api/view/view/:siteIndex` | 记录浏览 |
| 浏览 | GET | `/api/view/siteList` | 获取浏览记录 |
| 浏览 | GET | `/api/view/userList/:siteIndex` | 获取浏览用户列表 |
| OSS | GET | `/api/oss/touristType` | 获取游客类型饼图 |
| 推荐 | GET | `/api/rec/siteIdxList` | 获取推荐景点列表 |

## 快速开始

### 环境要求

- Go 1.25+
- Node.js 20.19+ / 22.12+
- Python 3.10+
- MySQL 8.0+
- Redis 6.0+
- Milvus 2.3+
- Nginx

### 本地开发

**1. 启动后端**

```bash
cd backend
# 配置 config.json（数据库、Redis、JWT、OSS、gRPC 等）
make depend    # 下载依赖
make build     # 构建
make run       # 运行（默认端口 9999）
```

**2. 启动前端**

```bash
cd fronted
npm install
npm run dev    # 启动开发服务器
```

**3. 启动推荐系统**

```bash
cd rec-sys/code/rec-sys
python rec_sys_server.py   # 启动 gRPC 服务（默认端口 7777）
```

**4. 一键启动/停止**

```bash
./startup.sh   # 构建并启动后端 + 前端（Nginx）
./stop.sh      # 停止所有服务
```

### 生产部署

前端使用 Vite 构建后通过 Nginx 托管静态资源，Nginx 同时反向代理 `/api/` 请求到后端服务。详见 [fronted/run/nginx.conf](file:///Users/bytedance/projects/software_design/fronted/run/nginx.conf)。

## gRPC 通信协议

后端与推荐系统之间通过 gRPC 进行通信，协议定义见 [rec_sys.proto](file:///Users/bytedance/projects/software_design/backend/rpc/proto/rec_sys.proto)：

```protobuf
service RecSysService {
    rpc GetRecResult(GetRecResultReq) returns (GetRecResultResp);
}
```

请求参数包含用户 ID、地址、出游类型、价格敏感度、偏好类型、出游动机、关注点、是否更新向量、返回数量限制、已点赞景点列表等。

## 数据流

```
用户操作 → 前端 (Vue) → HTTP 请求 → 后端 (Gin)
    ├── 登录/注册 → JWT 签发 → 触发推荐 Pipeline
    ├── 景点查询 → MySQL 查询
    ├── 点赞/浏览 → Redis 缓存 + MySQL 持久化
    └── 推荐请求 → Redis 缓存读取 → (缓存未命中时) gRPC → 推荐系统
                                                          ├── 双塔模型推理
                                                          ├── Milvus 向量检索
                                                          └── 精排 + 重排 → 返回结果
```

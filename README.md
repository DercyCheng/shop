# Shop 电商系统

Shop 是一个基于 Go 语言和微服务架构的完整电商解决方案，包含用户、商品、库存、订单和用户操作等核心微服务，以及 Vue 3 构建的管理后台。

## 项目概览

该项目采用前后端分离的设计理念，后端基于微服务架构，前端使用 Vue 3 构建现代化的管理界面。整个系统涵盖了电商平台所需的全部核心功能，包括：

- 用户注册、登录和管理
- 商品管理、分类、品牌和搜索
- 购物车和订单管理
- 库存管理
- 支付集成
- 用户收藏、地址管理和留言

## 系统架构

系统采用典型的微服务架构，分为四个主要层次：

1. **基础设施层**：包括 MySQL、Redis、ElasticSearch、RocketMQ、Consul、Nacos 等基础组件
2. **服务层（SRV）**：核心业务逻辑实现，包括用户服务、商品服务、库存服务、订单服务和用户操作服务
3. **API 层（Web）**：对外提供 HTTP 接口，负责参数校验、权限验证等，调用 SRV 层完成业务逻辑
4. **前端层**：包括管理后台界面

![架构图](./image/README/1744869775075.png)

## 技术栈

### 后端技术栈

- **编程语言**：Go
- **微服务框架**：gRPC
- **API 框架**：Gin
- **ORM 框架**：GORM
- **服务发现**：Consul
- **配置中心**：Nacos
- **消息队列**：RocketMQ
- **缓存**：Redis
- **搜索引擎**：ElasticSearch
- **数据库**：MySQL
- **容器化**：Docker & Docker Compose
- **分布式追踪**：Jaeger

### 前端技术栈

- **框架**：Vue 3
- **语言**：TypeScript
- **UI 组件库**：Element Plus
- **状态管理**：Pinia
- **路由**：Vue Router
- **HTTP 客户端**：Axios
- **构建工具**：Vite

## 主要功能模块

### 用户服务 (User Service)

- 用户注册、登录和信息管理
- JWT 认证
- 权限管理
- 短信验证
- 图形验证码

### 商品服务 (Goods Service)

- 商品管理
- 分类管理（多级分类）
- 品牌管理
- 商品搜索（ElasticSearch）
- 轮播图管理

### 库存服务 (Inventory Service)

- 库存管理
- 库存锁定和释放
- 并发控制
- 库存操作历史记录

### 订单服务 (Order Service)

- 购物车管理
- 订单创建和管理
- 支付集成（支付宝）
- 订单状态流转
- 超时订单处理

### 用户操作服务 (UserOp Service)

- 商品收藏
- 地址管理
- 用户留言

### 管理后台 (Admin Dashboard)

- 用户管理
- 商品管理
- 订单管理
- 库存管理
- 系统监控

## 技术亮点

1. **微服务架构**：系统采用微服务架构，各服务独立部署，提高系统弹性和可扩展性。

2. **分布式事务**：使用基于消息队列的最终一致性方案处理跨服务事务，如订单创建与库存锁定。

3. **高并发处理**：库存服务采用乐观锁和缓存策略，能够应对高并发场景。

4. **全文搜索**：集成 ElasticSearch 实现高效的商品搜索功能。

5. **容器化部署**：使用 Docker 和 Docker Compose 简化部署流程。

6. **分布式配置**：使用 Nacos 实现配置的统一管理和动态更新。

7. **服务发现**：通过 Consul 实现服务注册与发现，提高系统的可用性。

8. **异步通信**：使用 RocketMQ 实现服务间的异步通信，提高系统吞吐量和容错性。

## 项目结构

```
shop/
├── docker-compose.yml  # Docker 部署配置
├── Dockerfile          # Docker 构建文件
├── init.sql            # 数据库初始化脚本
├── README.md           # 项目说明
├── run.py              # 运行脚本
├── doc/                # 项目文档
├── image/              # 图片资源
├── scripts/            # 脚本文件
├── shop_admin/         # 管理后台前端
├── shop_api/           # API 层
└── shop_srv/           # 服务层
```

## 快速开始

详细的环境搭建和启动指南请参考 [环境搭建文档](./doc/环境搭建.md)。

### 使用 Docker Compose 启动

```bash
# 启动所有服务
docker-compose up -d
```

### 本地开发启动

```bash
# 启动服务层
cd scripts
./start.sh start_srv

# 启动 API 层
./start.sh start_api

# 启动前端
cd ../shop_admin
npm run dev
```

## 文档

- [环境搭建指南](./doc/环境搭建.md)
- [用户服务文档](./doc/用户服务.md)
- [商品服务文档](./doc/商品服务.md)
- [库存服务文档](./doc/库存服务.md)
- [订单服务文档](./doc/订单服务.md)
- [用户操作服务文档](./doc/用户操作服务.md)
- [面试介绍指南](./doc/interview.md)

## 作者

- dercy - [邮箱](mailto:dercyc@example.com)

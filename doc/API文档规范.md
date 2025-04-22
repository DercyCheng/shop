# API文档规范指南

## Swagger文档规范

本项目使用Swagger提供API文档，所有微服务API必须按照以下规范完善Swagger文档。

### 安装与配置

在每个Web服务中安装Swagger依赖：

```bash
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
go get -u github.com/swaggo/swag/cmd/swag
```

确保在每个服务的`go.mod`文件中有这些依赖项。

### Swagger注释规范

#### 主文档注释

每个API服务的main.go文件顶部需添加以下注释：

```go
// @title Shop微服务系统 - XXX服务API文档
// @version 1.0
// @description Shop电商项目XXX服务接口文档
// @termsOfService https://github.com/username/shop

// @contact.name API Support
// @contact.url https://github.com/username/shop/issues
// @contact.email dercyc@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:80XX
// @BasePath /api/v1
// @schemes http https
```

#### API接口注释

每个API处理函数必须包含以下格式的注释：

```go
// @Summary 接口简短描述
// @Description 接口详细描述
// @Tags 接口分组名称
// @Accept json
// @Produce json
// @Param 参数名称 参数位置 参数类型 是否必须 描述 "示例值"
// @Success 200 {object} 返回结构体 "成功返回"
// @Failure 400 {object} 错误返回结构体 "参数错误"
// @Failure 500 {object} 错误返回结构体 "服务器内部错误"
// @Router /路径 [请求方法]
```

示例：

```go
// @Summary 获取用户信息
// @Description 根据用户ID获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID" "1"
// @Success 200 {object} models.UserResponse "成功返回"
// @Failure 400 {object} models.Response "参数错误"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /users/{id} [get]
```

### 生成与访问Swagger文档

在每个微服务目录下执行以下命令生成Swagger文档：

```bash
swag init
```

启动服务后，访问以下路径查看API文档：

- 用户服务: http://localhost:8021/swagger/index.html
- 商品服务: http://localhost:8022/swagger/index.html
- 订单服务: http://localhost:8023/swagger/index.html
- 用户操作: http://localhost:8024/swagger/index.html
- OSS 服务: http://localhost:8025/swagger/index.html

### 最佳实践

1. **保持更新**：代码修改后及时更新Swagger注释
2. **详细描述**：为每个参数提供清晰的说明和示例值
3. **错误码规范**：统一错误码和返回格式
4. **版本控制**：在文档中明确标注API版本

## API接口命名规范

遵循RESTful API设计规范：

| HTTP方法 | 用途               | 示例                       |
|---------|-------------------|----------------------------|
| GET     | 获取资源           | GET /users                 |
| POST    | 创建资源           | POST /users                |
| PUT     | 更新资源（全部字段）| PUT /users/1               |
| PATCH   | 更新资源（部分字段）| PATCH /users/1             |
| DELETE  | 删除资源           | DELETE /users/1            |

## 接口测试指南

每个接口必须提供测试用例，包括：

1. 正常情况下的请求与响应
2. 异常情况的处理
3. 边界条件测试

可使用Postman工具创建集合并分享给团队成员。

## 接口变更管理

1. 重大变更需要在ADR文档中记录
2. 保持向前兼容，不要轻易删除现有字段
3. 使用版本控制来管理接口变化

## 安全性考虑

1. 所有敏感操作必须进行认证和鉴权
2. 使用HTTPS保护数据传输
3. 实施请求限流和防止滥用措施
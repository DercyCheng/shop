# Shop 用户服务重构文档

## 架构概述

我们对 Shop 电商平台的用户服务进行了重构，采用了领域驱动设计(DDD)的分层架构、面向接口编程以及依赖注入设计模式，使代码更加模块化、可测试和可维护。

## 目录结构

```
user_srv/
├── application/            # 应用层：协调领域对象完成业务用例
│   └── service/            # 服务实现
├── domain/                 # 领域层：业务核心概念和规则
│   ├── repository/         # 存储库接口
│   └── service/            # 服务接口
├── infrastructure/         # 基础设施层：技术细节实现
│   └── persistence/        # 持久化实现
├── interfaces/             # 接口层：与外部系统交互
│   └── grpc/               # gRPC处理器
├── wire/                   # 依赖注入配置
│   ├── wire.go             # Wire集合定义
│   ├── inject.go           # 注入点定义
│   └── wire_gen.go         # 生成的注入实现
└── ...                     # 其他原有目录和文件
```

## 面向接口编程

我们采用了面向接口编程的设计理念，主要体现在以下几点：

1. **接口定义与实现分离**：
   - 在 `domain/repository` 目录中定义 `UserRepository` 接口
   - 在 `infrastructure/persistence` 目录中实现 `UserRepositoryImpl` 类
   - 在 `domain/service` 目录中定义 `UserService` 接口
   - 在 `application/service` 目录中实现 `UserServiceImpl` 类

2. **依赖倒置原则**：
   - 上层模块依赖抽象接口，而不是具体实现
   - 例如 `UserServiceImpl` 依赖 `UserRepository` 接口，而非具体的实现类

3. **松耦合设计**：
   - 各模块间通过接口交互，降低了耦合度
   - 可以轻松替换实现，例如切换数据库访问方式，只需提供不同的 `UserRepository` 实现

## 面向接口编程的优势

1. **可测试性**：
   - 可以轻松创建模拟（Mock）对象进行单元测试
   - 例如，测试 `UserServiceImpl` 时可以使用模拟的 `UserRepository` 实现

2. **可扩展性**：
   - 可以轻松添加新功能或替换现有实现
   - 例如，可以增加 `CachedUserRepository` 实现，添加缓存功能

3. **并行开发**：
   - 团队成员可以并行开发不同层级的代码
   - 只要遵循接口规范，不同模块的开发人员可以独立工作

4. **关注点分离**：
   - 业务逻辑与技术细节分离
   - 例如，`UserService` 专注于业务规则，而不关心数据如何持久化

## 依赖注入与 Wire

我们使用 Google 的 Wire 框架来实现依赖注入，主要优势包括：

1. **编译时依赖注入**：
   - Wire 在编译时生成依赖注入代码，避免了运行时反射带来的性能开销
   - 依赖关系错误在编译时就能被发现

2. **自动解析依赖图**：
   - Wire 自动解析构造函数之间的依赖关系
   - 声明性地定义组件之间的依赖关系，无需手动管理

3. **模块化配置**：
   - 使用 Wire 集合（Set）将相关组件分组
   - 例如，`RepositorySet`、`ServiceSet` 和 `HandlerSet`

## 依赖注入实现

在 `wire/wire.go` 文件中，我们定义了三个主要的依赖集合：

1. **RepositorySet**：
   ```go
   var RepositorySet = wire.NewSet(
       persistence.NewUserRepository,
       wire.Bind(new(repository.UserRepository), new(*persistence.UserRepositoryImpl)),
   )
   ```

2. **ServiceSet**：
   ```go
   var ServiceSet = wire.NewSet(
       service.NewUserService,
       wire.Bind(new(domainService.UserService), new(*service.UserServiceImpl)),
   )
   ```

3. **HandlerSet**：
   ```go
   var HandlerSet = wire.NewSet(
       grpc.NewUserHandler,
   )
   ```

在 `wire/inject.go` 文件中，我们定义了注入点：

```go
func ProvideUserHandler(db *gorm.DB) *grpc.UserHandler {
    wire.Build(
        RepositorySet,
        ServiceSet,
        HandlerSet,
    )
    return nil
}
```

Wire 生成的 `wire_gen.go` 文件中包含了实际的依赖注入实现代码：

```go
func ProvideUserHandler(db *gorm.DB) *grpc.UserHandler {
    userRepositoryImpl := persistence.NewUserRepository(db)
    userServiceImpl := service.NewUserService(userRepositoryImpl)
    userHandler := grpc.NewUserHandler(userServiceImpl)
    return userHandler
}
```

## 使用方式

在 `main.go` 中，我们只需调用一行代码即可获取完整的依赖注入链：

```go
userHandler := wire.ProvideUserHandler(global.DB)
proto.RegisterUserServer(server, userHandler)
```

这样，所有的依赖关系都由 Wire 自动处理，无需手动管理对象的创建和注入。

## 单元测试示例

基于面向接口编程和依赖注入模式，我们可以轻松编写单元测试：

```go
func TestUserService_GetUserById(t *testing.T) {
    // 创建UserRepository的模拟实现
    mockRepo := new(MockUserRepository)
    
    // 设置模拟行为
    mockUser := &model.User{ID: 1, NickName: "测试用户"}
    mockRepo.On("GetUserById", mock.Anything, int32(1)).Return(mockUser, nil)
    
    // 创建被测试的服务
    userService := service.NewUserService(mockRepo)
    
    // 执行测试
    resp, err := userService.GetUserById(context.Background(), &proto.IdRequest{Id: 1})
    
    // 验证结果
    assert.NoError(t, err)
    assert.Equal(t, "测试用户", resp.NickName)
    mockRepo.AssertExpectations(t)
}
```

## 后续扩展

1. **添加更多仓储实现**：
   - 可以实现 Redis 缓存版本的用户仓储
   - 可以实现只读版本的用户仓储用于查询优化

2. **添加横切关注点**：
   - 可以添加日志、缓存、指标收集等横切关注点
   - 利用修饰器模式包装现有接口实现

3. **微服务间通信**：
   - 可以定义服务客户端接口，隐藏 gRPC 调用细节
   - 提供统一的错误处理和重试逻辑
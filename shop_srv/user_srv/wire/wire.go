package wire

import (
	"github.com/google/wire"
	"gorm.io/gorm"
	"nd/user_srv/application/service"
	"nd/user_srv/domain/repository"
	domainService "nd/user_srv/domain/service"
	"nd/user_srv/infrastructure/persistence"
	"nd/user_srv/interfaces/grpc"
)

// RepositorySet 存储库依赖项集合
var RepositorySet = wire.NewSet(
	persistence.NewUserRepository,
	wire.Bind(new(repository.UserRepository), new(*persistence.UserRepositoryImpl)),
)

// ServiceSet 服务依赖项集合
var ServiceSet = wire.NewSet(
	service.NewUserService,
	wire.Bind(new(domainService.UserService), new(*service.UserServiceImpl)),
)

// HandlerSet 处理器依赖项集合
var HandlerSet = wire.NewSet(
	grpc.NewUserHandler,
)

// ProvideUserHandler 提供用户处理器
func ProvideUserHandler(db *gorm.DB) *grpc.UserHandler {
	wire.Build(
		RepositorySet,
		ServiceSet,
		HandlerSet,
	)
	return nil
}
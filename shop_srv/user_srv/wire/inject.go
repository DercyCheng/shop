//go:build wireinject
// +build wireinject

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

// ProvideUserHandler 提供用户处理器，这是一个用于Wire的注入函数
func ProvideUserHandler(db *gorm.DB) *grpc.UserHandler {
	wire.Build(
		RepositorySet,
		ServiceSet,
		HandlerSet,
	)
	return nil
}
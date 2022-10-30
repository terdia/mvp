package userservice

import (
	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/internal/repository"
	"github.com/terdia/mvp/internal/service/auth"
	"github.com/terdia/mvp/pkg/dto"
)

type UserService interface {
	Create(request dto.CreateUserRequest) (*data.User, ValidationErrors, error)
	GetPermissions(userID int64) (data.Permissions, error)
	GetUserByToken(tokenPlainText, scope string) (*data.User, error)
	CreateAuthenticationToken(request dto.AuthTokenRequest, scope string) (*data.Token, ValidationErrors, error)
}

type (
	ValidationErrors map[string]string

	userService struct {
		repo           repository.UserRepository
		tokenService   auth.TokenService
		permissionRepo repository.PermissionRepository
	}
)

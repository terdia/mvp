package userservice

import (
	"errors"
	"time"

	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/internal/repository"
	"github.com/terdia/mvp/internal/service/auth"
	"github.com/terdia/mvp/pkg/dto"
	"github.com/terdia/mvp/pkg/validator"
)

func NewUserService(
	repo repository.UserRepository,
	tokenService auth.TokenService,
	permissionRepo repository.PermissionRepository,
) UserService {
	return &userService{
		repo:           repo,
		tokenService:   tokenService,
		permissionRepo: permissionRepo,
	}
}

func (srv *userService) Create(request dto.CreateUserRequest) (*data.User, data.ValidationErrors, error) {

	user := &data.User{
		Role:     request.Role,
		Username: request.Username,
	}

	err := user.Password.Set(request.Password)
	if err != nil {
		return nil, nil, err
	}

	//validate request
	v := validator.New()
	if user.Validate(v); !v.Valid() {
		return nil, v.Errors, nil
	}

	err = srv.repo.Insert(user)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateUsername) {
			v.AddError("username", "a user with this username already exists")
			return nil, v.Errors, nil
		}

		return nil, nil, err
	}

	if err = srv.permissionRepo.AddForUser(user.ID, user.GetRolePermissions()...); err != nil {
		return nil, nil, err
	}

	return user, nil, nil
}

func (srv *userService) CreateAuthenticationToken(
	request dto.AuthTokenRequest,
	scope string,
) (*data.Token, data.ValidationErrors, error) {

	v := validator.New()
	v.Check(len(request.Username) > 0, "username", "must not be empty")
	v.Check(len(request.Password) > 0, "password", "must not be empty")
	if !v.Valid() {
		return nil, v.Errors, nil
	}

	user, err := srv.repo.Get(request.Username)
	if err != nil {
		return nil, nil, err
	}

	matchPassword, err := user.Password.Matches(request.Password)
	if err != nil {
		return nil, nil, err
	}

	if !matchPassword {
		return nil, nil, data.ErrInvalidCredentials
	}

	token, err := srv.tokenService.CreateNew(user.ID, 24*time.Hour, scope)

	return token, nil, err
}

func (srv *userService) GetPermissions(userID int64) (data.Permissions, error) {
	return srv.permissionRepo.GetAllForUser(userID)
}

func (srv *userService) GetUserByToken(tokenPlainText, scope string) (*data.User, error) {
	return srv.repo.GetForToken(tokenPlainText, scope)
}

func (srv *userService) UpdateUser(user *data.User) error {
	return srv.repo.Update(user)
}

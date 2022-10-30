package repository

import (
	"time"

	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/pkg/dto"
)

const (
	QueryTimeout = 3 * time.Second
)

type (
	Repository interface {
		Delete(id int64) error
	}

	UserRepository interface {
		Repository
		Insert(user *data.User) error
		Get(username string) (*data.User, error)
		Update(user *data.User) error
		GetForToken(tokenPlainText, scope string) (*data.User, error)
	}

	ProductRepository interface {
		Repository
		Insert(product *data.Product) error
		Get(id int64) (*data.Product, error)
		Update(product *data.Product) error
		GetAll(request dto.ListProductRequest) ([]*data.Product, data.Metadata, error)
	}

	PermissionRepository interface {
		GetAllForUser(userID int64) (data.Permissions, error)
		AddForUser(userID int64, codes ...string) error
	}

	TokenRepository interface {
		Create(token *data.Token) error
		DeleteAllForUserByScope(scope string, userID int64) error
	}
)

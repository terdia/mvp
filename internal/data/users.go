package data

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/terdia/mvp/pkg/validator"
)

const (
	seller = "seller"
)

var AnonymousUser = &User{}

type User struct {
	ID        int64
	Role      string
	Deposit   int
	Username  string
	Password  Password
	CreatedAt time.Time
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (u *User) GetRolePermissions() []string {
	permissions := []string{"products:read"}

	if u.Role == seller {
		permissions = append(permissions, "products:write")
	}

	return permissions
}

type Password struct {
	Plaintext *string
	Hash      []byte
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Username != "", "username", "must be provided")
	v.Check(len(user.Username) <= 500, "username", "must not be more than 500 bytes long")

	if user.Password.Plaintext != nil {
		validatePasswordPlaintext(v, *user.Password.Plaintext)
	}

	if user.Deposit != 0 {
		validateDeposit(v, user.Deposit)
	}

	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}

func validatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 6, "password", "must be at least 6 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func validateDeposit(v *validator.Validator, amount int) {
	allowed := []int{5, 10, 20, 50, 100}

	v.Check(
		validator.In(amount, allowed),
		"deposit",
		"you can oly deposit only 5, 10, 20, 50 and 100 cent",
	)
}

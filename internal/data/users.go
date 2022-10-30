package data

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/terdia/mvp/pkg/validator"
)

const (
	roleSeller = "seller"
	roleBuyer  = "buyer"
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
	permissions := []string{PermissionProductsRead}

	if u.Role == roleSeller {
		permissions = append(permissions, PermissionProductsWrite)
	}

	if u.Role == roleBuyer {
		permissions = append(permissions, PermissionProductsBuy)
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

func (u *User) Validate(v *validator.Validator) {
	v.Check(u.Username != "", "username", "must be provided")
	v.Check(len(u.Username) <= 500, "username", "must not be more than 500 bytes long")

	if u.Password.Plaintext != nil {
		validatePasswordPlaintext(v, *u.Password.Plaintext)
	}

	if u.Deposit != 0 {
		ValidateDeposit(v, u.Deposit)
	}

	v.Check(
		validator.In(u.Role, []string{roleBuyer, roleSeller}),
		"role",
		"must be seller or buyer",
	)

	if u.Password.Hash == nil {
		panic("missing password hash for user")
	}
}

func validatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 6, "password", "must be at least 6 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateDeposit(v *validator.Validator, amount int) {
	allowed := []int{CoinFiveCent, CoinTenCent, CoinTwentyCent, CoinFiftyCent, CoinHundredCent}

	v.Check(
		validator.In(amount, allowed),
		"deposit",
		"you can oly deposit only 5, 10, 20, 50 and 100 cent",
	)
}

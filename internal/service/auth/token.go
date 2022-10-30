package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/internal/repository"
)

type TokenService interface {
	CreateNew(userId int64, ttl time.Duration, scope string) (*data.Token, error)
	DeleteByUserIdAndScope(userId int64, scope string) error
}

type tokenService struct {
	repo repository.TokenRepository
}

func NewTokenService(tokenRepository repository.TokenRepository) TokenService {
	return &tokenService{repo: tokenRepository}
}

func (tsrv tokenService) CreateNew(userId int64, ttl time.Duration, scope string) (*data.Token, error) {
	token, err := generateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = tsrv.repo.Create(token)

	return token, err
}

func (tsrv tokenService) DeleteByUserIdAndScope(userId int64, scope string) error {
	return tsrv.repo.DeleteAllForUserByScope(scope, userId)
}

func generateToken(userId int64, ttl time.Duration, scope string) (*data.Token, error) {

	token := &data.Token{
		UserId: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	// fill the byte slice with random bytes from your os CSPRNG.
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))

	//convert it to a slice using the [:] operator
	token.Hash = hash[:]

	return token, nil
}

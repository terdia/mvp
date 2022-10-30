package transaction

import (
	"errors"
	"testing"
	"testing/quick"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"

	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/internal/service/productservice"
	"github.com/terdia/mvp/internal/service/userservice"
	repo "github.com/terdia/mvp/mocks/repository"
	"github.com/terdia/mvp/pkg/dto"
)

func TestTransactionService_BuyProduct(t *testing.T) {

	ctrl := gomock.NewController(t)
	productRepo := repo.NewMockProductRepository(ctrl)
	userRepo := repo.NewMockUserRepository(ctrl)
	permissionRepo := repo.NewMockPermissionRepository(ctrl)

	newUserService := userservice.NewUserService(userRepo, nil, permissionRepo)
	newProductService := productservice.NewProductService(productRepo)
	tService := NewTransactionService(newUserService, newProductService)

	testCases := map[string]interface{}{
		"BuyProductSuccessful": func() bool {
			// arrange
			user := &data.User{
				ID:        1,
				Role:      "buyer",
				Deposit:   475,
				Username:  "tester",
				Password:  data.Password{Hash: []byte("password")},
				CreatedAt: time.Now(),
			}

			product := &data.Product{
				ID:              1,
				Cost:            100,
				Name:            "Lemonade",
				Seller:          data.User{ID: 2},
				CreatedAt:       time.Now(),
				AmountAvailable: 20,
			}

			userRepo.EXPECT().Update(gomock.Any()).Return(nil)
			productRepo.EXPECT().Get(gomock.Any()).Return(product, nil)
			productRepo.EXPECT().Update(gomock.Any()).Return(nil)

			expectedResponse := dto.BuyProductResponse{
				AmountSpent: 200,
				Product: struct {
					Name     string `json:"name"`
					Cost     int    `json:"cost"`
					Quantity int    `json:"quantity_purchased"`
				}{
					Name:     "Lemonade",
					Cost:     100,
					Quantity: 2,
				},
				Change: []int{100, 100, 50, 20, 5},
			}

			// act
			purchase, validationErrs, err := tService.BuyProduct(user, product, 2)

			//assert
			if validationErrs != nil {
				t.Errorf("unexpected validation errors: %+v", validationErrs)
				return false
			}

			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return false
			}

			if !cmp.Equal(expectedResponse, *purchase) {
				t.Errorf("expected: %+v; got: %+v", expectedResponse, *purchase)
				return false
			}

			return true
		},
		"BuyProductValidationErrors": func() bool {
			// arrange
			user := &data.User{Deposit: 0}

			product := &data.Product{Cost: 100, AmountAvailable: 1}

			// act
			_, validationErrs, err := tService.BuyProduct(user, product, 2)

			//assert
			if validationErrs == nil {
				t.Errorf("expected validation errors, got: %+v", validationErrs)
				return false
			}

			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return false
			}

			return true
		},
		"BuyProductDatabaseErrors": func() bool {
			// arrange
			user := &data.User{Deposit: 200}
			product := &data.Product{
				ID:              1,
				Cost:            100,
				Name:            "Lemonade",
				Seller:          data.User{ID: 2},
				CreatedAt:       time.Now(),
				AmountAvailable: 3,
			}

			productRepo.EXPECT().Get(gomock.Any()).Return(nil, errors.New("database error"))

			// act
			_, validationErrs, err := tService.BuyProduct(user, product, 2)

			//assert
			if validationErrs != nil {
				t.Errorf("unexpected validation errors: %+v", validationErrs)
				return false
			}

			if err == nil {
				t.Errorf("expected error, got: %s", err)
				return false
			}

			return true
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if err := quick.Check(tc, nil); err != nil {
				t.Errorf("%v case failed with an error: %+v", name, err)
			}
		})
	}
}

func TestTransactionService_DepositCoin(t *testing.T) {

	ctrl := gomock.NewController(t)
	productRepo := repo.NewMockProductRepository(ctrl)
	userRepo := repo.NewMockUserRepository(ctrl)
	permissionRepo := repo.NewMockPermissionRepository(ctrl)

	newUserService := userservice.NewUserService(userRepo, nil, permissionRepo)
	newProductService := productservice.NewProductService(productRepo)
	tService := NewTransactionService(newUserService, newProductService)

	testCases := map[string]interface{}{
		"DepositSuccessful": func() bool {
			// arrange
			user := &data.User{
				ID:        1,
				Role:      "buyer",
				Deposit:   0,
				Username:  "tester",
				Password:  data.Password{Hash: []byte("password")},
				CreatedAt: time.Now(),
			}

			userRepo.EXPECT().Update(gomock.Any()).Return(nil)

			// act
			validationErrs, err := tService.DepositCoin(user, 100)

			//assert
			if validationErrs != nil {
				t.Errorf("unexpected validation errors: %+v", validationErrs)
				return false
			}

			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return false
			}

			expectedBalance := 100
			if user.Deposit != expectedBalance {
				t.Errorf("expected: %+v; got: %+v", expectedBalance, user.Deposit)
				return false
			}

			return true
		},
		"DepositValidationErrors": func() bool {
			// arrange
			user := &data.User{Deposit: 0}

			// act
			validationErrs, err := tService.DepositCoin(user, 560)

			//assert
			if validationErrs == nil {
				t.Errorf("expected validation errors, got: %+v", validationErrs)
				return false
			}

			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return false
			}

			return true
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if err := quick.Check(tc, nil); err != nil {
				t.Errorf("%v case failed with an error: %+v", name, err)
			}
		})
	}

}

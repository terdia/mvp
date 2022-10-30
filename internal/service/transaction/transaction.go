package transaction

import (
	"fmt"

	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/internal/service/productservice"
	"github.com/terdia/mvp/internal/service/userservice"
	"github.com/terdia/mvp/pkg/dto"
	"github.com/terdia/mvp/pkg/validator"
)

type Service interface {
	BuyProduct(*data.User, *data.Product, int) (*dto.BuyProductResponse, data.ValidationErrors, error)
	DepositCoin(*data.User, int) (data.ValidationErrors, error)
	DepositReset(*data.User) (data.ValidationErrors, error)
}

type transactionService struct {
	userService    userservice.UserService
	productService productservice.ProductService
}

func NewTransactionService(userService userservice.UserService, productService productservice.ProductService) Service {

	return &transactionService{
		userService:    userService,
		productService: productService,
	}
}

func (t *transactionService) BuyProduct(user *data.User, product *data.Product, quantity int) (*dto.BuyProductResponse, data.ValidationErrors, error) {

	v := validator.New()

	// check quantity
	v.Check(quantity > 0, "product", "purchase quantity must be greater zero")
	v.Check(
		product.AmountAvailable >= quantity,
		"product",
		fmt.Sprintf("not enough quantity only %d remaining", product.AmountAvailable),
	)
	// check if user has enough money for this transaction
	v.Check(user.Deposit >= (product.Cost*quantity), "product", "you do not have sufficient balance")
	if !v.Valid() {
		return nil, v.Errors, nil
	}

	//reduce product quantity
	product.AmountAvailable = product.AmountAvailable - quantity
	updatedProduct, validationErrs, err := t.productService.Update(*product)
	if validationErrs != nil || err != nil {
		return nil, validationErrs, err
	}

	//spent
	cost := product.Cost * quantity
	user.Deposit = user.Deposit - cost
	if err = t.userService.UpdateUser(user); err != nil {
		return nil, nil, err
	}

	purchase := &dto.BuyProductResponse{
		AmountSpent: cost,
		Product: struct {
			Name     string `json:"name"`
			Cost     int    `json:"cost"`
			Quantity int    `json:"quantity_purchased"`
		}{
			Name:     updatedProduct.Name,
			Cost:     updatedProduct.Cost,
			Quantity: quantity,
		},
		Change: getChange(user.Deposit),
	}

	return purchase, nil, nil
}

func (t *transactionService) DepositCoin(user *data.User, deposit int) (data.ValidationErrors, error) {

	v := validator.New()
	v.Check(deposit > 0, "deposit", "must be greater than zero")
	if data.ValidateDeposit(v, deposit); !v.Valid() {
		return v.Errors, nil
	}

	user.Deposit = user.Deposit + deposit

	return nil, t.userService.UpdateUser(user)
}

func (t *transactionService) DepositReset(user *data.User) (data.ValidationErrors, error) {

	user.Deposit = 0
	v := validator.New()
	if user.Validate(v); !v.Valid() {
		return v.Errors, nil
	}
	return nil, t.userService.UpdateUser(user)
}

func getChange(balance int) (change []int) {

	for balance > 0 {
		if balance >= data.CoinHundredCent {
			change = append(change, data.CoinHundredCent)
			balance = balance - data.CoinHundredCent
			continue
		} else if balance >= data.CoinFiftyCent {
			change = append(change, data.CoinFiftyCent)
			balance = balance - data.CoinFiftyCent
			continue
		} else if balance >= data.CoinTwentyCent {
			change = append(change, data.CoinTwentyCent)
			balance = balance - data.CoinTwentyCent
			continue
		} else if balance >= data.CoinTenCent {
			change = append(change, data.CoinTenCent)
			balance = balance - data.CoinTenCent
			continue
		} else if balance >= data.CoinFiveCent {
			change = append(change, data.CoinFiveCent)
			balance = balance - data.CoinFiveCent
			continue
		} else {
			panic(fmt.Sprintf("system error balance %d is not valid", balance))
		}
	}

	return
}

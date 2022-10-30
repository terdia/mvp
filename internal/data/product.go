package data

import (
	"time"

	"github.com/terdia/mvp/pkg/validator"
)

const (
	CoinFiveCent    = 5
	CoinTenCent     = 10
	CoinTwentyCent  = 20
	CoinFiftyCent   = 50
	CoinHundredCent = 100
)

type Product struct {
	ID              int64
	Cost            int
	Name            string
	Seller          User
	CreatedAt       time.Time
	AmountAvailable int
}

func (p *Product) Validate(v *validator.Validator) {

	v.Check(p.Name != "", "name", "must be provided")
	v.Check(len(p.Name) <= 255, "name", "must not be more than 255 bytes long")

	v.Check(p.Cost > 5, "cost", "must be greater than 5")
	v.Check((p.Cost%5) == 0, "cost", "must be a multiple of 5")
	v.Check(p.AmountAvailable >= 0, "amount_available", "must be a greater than 0")
}

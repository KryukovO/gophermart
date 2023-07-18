package entities

import (
	"errors"
	"time"

	"github.com/KryukovO/gophermart/internal/utils"
)

var ErrNotEnoughFunds = errors.New("not enough funds")

const (
	BalanceOperationRefill     string = "refill"
	BalanceOperationWithdrawal string = "withdrawal"
)

// @Description User's loyalty points account balance.
type Balance struct {
	UserID    int64   `json:"-"         swaggerignore:"true"`
	Current   float64 `json:"current"   swaggerignore:"false"`
	Withdrawn float64 `json:"withdrawn" swaggerignore:"false"`
} // @name Balance

// @Description Change of the user's loyalty points account balance.
type BalanceChange struct {
	UserID      int64     `json:"-"                      swaggerignore:"true"`
	Operation   string    `json:"-"                      swaggerignore:"true"`
	Order       string    `json:"order"                  swaggerignore:"false"`
	Sum         float64   `json:"sum"                    swaggerignore:"false"`
	ProcessedAt time.Time `json:"processed_at,omitempty" swaggerignore:"false"`
} // @name BalanceChange

func (operation *BalanceChange) Validate() error {
	if ok := utils.LuhnCheck(operation.Order); !ok {
		return ErrInvalidOrderNumber
	}

	return nil
}

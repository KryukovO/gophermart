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

type Balance struct {
	UserID    int64 `json:"-"`
	Current   int   `json:"current"`
	Withdrawn int   `json:"withdrawn"`
}

type BalanceChange struct {
	UserID      int64     `json:"-"`
	Operation   string    `json:"-"`
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}

func (operation *BalanceChange) Validate() error {
	if ok := utils.LuhnCheck(operation.Order); !ok {
		return ErrInvalidOrderNumber
	}

	return nil
}

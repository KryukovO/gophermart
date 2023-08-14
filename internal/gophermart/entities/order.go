package entities

import (
	"errors"
	"time"

	"github.com/KryukovO/gophermart/internal/utils"
)

var (
	ErrInvalidOrderNumber = errors.New("invalid order number")
	ErrOrderAlreadyAdded  = errors.New("order has already been added")
	ErrOrderAddedByOther  = errors.New("order has already been added by another user")
)

const (
	OrderStatusNew        string = "NEW"
	OrderStatusProcessing string = "PROCESSING"
	OrderStatusInvalid    string = "INVALID"
	OrderStatusProcessed  string = "PROCESSED"

	AccrualStatusRegistered string = "REGISTERED"
	AccrualStatusProcessing string = "PROCESSING"
	AccrualStatusInvalid    string = "INVALID"
	AccrualStatusProcessed  string = "PROCESSED"
)

func AccrualToOrderStatus(status string) string {
	mapping := map[string]string{
		AccrualStatusRegistered: OrderStatusNew,
		AccrualStatusProcessing: OrderStatusProcessing,
		AccrualStatusInvalid:    OrderStatusInvalid,
		AccrualStatusProcessed:  AccrualStatusProcessed,
	}

	return mapping[status]
}

// @Description Order data.
type Order struct {
	UserID     int64     `json:"-"                 swaggerignore:"true"`
	Number     string    `json:"number"            swaggerignore:"false"`
	Status     string    `json:"status"            swaggerignore:"false"`
	Accrual    float64   `json:"accrual,omitempty" swaggerignore:"false"`
	UploadedAt time.Time `json:"uploaded_at"       swaggerignore:"false"`
} // @name Order

func NewOrder(number string, userID int64) *Order {
	return &Order{
		UserID: userID,
		Number: number,
		Status: OrderStatusNew,
	}
}

func (order *Order) Validate() error {
	if ok := utils.LuhnCheck(order.Number); !ok {
		return ErrInvalidOrderNumber
	}

	return nil
}

// @Description Order data from the Accrual service.
type AccrualOrder struct {
	Order   string  `json:"order"   swaggerignore:"false"`
	Status  string  `json:"status"  swaggerignore:"false"`
	Accrual float64 `json:"accrual" swaggerignore:"false"`
} // @name AccrualOrder

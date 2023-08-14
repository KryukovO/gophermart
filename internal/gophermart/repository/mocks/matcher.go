package mocks

import (
	"fmt"
	reflect "reflect"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
	"github.com/golang/mock/gomock"
)

type orderMatcher struct {
	order *entities.Order
}

func OrderMatcher(order *entities.Order) gomock.Matcher {
	return &orderMatcher{order: order}
}

func (m orderMatcher) Matches(x interface{}) bool {
	order, ok := x.(*entities.Order)
	if !ok {
		return false
	}

	return reflect.DeepEqual(m.order, order)
}

func (m orderMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.order)
}

type balanceChangeMatcher struct {
	balanceChange *entities.BalanceChange
}

func BalanceChangeMatcher(balanceChange *entities.BalanceChange) gomock.Matcher {
	return &balanceChangeMatcher{balanceChange: balanceChange}
}

func (m balanceChangeMatcher) Matches(x interface{}) bool {
	balanceChange, ok := x.(*entities.BalanceChange)
	if !ok {
		return false
	}

	return reflect.DeepEqual(m.balanceChange, balanceChange)
}

func (m balanceChangeMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.balanceChange)
}

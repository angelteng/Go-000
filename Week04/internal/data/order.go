package data

import (
	"app/internal/biz"
	"fmt"
)

var _ biz.OrderRepo = (biz.OrderRepo)(nil)

func NewOrderRepo() biz.OrderRepo {
	fmt.Println("new order repo")
	return new(orderRepo)
}

type orderRepo struct {
	Item string
}

func (or *orderRepo) SaveOrder(o *biz.Order) {
	fmt.Println("save order to db")
	// db logic
}

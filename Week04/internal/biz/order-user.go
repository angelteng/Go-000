package biz

import "fmt"

type Order struct {
	Item string
}

// PO 持久化的interface
type OrderRepo interface {
	SaveOrder(*Order)
}

// 依赖注入
func NewOrderUser(repo OrderRepo) *OrderUser {
	fmt.Println("new order biz")
	return &OrderUser{repo: repo}
}

type OrderUser struct {
	repo OrderRepo
}

// 具体实现逻辑
func (uc *OrderUser) Buy(o *Order) {
	uc.repo.SaveOrder(o)
}

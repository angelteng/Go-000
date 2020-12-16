package service

import (
	"app/api"
	"app/internal/biz"
	"context"
	"fmt"
)

type ShopService struct {
	ouc *biz.OrderUser
}

// 依赖注入do层
func NewShopService(ouc *biz.OrderUser) *ShopService {
	fmt.Println("new order service")
	return &ShopService{ouc: ouc}
}

// 关注api实现
func (svr *ShopService) CreateOrder(ctx context.Context, r *api.CreateOrderRequest) (*api.CreateOrderReply, error) {
	// dto -> do
	o := new(biz.Order)
	o.Item = r.Name

	// 编排逻辑
	svr.ouc.Buy(o)
	return &api.CreateOrderReply{Message: "ok"}, nil
}

// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package wire

import (
	"app/internal/biz"
	"app/internal/data"
	"app/internal/service"
)

// Injectors from wire.go:

func InitializeShopService() *service.ShopService {
	orderRepo := data.NewOrderRepo()
	orderUser := biz.NewOrderUser(orderRepo)
	shopService := service.NewShopService(orderUser)
	return shopService
}
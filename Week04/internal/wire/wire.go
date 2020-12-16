package wire

import (
	"app/internal/biz"
	"app/internal/data"
	"app/internal/service"

	"github.com/google/wire"
)

func InitializeShopService() *service.ShopService {
	wire.Build(service.NewShopService, biz.NewOrderUser, data.NewOrderRepo)
	return &service.ShopService{}
}

# 工程化
1. 目录https://github.com/golang-standards/project-layout/blob/master/README_zh.md 

## model -dao - service - server - api
1. model 面向数据库的结构体，参考贫血模型，定义相关方法
2. cache miss 在dao层处理
3. service 层与grpc / 路由 绑定
4. service 负责 dto -> do 转换的deep copy

## data - biz - service -api
1. service 是api的实现类， 实现了dto - do的转换、编排逻辑。
``` golang
    type ShopService struct{
        UnimplmentedShopService
        ouc *biz.OrderUser
    }
    // 依赖注入do层
    func NewShopService(ouc &biz.OrderUser) ShopService{
        return &ShopService{ouc:ouc}
    }

    // 关注api实现
    func(svr *ShopService) CreateOrder(ctx context.Context, r *CreateOrderRequest)  (*CreateOrderReply, error){
        // dto -> do
        o:= new(biz.Order)
        o.Item= r.Name

        // 编排逻辑
        svr.ouc.Buy(o)
        return &CreateOrderReply{Message:"ok"}, nil
    }
```
2. biz是业务逻辑层，定义领域对象do，及业务逻辑的具体实现
``` golang
    // 定义do 
    type Order struct{
        Item string
    }
    // PO 持久化的interface
    type OrderRepo interface{
        SaveOrder(*Order)
    }

    // 依赖注入
    func NewOrderUser(repo OrderRepo) *OrderUser{
        return &OrderUser{repo: repo}
    }
    type OrderUser struct{
        repo OrderRepo
    }
    // 具体实现逻辑
    func (uc *OrderUser) Buy(o *Order){
        uc.repo.SaveOrder(o)
    }
```
3. data 定义持久化对象，实现biz层定义po的interface， 存储具体逻辑，包括cache miss
``` golang
    // 告诉编译器实现了接口
    var _ biz.OrderRepo = (biz.OrderRepo)(nil)
    func NewOrderRepo() biz.OrderRepo{
        return new(orderRepo)
    }
    type orderRepo struct{}
    func (or *orderRepo) SaveOrder(o *biz.Order){

    }
```
4. 持久化对象，与持久层的数据结构形成对应关系: https://github.com/facebook/ent  

## Lifecycle：
1. 所有http/grpc依赖的前置资源初始化（data, biz, service)，之后再启动监听
2. 目的：
   1. 方便测试
   2. 单次初始化，多次复用
3. 管理依赖 https://github.com/google/wire
4. 参考 https://github.com/go-kratos/kratos/blob/v2/app.go 
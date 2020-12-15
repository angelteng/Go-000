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

## API规范
1. https://github.com/googleapis/googleapis https://github.com/envoyproxy/data-plane-api https://github.com/istio/api 
2. 命名
    1. package <package_name>.<version>; 
    2. RequestURL: /<package_name>.<version>.<service_name >/{method} 
3. grpc中定义需要零值的字段使用 Wrapper类型 https://github.com/protocolbuffers/protobuf/blob/master/ src/google/protobuf/wrappers.proto 
4. 错误
    1. httpcode 使用标准httpcode 而不是全使用200
    2. 当依赖的服务返回的错误时，应自行封装而不是直接往上传播
```golang
    // StatusError contains an error response from the server.
    type StatusError struct {
        // Code is the gRPC response status code and will always be populated.
        Code int `json:"code"`
        // Message is the server response message and is only populated when
        // explicitly referenced by the JSON server response.
        Message string `json:"message"`
        // Details provide more context to an error.
        Details []interface{} `json:"details"`
    }

    // ErrorInfo is a detailed error code & message from the API frontend.
    type ErrorInfo struct {
        // Reason is the typed error code. For example: "some_example".
        Reason string `json:"reason"`
        // Message is the human-readable description of the error.
        Message string `json:"message"`
    }

    // Reason returns the gRPC status for a particular error.
    // It supports wrapped errors.
    func Reason(err error) *ErrorInfo {
        if se := new(StatusError); errors.As(err, &se) {
            for _, d := range se.Details {
                if e, ok := d.(*ErrorInfo); ok {
                    return e
                }
            }
        }
        return &ErrorInfo{Reason: UnknownReason}
    }
```
5. grpc 定义部分更新使用FiledMask

配置
1. 环境变量
2. 静态配置
3. 动态配置: https://pkg.go.dev/expvar
4. 全局配置: 通过配置模版在配置中心配置
5. Functional options
```golang
    // DialOption specifies an option for dialing a Redis server.
    type DialOption struct {
    f func(*dialOptions)
    }

    // Dial connects to the Redis server at the given network and
    // address using the specified options.
    // 区分必须参数、可选参数
    func Dial(network, address string, options ...DialOption) (Conn, error) {
    do := dialOptions{
        dial: net.Dial,
    }
    for _, option := range options {
        option.f(&do)
    } // ...
    }

    // 可选参数修改配置方法
    func DialReadTimeout(d time.Duration) DialOption{
        return DialOption(func(do *dialOptions){
            do.readTimeout = d
        })
    }
```
6. 通过配置文件时：
``` golang
    // 1. 通过pb定义配置文件字段，区分可选字段、必须字段
    syntax = "proto3";
    import "google/protobuf/duration.proto";
    package config.redis.v1;
    // redis config.
    message redis {
    string network = 1;
    string address = 2;
    int32 database = 3;
    string password = 4;
    google.protobuf.Duration read_timeout = 5;
    }
    // 2. 序列化yml文件
    func ApplyYAML(s *redis.Config, yml string) error {
        js, err := yaml.YAMLToJSON([]byte(yml))
        if err != nil {
            return err
        }
        return ApplyJSON(s, string(js))
    }
    // 3. 生成option列表
    // Options apply config to options.
    func Options(c *redis.Config) []redis.Options {
        return []redis.Options{
            redis.DialDatabase(c.Database),
            redis.DialPassword(c.Password),
            redis.DialReadTimeout(c.ReadTimeout),
        }
    }

    func main() {
        // load config file from yaml.
        c := new(redis.Config)
        _ = ApplyYAML(c, loadConfig())
        r, _ := redis.Dial(c.Network, c.Address, Options(c)...)
    }

```
## 测试
1. 单元测试的基本要求：
    1. 快速
    2. 环境一致：跑完之后清理资源
    3. 任意顺序
    4. 并行 
2. 利用 go 官方提供的: Subtests  + Gomock 完成整个单元测试。
   1. /api：比较适合进行集成测试，直接测试 API，使用 API 测试框架(例如: yapi)，维护大量业务测试 case。
   2. /data：docker compose 把底层基础设施真实模拟，因此可以去掉 infra 的抽象层。
   3. /biz： 依赖  repo、rpc client，利用 gomock 模拟 interface 的实现，来进行业务单元测试。
   4. /service： 依赖 biz 的实现，构建 biz 的实现类传入，进行单元测试。
3. 基于 git branch 进行 feature 开发，本地进行 unittest，之后提交 gitlab merge request 进行 CI 的单元测试，基于 feature branch 进行构建，完成功能测试，之后合并 master，进行集成测试，上线后进行回归测试。

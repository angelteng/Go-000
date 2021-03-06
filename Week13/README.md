# 毕业总结

## 微服务项目搭建
以下从过往项目踩坑经历结合毛老师的授课内容，阐述一下一个项目从0到1的主要流程以及需要思考规划的地方。
### 1. 需求梳理、服务划分
1. 接到新需求后，先进行需求梳理，先思考这个需求是否是一个独立服务，如果是一个新的产品线或者一个快速试错的项目，很可能会拆分两、三个的服务（比如登录服务、接入服务等）。
2. 一个服务也可能需要提供给不同的目标用户，比如我们需要提供给业务层一个rpc接口，也需要给运营同学一个后台管理系统，同时我们还可能需要一个定时任务管理去做一些统计任务。
3. 因此一个服务可能分为
   1. BFF：业务编排逻辑，组合内部、外部接口数据返回给前端。
   2. Service：主要的API服务
   3. Admin：管理后台服务
   4. Job：定时任务
   5. Task
### 2. 项目分层架构
1. 然后就可以进行项目代码分层架构设计，此时可以选用开源的框架（比如Kratos），也可以根据需求规划自己的框架目录结构。
2. 以下是一个例子
```
├── api             
│   └── errors
│       ├── canaan.proto
│       ├── canaan.pb.go
│       ├── canaan_errors.pb.go
│       └── v1
│           ├── video.pb.go
│           ├── video.proto
│           ├── video.pb.go
│           └── video_http.pb.go
├── cmd // 服务模块入口
│   ├── service.go
|   ├── admin.go
├── configs 
│   └── config.yaml
├── generate.go
├── go.mod
├── go.sum
└── internal 
    ├── biz //业务逻辑层
    │   ├── camera.go
    │   ├── escalator.go
    │   └── floor.go
    ├── data //持久化逻辑层
    │   ├── camera.go
    │   ├── escalator.go
    │   └── floor.go
    ├── server //定义grpc、http服务
    │   ├── grpc.go
    │   ├── http.go
    │   └── server.go
    └── service //api实现层
        ├── camera.go
        ├── escalator.go
        └── floor.go
```
### 3. 接口文档
1. 以往我们接口可能通过公司内部wiki、swagger等生成，这时候由于更新不及时，可能会产生接口文档与实际接口不一致的情况，失去维护的文档对后续代码维护及迭代开发的人员十分不友好。
2. 通过api目录定义的proto文件作为接口文档，可以保证在版本一致的情况下，接口文档与实际代码肯定是一致，甚至可以在api文档约定服务的时效性等其他信息。
3. 在错误码约定，也是通过proto文件定义，要同时考虑到错误码与错误的详细信息是否能帮助开发快速识别问题，可以参考Krato设计：
```
type Error struct {
    // 错误码
    code int `json:"code"`
    // 错误消息
    msg string `json:"msg"`
    // 详细信息
    details []string `json:"details"`
}
```
### 4. 服务间通讯
1. 服务间的rpc通讯框架在不同公司可能有不同的习惯，但是需要保证的是：
   1. 不同语言间通用性，比如我们之前的zero-ice框架就支持golang后续不得不换了grpc
   2. 支持meta data，以传递超时信息、鉴证信息等。
   3. 社区热度、后续可维护性也是我们选型的时候需要考虑的地方。
### 5. 代码开发
在实际代码开发中，我们还需要考虑：
1. 是否使用DI：依赖注入，但如果我们只是开发一个快速试错项目，或者一个生命周期很短但活动，考虑到学习成本与代码复杂性到增加，可以选择性使用。如果使用依赖注入，可以通过Wire实现模块依赖到定义。
2. ORM选型：用sql语句还是ORM框架，取决于团队规范。
3. 并发使用：使用goroutine提高并发度，但是需要关注 保证goroutine的生命周期，防止goroutine泄漏：
    1. 如果启动goroutine，要知道什么时候结束（退出应该发出信号）
    2. 控制gouroutine控制推出（chan/超时控制）
    3. 由调用者决定是否开启goroutine
    4. 用chan代替不断创建goroutine
    5. 常用sync.WaitGroup判断是否结束
4. 消息队列：常用消息队列实现削峰、异步处理，比较常用Kakfa做消息队列，在嵌入式端也可以通过Redis Stream实现简单的消息队列。
5. 缓存优化：分布式缓存中还需要考虑到缓存一致性处理；在大批量写入时可以考虑使用Pipeline进行优化。
### 6. 监控与交付
1. 监控：Opentracing、Prometheus
2. 运营数据：可以通过Canal订阅Mysql binlog写入ELK、Influxdb等数据库，以提供数据分析服务。

## 总结及思考
通过13周的课程，我们可以知道，每种语言的产生都是为了解决不同的问题，不同语言有不同的特性又有相通的地方，正如我们每个项目，虽然业务痛点不一样，但是解决痛点的套路可能是相通的。从项目分层、项目结构、代码依赖、api定义，到并发、性能优化，都可以从计算机组成的发展历史找到优化思想的蛛丝马迹。
从微服务概览到内存模型，不仅收获了语言的“为什么这样设计”，还收获了“（目前）最优设计思路”，感谢毛老师的谆谆教诲，终身学习的精神影响着我们。
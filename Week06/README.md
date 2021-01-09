# 隔离
1. 目的：限制故障的传播范围
2. 服务隔离
    1. 动静分离
        1. CPU cacheline false sharing
        2. Mysql中避免bufferpool频繁过期隔离动静表/拆分索引表
        3. CDN加速节点
    2. 读写分离 
        1. 主从
        2. Repilicaset
        3. CQRS https://zhuanlan.zhihu.com/p/115685384 
3. 轻重隔离
    1. 核心、非核心
    2. 快慢：按不同维度（重要性/业务）隔离处理
    3. 热点隔离
        1. 小表广播: 从 remotecache 提升为 localcache，app 定时更新，甚至可以让运营平台支持广播刷新 localcache。atomic.Value
        2. 主动预热: 比如直播房间页高在线情况下bypass 监控主动防御。
4. 物理隔离
    1. 线程
    2. 进程：容器化
    3. 集群：逻辑上是一个应用，物理上部署多套应用，通过 cluster 区分
    4. 机房
5. case：
    1. info、error日志隔离

# 超时控制
1. 网络超时：连接超时、写超时、读超时
    1. http超时https://colobu.com/2016/07/01/the-complete-guide-to-golang-net-http-timeouts/ 
2. 服务控制：
    1. 服务提供者定义好 latency SLO，更新到 gRPC Proto 定义中 
    2. kit 基础库兜底默认超时，比如 100ms，进行配置防御保护 
    3. 配置中心设置默认超时配置
3. 超时传递：把当前服务的剩余时间传递到下游服务中，继承超时策略，控制请求级别的全局超时控制。
    1. 进程内超时控制：每个请求在每个阶段(网络请求)开始前，就要检查是否还有足够的剩余来处理请求，以及继承他的超时策略（min（配置超时事件，上一阶段结束后剩下的超时事件）），使用 Go 标准库的 context.WithTimeout。
    2. 服务间的超时控制：
        1. grpc天然支持 超时传递、级联取消
        2. grpc metadata exchange，基于http2 传递grpc-timeout。
4. 监控：
    1. 超时意味服务线程耗尽、服务崩溃
    2. 看95th 99th耗时分布统计
    3. 设置合理的超时，拒绝超长请求，Server不可用时主动失败
5. case：
    1. nginx proxy timeout
    2. db连接池超时

# 过载保护（自己保护自己的服务）
1. 单机版限流（被动限流，不能快速适应）
    1. 令牌桶算法 https://pkg.go.dev/golang.org/x/time/rate
    2. 漏桶算法 go.uber.org/ratelimit 
2. 目的：计算系统临近过载时的峰值吞吐作为限流的阈值来进行流量控制。
    1. 利特尔法则
    2. 参考算法：tcp bbr算法、vegas、CoDel
    3. CPU（内存增长时->触犯GC->CPU增加）作为信号量进行节流
3. 计算峰值吞吐
    1. 独立线程采样CPU，每个250ms触发一次，计算均值时使用滑动均值去除峰值影响Moving average
    2. infight：当前服务正在请求的数量 atmoic.add
    3. Pass&RT,pass为每100ms采样窗口成功的数量，rt为单个采样窗口平均响应时间
    4. 当CPU>80% 作为启发阈值，持续1s～2s冷却时间，冷却时间结束后如果CPU仍然>80%， pass*rt < inflight 放行 , pass*rt > infight丢弃
4. 滑动窗口https://github.com/go-kratos/kratos/tree/master/okg/stat/metric
5. 限流算法 GitHub.com/go-kratos/kratos/pkg/ratelimit/bbr

# 限流
1. 定义：在一定时间内，定义某个客户/应用可以接收或处理多少个请求的技术。
2. QPS限流：不同请求可能需要cpu核心等资源不同，静态QPS限流不准
3. 限制异常用户
4. 按优先级丢弃
5. 拒绝请求也需要成本
   
## 分布式限流
1. 每次心跳后，异步批量获取quota，减少请求redis等频次，获取完以后本地消费，基于令牌桶拦截。
2. 每次请求的配额：基于历史窗口数据预测 （y=ax+b)
3. 假设有N个节点，如何分配资源：最大最小公平分享算法 Max-Min Faireness
4. 步骤：
   1. 节点A，节点B每5s 发送心跳到令牌服务，带上自己根据历史值所计算出来的配额。Token Server持有配置的全局QPS配额。
   2. Token Server根据Max-Min Faireness分配令牌。
   3. A、B获取令牌后，根据令牌桶算法进行拦截。
5. 如何配置阈值：按照重要级:最重要、重要、可丢弃（批量任务会重试）、可丢弃（可偶尔不可用）
6. 客户端侧限流：熔断
    1. 计算需要熔断的概率： max(0, (requests - K*accepts) / (requests + 1))
    2. 计算是否进入熔断（Google SRE）
    3. 客户端需要限制请求频次应从接口返回解析出
7. 双熔断（Gutter）：基于熔断的 gutter kafka ，用于接管自动修复系统运行过程中的负载，这样只需要付出10%的资源就能解决部分系统可用性问题。如果 gutter 也接受不住的流量，重新回抛到主集群，最大力度来接受。 
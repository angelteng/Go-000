# 隔离
1. 服务隔离
    1. 动静分离
    2. 读写分离
2. 轻重隔离
    1. 核心
    2. 快慢
    3. 热点
        1. 小表广播: 从 remotecache 提升为 localcache，app 定时更新，甚至可以让运营平台支持广播刷新 localcache。atomic.Value
        2. 主动预热: 比如直播房间页高在线情况下bypass 监控主动防御。
3. 物理隔离
    1. 线程
    2. 进程：容器化
    3. 集群：逻辑上是一个应用，物理上部署多套应用，通过 cluster 区分
    4. 机房
4. case：
    1. info、error日志隔离

# 超时控制
1. 连接超时、写超时、读超时
2. 如何定义：
    1. 服务提供者定义好 latency SLO，更新到 gRPC Proto 定义中 
    2. kit 基础库兜底默认超时，比如 100ms，进行配置防御保护 
3. 超时传递：每个请求在每个阶段(网络请求)开始前，就要检查是否还有足够的剩余来处理请求，以及继承他的超时策略，使用 Go 标准库的 context.WithTimeout。
4. grpc支持 超时传递、级联取消
5. case：
    1. nginx proxy timeout
    2. db连接池超时

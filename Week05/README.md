
# 概览

￼1. 结构
    1. BFF 编排逻辑，包括调用其他rpc服务，降级
    2. comment-service 评论的api服务
    3. comment-job 削峰
    4. comment-admin 管理后台，运营侧
2. canal订阅mysql 写入es，es只做查询
3. 读：cache miss 的回源由job执行
4. 写：kafka  + sharding

# 存储设计

1. 先读后写for update ，避免死锁：先排序
2. 分离索引表、内容表
3. 用位存储属性
4. 缓存：增量加载的方式逐渐预热填充缓存 （发现redis没有再去填）。redis sortset必须保证时增量的，使用expire 而不是exisit判断key是否存在（因为可能马上过期）

# 可用性设计

1. 缓存穿透
    1. 不合适处理方式：setnx + 轮询，导致浪费线程，锁超时但是又没处理完会导致堆积请求。
    2. 单节点使用sync.siglefight（归并回源)，单机器多节点使用local cache 5s缓存减少回源请求，这里只拿一页的缓存。
    3. 其他页发送rebuild cache到kafka传递回源消息，消费进程处理回源逻辑。
    4. 可参考cdn回源处理
2. 热点：
    1. 多replica
    2. local cache
    3. 热点识别：环形数组统计key调用次数，最小堆+topk计算识别。
    4. 调用时机：接口时间与包含时间的原子值比对，大于间隔模式使用siglefight触发计算
￼


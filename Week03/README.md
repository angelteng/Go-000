# Goroutine生命周期管理
1. 如果启动goroutine，要知道什么时候结束（退出应该发出信号）
2. 控制gouroutine控制推出（chan/超时控制）
3. 由调用者决定是否开启goroutine
4. 用chan代替不断创建goroutine
5. 常用sync.WaitGroup判断是否结束
6. 内存模型：
    1. 原子性、可见性
    2. go memory model（了解happen-before）
    3. 底层的 memory reordering（可以挖一挖 cpu cacline、锁总线、mesi、memory barrier）
   
# Data race
1. 查看汇编go tool compile -S main.go
2. 查看是否有竞争行为 go build -race main.go
3. go 同步语义 mutex atomic
4. single machine word 赋值是原子的
5. data race：原子性、可见性
6. 写时复制（无锁访问共享数据）：
    1. 常用于热门信息缓存、服务将降级
    2. 写操作时复制全量老数据到新对象，更新新数据，然后atomic.value原子替换
7. Mutex锁实现：
    1. Barging：锁被释放时唤醒第一个等待者，但是锁给第一个等待的人或者第一个请求锁的人
    2. Handsoff：锁释放的时候会持有知道第一个等待者准备好活去锁(解锁锁饥饿的问题）
    3. Spinning：不需要park goroutine，自己在循环等待锁，自旋在等待队列为空||饮用程序重度使用锁
    4. go 1.8 Barging+Spining , go 1.9 Handsoff
8. errgroup 包 
    1. 并行工作流的错误处理||优雅降级，context传播取消，利用局部变量+闭包
    2. 野生goroutine
    3. context不能用于其他函数
9. sync.Pool 保存复用临时对象（可被任意时候回收），比如每次请求的trace对象，比如解析http包的对象
10. Context
    1. 显式传递context，不要包在结构体里，除非用于chan传递
    2. Context.Value是并发安全的，只读，新增的时候只会新增context对象
    3. 适合放请求的不能被变更的元数据，比如染色、trace
    4. 如果value是个map，当修改时，应该cow，然后创建新的context进行传递
11. context超时覆盖
12. Context.WithCancel的cancel方法一定要调用 否则会泄露
13. Chan
    1. 只有发送者可以close chan
    2. buffer可以减少延迟，不能增加吞吐，吞吐需要增加goroutine去竞争消费

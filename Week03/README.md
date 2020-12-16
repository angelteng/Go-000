# Goroutine生命周期管理
1. 控制goroutine的生命周期
    1. 由调用者决定是否开启goroutine
    2. 如果启动goroutine，要知道什么时候结束（退出应该发出信号）
    3. 控制gouroutine控制推出（chan/context）
2. 如果你的 goroutine 在从另一个 goroutine 获得结果之前无法取得进展，那么通常情况下，你自己去做这项工作比委托它( go func() )更简单。
3. 用chan+消费池代替不断创建goroutine
4. 常用sync.WaitGroup等待所有goroutine结束。也可以将 wg.Wait() 操作托管到其他 goroutine，owner goroutine 使用 context 处理超时。
```golang 
    ch := make(chan struct{})
    go func(){
        Wg.Wait()
        close(ch)
    }()
    select{
        case <-ch:
            return nil
        case <- ctx.Done()
            return errors.New("timeout")
    }
```
5. log.Fatal会调用os.Exit(),导致defer函数无法被调用。

# Go的内存模型
1. Happend Before：在一个 goroutine 中，读和写一定是按照程序中的顺序执行的。即编译器和处理器只有在不会改变这个 goroutine 的行为时才可能修改读和写的执行顺序。由于重排（CPU 重排、编译器重排），不同的goroutine 可能会看到不同的执行顺序。
2. 多个 goroutine 访问共享变量 v 时，它们必须使用同步事件来建立先行发生这一条件来保证读操作能看到需要的写操作。 
    1. 对变量v的零值初始化在内存模型中表现的与写操作相同。
    2. 对大于 single machine word 的变量的读写操作表现的像以不确定顺序对多个 single machine word的变量的操作。
    3. 参考 https://www.jianshu.com/p/5e44168f47a3
3. 因此，多个goroutine要确保没有data race 需要保证：原子性（以singel machine word操作或者利用同步语义）、可见性（消除内存屏障）
4. 关于底层的 memory reordering，可以挖一挖 cpu cacline、锁总线、mesi、memory barrier。
5. 编译器重排、内存重排，(重排是指程序在实际运行时对内存的访问顺序和代码编写时的顺序不一致)，目的都是为了减少程序指令数，最大化提高CPU利用率。
6. Memory Barrier: 现代 CPU 为了“抚平”内核、内存、硬盘之间的速度差异，搞出了各种策略，例如三级缓存等。每个线程的操作结果可能先缓存在自己内核的L1 L2cache，此时（还没刷到内存）别的内核线程是看不到的。因此，对于多线程的程序，所有的 CPU 都会提供“锁”支持，称之为 barrier，或者 fence。它要求：barrier 指令要求所有对内存的操作都必须要“扩散”到 memory 之后才能继续执行其他对 memory 的操作。 参考 https://cch123.github.io/ooo/，https://blog.csdn.net/qcrao/article/details/92759907 
   

# SYNC 包
1. Share Memory By Communicating： Go 没有显式地使用锁来协调对共享数据的访问，而是鼓励使用 chan 在 goroutine 之间传递对数据的引用。这种方法确保在给定的时间只有一个goroutine 可以访问数据。
2. 查看是否有竞争行为 go build -race main.go，查看汇编go tool compile -S main.go 
3. go 常用同步语义：mutex、RWMutex、atomic
6. 写时复制（无锁访问共享数据）：
    1. 常用于热门信息缓存、服务将降级
    2. 写操作时复制全量老数据到新对象，更新新数据，然后atomic.Value原子替换
7. Mutex锁实现：
    1. Barging：锁被释放时唤醒第一个等待者，但是锁给第一个等待的人或者第一个请求锁的人
    2. Handsoff：锁释放的时候会持有知道第一个等待者准备好活去锁(解锁锁饥饿的问题）
    3. Spinning：不需要park goroutine，自己在循环等待锁，自旋在等待队列为空||饮用程序重度使用锁
    4. go 1.8 Barging+Spining , go 1.9 Handsoff
8. errgroup 包 
    1. 并行工作流的错误处理
    2. 优雅降级（当错误发生时，在该goroutine内部进行降级逻辑处理）
    3. context传播取消，当有一个错误发生时，其内部当context会cancel。父context取消也一样会传播。
    2. 需要解决野生goroutine的问题。
    3. 当Wait函数返回后，其context已经被取消了，之后context不能用于其他函数。
9. sync.Pool 保存复用临时对象（可被任意时候回收），比如每次请求的trace对象，比如解析http包的对象

# Context
1. 显式传递context，不要包在结构体里，除非用于chan传递
2. Context.Value是并发安全的，只读，新增的时候只会新增context对象
3. 适合放请求的不能被变更的元数据，比如染色、trace
4. 如果value是个map，当修改时，应该cow，然后创建新的context进行传递
5. context超时覆盖
```golang
    func shrinkDeadline(ctx context.Context, timeout time.duration) time.Time{
        timeoutTime := time.Now().Add(timeout)
        if deadline,ok:= ctx.Deadline(); ok && timeoutTime.After(deadline){
            return deadline
        }
        return timeoutTime
    }
```
6. 重要：Context.WithCancel的cancel方法一定要调用，否则会造成泄露

# Chan
1. 只有发送者可以close chan
2. buffer可以减少延迟，不能增加吞吐，吞吐需要增加goroutine去竞争消费

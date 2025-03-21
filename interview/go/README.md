- [Golang面试题](#golang面试题)
- [1 基础问题](#1-基础问题)
- [2 context](#2-context)
- [3 channel相关](#3-channel相关)
- [4 map相关](#4-map相关)
- [5 GMP相关](#5-gmp相关)
- [6 锁相关](#6-锁相关)
- [7 并发相关](#7-并发相关)
- [8 GC相关](#8-gc相关)
- [9 内存相关](#9-内存相关)
- [10 其它](#10-其它)

# Golang面试题


# 1 基础问题
**1、golang 中 make 和 new 的区别？**
- make和new都是内存的分配（堆上）， make只用于slice、map以及channel的初始化(非零值);
- new用于为类型分配一块零值化的内存，并返回指向这块内存的指针。

**2、数组和切片的区别**
- 数组指长度固定的数据结构，数组的初始化和切片不一样，数组是值传递
- 切片的长度不固定，可变化；切片是地址传递

**3、for range 的时候它的地址会发生变化么？**
- 不会

**4、go defer，多个 defer 的顺序，defer 在什么时机会修改返回值？**
- defer是压栈的形式，多个defer按顺序入栈，先进后出
- defer会在程序发生panic的时候修改返回值

**5、 uint 类型溢出**
- 超过255就是溢出了

**6、介绍 rune 类型**
- rune其实就是int32

**7、 golang 中解析 tag 是怎么实现的？反射原理是什么？**
- 通过反射机制
- 先获取reflect.type，通过type Field获取tag

**8、调用函数传入结构体时，应该传值还是指针？**
- Golang 都是传值

**9、切片扩容机制**
> 内存分配：扩容时会分配新数组，并将原数据复制到新数组，原数组由垃圾回收器处理。
> 
> 性能影响：频繁扩容可能导致性能下降，可通过 make 函数预设足够容量来优化。

1. go 1.18之前
- 小切片：容量小于 1024 时，每次扩容为原来的 2 倍。
- 大切片：容量达到或超过 1024 时，每次扩容增加 25%，直到满足需求。

2. go 1.18之后
- 小切片：临界值换成了256，小于256的时候，切片先两倍扩容，如果两倍扩容后的容量还是不够，就直接以切片需要的容量作为容量。
- 大切片：大于256公式变为`(oldCap+3*256)/4`这个公式的值随着oldcap的越来越大，从2一直接近1.25，相对于1.18之前可以更平滑的过渡。

# 2 context
**1、context 结构是什么样的？**
```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```
**2、context 使用场景和用途**
- 主要用于控制并发，可以有效避免goroutine泄漏（goroutine资源一直得不到释放）
- 在 Go http 包的 Server 中，每一个请求在都有一个对应的goroutine去处理。请求处理函数通常会启动额外的goroutine用来访问后端服务，比如数据库和 RPC 服务。用来处理一个请求的goroutine通常需要访问一些与请求特定的数据，比如终端用户的身份认证信息、验证相关的 token、请求的截止时间。当一个请求被取消或超时时，所有用来处理该请求的goroutine都应该迅速退出，然后系统才能释放这些goroutine占用的资源

# 3 channel相关
**channel 数据结构：**
```go
type hchan struct {
    qcount   uint   // 队列中当前的元素个数
    dataqsiz uint   // 环形队列的大小（缓冲区大小）
    buf      unsafe.Pointer // 环形队列的指针（存储数据）
    sendx    uint   // 发送索引（用于环形队列）
    recvx    uint   // 接收索引（用于环形队列）
    recvq    waitq  // 等待接收的 Goroutine 队列
    sendq    waitq  // 等待发送的 Goroutine 队列
    lock     mutex  // 互斥锁，保证并发安全
}
```

**1、channel 是否线程安全？锁用在什么地方？**
- 发送一个数据到channel和从channel接收一个数据都是原子操作，而且Go的设计思想就是:`不要通过共享内存来通信（加锁），而是通过通信来共享内存（channel）`，前者就是传统的加锁，后者就是Channel，也就是说，设计Channel的主要目的就是在多任务间传递数据的，这当然是安全的。
- 锁用在数据入队列和出队列

**2、go channel 的底层实现原理（数据结构）**
- 是一个环形队列。优点如下
- 大小固定的，高效的内存利用率，在内存中是一片连续的内存空间，无需动态分配和释放内存
- 时间复杂度为O（1）的入队和出队
- 支持并发场景的锁优化，使用原子操作atomic是最轻量级的锁

**3、nil、关闭的 channel、有数据的 channel，再进行读、写、关闭会怎么样？**
- nil channel读写数据会死锁，永久阻塞
- 关闭的channel若有数据可以正常读取，不能写入
- 有数据的channel若空间占满，则后续写入数据会被阻塞，等待数据被读取后释放空间

**4、向 channel 发送数据和从 channel 读数据的流程是什么样的？**
- 分为有缓冲和无缓冲通道看待

**5、channel为什么是go中唯一满足并发安全的基础数据类型？**
1. 基于原子锁原生支持并发安全
2. 遵循CSP并发模型，强调通过消息传递而非共享内存来实现并发。
3. 自动阻塞机制，当写入数据到已满的缓冲 channel 时，写操作会阻塞，直到有空间

**6、channel 分配在栈上还是堆上？哪些对象分配在堆上，哪些对象分配在栈上？**
- 根据channel中数据类型而定

**7、channel 调度策略**
- FIFO 调度：发送 Goroutine 和接收 Goroutine 按顺序排队（先进先出）。
- 无锁（读写分离）：有缓冲 channel 读写不同 buf 索引，减少锁冲突。
- 超时机制：select 语句支持 timeout，防止 Goroutine 永远阻塞。

# 4 map相关
**1、map 使用注意的点，并发安全？**
- map的并发读写是不安全的，map属于引用类型，并发读写时多个协程之间是通过指针访问同一个地址，即访问共享变量，此时同时读写资源存在竞争关系。利用读写锁可实现对map的安全访问。

**2、map 循环是有序的还是无序的？为什么是无序的？**
- 无序的；底层是hash表，随即获取数据

**3、 map 中删除一个 key，它的内存会释放么？**
- 如果删除的是值类型的，那么不会释放
- 如果是引用类型的，那么会释放被删除元素的内存空间

**4、怎么处理对 map 进行并发访问？有没有其他方案？ 区别是什么？**
- 互斥锁，同步锁
- 启动单个协程对map进行读写，当其他协程需要读写map时，通过channel向这个协程发送信号即可。

**5、 nil map 和空 map 有何不同？**
- nil map 还没初始化，不能赋值（new map）
- 空map是已经初始化了内存空间，可以赋值，只不过默认值为对应map数据类型的默认空值（make map）

**6、map 的数据结构是什么？是怎么实现扩容？**

> map的数据结构是hash表
> 
> map使用的桶就是bmap，一个桶里可以放8个键值对，但是为了让内存排列更加紧凑，8个key放一起，8个value放一起，8个key的前面则是8个tophash，每个tophash都是对应哈希值的高8位。

```go
type hmap struct {
    count     int          // 当前map中元素的数量
    flags     uint8        // 状态标志（例如，map是否已经扩容）
    B         uint8        // 当前map的桶数量的指数（桶的数量是2的B次方）
    hash0     uint32       // 随机哈希种子
    buckets   unsafe.Pointer // 指向桶的指针
    oldbuckets unsafe.Pointer // 指向旧桶的指针（用于扩容时）
    nelems    int32        // 存储键值对的数量
    overflow  *bucket      // 溢出桶，用于存放哈希冲突较多的情况
}
```

Go map 的扩容主要有两种情况触发：
- 装载因子超过阈值(LoadFactor > 6.5) 装载因子计算公式为：LoadFactor = count / 2^B
其中 count 是 map 中的元素个数，2^B 是当前桶数量
- overflow buckets 过多，当 overflow buckets 数量超过 `2^B` 时，即使装载因子未达到阈值也会触发扩容

扩容方式分两种：

1. 翻倍扩容(grow)

- 当装载因子过大时使用
- 新建一个是原来 2 倍大小的 hash table
- 将原数据迁移到新表
- 桶的数量翻倍，B + 1


2. 等量扩容(sameSizeGrow)
- 当 overflow buckets 过多时使用
- 不改变大小，仅整理碎片
- 重新做一遍 hash，使得元素排列更紧密
- 桶数量不变，保持 B 不变

扩容过程的特点：

- 采用渐进式扩容，每次操作时只迁移 1-2 个桶
- 扩容期间，新旧 bucket 同时存在
- 查找时会同时查看两个 bucket
- 一个 bucket 迁移完成后才删除旧的
- 所有 bucket 都迁移完成后完成扩容

这种机制保证了 map 操作的性能平稳，避免了因一次性扩容导致的性能抖动。

**7、map冲突如何解决**
- 使用拉链法，每个哈希桶（hash bucket）在初始化时会有一个链表或者其他形式的数据结构来存储具有相同哈希值的键值对。当发生哈希冲突时，新的键值对会被添加到对应哈希桶的链表（或其他数据结构）的末尾，而不是直接覆盖原有的键值对。这样，即使发生哈希冲突，也能够保证所有的键值对都能够被正确地存储和检索。


# 5 GMP相关
**1、什么是 GMP？**
- M指的是Machine，一个M直接关联了一个内核线程，用以执行G，维护自己的本地内存以减少执行期间的争用。
- P指的是"processor"，代表了M所需的上下文环境，也是处理用户级代码逻辑的处理器，可以控制并发。
- G指的是Goroutine（协程），Go 运行时管理的轻量级用户空间线程，堆栈从小处开始（例如 2 KB），然后动态增大/缩小。与 OS 线程相比，创建和切换成本更低

**2、GMP如何工作？**
1. 创建G并将其分配给P。
2. P从其本地队列中选择G，并将它们分配给M （OS线程）。
3. M执行G，直到它阻塞、完成或被抢占。
4. P个实例的数量由GOMAXPROCS控制，它定义了并发运行的程序的最大数量。

**3、M、P、G 的数量问题？**
- G 的数量理论上可以 无限多（由 Go 运行时管理）。
- P 的数量由 GOMAXPROCS 决定，最多等于 CPU 核心数（默认 runtime.NumCPU()）。
- M 的数量不固定，Go 运行时会根据 G 的调度动态创建/回收 M（通常 M 的数量 不会超过 P 数量太多）。

**4、进程、线程、协程有什么区别？**
- 进程是资源（CPU、内存等）分配的基本单位，在程序运行时创建，它是程序执行时的一个实例
- 线程则是程序执行的最小单位，是进程的一个执行流，一个进程由多个线程组成
- 协程则是一种轻量级的线程

扩展
- **进程** 比较重量，占据独立的内存，所以上下文进程间的切换开销（栈、寄存器、虚拟内存、文件句柄等）比较大，但相对比较稳定安全。
- **线程** 间通信主要通过共享内存，上下文切换很快，资源开销较少，但相比进程不够稳定，容易丢失数据。线程的调度是抢占式的。线程的栈大小默认为2MB。
- **协程** 拥有自己的寄存器上下文和栈。协程调度切换时，将寄存器上下文和栈保存到其他地方，在切回来的时候，恢复先前保存的寄存器上下文和栈，直接操作栈则基本没有内核切换的开销，可以不加锁的访问全局变量，所以上下文的切换非常快。协程的调度是协作式的，协程执行完任务会主动让给其他协程。协程的栈大小默认为2KB。

**4、什么是抢占式调度？是如何抢占的？**
- 抢占式调度主要用于防止某些 Goroutine 长时间占用 CPU，而不让出执行权，影响并发性。
- Go 采用了一种 混合调度策略，结合了 协作式调度（Cooperative Scheduling） 和 抢占式调度（Preemptive Scheduling）
- 抢占式调度的触发时机：系统调用时（Syscall Preemption）、运行时间太长（Safe Point Preemption）
- Go 不会 在任意指令执行时抢占 Goroutine，而是在 函数调用 处插入检查点。这样可以确保 Goroutine 在安全的状态下被调度。
- Go 运行时会向长时间运行的 Goroutine 发送 `SIGURG` 信号来触发调度。

**5、 Goroutine进行一个读写，然后阻塞了它，逻辑处理器P和操作系统线程M会发生什么样的一个对应的一个操作**
- 当一个Goroutine阻塞时，逻辑处理器P会寻找其他可运行的Goroutine来填充其空闲时间，而操作系统线程M仍然会保持活动状态以继续执行其他Goroutine。这是Go并发模型的一个关键特点，可以有效地管理大量Goroutine，确保程序在并发执行中高效运行。

**6、进程间通信IPC**
- 管道：速度慢，容量有限，只有父子进程能通讯
- FIFO：任何进程间都能通讯，但速度慢
- 消息队列：容量受到系统限制，且要注意第一次读的时候，要考虑上一次没有读完数据的问题
- 信号量：不能传递复杂消息，只能用来同步
- 共享内存区：能够很容易控制容量，速度快，但要保持同步，比如一个进程在写的时候，另一个进程要注意读写的问题，相当于线程中的线程安全，当然，共享内存区同样可以用作线程间通讯，不过没这个必要，线程间本来就已经共享了同一进程内的一块内存

# 6 锁相关
加锁的目的就是保证共享资源在任意时间里，只有一个线程访问，这样就可以避免多线程导致共享数据错乱的问题

**1、除了 mutex 以外还有那些方式安全读写共享变量？**
- channel

**2、什么是原子操作？Go 如何实现原子操作？**
- 原子操作即是进行过程中不能被中断的操作，针对某个值的原子操作在被进行的过程中，CPU绝不会再去进行其他的针对该值的操作。为了实现这样的严谨性，原子操作仅会由一个独立的CPU指令代表和完成。原子操作是无锁的，常常直接通过CPU指令直接实现。 事实上，其它同步技术的实现常常依赖于原子操作。golang原子操作可确保这些goroutine之间不存在数据竞争
- Go 语言的sync/atomic包提供了对原子操作的支持，用于同步访问整数和指针。
- Go语言提供的原子操作都是非入侵式的。这些函数提供的原子操作共有五种：增减、比较并交换、载入、存储、交换。原子操作支持的类型类型包括int32、int64、uint32、uint64、uintptr、unsafe.Pointer。

**3、Mutex 是悲观锁还是乐观锁？悲观锁、乐观锁是什么？**
- 悲观锁：顾名思义，就是很悲观，每次去拿数据的时候都认为别人会修改，所以每次在拿数据的时候都会上锁，这样别人想拿这个数据就会block直到它拿到锁；互斥锁、自旋锁和读写锁都是悲观锁。
- 乐观锁：顾名思义，就是很乐观，每次去拿数据的时候都认为别人不会修改，所以不会上锁，但是在更新的时候会判断一下在此期间别人有没有去更新这个数据，可以使用版本号等机制。sync/atomic是乐观锁。

**4、Mutex 有几种模式？**
- 互斥锁：就是互相排斥
- 自旋锁：自旋锁是指当一个线程在获取锁的时候，如果锁已经被其他线程获取，那么该线程将循环等待，然后不断地判断是否能够被成功获取，直到获取到锁才会退出循环
- 读写锁：实际是一种特殊的自旋锁，它把对共享资源的访问者划分成读者和写者，读者只对共享资源进行读访问，写者则需要对共享资源进行写操作
- Map锁
无论是互斥锁，还是自旋锁，在任何时刻，最多只能有一个保持者，也就说，在任何时刻最多只能有一个执行单元获得锁。

**5、goroutine 的自旋占用资源如何解决**
- 自旋锁是指，当一个 Goroutine 尝试获取锁时，如果该锁被其他 Goroutine 持有，它不会进入阻塞状态，而是会不断地检查锁的状态，直到锁被释放为止。

# 7 并发相关
**1、怎么控制并发数？**
- 使用channel阻塞
- 使用计数器，sync.WaitGroup

**2、怎么控制并发使用CPU核数**
- runtime.GOMAXPROCS(runtime.NumCPU())

**3、多个 goroutine 对同一个 map 写会 panic，异常是否可以用 defer 捕获？**
- 不可以吧

**4、如何优雅的实现一个 goroutine 池**
```go
package workerpool

import (
    "context"
    "errors"
    "sync"
)

// Task represents a unit of work
type Task struct {
    handler func() error    // the actual work to be done
    done    chan error     // channel to signal completion
}

// Pool represents a goroutine pool
type Pool struct {
    workers    int           // number of workers
    taskQueue  chan Task     // channel for tasks
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup
}

// NewPool creates a new worker pool with specified number of workers
func NewPool(workers int) *Pool {
    if workers <= 0 {
        workers = 1
    }
    
    ctx, cancel := context.WithCancel(context.Background())
    
    p := &Pool{
        workers:    workers,
        taskQueue:  make(chan Task, workers*2), // buffer size = 2 * workers
        ctx:        ctx,
        cancel:     cancel,
    }
    
    p.start()
    return p
}

// start launches the workers
func (p *Pool) start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go func() {
            defer p.wg.Done()
            for {
                select {
                case task, ok := <-p.taskQueue:
                    if !ok {
                        return
                    }
                    err := task.handler()
                    task.done <- err
                    close(task.done)
                case <-p.ctx.Done():
                    return
                }
            }
        }()
    }
}

// Submit submits a task to the pool
func (p *Pool) Submit(handler func() error) error {
    select {
    case <-p.ctx.Done():
        return errors.New("pool is closed")
    default:
        done := make(chan error, 1)
        task := Task{
            handler: handler,
            done:    done,
        }
        
        select {
        case p.taskQueue <- task:
            return <-done
        case <-p.ctx.Done():
            return errors.New("pool is closed")
        }
    }
}

// Close gracefully shuts down the pool
func (p *Pool) Close() {
    p.cancel()
    close(p.taskQueue)
    p.wg.Wait()
}
```

调用
```go
func main() {
    // 创建一个包含 3 个 worker 的池
    pool := NewPool(3)
    defer pool.Close()
    
    // 提交任务
    err := pool.Submit(func() error {
        // 执行具体工作
        time.Sleep(time.Second)
        return nil
    })
    
    if err != nil {
        log.Printf("task error: %v", err)
    }
}
```

**5、Golang中常见的并发模型**
- 通过channel通知实现并发控制
- 通过sync包中的WaitGroup实现并发控制
- 互斥锁

**6、获取goroutine数量**
- `runtime.NumGoroutine()`

# 8 GC相关
**1、go gc 是怎么实现的？**
- [GC](./gc.md)

**2、GC 中 stw 时机，各个阶段是如何解决的？**
- 标记准备阶段（使用懒清扫算法，会为每个P启动一个标记协程，但不是所有标记协程都有执行的机会，开启协程的数量为0.25P）
- 标记终止阶段

**3、GC 的触发时机？**
- 一般是当 Heap 上的内存达到一定数值后，会触发一次 GC，这个数值我们可以通过环境变量 GOGC 或者 debug.SetGCPercent() 设置，默认是 100，表示当内存增长 100% 执行一次 GC。如果当前堆内存使用了 10MB，那么等到它涨到 20MB 的时候就会触发 GC。
- 再就是每隔 2 分钟，如果期间内没有触发 GC，也会强制触发一次。
- 最后就是用户手动触发了，也就是调用 runtime.GC() 强制触发一次。
- 扫描过程最多使用 25% 的 CPU 进行标记，这是为了尽可能降低 GC 过程对用户的影响。而如果 GC 未完成，下一轮 GC 又触发了，系统会等待上一轮 GC 结束。

**4、GC的优化**
- 内存分配过快，导致频繁gc，如果从代码层面没有很好的解决的办法，那么可以设置GOGC的值，降低清扫频率比如经常有大量IO的程序；
- 小对象越多，对gc会造成标记清扫压力，尽量避免过多小对象的创建；
- 尽可能重用对象，减少对象的创建和销毁，可以使用对象池sync.pool；
- 避免内存泄漏，及时释放无用的内存；
- 使用gc调优工具，`export GODEBUG=gctrace=1,gcpacertrace=1`；其中，`gctrace=1`表示输出gc的执行日志，`gcpacertrace=1`表示输出内存分配器的执行日志；
- 编译器优化

# 9 内存相关
**1、谈谈内存泄露，什么情况下内存会泄露？怎么定位排查内存泄漏问题？**
- 内存泄漏（Memory Leak）是指程序中已动态分配的堆内存由于某种原因程序未释放或无法释放，造成系统内存的浪费，导致程序运行速度减慢甚至系统崩溃等严重后果。
- 代码里：使用pprof工具或者runtime.ReadMemStats()方法
- 做好code review

**2、知道 golang 的内存逃逸吗？什么情况下会发生内存逃逸？**
- 本质上就是数据在堆中还是在栈中的选择
- 多级间接赋值容易导致逃逸，如[]interface{} —> data[0]=100
- 发送指针或带有指针的值到 channel 中
- 在一个切片上存储指针或带指针的值
- slice 的背后数组被重新分配了，因为 append 时可能会超出其容量(cap)
- 在 interface 类型上调用方法。

**3、请简述 Go 是如何分配内存的？**
- Go将内存分为大小67个span，每个span所容纳的元素大小和元素个数不相同，class1 8k 1024个，class2 16k 512个…以此类推。为了方便对这些span进行管理，Go采用了三级内存管理模型，分别为mcache、mcentral、mheap。Go的内存分配采用了现代TCMalloc内存分配思想，每个逻辑处理器都存储了一个本地span，称作mcache，每个mcache包含了完整的不同类别的span，协程如果需要内存，就从mcache中取，mcache中没有就会往上一级mcentral中取，mcentral需要加锁处理，因为mcentral中包含了两个mspan链表，需要对这些链表进行循环遍历，如果mcentral中也没有，那么就会往mhead中申请，这三级是在用户态操作的，如果mheap中也没有空闲的内存空间，那么就会向操作系统中申请资源（内核态）。

**4、介绍一下大对象小对象，为什么小对象多了会造成 gc 压力？**
- 大对象大于32字节的，小对象小于32字节的

- GC 需要标记每个对象，对象数量越多，需要标记的次数越多
- 内存碎片，小对象分散在内存中
- 小对象之间可能有复杂的引用关系，追踪这些引用关系会消耗更多 CPU

**5、使用多级内存管理的优势**
1. 内存分配大多时候都是在用户态完成的，不需要频繁进入内核态。
2. 每个 P 都有独立的 span cache，多个 CPU 不会并发读写同一块内存，进而减少 CPU L1 cache 的 cacheline 出现 dirty 情况，增大 cpu cache 命中率。
3. 内存碎片的问题，Go 是自己在用户态管理的，在 OS 层面看是没有碎片的，使得操作系统层面对碎片的管理压力也会降低。
4. mcache 的存在使得内存分配不需要加锁。

参考：
http://ms.gtalent.cn/pin/d0PJ_1hYQ7J9mkia84KmBU

# 10 其它

**1、中间状态怎么记录**
- 在Go语言中，要记录中间状态并进行对比通常可以通过变量、数据结构或者函数返回值来实现

相关文章
- [Mastering Golang’s Concurrency and Memory Management: GMP, Garbage Collection, and Channel Handling](https://charleswan111.medium.com/mastering-golangs-concurrency-and-memory-management-gmp-garbage-collection-and-channel-handling-212dea055961)
- [(Golang Triad)-II-Comprehensive Analysis of Go's Hybrid Write Barrier Garbage Collection](https://dev.to/aceld/golang-triad-ii-comprehensive-analysis-of-gos-hybrid-write-barrier-garbage-collection-3knh)


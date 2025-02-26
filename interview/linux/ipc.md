# 进程间通信

汇总一下关于进程间通信（IPC）的知识。

## 进程通信的目的

- **数据传输**
  一个进程需要将它的数据发送给另一个进程，发送的数据量在一个字节到几M字节之间
- **共享数据**
  多个进程想要操作共享数据，一个进程对共享数据
- **通知事件**
  一个进程需要向另一个或一组进程发送消息，通知它（它们）发生了某种事件（如进程终止时要通知父进程）。
- **资源共享**
  多个进程之间共享同样的资源。为了作到这一点，需要内核提供锁和同步机制。
- **进程控制**
  有些进程希望完全控制另一个进程的执行（如Debug进程），此时控制进程希望能够拦截另一个进程的所有陷入和异常，并能够及时知道它的状态改变。



## 进程通信的发展

> linux下的进程通信手段基本上是从Unix平台上的进程通信手段继承而来的。而对Unix发展做出重大贡献的两大主力AT&T的贝尔实验室及BSD（加州大学伯克利分校的伯克利软件发布中心）在进程间通信方面的侧重点有所不同。
>
> 前者对Unix早期的进程间通信手段进行了系统的改进和扩充，形成了“system V IPC”，通信进程局限在单个计算机内；
>
> 后者则跳过了该限制，形成了基于套接口（socket）的进程间通信机制。
>
> Linux则把两者继承了下来

- 早期UNIX进程间通信
- 基于System V进程间通信
- 基于Socket进程间通信
- POSIX进程间通信。

UNIX进程间通信方式包括：管道、FIFO、信号。

System V进程间通信方式包括：System V消息队列、System V信号灯、System V共享内存

POSIX进程间通信包括：posix消息队列、posix信号灯、posix共享内存。

> 由于Unix版本的多样性，电子电气工程协会（IEEE）开发了一个独立的Unix标准，这个新的ANSI Unix标准被称为计算机环境的可移植性操作系统界面（PSOIX）。现有大部分Unix和流行版本都是遵循POSIX标准的，而Linux从一开始就遵循POSIX标准；
>
> BSD并不是没有涉足单机内的进程间通信（socket本身就可以用于单机内的进程间通信）。事实上，很多Unix版本的单机IPC留有BSD的痕迹，如4.4BSD支持的匿名内存映射、4.3+BSD对可靠信号语义的实现等等。



## linux使用的进程间通信方式

1. 管道（pipe）,流管道(s_pipe)和有名管道（FIFO）
2. 信号（signal）
3. 消息队列
4. 共享内存
5. 信号量
6. 套接字（socket)



#### 管道( pipe )

管道这种通讯方式有两种限制，一是半双工的通信，数据只能单向流动，二是只能在具有亲缘关系的进程间使用。进程的亲缘关系通常是指父子进程关系。

流管道s_pipe: 去除了第一种限制,可以双向传输.

管道可用于具有亲缘关系进程间的通信，命名管道:name_pipe克服了管道没有名字的限制，因此，除具有管道所具有的功能外，它还允许无亲缘关系进程间的通信；



#### 信号量( semophore )

信号量是一个计数器，可以用来控制多个进程对共享资源的访问。它常作为一种锁机制，防止某进程正在访问共享资源时，其他进程也访问该资源。因此，主要作为进程间以及同一进程内不同线程之间的同步手段。

信号是比较复杂的通信方式，用于通知接受进程有某种事件发生，除了用于进程间通信外，进程还可以发送信号给进程本身；linux除了支持Unix早期信号语义函数sigal外，还支持语义符合Posix.1标准的信号函数sigaction（实际上，该函数是基于BSD的，BSD为了实现可靠信号机制，又能够统一对外接口，用sigaction函数重新实现了signal函数）；



#### 消息队列( message queue )

消息队列是由消息的链表，存放在内核中并由消息队列标识符标识。消息队列克服了信号传递信息少、管道只能承载无格式字节流以及缓冲区大小受限等缺点。

消息队列是消息的链接表，包括Posix消息队列system V消息队列。有足够权限的进程可以向队列中添加消息，被赋予读权限的进程则可以读走队列中的消息。消息队列克服了信号承载信息量少，管道只能承载无格式字节流以及缓冲区大小受限等缺点。



#### 信号 ( singal )

信号是一种比较复杂的通信方式，用于通知接收进程某个事件已经发生。

主要作为进程间以及同一进程不同线程之间的同步手段。



#### 共享内存( shared memory )

共享内存就是映射一段能被其他进程所访问的内存，这段共享内存由一个进程创建，但多个进程都可以访问。共享内存是最快的 IPC 方式，它是针对其他进程间通信方式运行效率低而专门设计的。它往往与其他通信机制，如信号量，配合使用，来实现进程间的同步和通信。

使得多个进程可以访问同一块内存空间，是最快的可用IPC形式。是针对其他通信机制运行效率较低而设计的。往往与其它通信机制，如信号量结合使用，来达到进程间的同步及互斥。

共享内存针对消息缓冲的缺点改而利用内存缓冲区直接交换信息，无须复制，快捷、信息量大是其优点。

但是共享内存的通信方式是通过将共享的内存缓冲区直接附加到进程的虚拟地址空间中来实现的，因此，这些进程之间的读写操作的同步问题操作系统无法实现。必须由各进程利用其他同步工具解决。另外，由于内存实体存在于计算机系统中，所以只能由处于同一个计算机系统中的诸进程共享。不方便网络通信。

共享内存块提供了在任意数量的进程之间进行高效双向通信的机制。每个使用者都可以读取写入数据，但是所有程序之间必须达成并遵守一定的协议，以防止诸如在读取信息之前覆写内存空间等竞争状态的出现。

不幸的是，Linux无法严格保证提供对共享内存块的独占访问，甚至是在您通过使用IPC_PRIVATE创建新的共享内存块的时候也不能保证访问的独占性。 同时，多个使用共享内存块的进程之间必须协调使用同一个键值。



#### 套接字( socket )

套解字也是一种进程间通信机制，与其他通信机制不同的是，它可用于不同机器间的进程通信

更为一般的进程间通信机制，可用于不同机器之间的进程间通信。起初是由Unix系统的BSD分支开发出来的，但现在一般可以移植到其它类Unix系统上：Linux和System V的变种都支持套接字。



## 比较

| 类型             | 无连接 | 可靠 | 流控制 | 记录消息类型 | 优先级 |
| ---------------- | ------ | ---- | ------ | ------------ | ------ |
| 普通PIPE         | N      | Y    | Y      |              | N      |
| 流PIPE           | N      | Y    | Y      |              | N      |
| 命名PIPE(FIFO)   | N      | Y    | Y      |              | N      |
| 消息队列         | N      | Y    | Y      |              | Y      |
| 信号量           | N      | Y    | Y      |              | Y      |
| 共享存储         | N      | Y    | Y      |              | Y      |
| UNIX流SOCKET     | N      | Y    | Y      |              | N      |
| UNIX数据包SOCKET | Y      | Y    | N      |              | N      |

> 注:无连接: 指无需调用某种形式的OPEN,就有发送消息的能力流控制:
>
> 如果系统资源短缺或者不能接收更多消息,则发送进程能进行流量控制

各种通信方式的比较和优缺点

1. 管道：速度慢，容量有限，只有父子进程能通讯
2. FIFO：任何进程间都能通讯，但速度慢
3. 消息队列：容量受到系统限制，且要注意第一次读的时候，要考虑上一次没有读完数据的问题
4. 信号量：不能传递复杂消息，只能用来同步
5. 共享内存区：能够很容易控制容量，速度快，但要保持同步，比如一个进程在写的时候，另一个进程要注意读写的问题，相当于线程中的线程安全，当然，共享内存区同样可以用作线程间通讯，不过没这个必要，线程间本来就已经共享了同一进程内的一块内存

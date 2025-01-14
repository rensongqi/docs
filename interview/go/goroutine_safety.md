## 第一种方式：

对slice加锁，进行保护

```go
num := 10000

var a []int
var l sync.Mutex

var wg sync.WaitGroup
wg.Add(num)

for i := 0; i < num; i++ {
    go func() {
        l.Lock() // 加锁
        a = append(a, 1)
        l.Unlock() // 解锁
        wg.Done()
    }()
}

wg.Wait()

fmt.Println(len(a))
```

缺点：锁会影响性能

 

## 第二种方式：

使用channel的传递数据

```go
num := 10000

var wg sync.WaitGroup
wg.Add(num)

c := make(chan int)
for i := 0; i < num; i++ {
    go func() {
        c <- 1 // channl是协程安全的
        wg.Done()
    }()
}

// 等待关闭channel
go func() {
    wg.Wait()
    close(c)
}()

// 读取数据
var a []int
for i := range c {
    a = append(a, i)
}

fmt.Println(len(a))
```



## 第三种方式：

使用索引

```go
num := 10000

a := make([]int, num, num)

var wg sync.WaitGroup
wg.Add(num)

for i := 0; i < num; i++ {
    i := i // 必须使用局部变量
    go func() {
        a[i] = 1
        wg.Done()
    }()
}

wg.Wait()

count := 0
for i := range a {
    if a[i] != 0 {
        count++
    }
}
fmt.Println(count)
```

优点：无锁，不影响性能
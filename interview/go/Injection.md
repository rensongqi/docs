
# 依赖注入

> Golang 依赖注入（Dependency Injection, DI）是一种设计模式，它用于将依赖关系（即一个组件需要的外部对象）从组件内部的创建过程剥离出来，而是在外部构造并传递给该组件。这种方式可以提高代码的可测试性、可维护性和灵活性。

为什么需要依赖注入？

- 降低耦合：组件不需要关心依赖对象的创建和管理，而是由外部提供。
- 提高可测试性：可以在测试时轻松替换真实依赖，使用 mock/stub 进行单元测试。
- 增强灵活性：可以更方便地更换不同的实现（如使用不同的数据库、日志库等）。

## 1.1 手动依赖注入（最常见）

手动创建依赖并将其作为参数传递。
```go
package main

import "fmt"

// 定义一个 Logger 结构体
type Logger struct{}

func (l *Logger) Log(message string) {
    fmt.Println("Log:", message)
}

// Service 依赖 Logger
type Service struct {
    logger *Logger
}

// Service 需要 Logger 作为依赖
func NewService(logger *Logger) *Service {
    return &Service{logger: logger}
}

func (s *Service) DoSomething() {
    s.logger.Log("Doing something...")
}

func main() {
    logger := &Logger{}       // 由外部创建 Logger
    service := NewService(logger) // 通过构造函数注入 Logger
    service.DoSomething()
}
```
优势：
- 代码清晰，易于理解
- 依赖关系显式可见
- 不需要额外库，完全由 Go 语言原生支持

## 1.2 使用全局变量（不推荐）

```go
var logger = &Logger{}

type Service struct{}

func (s *Service) DoSomething() {
    logger.Log("Doing something...")
}

```
问题：
- 依赖是隐藏的，难以测试
- 代码可维护性差
## 1.3 使用 Google 的 wire 依赖注入框架

[Google Wire](https://github.com/google/wire) 是 Go 中一个流行的依赖注入工具，使用代码生成的方式实现依赖注入。[Wire tutorial](https://github.com/google/wire/blob/main/_tutorial/README.md)

```go
// wire.go
package main

import "github.com/google/wire"

// 提供 Logger 实例
func NewLogger() *Logger {
    return &Logger{}
}

// 生成 Service 实例
func NewService(logger *Logger) *Service {
    return &Service{logger: logger}
}

// 使用 wire 绑定依赖
var ProviderSet = wire.NewSet(NewLogger, NewService)
```

优势：

- 代码清晰，无需手动管理依赖树
- 自动生成依赖图，避免手写工厂方法

劣势：

- 需要额外的工具
- 代码生成带来一定复杂性
# DelayTask 包

延时任务处理包，基于Redis实现的分布式延时任务队列。

## 功能特性

- **任务推送**: 将任务推送到指定时间点的队列中
- **任务拉取**: 按时间顺序拉取到期的任务
- **分布式支持**: 基于Redis，支持多实例部署
- **灵活的时间间隔**: 支持毫秒、秒、分钟、小时级别的时间间隔

## 主要函数

### 基础函数
- `PushTaskInternal(r, ctx, taskname, tick_time, content, interval)` - 推送任务到指定间隔的队列
- `PullTaskInternal(r, ctx, taskname, interval)` - 从指定间隔的队列拉取任务

### 便捷函数
- `PushTask(r, ctx, taskname, tick_time, content)` - 推送任务（使用秒级间隔）
- `PullTask(r, ctx, taskname)` - 拉取任务（使用秒级间隔）

### 时间间隔常量
```go
const (
    INTERVAL_MILLISECONDS = 1          // 毫秒级
    INTERVAL_SECONDS      = 1000       // 秒级
    INTERVAL_MUNITES      = 60 * 1000  // 分钟级
    INTERVAL_HOUR         = 60 * 60 * 1000 // 小时级
)
```

## 使用示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/daozhonglee/go-util/delaytask"
    times "github.com/daozhonglee/go-util/time"
    "github.com/redis/go-redis/v9"
)

func main() {
    // 创建Redis客户端
    r := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        Password: "",
        DB: 0,
    })
    
    ctx := context.Background()
    taskname := "example_task"
    
    // 推送任务
    currentTime := times.GetCurrentMilliUnix()
    err := delaytask.PushTask(r, ctx, taskname, currentTime, "Hello, World!")
    if err != nil {
        log.Fatal(err)
    }
    
    // 拉取任务
    tasks, err := delaytask.PullTask(r, ctx, taskname)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, task := range tasks {
        fmt.Println("处理任务:", task)
    }
}
```

## 运行测试


在运行测试之前，请确保Redis服务器正在运行：


### 运行测试

```bash
# 运行所有测试
go test -v ./delaytask

# 运行特定测试
go test -v ./delaytask -run TestPushTask

# 使用自定义Redis地址
REDIS_ADDR=localhost:6379 go test -v ./delaytask

# 运行基准测试
go test -v ./delaytask -bench=.
```

### 测试说明

测试包含以下场景：

1. **基础功能测试**
   - `TestPushTaskInternal` - 测试内部推送函数
   - `TestPullTaskInternal` - 测试内部拉取函数
   - `TestPushTask` - 测试便捷推送函数
   - `TestPullTask` - 测试便捷拉取函数

2. **边界情况测试**
   - `TestPullTaskEmpty` - 测试拉取空队列
   - `TestConstants` - 测试时间间隔常量

3. **工作流程测试**
   - `TestDelayTaskWorkflow` - 测试完整的推送-拉取工作流
   - `TestDifferentIntervals` - 测试不同时间间隔

4. **连接处理**
   - 如果Redis不可用，相关测试会被自动跳过
   - 可以通过环境变量 `REDIS_ADDR` 配置Redis地址

## 注意事项

1. **Redis键命名规则**：
   - 任务队列键：`dtaskq:{taskname}:{tick}`
   - 时间戳键：`dtaskt:{taskname}`

2. **任务过期时间**：
   - 任务队列默认过期时间：24小时 (86400秒)
   - 时间戳键过期时间：2小时 (7200秒)

3. **任务拉取限制**：
   - 每次最多拉取100个任务

4. **并发安全**：
   - 使用Redis的原子操作确保并发安全
   - 支持多实例同时处理任务
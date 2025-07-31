# AT-GoLib

Go语言工具库集合，提供常用的工具函数和组件。

## 🏗️ 项目结构

```
go-util/
├── id/           # 身份证相关
├── crypto/       # 加密和哈希相关  
├── json/         # JSON处理
├── random/       # 随机数生成
├── xss/          # XSS防护
├── async/        # 异步操作
├── time/         # 时间处理（包名：times）
├── collection/   # 集合类型
├── api/          # HTTP响应
├── delaytask/    # 延时任务队列
├── errorutil/    # 错误处理和panic恢复
├── pkg/          # 统一入口包
│   ├── util.go   # 主入口，重新导出常用函数
│   └── util_test.go # 入口测试文件
└── example/      # 使用示例
```

## 📦 功能模块

| 包名 | 功能描述 | 主要函数 |
|------|----------|----------|
| **id** | 中国身份证验证 | `ValidateChinese()`, `CheckAge16To18()` |
| **crypto** | 加密和哈希计算 | `MD5Hash()`, `SHA256Hash()`, `HMACSign()` |
| **json** | JSON处理 | `Marshal()`, `Unmarshal()`, `Pretty()` |
| **random** | 随机数生成 | `Int(min, max)` |
| **xss** | XSS防护 | `Clean()` |
| **async** | 异步操作 | `Go()`, `Timeout()`, `Safe()` |
| **time** | 时间处理（包名times） | `Relative()` |
| **collection** | 集合数据结构 | `NewSet()` |
| **api** | HTTP响应 | `Success()`, `Error()` |
| **errorutil** | 错误处理和panic恢复 | `PanicIf()`, `PanicWithStack()`, `Recover()` |

## 🚀 使用方式

### 方式1：使用子包（推荐）

```go
import (
    "github.com/daozhonglee/go-util/id"
    "github.com/daozhonglee/go-util/crypto"
    "github.com/daozhonglee/go-util/json"
    times "github.com/daozhonglee/go-util/time"  // 注意包名冲突
)

// 身份证验证
isValid := id.ValidateChinese("11010119900307001X")

// 哈希计算
hash := crypto.MD5Hash([]byte("test"))

// JSON处理
jsonStr := json.Marshal(map[string]string{"key": "value"})

// 时间处理
relativeTime := times.Relative(pastTimestamp)
```

### 方式2：使用统一入口（pkg包）

```go
import "github.com/daozhonglee/go-util/pkg"

// 身份证验证
isValid := util.ValidateChineseID("11010119900307001X")

// 哈希计算
hash := util.MD5Hash([]byte("test"))

// JSON处理
jsonStr := util.JSONMarshal(map[string]string{"key": "value"})

// 时间处理
relativeTime := util.RelativeTime(pastTimestamp)
```

## 🧪 测试

运行所有测试：

```bash
go test -v ./...
```

运行特定包的测试：

```bash
go test -v ./id
go test -v ./crypto
```

## 💡 设计理念

- **模块化**：功能按领域分包，清晰明确
- **按需导入**：只导入需要的功能，减少依赖
- **简洁API**：函数命名简洁直观
- **性能优先**：考虑高性能场景的优化
- **错误友好**：提供详细的错误信息和恢复机制

## 📋 依赖管理

- 最小化外部依赖
- 每个子包的依赖独立管理
- 支持Go modules
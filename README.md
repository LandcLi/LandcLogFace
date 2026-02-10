# LandcLogFace - Go语言日志门面项目

LandcLogFace是一个Go语言的日志门面（Logging Facade）项目，提供统一的日志接口，支持多种常用日志库的适配，帮助开发者在不同日志实现之间无缝切换，同时保持代码的一致性和可维护性。

## 项目特性

- **统一的日志接口**：定义标准的Logger接口，包含各种日志级别和功能
- **按需导入**：使用Build Tags实现依赖分离，只引入需要的日志库
- **多日志库支持**：实现多种常用日志库的适配器
  - 控制台日志（默认，无第三方依赖）
  - zap日志库（高性能）
  - logrus日志库（功能丰富）
  - 标准库log（轻量，无第三方依赖）
- **框架适配器**：支持Gin和GoFrame框架的日志集成
- **灵活的配置管理**：支持通过选项函数和配置map进行灵活配置
- **日志工厂**：提供统一的日志实例创建和管理功能
- **全局日志**：提供便捷的全局日志函数
- **结构化字段**：支持添加结构化日志字段
- **上下文支持**：支持添加上下文信息
- **错误处理**：支持添加错误信息
- **时间管理**：支持自定义时间字段
- **日志级别控制**：支持细粒度的日志级别控制
- **日志文件轮转**：支持根据文件大小自动切分日志文件
- **日志保留策略**：支持设置日志文件的最大保留时间和数量
- **日志压缩**：支持压缩旧日志文件以节省空间
- **单条日志大小限制**：支持限制单条日志的最大大小
- **可扩展性**：支持自定义日志提供者

## 安装

### 前提条件

- Go 1.24.0或更高版本

### 安装步骤

1. **使用go get安装**

```bash
go get github.com/LandcLi/LandcLogFace
```

2. **或直接克隆仓库**

```bash
git clone https://github.com/LandcLi/LandcLogFace.git
cd LandcLogFace
go mod tidy
```

## 快速开始

### 最小化使用（无第三方依赖）

```go
package main

import "github.com/LandcLi/LandcLogFace"

func main() {
	// 获取全局日志实例
	logger := LandcLogFace.GetLogger()

	// 输出不同级别的日志
	logger.Debug("这是一条调试日志")
	logger.Info("这是一条信息日志")
	logger.Warn("这是一条警告日志")
	logger.Error("这是一条错误日志")

	// 输出格式化日志
	logger.Infof("Hello, %s!", "world")
	logger.Errorf("Error: %d", 404)

	// 使用字段
	logger.Info("用户登录", 
		LandcLogFace.Field{Key: "user_id", Value: 123},
		LandcLogFace.Field{Key: "ip", Value: "192.168.1.1"},
	)
}
```

### 使用全局函数

```go
package main

import "github.com/LandcLi/LandcLogFace"

func main() {
	// 使用全局函数输出日志
	LandcLogFace.Debug("全局调试日志")
	LandcLogFace.Info("全局信息日志")
	LandcLogFace.Warn("全局警告日志")
	LandcLogFace.Error("全局错误日志")

	// 使用全局格式化函数
	LandcLogFace.Infof("全局格式化日志: %s", "test")
	LandcLogFace.Errorf("全局错误格式化日志: %d", 500)
}
```

## 按需导入日志库

LandcLogFace 使用 Go 的 Build Tags 特性实现依赖分离，让你可以按需导入日志库，避免不必要的依赖。

### 核心包（无第三方依赖）

核心包不包含任何第三方日志库依赖，只提供基础的日志接口和功能：

```go
import "github.com/LandcLi/LandcLogFace"
```

核心包包含：
- 基础日志接口定义
- 控制台日志实现（默认）
- 标准库日志实现
- 日志工厂和配置管理
- 全局日志函数

### 使用 Zap 日志库

如果你只想使用 Zap 日志库，只需要导入 Zap 提供者包：

```go
import (
    "github.com/LandcLi/LandcLogFace"
    _ "github.com/LandcLi/LandcLogFace/providers/zap" // 导入并注册 Zap 提供者
)

func main() {
    // Zap 提供者已自动注册，可以直接使用
    logger := LandcLogFace.GetLoggerWithProvider("app", "zap",
        LandcLogFace.WithLevel(LandcLogFace.InfoLevel),
        LandcLogFace.WithFormat("json"),
    )
    logger.Info("使用 Zap 日志")
}
```

**依赖**：只会引入 `go.uber.org/zap` 和 `gopkg.in/natefinch/lumberjack.v2`

### 使用 Logrus 日志库

如果你只想使用 Logrus 日志库，只需要导入 Logrus 提供者包：

```go
import (
    "github.com/LandcLi/LandcLogFace"
    _ "github.com/LandcLi/LandcLogFace/providers/logrus" // 导入并注册 Logrus 提供者
)

func main() {
    // Logrus 提供者已自动注册，可以直接使用
    logger := LandcLogFace.GetLoggerWithProvider("app", "logrus",
        LandcLogFace.WithLevel(LandcLogFace.InfoLevel),
        LandcLogFace.WithFormat("text"),
    )
    logger.Info("使用 Logrus 日志")
}
```

**依赖**：只会引入 `github.com/sirupsen/logrus` 和 `gopkg.in/natefinch/lumberjack.v2`

### 组合使用

你可以同时导入多个提供者，根据需要选择使用：

```go
import (
    "github.com/LandcLi/LandcLogFace"
    _ "github.com/LandcLi/LandcLogFace/providers/zap"
    _ "github.com/LandcLi/LandcLogFace/providers/logrus"
)

func main() {
    // 使用 Zap 日志
    zapLogger := LandcLogFace.GetLoggerWithProvider("app", "zap")
    zapLogger.Info("使用 Zap")

    // 使用 Logrus 日志
    logrusLogger := LandcLogFace.GetLoggerWithProvider("app", "logrus")
    logrusLogger.Info("使用 Logrus")

    // 使用默认的控制台日志
    consoleLogger := LandcLogFace.GetLoggerWithProvider("app", "console")
    consoleLogger.Info("使用控制台")
}
```

## 详细使用指南

### 1. 选择不同的日志库

LandcLogFace支持多种日志库，你可以根据需要选择合适的日志实现：

```go
package main

import (
	"github.com/LandcLi/LandcLogFace"
)

func main() {
	// 使用console提供者（默认）
	consoleLogger := LandcLogFace.GetLoggerWithProvider("app", "console")
	consoleLogger.Info("使用控制台日志")

	// 使用std提供者（轻量）
	stdLogger := LandcLogFace.GetLoggerWithProvider("app", "std")
	stdLogger.Info("使用标准库日志")
}
```

**注意**：要使用 `zap` 或 `logrus` 提供者，需要先导入对应的提供者包。

### 2. 配置日志实例

你可以通过选项函数或配置map来配置日志实例：

#### 使用选项函数

```go
package main

import (
	"github.com/LandcLi/LandcLogFace"
    _ "github.com/LandcLi/LandcLogFace/providers/zap"
)

func main() {
	// 使用选项函数配置
	logger := LandcLogFace.GetLoggerWithProvider("app", "zap",
		LandcLogFace.WithLevel(LandcLogFace.DebugLevel),
		LandcLogFace.WithFormat("json"),
		LandcLogFace.WithOutputPath("stdout"),
	)

	logger.Debug("调试日志")
}
```

#### 使用配置map

```go
package main

import (
	"github.com/LandcLi/LandcLogFace"
)

func main() {
	// 使用配置map
	config := map[string]interface{}{
		"provider":    "zap",
		"level":       LandcLogFace.DebugLevel,
		"format":      "json",
		"outputPath":  "stdout",
	}

	logger := LandcLogFace.GetLoggerWithMap("app", config)
	logger.Info("带配置的日志")
}
```

### 3. 高级功能

#### 字段使用

```go
package main

import (
	"github.com/LandcLi/LandcLogFace"
)

func main() {
	logger := LandcLogFace.GetLogger()

	// 创建带字段的日志实例
	userLogger := logger.WithFields(
		LandcLogFace.Field{Key: "service", Value: "user"},
		LandcLogFace.Field{Key: "version", Value: "1.0.0"},
	)

	// 使用带字段的日志实例
	userLogger.Info("用户注册", 
		LandcLogFace.Field{Key: "username", Value: "john"},
		LandcLogFace.Field{Key: "email", Value: "john@example.com"},
	)

	// 链式添加字段
	userLogger.WithField("status", "success").Info("注册成功")
}
```

#### 上下文支持

```go
package main

import (
	"github.com/LandcLi/LandcLogFace"
	"context"
)

func main() {
	logger := LandcLogFace.GetLogger()

	// 创建上下文
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "123456")

	// 添加上下文到日志
	ctxLogger := logger.WithContext(ctx)
	ctxLogger.Info("带上下文的日志")
}
```

#### 错误处理

```go
package main

import (
	"github.com/LandcLi/LandcLogFace"
	"errors"
)

func main() {
	logger := LandcLogFace.GetLogger()

	// 创建错误
	err := errors.New("数据库连接失败")

	// 添加错误到日志
	errLogger := logger.WithError(err)
	errLogger.Error("操作失败")

	// 直接在日志中添加错误字段
	logger.Error("操作失败", LandcLogFace.Field{Key: "error", Value: err})
}
```

#### 时间管理

```go
package main

import (
	"github.com/LandcLi/LandcLogFace"
	"time"
)

func main() {
	logger := LandcLogFace.GetLogger()

	// 添加自定义时间
	now := time.Now()
	logger.WithTime(now).Info("带自定义时间的日志")
}
```

#### 日志级别检查

```go
package main

import (
	"github.com/LandcLi/LandcLogFace"
)

func main() {
	logger := LandcLogFace.GetLogger()

	// 检查日志级别是否启用
	if logger.IsDebugEnabled() {
		// 执行耗时操作
		data := "耗时操作的结果"
		logger.Debug("调试信息", LandcLogFace.Field{Key: "data", Value: data})
	}

	if logger.IsInfoEnabled() {
		logger.Info("信息日志已启用")
	}
}
```

### 4. 日志文件轮转配置

LandcLogFace支持详细的日志文件轮转配置，包括文件大小限制、保留时间、文件数量等参数：

#### 配置选项

| 配置项 | 类型 | 默认值 | 描述 |
|-------|------|-------|------|
| `MaxLogSize` | `int64` | 100 | 单个日志文件最大大小（MB） |
| `MaxLogAge` | `time.Duration` | 7*24*time.Hour | 日志文件最大保留时间 |
| `MaxLogFiles` | `int` | 10 | 最大保留日志文件数量 |
| `CompressLogs` | `bool` | false | 是否压缩旧日志 |
| `MaxMessageSize` | `int` | 0 | 单条日志最大大小（KB），0表示不限制 |

#### 使用示例

```go
package main

import (
	"github.com/LandcLi/LandcLogFace"
	"time"
    _ "github.com/LandcLi/LandcLogFace/providers/zap"
)

func main() {
	// 配置文件轮转
	logger := LandcLogFace.GetLoggerWithProvider("app", "zap",
		LandcLogFace.WithLevel(LandcLogFace.InfoLevel),
		LandcLogFace.WithFormat("json"),
		LandcLogFace.WithOutputPath("app.log"),
		LandcLogFace.WithMaxLogSize(50),        // 单个日志文件最大50MB
		LandcLogFace.WithMaxLogAge(7*24*time.Hour), // 日志文件保留7天
		LandcLogFace.WithMaxLogFiles(5),        // 最多保留5个日志文件
		LandcLogFace.WithCompressLogs(true),    // 压缩旧日志
	)

	logger.Info("使用文件轮转的zap日志")

	// 配置单条日志大小限制
	sizeLogger := LandcLogFace.GetLoggerWithProvider("app", "logrus",
		LandcLogFace.WithLevel(LandcLogFace.InfoLevel),
		LandcLogFace.WithFormat("text"),
		LandcLogFace.WithOutputPath("app.log"),
		LandcLogFace.WithMaxMessageSize(10), // 单条日志最大10KB
	)

	sizeLogger.Info("限制单条日志大小的logrus日志")
}
```

### 5. 使用统一配置类

LandcLogFace提供了`LogConfig`统一配置类，用于集中管理所有日志配置选项：

#### 配置示例

```go
package main

import (
	"github.com/LandcLi/LandcLogFace"
	"time"
    _ "github.com/LandcLi/LandcLogFace/providers/zap"
)

func main() {
	// 创建并配置LogConfig
	config := LandcLogFace.NewLogConfig()
	config.WithProvider("zap").
		WithName("myapp").
		WithLevel(LandcLogFace.DebugLevel).
		WithFormat("json").
		WithOutputPath("app.log").
		WithMaxLogSize(50).        // 单个日志文件最大50MB
		WithMaxLogAge(7*24*time.Hour). // 日志文件保留7天
		WithMaxLogFiles(5).        // 最多保留5个日志文件
		WithCompressLogs(true).    // 压缩旧日志
		WithMaxMessageSize(10).    // 单条日志最大10KB
		WithExtraConfig("key", "value") // 添加额外配置

	// 根据LogConfig创建日志实例
	logger := LandcLogFace.GetLoggerWithLogConfig(config)
	logger.Info("使用LogConfig的日志")
}
```

#### 配置项说明

| 配置项 | 类型 | 默认值 | 描述 |
|-------|------|-------|------|
| `Provider` | `string` | "console" | 日志提供者名称 |
| `Name` | `string` | "app" | 日志名称 |
| `Level` | `LogLevel` | `InfoLevel` | 日志级别 |
| `Format` | `string` | "text" | 日志格式（text/json） |
| `OutputPath` | `string` | "stdout" | 日志输出路径 |
| `MaxLogSize` | `int64` | 100 | 单个日志文件最大大小（MB） |
| `MaxLogAge` | `time.Duration` | 7*24*time.Hour | 日志文件最大保留时间 |
| `MaxLogFiles` | `int` | 10 | 最大保留日志文件数量 |
| `CompressLogs` | `bool` | false | 是否压缩旧日志 |
| `MaxMessageSize` | `int` | 0 | 单条日志最大大小（KB），0表示不限制 |
| `ExtraConfig` | `map[string]interface{}` | 空 | 额外的提供者特定配置 |

### 6. 框架适配器

LandcLogFace提供了常用Web框架的日志适配器，方便在框架中使用统一的日志系统。

#### 6.1 Gin框架适配器

**注意：使用Gin适配器前，需要先安装Gin框架依赖：**

```bash
go get github.com/gin-gonic/gin
```

使用Gin适配器可以将Gin框架的请求日志集成到统一的日志系统中：

```go
package main

import (
	"github.com/LandcLi/LandcLogFace"
	"github.com/LandcLi/LandcLogFace/providers/zap"
	"github.com/LandcLi/LandcLogFace/adapter/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建日志实例
	config := LandcLogFace.NewLogConfig()
	config.WithProvider("zap").
		WithName("gin-app").
		WithLevel(LandcLogFace.InfoLevel).
		WithFormat("json").
		WithOutputPath("gin.log")

	logger := LandcLogFace.GetLoggerWithLogConfig(config)

	// 创建Gin引擎
	r := gin.Default()

	// 使用LandcLogFace的Gin适配器
	ginProvider.UseWithGin(r, logger)

	// 定义路由
	r.GET("/", func(c *gin.Context) {
		// 在处理函数中使用统一日志
		LandcLogFace.Info("处理根路径请求")
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// 启动服务器
	r.Run(":8080")
}
```

#### 6.2 GoFrame框架适配器

**注意：使用GoFrame适配器前，需要先安装GoFrame框架依赖：**

```bash
go get github.com/gogf/gf/v2
```

使用GoFrame适配器可以将GoFrame框架的日志系统替换为统一的日志系统：

```go
package main

import (
	"github.com/LandcLi/LandcLogFace"
	"github.com/LandcLi/LandcLogFace/providers/logrus"
	"github.com/LandcLi/LandcLogFace/adapter/gf"
	"context"

	"github.com/gogf/gf/v2/os/glog"
)

func main() {
	// 创建日志实例
	config := LandcLogFace.NewLogConfig()
	config.WithProvider("logrus").
		WithName("gf-app").
		WithLevel(LandcLogFace.InfoLevel).
		WithFormat("text").
		WithOutputPath("gf.log")

	logger := LandcLogFace.GetLoggerWithLogConfig(config)

	// 创建GoFrame日志适配器
	gfLogger := gf.NewGFLogger(logger)

	// 设置GoFrame的全局日志器
	glog.SetLogger(gfLogger)

	// 使用GoFrame的日志系统（会自动使用我们的适配器）
	ctx := context.Background()
	glog.Info(ctx, "GoFrame日志测试")
	glog.Debug(ctx, "调试信息")
	glog.Error(ctx, "错误信息")
}
```

### 7. 自定义日志提供者

如果你需要使用项目未内置的日志库，可以通过实现`LoggerProvider`接口来添加自定义日志提供者：

```go
package main

import (
	"github.com/LandcLi/LandcLogFace"
	"fmt"
)

// CustomLogger 自定义日志实现
type CustomLogger struct {
	name string
}

// NewCustomLogger 创建自定义日志实例
func NewCustomLogger(name string) *CustomLogger {
	return &CustomLogger{name: name}
}

// 实现Logger接口的所有方法...
// （具体实现可参考项目中的示例代码）

// CustomLoggerProvider 自定义日志提供者
type CustomLoggerProvider struct{}

// Create 创建日志实例
func (p *CustomLoggerProvider) Create(name string, opts ...LandcLogFace.Option) LandcLogFace.Logger {
	return NewCustomLogger(name)
}

// CreateWithConfig 根据配置创建日志实例
func (p *CustomLoggerProvider) CreateWithConfig(name string, config map[string]interface{}) LandcLogFace.Logger {
	return NewCustomLogger(name)
}

func main() {
	// 注册自定义提供者
	LandcLogFace.RegisterProvider("custom", &CustomLoggerProvider{})

	// 使用自定义提供者
	customLogger := LandcLogFace.GetLoggerWithProvider("app", "custom")
	customLogger.Info("使用自定义日志提供者")

	// 注销自定义提供者
	LandcLogFace.UnregisterProvider("custom")
}
```

## 依赖对比

| 使用场景 | 必需依赖 |
|---------|----------|
| 仅核心包 | 无 |
| Zap | `go.uber.org/zap`, `gopkg.in/natefinch/lumberjack.v2` |
| Logrus | `github.com/sirupsen/logrus`, `gopkg.in/natefinch/lumberjack.v2` |
| Gin | `github.com/gin-gonic/gin` |
| GoFrame | `github.com/gogf/gf/v2` |

## 最佳实践

1. **按需导入**：只导入你实际需要的日志库提供者
2. **使用空白导入**：使用 `_` 前缀导入提供者包，自动注册提供者
3. **核心包优先**：优先使用核心包提供的功能，减少依赖
4. **明确依赖**：在 `go.mod` 中明确指定需要的依赖

## 项目结构

```
LandcLogFace/
├── landc_logface.go      # 主包入口文件
├── go.mod                # 项目依赖管理
├── go.sum                # 依赖版本锁定
├── README.md             # 项目文档
├── LICENSE               # 许可证文件
├── internal/             # 核心代码目录
│   └── logger/           # 日志核心实现
│       ├── logger.go         # 核心接口定义
│       ├── config.go         # 统一配置类
│       ├── log_factory.go    # 日志工厂和配置管理
│       ├── console_logger.go # 控制台日志适配器
│       ├── std_logger.go     # 标准库log适配器
│       ├── zap_logger.go     # zap日志库适配器
│       └── logrus_logger.go  # logrus日志库适配器
├── providers/           # 按需导入的提供者包
│   ├── zap/             # Zap日志库提供者
│   └── logrus/          # Logrus日志库提供者
├── adapter/             # 框架适配器目录
│   ├── gin/             # Gin框架适配器
│   └── gf/              # GoFrame框架适配器
├── examples/             # 示例代码目录
│   └── example.go        # 使用示例
└── tests/                # 测试目录
    ├── logger_test.go    # 测试用例
    ├── options_test.go   # 选项函数测试
    ├── custom_provider_test.go # 自定义提供者测试
    ├── core_only_test.go # 核心包测试
    ├── zap_test.go      # Zap提供者测试
    └── logrus_test.go   # Logrus提供者测试
```

## 依赖管理

项目使用Go 1.24.0版本，依赖以下第三方库：

| 依赖库 | 版本 | 用途 |
|-------|------|------|
| `go.uber.org/zap` | v1.26.0 | 高性能日志库（可选） |
| `github.com/sirupsen/logrus` | v1.9.3 | 功能丰富的日志库（可选） |
| `github.com/gin-gonic/gin` | v1.9.1 | Gin框架，用于实现Gin适配器（可选） |
| `gopkg.in/natefinch/lumberjack.v2` | v2.2.1 | 日志文件轮转库（可选） |

**可选依赖**
| 依赖库 | 版本 | 用途 |
|-------|------|------|
| `github.com/gogf/gf/v2` | v2.9.7 | GoFrame框架，用于实现GoFrame适配器（需要手动安装） |

## 测试

项目包含完整的测试用例，验证了所有核心功能的正确性：

```bash
# 测试核心包（无第三方依赖）
go test -v -tags=core_only ./tests/core_only_test.go

# 测试Zap提供者
go test -v -tags=zap_provider ./tests/zap_test.go

# 测试Logrus提供者
go test -v -tags=logrus_provider ./tests/logrus_test.go

# 运行所有测试
go test -v ./...
```

## 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 贡献

欢迎贡献代码、报告问题或提出改进建议。请通过以下方式参与：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 联系方式

- 项目主页：https://github.com/LandcLi/LandcLogFace
- 问题反馈：https://github.com/LandcLi/LandcLogFace/issues

## 致谢

感谢所有为本项目做出贡献的开发者！

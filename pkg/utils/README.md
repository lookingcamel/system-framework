# 工具函数模块

提供项目公共工具函数，包括环境变量、请求 ID 生成、通用辅助函数等。

## 功能特性

- 🔧 **环境变量**: 环境变量读取和默认值
- 🆔 **请求 ID**: 唯一请求标识生成和追踪
- 🔄 **类型转换**: 通用类型转换工具
- 📝 **字符串处理**: 字符串工具函数
- ⏰ **时间处理**: 时间格式化工具
- ✅ **数据验证**: 通用验证函数

## 快速开始

### 环境变量

```go
import "github.com/lookingcamel/system-framework/pkg/utils"

func main() {
    // 获取环境变量，带默认值
    dbHost := utils.GetEnv("DB_HOST", "localhost")
    dbPort := utils.GetEnvAsInt("DB_PORT", "3306")

    // 获取布尔值
    debug := utils.GetEnvAsBool("DEBUG", "false")

    fmt.Printf("DB: %s:%d, Debug: %v\n", dbHost, dbPort, debug)
}
```

### 请求 ID

```go
import "github.com/lookingcamel/system-framework/pkg/utils"

// 生成请求 ID
requestID := utils.GenerateRequestID()

// 设置到 Context
ctx := utils.SetRequestID(context.Background(), requestID)

// 从 Context 获取
requestID = utils.GetRequestID(ctx)

// Gin Context 操作
utils.GinSetRequestID(c, requestID)
requestID = utils.GinGetRequestID(c)
```

## 工具函数

### 1. 环境变量

```go
// 获取字符串值
func GetEnv(key, defaultValue string) string

// 获取整数值
func GetEnvAsInt(key string, defaultValue int) int

// 获取浮点数值
func GetEnvAsFloat(key string, defaultValue float64) float64

// 获取布尔值
func GetEnvAsBool(key string, defaultValue bool) bool

// 获取 Duration
func GetEnvAsDuration(key string, defaultValue time.Duration) time.Duration
```

### 2. 请求 ID

```go
// 生成 UUID
func GenerateRequestID() string

// Context 操作
func SetRequestID(ctx context.Context, requestID string) context.Context
func GetRequestID(ctx context.Context) string

// Gin Context 操作
func GinSetRequestID(c *gin.Context, requestID string)
func GinGetRequestID(c *gin.Context) string
```

### 3. 类型转换

```go
import "github.com/lookingcamel/system-framework/pkg/utils"

// 字符串转换
utils.ToString(value)           // 任意类型转字符串
utils.ToInt(value)              // 转整数
utils.ToInt64(value)            // 转 64 位整数
utils.ToFloat64(value)          // 转浮点数
utils.ToBool(value)             // 转布尔值

// 结构体转换
utils.StructToMap(obj)          // 结构体转 Map
utils.MapToStruct(m, obj)       // Map 转结构体

// JSON 操作
utils.JSONMarshal(obj)          // JSON 序列化
utils.JSONUnmarshal(data, obj)  // JSON 反序列化
```

### 4. 字符串处理

```go
// 字符串工具
utils.IsEmpty(str)             // 判断空字符串
utils.TrimSpace(str)            // 去除空格
utils.Truncate(str, length)     // 截断字符串
utils.SnakeToCamel(str)        // 蛇形转驼峰
utils.CamelToSnake(str)         // 驼峰转蛇形

// 字符串构建
utils.Join(list, separator)     // 拼接字符串
utils.BuildQuery(params)        // 构建查询字符串
```

### 5. 时间处理

```go
// 时间格式化
utils.FormatTime(t, layout)     // 格式化时间
utils.ParseTime(str, layout)    // 解析时间
utils.Now()                    // 当前时间
utils.Today()                   // 今天零点

// 时间计算
utils.AddDays(t, days)          // 加减天数
utils.AddHours(t, hours)        // 加减小时
utils.DaysBetween(t1, t2)       // 计算天数差

// 时间戳
utils.Timestamp()               // 秒时间戳
utils.Milliseconds()            // 毫秒时间戳
utils.Nanoseconds()             // 纳秒时间戳
```

### 6. 数据验证

```go
// 验证函数
utils.IsEmail(email)           // 验证邮箱
utils.IsURL(url)               // 验证 URL
utils.IsPhone(phone)           // 验证手机号
utils.IsIP(ip)                 // 验证 IP 地址
utils.InRange(value, min, max) // 验证范围

// 验证并返回错误
func ValidateEmail(email string) error
func ValidateURL(url string) error
```

### 7. 加密和哈希

```go
// MD5
utils.MD5(str)                 // MD5 哈希
utils.MD5Salt(str, salt)       // MD5 加盐

// SHA256
utils.SHA256(str)              // SHA256 哈希

// Base64
utils.Base64Encode(data)       // Base64 编码
utils.Base64Decode(str)        // Base64 解码

// UUID
utils.NewUUID()                // 生成 UUID
utils.IsValidUUID(str)         // 验证 UUID
```

### 8. 文件操作

```go
// 文件工具
utils.FileExists(path)         // 文件是否存在
utils.IsDir(path)              // 是否是目录
utils.CreateDir(path)          // 创建目录
utils.ReadFile(path)           // 读取文件
utils.WriteFile(path, data)    // 写入文件
utils.CopyFile(src, dst)       // 复制文件

// 路径工具
utils.JoinPath(parts...)       // 连接路径
utils.BaseName(path)           // 获取文件名
utils.DirName(path)            // 获取目录名
utils.ExtName(path)            // 获取扩展名
```

### 9. 网络工具

```go
// IP 地址
utils.GetLocalIP()            // 获取本机 IP
utils.IsPrivateIP(ip)          // 是否私有 IP
utils.IPToInt(ip)             // IP 转整数
utils.IntToIP(num)            // 整数转 IP

// HTTP 工具
utils.BuildURL(base, path, params)  // 构建 URL
utils.ParseQuery(query)        // 解析查询字符串
```

### 10. 随机数据

```go
// 随机字符串
utils.RandomString(length)     // 随机字符串
utils.RandomAlphanumeric(length) // 随机字母数字

// 随机数
utils.RandomInt(min, max)      // 随机整数
utils.RandomFloat(min, max)    // 随机浮点数

// 随机选择
utils.RandomChoice(list)       // 随机选择一个
```

## 使用示例

### 环境变量配置

```go
func loadConfig() *Config {
    return &Config{
        Host:     utils.GetEnv("APP_HOST", "0.0.0.0"),
        Port:     utils.GetEnvAsInt("APP_PORT", 8080),
        Database: utils.GetEnv("DB_URL", "sqlite://./data/app.db"),
        Redis: utils.RedisConfig{
            Host:     utils.GetEnv("REDIS_HOST", "localhost"),
            Port:     utils.GetEnvAsInt("REDIS_PORT", 6379),
            Password: utils.GetEnv("REDIS_PASSWORD", ""),
        },
        Debug: utils.GetEnvAsBool("DEBUG", false),
    }
}
```

### 请求追踪

```go
func loggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 生成或获取请求 ID
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = utils.GenerateRequestID()
        }

        // 设置到 Context
        utils.GinSetRequestID(c, requestID)
        c.Writer.Header().Set("X-Request-ID", requestID)

        // 记录日志
        start := time.Now()
        c.Next()

        logger.Info("HTTP request",
            zap.String("request_id", requestID),
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("latency", time.Since(start)),
        )
    }
}
```

### 数据验证

```go
func validateUser(user *User) error {
    if !utils.IsEmail(user.Email) {
        return errors.New("invalid email format")
    }

    if !utils.IsPhone(user.Phone) {
        return errors.New("invalid phone format")
    }

    if utils.IsEmpty(user.Name) {
        return errors.New("name is required")
    }

    if !utils.InRange(len(user.Password), 6, 32) {
        return errors.New("password must be 6-32 characters")
    }

    return nil
}
```

### 文件处理

```go
func saveUserAvatar(userID string, file io.Reader) error {
    // 生成文件名
    ext := utils.ExtName(file.Filename)
    filename := fmt.Sprintf("%s%s", utils.GenerateRequestID(), ext)

    // 创建目录
    dir := fmt.Sprintf("./uploads/avatars/%s", userID)
    if err := utils.CreateDir(dir); err != nil {
        return err
    }

    // 保存文件
    path := utils.JoinPath(dir, filename)
    return utils.WriteFile(path, file)
}
```

### 安全处理

```go
func hashPassword(password string) string {
    salt := utils.RandomString(16)
    return utils.MD5Salt(password, salt)
}

func generateToken(userID string) string {
    data := fmt.Sprintf("%s:%s:%d", userID, utils.RandomString(32), utils.Timestamp())
    return utils.SHA256(data)
}
```

## 最佳实践

1. **错误处理**
   - 所有函数返回 error
   - 记录详细错误信息
   - 优雅降级

2. **性能优化**
   - 避免重复创建对象
   - 使用对象池
   - 批量处理

3. **安全考虑**
   - 不记录敏感信息
   - 使用安全的随机数
   - 加密存储敏感数据

4. **代码复用**
   - 统一工具函数
   - 避免重复代码
   - 保持函数简洁

## 测试

```go
func TestEnv(t *testing.T) {
    os.Setenv("TEST_VAR", "test_value")
    defer os.Unsetenv("TEST_VAR")

    value := utils.GetEnv("TEST_VAR", "default")
    assert.Equal(t, "test_value", value)
}

func TestRequestID(t *testing.T) {
    id1 := utils.GenerateRequestID()
    id2 := utils.GenerateRequestID()

    assert.NotEmpty(t, id1)
    assert.NotEmpty(t, id2)
    assert.NotEqual(t, id1, id2)

    ctx := utils.SetRequestID(context.Background(), id1)
    assert.Equal(t, id1, utils.GetRequestID(ctx))
}
```

## 相关文档

- [配置模块](../../internal/config/README.md)
- [日志模块](../../internal/logger/README.md)
- [项目总览](../../README.md)

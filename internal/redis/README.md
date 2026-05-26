# Redis 缓存模块

高性能 Redis 缓存客户端，支持连接池管理和分布式缓存。

## 功能特性

- 🔗 **连接池**: 高性能连接池管理
- 🔄 **自动重连**: 网络异常自动重连
- 📊 **连接监控**: 实时连接状态监控
- ⏰ **超时控制**: 请求超时设置
- 🔒 **密码认证**: 支持 Redis 密码认证
- 🗃️ **多数据库**: 支持 Redis 多数据库
- 🎯 **Pipeline**: 支持管道批量操作
- 📈 **发布订阅**: 支持 Pub/Sub 模式

## 快速开始

### 初始化

```go
import "github.com/lookingcamel/system-framework/internal/redis"

func main() {
    err := redis.Init(cfg.Redis)
    if err != nil {
        log.Fatal(err)
    }
    defer redis.Close()
}
```

### 基本操作

```go
import "github.com/lookingcamel/system-framework/internal/redis"

// 设置字符串
err := redis.Set("key", "value", 3600)
if err != nil {
    log.Fatal(err)
}

// 获取字符串
value, err := redis.Get("key")
if err != nil {
    log.Fatal(err)
}
fmt.Println(value)

// 删除键
err := redis.Del("key")

// 设置 Hash
err := redis.HSet("user:1", map[string]interface{}{
    "name":  "john",
    "email": "john@example.com",
    "age":   30,
})

// 获取 Hash
user, err := redis.HGetAll("user:1")
```

## 配置说明

```yaml
redis:
  enabled: true
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10
  min_idle_conns: 5
```

### 配置参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `host` | Redis 服务器地址 | localhost |
| `port` | Redis 端口 | 6379 |
| `password` | Redis 密码 | "" |
| `db` | Redis 数据库编号 | 0 |
| `pool_size` | 连接池大小 | 10 |
| `min_idle_conns` | 最小空闲连接数 | 5 |

## String 操作

### 基本操作

```go
// 设置值
redis.Set("name", "john", 3600)

// 获取值
value, _ := redis.Get("name")

// 设置多个值
redis.MSet(map[string]string{
    "name": "john",
    "city": "beijing",
})

// 获取多个值
values, _ := redis.MGet("name", "city")

// 删除键
redis.Del("key")

// 设置过期时间
redis.Expire("key", 3600)

// 获取剩余时间
ttl, _ := redis.TTL("key")
```

### 计数操作

```go
// 递增
redis.Incr("counter")
redis.IncrBy("counter", 10)

// 递减
redis.Decr("counter")
redis.DecrBy("counter", 5)

// 自增浮点数
redis.IncrByFloat("price", 0.5)
```

## Hash 操作

```go
// 设置 Hash 字段
redis.HSet("user:1", "name", "john")

// 获取 Hash 字段
name, _ := redis.HGet("user:1", "name")

// 获取所有字段
user, _ := redis.HGetAll("user:1")

// 设置多个字段
redis.HMSet("user:1", map[string]interface{}{
    "name":  "john",
    "email": "john@example.com",
})

// 获取多个字段
fields, _ := redis.HMGet("user:1", "name", "email")

// 删除字段
redis.HDel("user:1", "name")

// 字段是否存在
exists, _ := redis.HExists("user:1", "name")

// 字段递增
redis.HIncrBy("user:1", "age", 1)
```

## List 操作

```go
// 左边入队
redis.LPush("queue", "task1")
redis.LPush("queue", "task2", "task3")

// 右边入队
redis.RPush("queue", "task4")

// 左边出队
task, _ := redis.LPop("queue")

// 右边出队
task, _ := redis.RPop("queue")

// 获取列表长度
length, _ := redis.LLen("queue")

// 获取列表范围
tasks, _ := redis.LRange("queue", 0, 9)

// 裁剪列表
redis.LTrim("queue", 0, 99)
```

## Set 操作

```go
// 添加成员
redis.SAdd("tags", "golang", "redis", "cache")

// 获取所有成员
tags, _ := redis.SMembers("tags")

// 成员数量
count, _ := redis.SCard("tags")

// 成员是否存在
exists, _ := redis.SIsMember("tags", "golang")

// 随机获取成员
member, _ := redis.SRandMember("tags")

// 删除成员
redis.SRem("tags", "cache")
```

## Sorted Set 操作

```go
// 添加成员
redis.ZAdd("leaderboard", map[string]float64{
    "john":  100,
    "alice": 200,
    "bob":   150,
})

// 获取成员分数
score, _ := redis.ZScore("leaderboard", "john")

// 获取排名
rank, _ := redis.ZRank("leaderboard", "john")

// 获取逆序排名
revRank, _ := redis.ZRevRank("leaderboard", "john")

// 获取分数范围内的成员
members, _ := redis.ZRangeByScore("leaderboard", "100", "200")

// 获取逆序分数范围内的成员
members, _ := redis.ZRevRangeByScore("leaderboard", "200", "100")

// 增加分数
redis.ZIncrBy("leaderboard", 50, "john")
```

## Pub/Sub 发布订阅

```go
// 订阅消息
pubsub := redis.Subscribe("channel1", "channel2")
defer pubsub.Close()

ch := pubsub.Channel()
for msg := range ch {
    fmt.Printf("Received: %s from %s\n", msg.Payload, msg.Channel)
}

// 发布消息
redis.Publish("channel1", "Hello, Redis!")
```

## Pipeline 管道

```go
pipe := redis.Pipeline()

pipe.Set("key1", "value1", 0)
pipe.Get("key2")
pipe.HSet("key3", "field1", "value3")

results, err := pipe.Exec()
```

## 分布式锁

```go
// 获取锁
lock, err := redis.AcquireLock("resource:1", 10*time.Second)
if err != nil {
    return errors.New("failed to acquire lock")
}
defer lock.Release()

// 业务逻辑
doSomething()
```

## 缓存模式

### Cache-Aside

```go
func GetUser(id string) (*User, error) {
    // 1. 检查缓存
    cached, _ := redis.Get("user:" + id)
    if cached != "" {
        return unmarshalUser(cached)
    }

    // 2. 查询数据库
    user, err := db.GetUser(id)
    if err != nil {
        return nil, err
    }

    // 3. 更新缓存
    redis.Set("user:"+id, marshalUser(user), 3600)

    return user, nil
}
```

### Write-Through

```go
func CreateUser(user *User) error {
    // 1. 写入数据库
    err := db.CreateUser(user)
    if err != nil {
        return err
    }

    // 2. 写入缓存
    redis.Set("user:"+user.ID, marshalUser(user), 3600)

    return nil
}
```

### Write-Behind

```go
// 使用 Channel 异步写入
go func() {
    for user := range userChan {
        redis.Set("user:"+user.ID, marshalUser(user), 3600)
    }
}()
```

## 连接池监控

```go
// 获取连接池统计
stats := redis.PoolStats()
fmt.Printf("Hits: %d, Misses: %d\n", stats.Hits, stats.Misses)
fmt.Printf("Timeouts: %d, TotalConns: %d\n", stats.Timeouts, stats.TotalConns)
fmt.Printf("IdleConns: %d, StaleConns: %d\n", stats.IdleConns, stats.StaleConns)
```

## 最佳实践

1. **键命名规范**
   ```
   user:{user_id}:profile
   session:{session_id}:data
   cache:product:{product_id}
   lock:resource:{resource_id}
   ```

2. **过期时间设置**
   - 热点数据: 5-30 分钟
   - 会话数据: 根据会话超时
   - 配置缓存: 1-24 小时

3. **内存管理**
   - 定期清理过期键
   - 使用压缩减少内存
   - 合理设置 maxmemory

4. **性能优化**
   - 使用 Pipeline 批量操作
   - 避免大 Key
   - 使用连接池

## 故障排查

### 常见问题

**Q: 连接被拒绝**
```
A: 检查 Redis 服务状态，确认网络连通性
```

**Q: 缓存穿透**
```
A: 使用布隆过滤器或设置短过期时间
```

**Q: 缓存雪崩**
```
A: 随机过期时间，使用多级缓存
```

## 相关文档

- [数据库模块](../database/README.md)
- [缓存策略](../cache/README.md)
- [项目总览](../../README.md)

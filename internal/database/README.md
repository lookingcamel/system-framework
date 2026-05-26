# 数据库模块

支持 MySQL、PostgreSQL、SQLite 的多数据库抽象层，提供统一的数据库操作接口。

## 功能特性

- 🗄️ **多数据库支持**: MySQL、PostgreSQL、SQLite
- 🔄 **连接池管理**: 自动管理连接池
- 🛡️ **熔断保护**: 集成熔断器防止雪崩
- 🔄 **自动迁移**: 数据库版本化迁移
- 💾 **自动备份**: 定时数据库备份
- 📊 **监控指标**: 数据库操作监控
- 🔒 **安全连接**: 支持 TLS/SSL
- ⏰ **超时控制**: 请求超时设置

## 快速开始

### 初始化

```go
import "github.com/lookingcamel/system-framework/internal/database"

func main() {
    err := database.Init(cfg.Database)
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()
}
```

### 执行查询

```go
import "github.com/lookingcamel/system-framework/internal/database"

// 查询单行
var name string
err := database.QueryRow("SELECT name FROM users WHERE id = ?", 1).Scan(&name)

// 查询多行
rows, err := database.Query("SELECT id, name FROM users WHERE status = ?", "active")
defer rows.Close()

for rows.Next() {
    var id int
    var name string
    rows.Scan(&id, &name)
    fmt.Printf("User: %d - %s\n", id, name)
}

// 执行语句
result, err := database.Exec("UPDATE users SET status = ? WHERE id = ?", "inactive", 1)
affected, _ := result.RowsAffected()
```

## 数据库配置

### MySQL 配置

```yaml
database:
  type: "mysql"
  host: "localhost"
  port: 3306
  name: "app"
  username: "admin"
  password: "password"
  max_open_connections: 10
  max_idle_connections: 5
  connection_max_lifetime: 300
  migration_path: "./migrations"
  backup_enabled: true
  backup_dir: "./backups"
```

### PostgreSQL 配置

```yaml
database:
  type: "postgres"
  host: "localhost"
  port: 5432
  name: "app"
  username: "admin"
  password: "password"
  max_open_connections: 10
  max_idle_connections: 5
  connection_max_lifetime: 300
```

### SQLite 配置

```yaml
database:
  type: "sqlite"
  sqlite_path: "./data/app.db"
  max_open_connections: 10
  max_idle_connections: 5
  connection_max_lifetime: 300
```

## 连接池配置

```yaml
database:
  max_open_connections: 10      # 最大打开连接数
  max_idle_connections: 5       # 最大空闲连接数
  connection_max_lifetime: 300  # 连接最大生存时间（秒）
```

### 连接池调优

```go
// 根据 QPS 调整连接数
// 公式: max_open_connections = (核心数 * 2) + 磁盘数

// 高并发场景
max_open_connections: 50
max_idle_connections: 25

// 低并发场景
max_open_connections: 10
max_idle_connections: 5
```

## 数据库操作

### 事务处理

```go
// 开始事务
tx, err := database.Begin()
if err != nil {
    return err
}
defer tx.Rollback()

// 执行事务操作
_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromID)
if err != nil {
    return err
}

_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toID)
if err != nil {
    return err
}

// 提交事务
return tx.Commit()
```

### 批量操作

```go
// 批量插入
stmt, _ := database.Prepare("INSERT INTO users (name, email) VALUES (?, ?)")
defer stmt.Close()

for _, user := range users {
    _, err = stmt.Exec(user.Name, user.Email)
}
```

### 错误处理

```go
err := database.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)
if err == sql.ErrNoRows {
    return nil, nil  // 未找到记录
}
if err != nil {
    return nil, err  // 其他错误
}
```

## SQL 方言差异

### MySQL

```sql
-- 分页
SELECT * FROM users LIMIT 10 OFFSET 20;

-- 模糊查询
SELECT * FROM users WHERE name LIKE '%john%';

-- 插入或更新
INSERT INTO users (id, name) VALUES (1, 'john') ON DUPLICATE KEY UPDATE name='john';
```

### PostgreSQL

```sql
-- 分页
SELECT * FROM users LIMIT 10 OFFSET 20;

-- 模糊查询
SELECT * FROM users WHERE name ILIKE '%john%';

-- 插入或更新
INSERT INTO users (id, name) VALUES (1, 'john') ON CONFLICT (id) DO UPDATE SET name='john';
```

### SQLite

```sql
-- 分页
SELECT * FROM users LIMIT 10 OFFSET 20;

-- 模糊查询
SELECT * FROM users WHERE name LIKE '%john%';

-- 插入或更新
INSERT OR REPLACE INTO users (id, name) VALUES (1, 'john');
```

## 数据库迁移

### 迁移文件命名

```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_add_email_column.up.sql
└── 000002_add_email_column.down.sql
```

### 创建迁移

```sql
-- 000003_create_orders_table.up.sql
CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    total DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
```

```sql
-- 000003_create_orders_table.down.sql
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_user_id;
DROP TABLE IF EXISTS orders;
```

### 运行迁移

```go
import "github.com/lookingcamel/system-framework/internal/migrate"

err := migrate.RunMigrations(cfg.Database)
if err != nil {
    log.Fatal(err)
}
```

## 数据库备份

### 自动备份配置

```yaml
database:
  backup_enabled: true
  backup_dir: "./backups"
```

### 手动备份

```go
import "github.com/lookingcamel/system-framework/internal/backup"

err := backup.BackupDatabase(cfg.Database, "manual_backup")
if err != nil {
    log.Fatal(err)
}
```

### 恢复备份

```go
err := backup.RestoreDatabase(cfg.Database, "backup_20240101_120000.sql")
if err != nil {
    log.Fatal(err)
}
```

## 监控指标

### Prometheus 指标

```
# 数据库连接池指标
database_connections_open{service="app"} 10
database_connections_in_use{service="app"} 5
database_connections_idle{service="app"} 5

# 查询延迟指标
database_query_duration_seconds{service="app", query_type="select"} 0.023
database_query_duration_seconds{service="app", query_type="insert"} 0.045
database_query_duration_seconds{service="app", query_type="update"} 0.034
```

## 最佳实践

1. **连接池配置**
   - 根据 QPS 调整连接数
   - 设置合理的连接超时
   - 监控连接泄漏

2. **查询优化**
   - 使用参数化查询防止 SQL 注入
   - 合理使用索引
   - 避免 SELECT *

3. **事务管理**
   - 保持事务简短
   - 合理设置超时时间
   - 避免嵌套事务

4. **错误处理**
   - 区分不同类型错误
   - 记录详细错误信息
   - 优雅处理断开连接

## 故障排查

### 常见问题

**Q: 连接池耗尽**
```
A: 检查是否有长事务，增加 max_open_connections
```

**Q: 数据库连接超时**
```
A: 检查网络连通性，增加 connection_max_lifetime
```

**Q: 迁移失败**
```
A: 检查迁移文件语法，确认数据库权限
```

## 相关文档

- [迁移模块](../migrate/README.md)
- [备份模块](../backup/README.md)
- [熔断器模块](../circuitbreaker/README.md)
- [项目总览](../../README.md)

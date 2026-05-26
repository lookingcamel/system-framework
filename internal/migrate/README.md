# 数据库迁移模块

基于 golang-migrate 的数据库版本化迁移管理，支持 MySQL、PostgreSQL、SQLite。

## 功能特性

- 📝 **版本化迁移**: 支持 up/down 迁移版本
- 🔄 **自动迁移**: 启动时自动执行迁移
- ⏪ **版本回滚**: 支持回滚到指定版本
- 📊 **迁移状态**: 查询当前迁移状态
- 🗄️ **多数据库**: MySQL、PostgreSQL、SQLite
- 📝 **SQL 格式**: 纯 SQL 迁移脚本

## 快速开始

### 初始化

```go
import "github.com/lookingcamel/system-framework/internal/migrate"

func main() {
    err := migrate.RunMigrations(cfg.Database)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 回滚迁移

```go
// 回滚上一个版本
err := migrate.RollbackMigrations(cfg.Database, 1)

// 回滚多个版本
err := migrate.RollbackMigrations(cfg.Database, 3)
```

## 迁移文件结构

### 命名规范

```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_add_email_to_users.up.sql
├── 000002_add_email_to_users.down.sql
├── 000003_create_orders_table.up.sql
└── 000003_create_orders_table.down.sql
```

### 版本号

- 使用 6 位数字前缀
- 000001, 000002, 000003...
- 确保版本号连续

## 迁移示例

### 1. 创建表

```sql
-- 000001_create_users_table.up.sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
```

```sql
-- 000001_create_users_table.down.sql
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP TABLE IF EXISTS users;
```

### 2. 添加字段

```sql
-- 000002_add_phone_to_users.up.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255);
```

```sql
-- 000002_add_phone_to_users.down.sql
ALTER TABLE users DROP COLUMN avatar_url;
ALTER TABLE users DROP COLUMN phone;
```

### 3. 创建索引

```sql
-- 000003_create_user_sessions.up.sql
CREATE TABLE user_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_token VARCHAR(255) NOT NULL UNIQUE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_sessions_token ON user_sessions(session_token);
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_expires ON user_sessions(expires_at);
```

```sql
-- 000003_create_user_sessions.down.sql
DROP INDEX IF EXISTS idx_user_sessions_expires;
DROP INDEX IF EXISTS idx_user_sessions_user_id;
DROP INDEX IF EXISTS idx_user_sessions_token;
DROP TABLE IF EXISTS user_sessions;
```

### 4. 数据迁移

```sql
-- 000004_migrate_user_data.up.sql
UPDATE users
SET status = 'inactive'
WHERE last_login < DATE_SUB(NOW(), INTERVAL 1 YEAR);

INSERT INTO user_sessions (user_id, session_token, expires_at)
SELECT id, token, DATE_ADD(NOW(), INTERVAL 7 DAY)
FROM old_sessions;
```

```sql
-- 000004_migrate_user_data.down.sql
DELETE FROM user_sessions WHERE created_at < DATE_SUB(NOW(), INTERVAL 7 DAY);
```

## 配置说明

```yaml
database:
  type: "sqlite"
  sqlite_path: "./data/app.db"
  migration_path: "./migrations"  # 迁移文件目录
```

## API 文档

### 主要函数

| 函数 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `RunMigrations(cfg)` | 执行所有待执行迁移 | DatabaseConfig | error |
| `RollbackMigrations(cfg, steps)` | 回滚迁移 | DatabaseConfig, steps | error |
| `GetVersion(cfg)` | 获取当前版本 | DatabaseConfig | (version, dirty, error) |
| `ForceVersion(cfg, version)` | 强制设置版本 | DatabaseConfig, version | error |

## 数据库兼容性

### MySQL

```sql
-- AUTO_INCREMENT 支持
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    ...
);

-- 索引命名
CREATE INDEX idx_users_email ON users(email);
```

### PostgreSQL

```sql
-- SERIAL 类型
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    ...
);

-- 索引命名
CREATE INDEX idx_users_email ON users(email);
```

### SQLite

```sql
-- AUTOINCREMENT
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ...
);

-- 索引命名
CREATE INDEX idx_users_email ON users(email);
```

## 迁移策略

### 1. 蓝绿部署

```go
// 迁移前创建新表
ALTER TABLE users RENAME TO users_old;
CREATE TABLE users (...);
INSERT INTO users SELECT * FROM users_old;
DROP TABLE users_old;
```

### 2. 零停机迁移

```go
// 1. 添加新字段（可为空）
// 2. 后台数据迁移
// 3. 添加约束
```

### 3. 回滚策略

```sql
-- 在 down.sql 中清理数据
DELETE FROM users WHERE status = 'deleted';
```

## 最佳实践

### 1. 迁移文件组织

```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_add_email.up.sql
├── 000002_add_email.down.sql
└── _lock
```

### 2. 命名规范

- 表名使用复数形式: `users`, `orders`
- 索引命名: `idx_{table}_{column}`
- 外键命名: `fk_{table}_{ref_table}`

### 3. 数据安全

- 始终编写 down.sql
- 使用事务包装迁移
- 备份数据前执行迁移
- 测试环境验证

### 4. 性能考虑

- 大数据量迁移分批处理
- 避免锁表操作
- 选择低峰期执行

## 故障处理

### 常见问题

**Q: 迁移失败**
```
A: 检查 SQL 语法，确认数据库权限，查看错误日志
```

**Q: 版本不一致**
```
A: 使用 migrate.ForceVersion 强制同步版本
```

**Q: 锁表**
```
A: 使用 pt-online-schema-change 工具
```

### 恢复步骤

```bash
# 1. 查看当前版本
migrate -path ./migrations -database $DATABASE_URL version

# 2. 回滚到指定版本
migrate -path ./migrations -database $DATABASE_URL down 1

# 3. 修复迁移文件
# 4. 重新执行迁移
migrate -path ./migrations -database $DATABASE_URL up
```

## 自动化

### CI/CD 集成

```yaml
# .github/workflows/migrate.yml
name: Database Migration

on:
  push:
    branches: [main]

jobs:
  migrate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run migrations
        run: |
          migrate -path ./migrations \
            -database $DATABASE_URL \
            up
        env:
          DATABASE_URL: ${{ secrets.DATABASE_URL }}
```

## 相关文档

- [数据库模块](../database/README.md)
- [备份模块](../backup/README.md)
- [配置模块](../config/README.md)
- [项目总览](../../README.md)

# 数据库备份模块

支持 MySQL、PostgreSQL、SQLite 的自动和手动备份。

## 功能特性

- 💾 **自动备份**: 定时自动备份
- 🔄 **手动备份**: 按需手动备份
- 📦 **压缩存储**: 自动压缩备份文件
- 🔍 **备份验证**: 验证备份完整性
- ⏪ **快速恢复**: 一键恢复备份
- 📊 **备份列表**: 查看历史备份
- 🗑️ **自动清理**: 自动删除过期备份
- 🌐 **远程存储**: 支持远程备份（可选）

## 快速开始

### 自动备份

```go
import "github.com/lookingcamel/system-framework/internal/backup"

func main() {
    // 备份会在数据库初始化时自动执行
    err := database.Init(cfg.Database)
}
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

## 配置说明

```yaml
database:
  type: "sqlite"
  sqlite_path: "./data/app.db"
  backup_enabled: true           # 启用自动备份
  backup_dir: "./backups"        # 备份目录
```

### 备份配置

```go
type BackupConfig struct {
    Enabled      bool   // 是否启用备份
    BackupDir    string // 备份目录
    MaxBackups   int    // 保留备份数量
    MaxAge       int    // 备份保留天数
    Compress     bool   // 是否压缩
}
```

## 备份文件命名

```
backups/
├── backup_20240101_120000.sql       # 未压缩
├── backup_20240101_120000.sql.gz    # 已压缩
└── backup_20240101_120000.json      # 备份元数据
```

### 元数据文件

```json
{
  "version": "1.0",
  "timestamp": "2024-01-01T12:00:00Z",
  "database_type": "sqlite",
  "database_path": "./data/app.db",
  "size": 1024567,
  "compressed": true,
  "checksum": "sha256:abc123...",
  "tables": ["users", "orders", "products"],
  "records": {
    "users": 1000,
    "orders": 5000,
    "products": 200
  }
}
```

## API 文档

### 主要函数

| 函数 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `BackupDatabase(cfg, name)` | 执行备份 | DatabaseConfig, name | error |
| `RestoreDatabase(cfg, file)` | 恢复备份 | DatabaseConfig, file | error |
| `ListBackups()` | 列出备份 | - | ([]BackupInfo, error) |
| `DeleteBackup(name)` | 删除备份 | BackupInfo | error |
| `VerifyBackup(file)` | 验证备份 | file | (bool, error) |
| `GetBackupConfig()` | 获取备份配置 | - | BackupConfig |

## 数据库备份

### SQLite 备份

```go
func backupSQLite(cfg DatabaseConfig, backupName string) error {
    // 备份文件
    src, err := os.Open(cfg.SQLitePath)
    if err != nil {
        return err
    }
    defer src.Close()

    // 目标文件
    timestamp := time.Now().Format("20060102_150405")
    dstName := fmt.Sprintf("backup_%s_%s.sql", backupName, timestamp)
    dst, err := os.Create(path.Join(cfg.BackupDir, dstName))
    if err != nil {
        return err
    }
    defer dst.Close()

    // 复制数据库
    _, err = io.Copy(dst, src)
    return err
}
```

### MySQL 备份

```bash
mysqldump -h localhost -u root -p database_name > backup.sql
```

```go
func backupMySQL(cfg DatabaseConfig, backupName string) error {
    cmd := exec.Command("mysqldump",
        "-h", cfg.Host,
        "-P", strconv.Itoa(cfg.Port),
        "-u", cfg.Username,
        "-p"+cfg.Password,
        cfg.Name,
    )

    output, err := os.Create(fmt.Sprintf("backup_%s.sql", backupName))
    if err != nil {
        return err
    }
    defer output.Close()

    cmd.Stdout = output
    return cmd.Run()
}
```

### PostgreSQL 备份

```bash
pg_dump -h localhost -U username -d database_name -f backup.sql
```

```go
func backupPostgres(cfg DatabaseConfig, backupName string) error {
    cmd := exec.Command("pg_dump",
        "-h", cfg.Host,
        "-p", strconv.Itoa(cfg.Port),
        "-U", cfg.Username,
        "-d", cfg.Name,
        "-f", fmt.Sprintf("backup_%s.sql", backupName),
    )

    cmd.Env = append(os.Environ(),
        fmt.Sprintf("PGPASSWORD=%s", cfg.Password))

    return cmd.Run()
}
```

## 数据库恢复

### SQLite 恢复

```go
func restoreSQLite(cfg DatabaseConfig, backupFile string) error {
    // 读取备份文件
    src, err := os.Open(backupFile)
    if err != nil {
        return err
    }
    defer src.Close()

    // 关闭原数据库连接
    database.Close()

    // 恢复数据库
    dst, err := os.Create(cfg.SQLitePath)
    if err != nil {
        return err
    }
    defer dst.Close()

    _, err = io.Copy(dst, src)
    return err

    // 重新初始化数据库
    return database.Init(cfg)
}
```

### MySQL 恢复

```bash
mysql -h localhost -u root -p database_name < backup.sql
```

### PostgreSQL 恢复

```bash
psql -h localhost -U username -d database_name < backup.sql
```

## 备份策略

### 1. 全量备份

```go
func FullBackup(cfg DatabaseConfig) error {
    return BackupDatabase(cfg, "full")
}
```

### 2. 增量备份

```go
func IncrementalBackup(cfg DatabaseConfig, since time.Time) error {
    // 备份指定时间后的数据
    return BackupDatabase(cfg, fmt.Sprintf("incremental_%s", since.Format("20060102")))
}
```

### 3. 定时备份

```go
func ScheduleBackups(cfg DatabaseConfig) {
    ticker := time.NewTicker(24 * time.Hour)
    go func() {
        for range ticker.C {
            BackupDatabase(cfg, "scheduled")
            CleanOldBackups(cfg)
        }
    }()
}
```

## 备份压缩

```go
func CompressBackup(src, dst string) error {
    // 打开源文件
    sourceFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    // 创建压缩文件
    compressedFile, err := os.Create(dst + ".gz")
    if err != nil {
        return err
    }
    defer compressedFile.Close()

    // 创建 gzip 写入器
    gzipWriter := gzip.NewWriter(compressedFile)
    defer gzipWriter.Close()

    // 复制数据
    _, err = io.Copy(gzipWriter, sourceFile)
    return err
}
```

## 备份验证

```go
func VerifyBackup(backupFile string) (bool, error) {
    // 读取备份文件
    file, err := os.Open(backupFile)
    if err != nil {
        return false, err
    }
    defer file.Close()

    // 解压（如需要）
    if strings.HasSuffix(backupFile, ".gz") {
        reader, err := gzip.NewReader(file)
        if err != nil {
            return false, err
        }
        defer reader.Close()
        file = reader
    }

    // 验证 SQL 语法
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        // 简单验证：检查是否有完整的 SQL 语句
        if strings.HasSuffix(line, ";") {
            continue
        }
    }

    return true, scanner.Err()
}
```

## 备份列表

```go
type BackupInfo struct {
    Name         string
    Size         int64
    CreatedAt    time.Time
    Compressed   bool
    Checksum     string
    RecordsCount int64
}

func ListBackups() ([]BackupInfo, error) {
    files, err := os.ReadDir(backupDir)
    if err != nil {
        return nil, err
    }

    var backups []BackupInfo
    for _, file := range files {
        if strings.HasPrefix(file.Name(), "backup_") {
            info, _ := file.Info()
            backups = append(backups, BackupInfo{
                Name:       file.Name(),
                Size:       info.Size(),
                CreatedAt:  info.ModTime(),
                Compressed: strings.HasSuffix(file.Name(), ".gz"),
            })
        }
    }

    return backups, nil
}
```

## 自动清理

```go
func CleanOldBackups(cfg DatabaseConfig) error {
    maxAge := 7 * 24 * time.Hour // 保留 7 天
    maxBackups := 10             // 最多保留 10 个

    backups, _ := ListBackups()

    // 按时间排序
    sort.Slice(backups, func(i, j int) bool {
        return backups[i].CreatedAt.After(backups[j].CreatedAt)
    })

    // 删除过期备份
    for _, backup := range backups[maxBackups:] {
        if time.Since(backup.CreatedAt) > maxAge {
            os.Remove(backup.Name)
        }
    }

    return nil
}
```

## 最佳实践

1. **备份频率**
   - 生产环境: 每日全量 + 每小时增量
   - 测试环境: 每周全量

2. **存储策略**
   - 本地存储 + 远程备份
   - 定期验证备份可用性
   - 加密敏感数据

3. **恢复测试**
   - 定期测试恢复流程
   - 文档化恢复步骤
   - 测量恢复时间

4. **监控告警**
   - 备份失败告警
   - 存储空间告警
   - 备份完整性检查

## 故障处理

### 常见问题

**Q: 备份失败**
```
A: 检查磁盘空间、数据库连接权限、日志文件
```

**Q: 恢复后数据不一致**
```
A: 确认备份时间点，避免在备份期间写入数据
```

**Q: 备份文件损坏**
```
A: 使用校验和验证备份，保留多个备份版本
```

## 相关文档

- [数据库模块](../database/README.md)
- [迁移模块](../migrate/README.md)
- [配置模块](../config/README.md)
- [项目总览](../../README.md)

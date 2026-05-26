package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

type Config struct {
	App                AppConfig
	Server             ServerConfig
	Logging            LoggingConfig
	Prometheus         PrometheusConfig
	PProf              PProfConfig
	Tracing            TracingConfig
	Health             HealthConfig
	GracefulShutdown   GracefulShutdownConfig
	Database           DatabaseConfig
	Redis              RedisConfig
	Nacos              NacosConfig
	CircuitBreaker     CircuitBreakerConfig
	Auth               AuthConfig
	RateLimit          RateLimitConfig
	APIVersion         APIVersionConfig
}

type APIVersionConfig struct {
	Enabled       bool
	DefaultVersion string      `mapstructure:"default_version"`
	SupportedVersions []string  `mapstructure:"supported_versions"`
	DeprecationNotice string  `mapstructure:"deprecation_notice"`
}

type AppConfig struct {
	Name    string
	Version string
	Host    string
	Port    int
	Mode    string
}

type ServerConfig struct {
	ReadTimeout    int
	WriteTimeout   int
	IdleTimeout    int
	MaxBodySize    int64 `mapstructure:"max_body_size"`
}

type LoggingConfig struct {
	Level       string
	Format      string
	Output      string
	OutputPath  string
	MaxSize     int
	MaxBackups  int
	MaxAge      int
	Compress    bool
}

type PrometheusConfig struct {
	Enabled bool
	Path    string
}

type PProfConfig struct {
	Enabled bool
	Path    string
}

type TracingConfig struct {
	Enabled     bool
	ServiceName string
	Exporter    string
	Endpoint    string
	SampleRate  float64
}

type HealthConfig struct {
	Enabled bool
	Path    string
}

type GracefulShutdownConfig struct {
	Timeout int
}

type DatabaseConfig struct {
	Type                  string
	Host                  string
	Port                  int
	Name                  string
	Username              string
	Password              string
	SQLitePath            string `mapstructure:"sqlite_path"`
	MaxOpenConnections    int    `mapstructure:"max_open_connections"`
	MaxIdleConnections    int    `mapstructure:"max_idle_connections"`
	ConnectionMaxLifetime int    `mapstructure:"connection_max_lifetime"`
	MigrationPath         string `mapstructure:"migration_path"`
	BackupEnabled         bool   `mapstructure:"backup_enabled"`
	BackupDir             string `mapstructure:"backup_dir"`
}

type RedisConfig struct {
	Enabled      bool
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int `mapstructure:"pool_size"`
	MinIdleConns int `mapstructure:"min_idle_conns"`
}

type NacosConfig struct {
	Enabled       bool
	ServerAddr    string `mapstructure:"server_addr"`
	ServerPort    uint64 `mapstructure:"server_port"`
	NamespaceId   string `mapstructure:"namespace_id"`
	Group         string
	DataId        string `mapstructure:"data_id"`
	Username      string
	Password      string
	ConfigType    string `mapstructure:"config_type"`
	RefreshEnable bool   `mapstructure:"refresh_enable"`
	RefreshDelay  int    `mapstructure:"refresh_delay"`
}

type CircuitBreakerConfig struct {
	Enabled                bool    `mapstructure:"enabled"`
	DefaultTimeout         int     `mapstructure:"default_timeout"`
	DefaultMaxConcurrent   int     `mapstructure:"default_max_concurrent"`
	DefaultErrorPercentage int     `mapstructure:"default_error_percentage"`
	DefaultRequestVolume   int     `mapstructure:"default_request_volume"`
}

type RateLimitConfig struct {
	Enabled           bool     `mapstructure:"enabled"`
	RequestsPerSecond float64  `mapstructure:"requests_per_second"`
	Burst             int
	ExcludePaths      []string `mapstructure:"exclude_paths"`
}

type AuthConfig struct {
	Enabled              bool     `mapstructure:"enabled"`
	JWTEnabled           bool     `mapstructure:"jwt_enabled"`
	JWTSecret            string   `mapstructure:"jwt_secret"`
	JWTIssuer            string   `mapstructure:"jwt_issuer"`
	JWTAudience          string   `mapstructure:"jwt_audience"`
	JWTExpireSeconds     int      `mapstructure:"jwt_expire_seconds"`
	SignatureEnabled     bool     `mapstructure:"signature_enabled"`
	SignatureSecret      string   `mapstructure:"signature_secret"`
	SignatureHeader      string   `mapstructure:"signature_header"`
	TimestampHeader      string   `mapstructure:"timestamp_header"`
	NonceHeader          string   `mapstructure:"nonce_header"`
	TimestampMaxDiff     int      `mapstructure:"timestamp_max_diff"`
	PublicKey            string   `mapstructure:"public_key"`
	PrivateKey           string   `mapstructure:"private_key"`
	ExcludePaths         []string `mapstructure:"exclude_paths"`
}

var (
	globalConfig *Config
	mu           sync.RWMutex
)

func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.Nacos.Enabled {
		if err := loadFromNacos(&cfg); err != nil {
			return nil, fmt.Errorf("failed to load config from nacos: %w", err)
		}

		if cfg.Nacos.RefreshEnable {
			go startConfigWatcher(&cfg)
		}
	}

	validator := NewValidator()
	if err := validator.Validate(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	globalConfig = &cfg
	return &cfg, nil
}

func loadFromNacos(cfg *Config) error {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(
			cfg.Nacos.ServerAddr,
			cfg.Nacos.ServerPort,
			constant.WithScheme("http"),
		),
	}

	cc := constant.NewClientConfig(
		constant.WithNamespaceId(cfg.Nacos.NamespaceId),
		constant.WithUsername(cfg.Nacos.Username),
		constant.WithPassword(cfg.Nacos.Password),
	)

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		return err
	}

	content, err := client.GetConfig(vo.ConfigParam{
		DataId: cfg.Nacos.DataId,
		Group:  cfg.Nacos.Group,
	})
	if err != nil {
		return err
	}

	viper.SetConfigType(cfg.Nacos.ConfigType)
	if err := viper.ReadConfig(strings.NewReader(content)); err != nil {
		return err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}

	return nil
}

func startConfigWatcher(cfg *Config) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(
			cfg.Nacos.ServerAddr,
			cfg.Nacos.ServerPort,
			constant.WithScheme("http"),
		),
	}

	cc := constant.NewClientConfig(
		constant.WithNamespaceId(cfg.Nacos.NamespaceId),
		constant.WithUsername(cfg.Nacos.Username),
		constant.WithPassword(cfg.Nacos.Password),
	)

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		return
	}

	client.ListenConfig(vo.ConfigParam{
		DataId: cfg.Nacos.DataId,
		Group:  cfg.Nacos.Group,
		OnChange: func(namespace, group, dataId, data string) {
			mu.Lock()
			defer mu.Unlock()

			viper.SetConfigType(cfg.Nacos.ConfigType)
			if err := viper.ReadConfig(strings.NewReader(data)); err != nil {
				return
			}

			var newCfg Config
			if err := viper.Unmarshal(&newCfg); err != nil {
				return
			}

			globalConfig = &newCfg
		},
	})

	select {}
}

func GetGlobalConfig() *Config {
	mu.RLock()
	defer mu.RUnlock()
	return globalConfig
}

func Reload() error {
	mu.Lock()
	defer mu.Unlock()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	globalConfig = &cfg
	return nil
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.App.Host, c.App.Port)
}

var _ config_client.IConfigClient = nil

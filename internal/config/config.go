// Package config 定义配置结构并提供全局配置对象。
// 配置分两层：base（框架层：app/mysql/redis）+ custom_conf（业务层：JWT 等）。
package config

// AppConfig 应用基础配置。
type AppConfig struct {
	ServiceName string `yaml:"service_name" env:"ServiceName" env-default:"pickcampus-backend"`
	LocalIP     string `yaml:"local_ip" env:"LocalIP" env-default:"0.0.0.0"`
	APIPort     int    `yaml:"api_port" env:"APIPort" env-default:"8080"`
	RunMode     string `yaml:"run_mode" env:"RunMode" env-default:"debug"`
	Cors        string `yaml:"cors" env:"Cors" env-default:"1"`
	// CorsAllowOrigins 允许的跨域来源(逗号分隔)。空=回显任意 Origin(本地开发);
	// 非空=仅放行清单内来源(生产收敛到前端域名)。
	CorsAllowOrigins string `yaml:"cors_allow_origins" env:"CORS_ALLOW_ORIGINS"`
}

// MySQLConfig MySQL 连接配置。
type MySQLConfig struct {
	Host            string `yaml:"write_db_host" env:"WriteDBHost" env-default:"127.0.0.1"`
	Port            int    `yaml:"write_db_port" env:"WriteDBPort" env-default:"3306"`
	User            string `yaml:"write_db_user" env:"WriteDBUser" env-default:"root"`
	Password        string `yaml:"write_db_password" env:"WriteDBPassword"`
	DBName          string `yaml:"write_db" env:"WriteDB" env-default:"pickcampus"`
	TablePrefix     string `yaml:"table_prefix" env:"TablePrefix" env-default:"tbl_"`
	MaxIdleConns    int    `yaml:"max_idle_conns" env:"MaxIdleConns" env-default:"10"`
	MaxOpenConns    int    `yaml:"max_open_conns" env:"MaxOpenConns" env-default:"100"`
	ConnMaxLifeTime int64  `yaml:"conn_max_life_time" env:"ConnMaxLifeTime" env-default:"3600"`
	// SkipEnsureDB 为 true 时跳过启动时的 CREATE DATABASE IF NOT EXISTS,
	// 仅连接已存在的库(生产用,应用账号只需库级权限)。
	SkipEnsureDB bool `yaml:"skip_ensure_db" env:"SkipEnsureDB" env-default:"false"`
}

// RedisConfig Redis 连接配置。
type RedisConfig struct {
	Addr     string `yaml:"host_and_port" env:"RedisHostAndPort" env-default:"127.0.0.1:6379"`
	Username string `yaml:"username" env:"RedisUsername"`
	Password string `yaml:"password" env:"RedisPassword"`
	DB       int    `yaml:"db" env:"RedisDB" env-default:"0"`
	PoolSize int    `yaml:"pool_size" env:"RedisPoolSize" env-default:"20"`
}

// BaseConfig 框架层配置。
type BaseConfig struct {
	App   AppConfig   `yaml:"app"`
	MySQL MySQLConfig `yaml:"mysql"`
	Redis RedisConfig `yaml:"redis"`
}

// CustomConf 业务层配置。
type CustomConf struct {
	JWTTokenSecret  string `yaml:"jwt_token_secret" env:"JWTTokenSecret"`
	JWTTokenExpires int64  `yaml:"jwt_token_expires" env:"JWTTokenExpires" env-default:"86400"`
}

// Config 顶层配置。
type Config struct {
	Base BaseConfig `yaml:"base"`
	Conf CustomConf `yaml:"custom_conf"`
}

// G 全局配置对象。
var G = &Config{}

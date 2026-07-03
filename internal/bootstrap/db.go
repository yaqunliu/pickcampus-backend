// Package bootstrap 负责初始化外部依赖（MySQL、Redis），直接用官方库建客户端，
// 不引入 TC-Backend 那套自研 lib 重封装，保持轻量自包含。
package bootstrap

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"pickcampus-backend/internal/config"
)

// _db MySQL 单例。
var _db *gorm.DB

// InitDB 建立 MySQL 连接：先自动建库（IF NOT EXISTS），再连库并配置连接池与命名策略。
// SkipEnsureDB=true 时跳过建库，只连已存在的库（生产：库预建、账号仅库级权限）。
func InitDB(cfg config.MySQLConfig) error {
	// ① 用无库名 DSN 建库（utf8mb4）；显式跳过时直接连库
	if !cfg.SkipEnsureDB {
		if err := ensureDatabase(cfg); err != nil {
			return fmt.Errorf("建库失败: %w", err)
		}
	}

	// ② 连接目标库
	dsn := buildDSN(cfg, cfg.DBName)
	gormCfg := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.TablePrefix, // tbl_
			SingularTable: true,            // 结构体 User -> 表 tbl_user
		},
		DisableForeignKeyConstraintWhenMigrating: true,
		// 把驱动的唯一键冲突等错误翻译成 gorm 哨兵错误（gorm.ErrDuplicatedKey），
		// 使 repo 层能可靠识别邮箱竞态冲突。
		TranslateError: true,
	}
	db, err := gorm.Open(mysql.Open(dsn), gormCfg)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// ③ 连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层连接失败: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	if cfg.ConnMaxLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifeTime) * time.Second)
	}

	_db = db
	return nil
}

// ensureDatabase 用无库名 DSN 执行 CREATE DATABASE IF NOT EXISTS。
func ensureDatabase(cfg config.MySQLConfig) error {
	rootDSN := buildDSN(cfg, "")
	db, err := gorm.Open(mysql.Open(rootDSN), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer func() { _ = sqlDB.Close() }()

	createSQL := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		cfg.DBName,
	)
	return db.Exec(createSQL).Error
}

// buildDSN 拼接 MySQL DSN；dbName 为空则连到无库（用于建库）。
func buildDSN(cfg config.MySQLConfig, dbName string) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, dbName,
	)
}

// AutoMigrate 执行自动迁移。
func AutoMigrate(dst ...interface{}) error {
	if _db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return _db.AutoMigrate(dst...)
}

// GetDB 返回 MySQL 单例。
func GetDB() *gorm.DB {
	return _db
}

// Cli 返回带 context 的 *gorm.DB，供 logic/repo 使用。
func Cli(ctx context.Context) *gorm.DB {
	return _db.WithContext(ctx)
}

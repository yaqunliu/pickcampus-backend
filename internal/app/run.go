// Package app 是启动引导：读配置 → 初始化 DB/Redis → 迁移 → 建 Gin → 起服务。
package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/cobra"

	"pickcampus-backend/internal/bootstrap"
	"pickcampus-backend/internal/common"
	"pickcampus-backend/internal/config"
	"pickcampus-backend/internal/handler"
	"pickcampus-backend/models"
)

const projectName = "pickcampus-backend"

var configFile string

var rootCmd = &cobra.Command{
	Use:   projectName,
	Short: "PickCampus 后端服务",
	RunE: func(*cobra.Command, []string) error {
		return run()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "configs/config.yaml", "配置文件路径")
}

// Execute Cobra 入口。
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	// ① 读配置（YAML + env 覆盖）
	if err := cleanenv.ReadConfig(configFile, config.G); err != nil {
		return fmt.Errorf("读取配置失败: %w", err)
	}

	// ② 初始化 MySQL 并自动迁移
	if err := bootstrap.InitDB(config.G.Base.MySQL); err != nil {
		return err
	}
	if err := bootstrap.AutoMigrate(models.AllTables...); err != nil {
		return fmt.Errorf("自动迁移失败: %w", err)
	}

	// ③ 初始化 Redis
	if err := bootstrap.InitRedis(config.G.Base.Redis); err != nil {
		return err
	}

	// ④ 建 Gin
	if config.G.Base.App.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.New()
	g.Use(gin.Recovery())
	if config.G.Base.App.Cors == "1" {
		g.Use(common.Cors())
	}

	// ⑤ 注册路由
	handler.RegisterRouter(g)

	// ⑥ 起服务（优雅关闭）
	addr := fmt.Sprintf("%s:%d", config.G.Base.App.LocalIP, config.G.Base.App.APIPort)
	srv := &http.Server{Addr: addr, Handler: g}

	go func() {
		log.Printf("%s 启动，监听 %s", projectName, addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("收到关闭信号，正在优雅关闭...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("优雅关闭失败: %w", err)
	}
	log.Println("服务已关闭")
	return nil
}

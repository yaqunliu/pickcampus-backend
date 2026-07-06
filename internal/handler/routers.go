package handler

import (
	"github.com/gin-gonic/gin"

	"pickcampus-backend/internal/common"
)

// RegisterRouter 注册所有路由。
func RegisterRouter(g *gin.Engine) {
	// 健康检查
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "message": "pong"})
	})

	userHandler := NewUserHandler()
	candidateHandler := NewCandidateHandler()
	routeHandler := NewRouteHandler()

	v1 := g.Group("/api/v1")
	{
		// 公开接口
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.GET("/route", routeHandler.Query) // 离家路程查询(高德代理),无需登录

		// 受保护接口（JWT + Redis 会话校验）
		auth := v1.Group("")
		auth.Use(common.AuthMiddleware())
		{
			auth.POST("/logout", userHandler.Logout)
			auth.GET("/user", userHandler.GetUserInfo) // 拿当前用户（me）

			// 候选院校 / 候选专业
			auth.GET("/candidates", candidateHandler.List)
			auth.POST("/candidates", candidateHandler.AddSchool)
			auth.DELETE("/candidates/:school_id", candidateHandler.RemoveSchool)
			auth.POST("/candidates/:school_id/majors", candidateHandler.AddMajor)
			auth.DELETE("/candidates/:school_id/majors", candidateHandler.RemoveMajor)
		}
	}
}

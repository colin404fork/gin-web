package initialize

import (
	"fmt"
	"gin-web/api"
	"gin-web/middleware"
	"gin-web/pkg/global"
	"gin-web/router"
	"github.com/gin-gonic/gin"
)

// 初始化总路由
func Routers() *gin.Engine {
	// 创建带有默认中间件的路由:
	// 日志与恢复中间件
	// r := gin.Default()
	// 创建不带中间件的路由:
	r := gin.New()

	// 添加速率访问中间件
	r.Use(middleware.RateLimiter())
	// 添加访问记录
	r.Use(middleware.AccessLog)
	// 添加全局异常处理中间件
	r.Use(middleware.Exception)
	// 添加全局事务处理中间件
	r.Use(middleware.Transaction)
	// 添加跨域中间件, 让请求支持跨域
	r.Use(middleware.Cors())
	global.Log.Debug("请求已支持跨域")

	// 初始化jwt auth中间件
	authMiddleware, err := middleware.InitAuth()

	if err != nil {
		panic(fmt.Sprintf("初始化jwt auth中间件失败: %v", err))
	}
	global.Log.Debug("初始化jwt auth中间件完成")

	apiGroup := r.Group(global.Conf.System.UrlPathPrefix)
	// ping
	apiGroup.GET("/ping", api.Ping)

	// 方便统一添加路由前缀
	v1Group := apiGroup.Group("v1")
	router.InitPublicRouter(v1Group)                   // 注册公共路由
	router.InitBaseRouter(v1Group, authMiddleware)     // 注册基础路由, 不会鉴权
	router.InitUserRouter(v1Group, authMiddleware)     // 注册用户路由
	router.InitMenuRouter(v1Group, authMiddleware)     // 注册菜单路由
	router.InitRoleRouter(v1Group, authMiddleware)     // 注册角色路由
	router.InitApiRouter(v1Group, authMiddleware)      // 注册接口路由
	router.InitWorkflowRouter(v1Group, authMiddleware) // 注册工作流路由

	global.Log.Debug("初始化路由完成")
	return r
}

// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	rediscache "github.com/golanguzb70/redis-cache"
	"github.com/golanguzb70/udevslabs-twitter/config"
	_ "github.com/golanguzb70/udevslabs-twitter/docs"
	"github.com/golanguzb70/udevslabs-twitter/internal/controller/http/v1/handler"
	"github.com/golanguzb70/udevslabs-twitter/internal/usecase"
	"github.com/golanguzb70/udevslabs-twitter/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description This is a sample server Go Clean Template server.
// @version     1.0
// @BasePath    /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(engine *gin.Engine, l *logger.Logger, config *config.Config, useCase *usecase.UseCase, redis rediscache.RedisCache) {
	// Options
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	handlerV1 := handler.NewHandler(l, config, useCase, redis)

	// Swagger - Place this before AuthMiddleware
	url := ginSwagger.URL("swagger/doc.json") // The URL pointing to API definition
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Initialize Casbin enforcer
	e := casbin.NewEnforcer("config/rbac.conf", "config/policy.csv")
	engine.Use(handlerV1.AuthMiddleware(e)) // Apply authentication middleware to all routes except Swagger

	// K8s probe
	engine.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routes
	v1 := engine.Group("/v1") 
	{
		v1.POST("/user", handlerV1.CreateUser)
		v1.GET("/user/list", handlerV1.GetUsers)
		v1.GET("/user/:id", handlerV1.GetUser)
		v1.PUT("/user", handlerV1.UpdateUser)
		v1.DELETE("/user/:id", handlerV1.DeleteUser)

		v1.GET("/session/list", handlerV1.GetSessions)
		v1.GET("/session/:id", handlerV1.GetSession)
		v1.PUT("/session", handlerV1.UpdateSession)
		v1.DELETE("/session/:id", handlerV1.DeleteSession)

		v1.POST("/auth/logout", handlerV1.Logout)
		v1.POST("/auth/register", handlerV1.Register)
		v1.POST("/auth/verify-email", handlerV1.VerifyEmail)
		v1.POST("/auth/login", handlerV1.Login)

		v1.POST("/tag", handlerV1.CreateTag)
		v1.GET("/tag/list", handlerV1.GetTags)
		v1.GET("/tag/:id", handlerV1.GetTag)
		v1.PUT("/tag", handlerV1.UpdateTag)
		v1.DELETE("/tag/:id", handlerV1.DeleteTag)

		v1.POST("/follower", handlerV1.FollowUnfollow)
		v1.GET("/follower/list", handlerV1.GetFollowers)

		v1.POST("/tweet", handlerV1.CreateTweet)
		v1.GET("/tweet/list", handlerV1.GetTweets)
		v1.GET("/tweet/:id", handlerV1.GetTweet)
		v1.PUT("/tweet", handlerV1.UpdateTweet)
		v1.DELETE("/tweet/:id", handlerV1.DeleteTweet)

		
	}

	// user := v1.Group("/user")
	// {
	// 	user.POST("/user", handlerV1.CreateUser)
	// 	user.GET("/user/list", handlerV1.GetUsers)
	// 	user.GET("/user/:id", handlerV1.GetUser)
	// 	user.PUT("/user/", handlerV1.UpdateUser)
	// 	user.DELETE("/user/:id", handlerV1.DeleteUser)
	// }

	// session := v1.Group("/session")
	// {
	// 	session.GET("/session/list", handlerV1.GetSessions)
	// 	session.GET("/session/:id", handlerV1.GetSession)
	// 	session.PUT("/session", handlerV1.UpdateSession)
	// 	session.DELETE("/session/:id", handlerV1.DeleteSession)
	// }

	// auth := v1.Group("/auth")
	// {
	// 	auth.POST("/auth/logout", handlerV1.Logout)
	// 	auth.POST("/auth/register", handlerV1.Register)
	// 	auth.POST("/auth/verify-email", handlerV1.VerifyEmail)
	// 	auth.POST("/auth/login", handlerV1.Login)
	// }

	// tag := v1.Group("/tag")
	// {
	// 	tag.POST("/tag", handlerV1.CreateTag)
	// 	tag.GET("/tag/list", handlerV1.GetTags)
	// 	tag.GET("/tag/:id", handlerV1.GetTag)
	// 	tag.PUT("/tag", handlerV1.UpdateTag)
	// 	tag.DELETE("/tag/:id", handlerV1.DeleteTag)
	// }

	// follower := v1.Group("/follower")
	// {
	// 	follower.POST("/follower", handlerV1.FollowUnfollow)
	// 	follower.GET("/follower/list", handlerV1.GetFollowers)
	// }

	// tweet := v1.Group("/tweet")
	// {
	// 	tweet.POST("/tweet", handlerV1.CreateTweet)
	// 	tweet.GET("/tweet/list", handlerV1.GetTweets)
	// 	tweet.GET("/tweet/:id", handlerV1.GetTweet)
	// 	tweet.PUT("/tweet", handlerV1.UpdateTweet)
	// 	tweet.DELETE("/tweet/:id", handlerV1.DeleteTweet)
	// }
}

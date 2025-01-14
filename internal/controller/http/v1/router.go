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
// @host        localhost:8080
// @BasePath    /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(engine *gin.Engine, l *logger.Logger, config *config.Config, useCase *usecase.UseCase, redis rediscache.RedisCache) {
	// Options
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	handlerV1 := handler.NewHandler(l, config, useCase, redis)

	// Initialize Casbin enforcer
	e := casbin.NewEnforcer("config/rbac.conf", "config/policy.csv")
	engine.Use(handlerV1.AuthMiddleware(e))

	// Swagger
	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// K8s probe
	engine.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routes
	v1 := engine.Group("/v1")

	user := v1.Group("/user")
	{
		user.POST("/", handlerV1.CreateUser)
		user.GET("/list", handlerV1.GetUsers)
		user.GET("/:id", handlerV1.GetUser)
		user.PUT("/", handlerV1.UpdateUser)
		user.DELETE("/:id", handlerV1.DeleteUser)
	}

	session := v1.Group("/session")
	{
		session.GET("/list", handlerV1.GetSessions)
		session.GET("/:id", handlerV1.GetSession)
		session.PUT("/", handlerV1.UpdateSession)
		session.DELETE("/:id", handlerV1.DeleteSession)
	}

	auth := v1.Group("/auth")
	{
		auth.POST("/logout", handlerV1.Logout)
		auth.POST("/register", handlerV1.Register)
		auth.POST("/verify-email", handlerV1.VerifyEmail)
		auth.POST("/login", handlerV1.Login)
	}

	tag := v1.Group("/tag")
	{
		tag.POST("/", handlerV1.CreateTag)
		tag.GET("/list", handlerV1.GetTags)
		tag.GET("/:id", handlerV1.GetTag)
		tag.PUT("/", handlerV1.UpdateTag)
		tag.DELETE("/:id", handlerV1.DeleteTag)
	}

	follower := v1.Group("/follower")
	{
		follower.POST("/", handlerV1.FollowUnfollow)
		follower.GET("/list", handlerV1.GetFollowers)
	}

	tweet := v1.Group("/tweet")
	{
		tweet.POST("/", handlerV1.CreateTweet)
		tweet.GET("/list", handlerV1.GetTweets)
		tweet.GET("/:id", handlerV1.GetTweet)
		tweet.PUT("/", handlerV1.UpdateTweet)
		tweet.DELETE("/:id", handlerV1.DeleteTweet)
	}

}

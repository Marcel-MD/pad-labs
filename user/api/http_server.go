package api

import (
	"net/http"

	"user/api/controllers"
	"user/api/middleware"
	"user/config"
	docs "user/docs"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHttpServer(cfg config.Config, userController controllers.UserController) *http.Server {
	log.Info().Msg("Creating new HTTP server")

	e := gin.Default()
	e.Use(middleware.CORS(cfg.AllowOrigin))

	r := e.Group("/api")

	// Register routes
	registerSwaggerRoutes(r, cfg)
	registerUserRoutes(r, cfg, userController)

	return &http.Server{
		Addr:    cfg.HttpPort,
		Handler: e,
	}
}

func registerSwaggerRoutes(router *gin.RouterGroup, cfg config.Config) {
	if cfg.Env == "prod" {
		return
	}

	docs.SwaggerInfo.Host = cfg.Host
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func registerUserRoutes(router *gin.RouterGroup, cfg config.Config, c controllers.UserController) {
	r := router.Group("/users")
	r.POST("/register", c.Register)
	r.POST("/login", c.Login)
	r.GET("/", c.GetAll)
	r.GET("/:id", c.GetById)

	pr := r.Use(middleware.JwtAuth(cfg.AccessTokenSecret))
	pr.GET("/current", c.GetCurrent)
}

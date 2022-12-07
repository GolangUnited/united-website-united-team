package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zhuravlev-pe/course-watch/api/swagger"
	v1 "github.com/zhuravlev-pe/course-watch/internal/adapter/http/v1"
	"net/http"
)

type Handler struct {
	handlerV1 *v1.Handler
}

func NewHandler(handlerV1 *v1.Handler) *Handler {
	return &Handler{
		handlerV1: handlerV1,
	}
}

func (h *Handler) Init() *gin.Engine {

	router := gin.New()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	swagger.SwaggerInfo.Host = "localhost:8080"
	// http://localhost:8080/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	api := router.Group("/api")
	{
		h.handlerV1.Init(api)
	}
}

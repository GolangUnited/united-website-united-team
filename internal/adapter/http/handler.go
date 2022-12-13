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
	services *v1.Services
	bearer   v1.BearerAuthenticator
}

func NewHandler(users v1.UserService, courses v1.CourseService, bearer v1.BearerAuthenticator) *Handler {
	return &Handler{
		services: &v1.Services{
			Users:   users,
			Courses: courses,
		},
		bearer: bearer,
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
	handlerV1 := v1.NewHandler(h.services, h.bearer)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}

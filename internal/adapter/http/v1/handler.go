package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"time"
)

type Handler struct {
	services *services
	bearer   BearerAuthenticator
	//TODO: logger
}

type BearerAuthenticator interface {
	Authenticate(ctx *gin.Context)
	Authorize(role security.Role) func(ctx *gin.Context)
	GenerateToken(principal *security.UserPrincipal) (string, error)
	GetTokenTtl() time.Duration
}

func NewHandler(userService UserService, courseService CourseService, bearer BearerAuthenticator) *Handler {
	return &Handler{
		services: &services{
			Users:   userService,
			Courses: courseService,
		},
		bearer: bearer,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initCoursesRoutes(v1)
		h.initUserRoutes(v1)
		h.initAuthRoutes(v1)
	}
}

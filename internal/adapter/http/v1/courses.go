package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhuravlev-pe/course-watch/internal/adapter/http/v1/utils"
	"github.com/zhuravlev-pe/course-watch/internal/core"
	"net/http"
)

func (h *Handler) initCoursesRoutes(api *gin.RouterGroup) {
	courses := api.Group("/courses")
	{
		//courses.GET("", h.getAllCourses)
		courses.POST("/", h.create)
		courses.GET("/:id", h.getCourseById)
	}
}

// @Summary Get Course By course id
// @Tags courses
// @Description  get course by id
// @ModuleID getCourseById
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Success 200 {object} core.Course
// @Failure 400,404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Failure default {object} utils.Response
// @Router /courses/{id} [get]
func (h *Handler) getCourseById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponseString(c, http.StatusBadRequest, "empty id param")
		return
	}
	
	course, err := h.services.Courses.GetById(c.Request.Context(), id)
	if err != nil {
		if err == core.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	
	c.JSON(http.StatusOK, course)
}

// @Summary Creates a new Course entity
// @Tags courses
// @Description Creates a new Course entity
// @ModuleID create
// @Accept  json
// @Produce  json
// @Param input body core.CreateCourseInput true "sign up info"
// @Success 201 "The generated id is returned in Location header"
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /courses/ [post]
func (h *Handler) create(c *gin.Context) {
	var input core.CreateCourseInput
	if err := c.BindJSON(&input); err != nil {
		utils.ErrorResponseString(c, http.StatusBadRequest, "invalid input body")
		return
	}
	course, err := h.services.Courses.Create(c.Request.Context(), input)
	if err != nil {
		//TODO: discriminate between validation errors, logic errors and internal server errors
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.Header("Location", "/"+course.Id)
	c.Status(http.StatusCreated)
}

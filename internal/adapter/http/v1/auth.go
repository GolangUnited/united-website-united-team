package v1

import (
	"github.com/zhuravlev-pe/course-watch/internal/core"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/zhuravlev-pe/course-watch/internal/adapter/http/v1/utils"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
)

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	courses := api.Group("/auth")
	{
		courses.POST("/signup", h.signupNewUser)
		courses.POST("/login", h.userLogin)
	}
}

// @Summary New user signup
// @Tags Authentication
// @Description Creates new user with the given details
// @ModuleID signupNewUser
// @Accept  json
// @Produce  json
// @Param input body core.SignupUserInput true "New user signup details"
// @Success 200
// @Failure 400,500 {object} utils.Response
// @Router /auth/signup [Post]
func (h *Handler) signupNewUser(ctx *gin.Context) {
	var input core.SignupUserInput
	if err := ctx.BindJSON(&input); err != nil {
		utils.ErrorResponseString(ctx, http.StatusBadRequest, "invalid input body")
		return
	}
	err := h.services.Users.Signup(ctx.Request.Context(), &input)
	
	if err != nil {
		if err == core.ErrUserAlreadyExist {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	
	ctx.Status(http.StatusOK)
}

// @Summary Authenticate user credentials
// @Tags Authentication
// @Description authenticates the user log-in credentials
// @ModuleID userLogin
// @Accept  json
// @Produce  json
// @Param input body core.LoginInput true "Login user details"
// @Success 200 {object} PostUserLoginOutput
// @Failure 400 {object} utils.Response
// @Router /auth/login [Post]
func (h *Handler) userLogin(ctx *gin.Context) {
	var input core.LoginInput
	if err := ctx.BindJSON(&input); err != nil {
		utils.ErrorResponseString(ctx, http.StatusBadRequest, "invalid input body")
		return
	}
	result, err := h.services.Users.Login(ctx.Request.Context(), &input)
	
	if err != nil {
		if err == core.ErrInvalidCredentials {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err)
			return
		}
		
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	
	up := security.UserPrincipal{UserId: result.Id, Roles: result.Roles}
	token, err := h.bearer.GenerateToken(&up)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	output := PostUserLoginOutput{
		UserId:      up.UserId,
		AccessToken: token,
		ExpiresIn:   int(h.bearer.GetTokenTtl().Seconds()),
	}
	ctx.JSON(http.StatusOK, output)
}

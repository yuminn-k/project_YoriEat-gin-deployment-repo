package controller

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/ryokushaka/project_YoriEat-gin-deployment-repo/bootstrap"
	"github.com/ryokushaka/project_YoriEat-gin-deployment-repo/domain"
)

// LoginController handles login-related HTTP requests.
type LoginController struct {
	LoginUsecase domain.LoginUsecase
	Env          *bootstrap.Env
}

// NewLoginController creates a new LoginController.
func NewLoginController(usecase domain.LoginUsecase) *LoginController {
	return &LoginController{
		LoginUsecase: usecase,
	}
}

// Login handles user login requests.
func (lc *LoginController) Login(c *gin.Context) {
	var request domain.LoginRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := lc.LoginUsecase.GetUserByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "User not found with the given email"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Invalid credentials"})
		return
	}

	accessToken, err := lc.LoginUsecase.CreateAccessToken(&user, lc.Env.AccessTokenSecret, lc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	refreshToken, err := lc.LoginUsecase.CreateRefreshToken(&user, lc.Env.RefreshTokenSecret, lc.Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	loginResponse := domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, loginResponse)
}

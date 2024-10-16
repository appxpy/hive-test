package v1

import (
	"net/http"

	"github.com/appxpy/hive-test/internal/usecase"
	"github.com/appxpy/hive-test/pkg/logger"
	"github.com/gin-gonic/gin"
)

type userRoutes struct {
	u usecase.UserUseCase
	l logger.Interface
}

func newUserRoutes(handler *gin.RouterGroup, u usecase.UserUseCase, l logger.Interface) {
	r := &userRoutes{u, l}

	h := handler.Group("/auth")
	{
		h.POST("/register", r.register)
		h.POST("/login", r.login)
	}
}

type userCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary     Register a new user
// @Description Registers a new user
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       credentials body userCredentials true "User Credentials"
// @Success     201
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /auth/register [post]
func (r *userRoutes) register(c *gin.Context) {
	var creds userCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		r.l.Error(err, "http - v1 - register")
		errorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := r.u.Register(c.Request.Context(), creds.Username, creds.Password)
	if err != nil {
		r.l.Error(err, "http - v1 - register")
		errorResponse(c, http.StatusInternalServerError, "Could not create user")
		return
	}

	c.Status(http.StatusCreated)
}

type loginResponse struct {
	Token string `json:"token" exaple:"your_jwt_token"`
}

// @Summary     Login
// @Description Authenticates a user and returns a JWT token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       credentials body userCredentials true "User Credentials"
// @Success     200 {object} loginResponse
// @Failure     400 {object} response
// @Failure     401 {object} response
// @Router      /auth/login [post]
func (r *userRoutes) login(c *gin.Context) {
	var creds userCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		r.l.Error(err, "http - v1 - login")
		errorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := r.u.Login(c.Request.Context(), creds.Username, creds.Password)
	if err != nil {
		r.l.Error(err, "http - v1 - login")
		errorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

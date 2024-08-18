package handler

import (
	"avito2024/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NoAuthService interface {
	DummyLoginService(status string) (string, error)
	LoginService(user domain.User) (string, error)
	RegisterService(user domain.User) (string, error)
}
type AuthOnlyService interface {
	CheckUserService(tokenString string) (domain.User, error)
	GetHouseFlatsService(id int, user domain.User) ([]domain.Flat, error)
	SubscribeHouseService(email string) error
	CreateFlatService(flat domain.Flat) (domain.Flat, error)
}

type ModerationOnlyService interface {
	CreateHouseService(house domain.House) (domain.House, error)
	UpdateFlatService(fl domain.Flat) (domain.Flat, error)
}
type Service interface {
	NoAuthService
	AuthOnlyService
	ModerationOnlyService
}

type Handler struct {
	Service
}

func NewHandler(s Service, engine *gin.Engine) *Handler {
	handler := &Handler{s}

	engine.GET("/dummyLogin", handler.DummyLogin)
	engine.POST("/login", handler.Login)
	engine.POST("/register", handler.Register)

	middleWareGroup := engine.Group("/").Use(handler.CheckUser).Use(handler.SetRequestID)
	{
		middleWareGroup.POST("/house/create", handler.CheckModerator, handler.CreateHouse)
		middleWareGroup.GET("/house/:id", handler.GetHouseFlats)
		middleWareGroup.POST("/house/:id/subscribe", handler.SubscribeHouse)
		middleWareGroup.POST("/flat/create", handler.CreateFlat)
		middleWareGroup.POST("/flat/update", handler.CheckModerator, handler.UpdateFlat)
	}

	return handler
}

func respondWithError(c *gin.Context, err error) {
	customErr, ok := err.(*domain.CustomError)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	switch httpStatus := customErr.GetHttpStatus(); httpStatus {
	case 400, 401, 404:
		c.AbortWithStatus(httpStatus)
		return
	}

	requestID, _ := c.Get("requestID")
	c.JSON(customErr.GetHttpStatus(), gin.H{
		"message":    customErr.GetUserLog(),
		"request_id": requestID,
		"code":       customErr.GetCode(),
	})
}

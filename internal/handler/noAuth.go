package handler

import (
	"avito2024/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) DummyLogin(c *gin.Context) {
	userStatus := c.Query("user_type")
	token, err := h.DummyLoginService(userStatus)
	if err != nil {
		respondWithError(c, err)
		return
	}
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) Login(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		respondWithError(c, domain.NewCustomError(domain.InvalidInputError()))
		return
	}
	token, err := h.LoginService(user)
	if err != nil {
		respondWithError(c, err)
		return
	}
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		respondWithError(c, domain.NewCustomError(domain.InvalidInputError()))
		return
	}
	token, err := h.RegisterService(user)
	if err != nil {
		respondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": token})
}

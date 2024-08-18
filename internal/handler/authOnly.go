package handler

import (
	"avito2024/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) CheckUser(c *gin.Context) {
	tokenString, err := c.Cookie("token")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user, err := h.CheckUserService(tokenString)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", user)
	c.Next()
}

func (h *Handler) GetHouseFlats(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user, _ := c.Get("user")
	flats, err := h.GetHouseFlatsService(id, user.(domain.User))
	if err != nil {
		respondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, flats)
}

// ???
func (h *Handler) SubscribeHouse(c *gin.Context) {
	email := c.PostForm("email")
	if err := h.SubscribeHouseService(email); err != nil {
		respondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "subscription successful"})
}

func (h *Handler) CreateFlat(c *gin.Context) {
	var flat domain.Flat
	if err := c.ShouldBindJSON(&flat); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	createdFlat, err := h.CreateFlatService(flat)
	if err != nil {
		respondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, createdFlat)
}

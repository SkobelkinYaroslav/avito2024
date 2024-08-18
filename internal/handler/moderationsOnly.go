package handler

import (
	"avito2024/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) CheckModerator(c *gin.Context) {
	user, _ := c.Get("user")
	if user.(domain.User).UserType != "moderator" {
		respondWithError(c, domain.NewCustomError(domain.AuthError()))
		return
	}

	c.Next()
}

func (h *Handler) CreateHouse(c *gin.Context) {
	var house domain.House
	if err := c.ShouldBindJSON(&house); err != nil {
		respondWithError(c, domain.NewCustomError(domain.InvalidInputError()))
		return
	}
	createdHouse, err := h.CreateHouseService(house)
	if err != nil {
		respondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, createdHouse)
}

func (h *Handler) UpdateFlat(c *gin.Context) {
	var flat domain.Flat
	if err := c.ShouldBindJSON(&flat); err != nil {
		respondWithError(c, domain.NewCustomError(domain.InvalidInputError()))
		return
	}
	updatedFlat, err := h.UpdateFlatService(flat)
	if err != nil {
		respondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, updatedFlat)
}

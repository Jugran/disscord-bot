package handlers

import (
	"diss-cord/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userSchema struct {
	Name              *string  `json:"name"`
	Disrespect        float32  `json:"disrespect,string,omitempty"`
	DiscordID         string   `json:"discord_id" binding:"required"`
	SeverityThreshold uint8    `json:"severity_threshold"`
	Roles             []string `json:"roles" binding:"required"`
}

type roleSchema struct {
	Name string `json:"name" uri:"name" binding:"required"`
}

type userRoleShema struct {
	UserID int      `json:"user_id" binding:"required"`
	Roles  []string `json:"roles" binding:"required"`
}

func AddRoleHandler(c *gin.Context) {
	roleData := roleSchema{}

	if err := c.ShouldBindUri(&roleData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid Input",
		})
		return
	}

	role := models.Role{
		Name: roleData.Name,
	}

	models.AddRoleAction(&role)

	c.JSON(http.StatusCreated, &role)
}

func FetchUsersHandler(c *gin.Context) {
	users, count := models.FetchAllUsersActions()

	fmt.Printf("Fetched %d user entries", count)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": count,
	})
}

func AddUserHandler(c *gin.Context) {
	var userData userSchema
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid Input",
		})
		return
	}

	user := models.User{
		Name:              userData.Name,
		Disrespect:        userData.Disrespect,
		DiscordID:         userData.DiscordID,
		SeverityThreshold: userData.SeverityThreshold,
	}

	models.AddUserAction(&user, &userData.Roles)

	c.JSON(http.StatusCreated, &user)
}

func AddUserRoleHandler(c *gin.Context) {
	var userData userRoleShema
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid Input",
		})
		return
	}

	models.AddUserRolesAction(userData.UserID, userData.Roles)

	c.JSON(http.StatusCreated, &userData)
}

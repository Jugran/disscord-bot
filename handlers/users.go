package handlers

import (
	"database/sql"
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
	UserID int    `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
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

	models.AddRole(&role)

	c.JSON(http.StatusCreated, &role)
}

func GetUsersHandler(c *gin.Context) {
	users, count := models.FetchAllUsers()

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

	name := sql.NullString{String: *userData.Name, Valid: true}
	if userData.Name == nil {
		name = sql.NullString{Valid: false}
	}

	var roles []models.Role

	for _, role := range userData.Roles {
		roles = append(roles, models.Role{
			Name: role,
		})
	}

	user := models.User{
		Name:              name,
		Disrespect:        userData.Disrespect,
		DiscordID:         userData.DiscordID,
		SeverityThreshold: userData.SeverityThreshold,
		Roles:             roles,
	}

	models.AddUser(&user)

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

	models.AddUserRole(userData.UserID, userData.Role)

	c.JSON(http.StatusCreated, &userData)
}

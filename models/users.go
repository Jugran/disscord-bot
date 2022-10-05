package models

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name              sql.NullString `json:"name" gorm:"type:varchar(30)"`
	Disrespect        float32        `json:"disrespect" gorm:"default:0.5"`
	DiscordID         string         `json:"discord_id" gorm:"uniqueIndex;not null;default:null"`
	SeverityThreshold uint8          `json:"severity_threshold" gorm:"default:2"`
	Roles             []Role         `gorm:"many2many:user_role"`
}

type Role struct {
	gorm.Model
	Name    string   `json:"name" gorm:"uniqueIndex;not null;default:null"`
	Insults []Insult `gorm:"many2many:role_insult"`
}

func FetchAllUsers() ([]User, int64) {
	var users []User
	result := DB.Find(&users)

	if result.Error != nil {
		fmt.Println("Cannot fetch user data:", result.Error)
		return []User{}, 0
	}

	return users, result.RowsAffected
}

func AddUser(user *User) bool {
	result := DB.Create(&user)

	if result.Error != nil {
		fmt.Println("Cannot add user data:", result.Error)
		return false
	}

	return true
}

func AddRole(role *Role) bool {
	result := DB.Create(&role)

	if result.Error != nil {
		fmt.Println("Cannot add role data:", result.Error)
		return false
	}

	return true
}

func AddUserRole(userID int, roleName string) bool {
	err := DB.Model(&User{Model: gorm.Model{ID: uint(userID)}}).Association("Roles").Append([]Role{
		{Name: roleName},
	})

	if err != nil {
		fmt.Println("Cannot add user role:", err)
		return false
	}

	return true
}

func RemoveUserRole(userID int, roleName string) bool {
	err := DB.Model(&User{Model: gorm.Model{ID: uint(userID)}}).Association("Roles").Delete([]Role{
		{Name: roleName},
	})

	if err != nil {
		fmt.Println("Cannot remove user role:", err)
		return false
	}

	return true
}

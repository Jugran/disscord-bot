package models

// https://foaas.com/

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
	result := DB.Preload("Roles").Find(&users)

	if result.Error != nil {
		fmt.Println("Cannot fetch user data:", result.Error)
		return []User{}, 0
	}

	return users, result.RowsAffected
}

func AddUser(user *User, roleNames *[]string) bool {
	roles := GetRolesByName(roleNames)
	result := DB.Debug().Where("name IN (?)", *roleNames).Find(&roles)

	if result.Error == nil {
		user.Roles = *roles
	}

	result = DB.Debug().Omit("Roles.*").Create(&user)

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

func AddUserRoles(userID int, roleNames []string) bool {
	roles := GetRolesByName(&roleNames)

	err := DB.Debug().
		Model(&User{Model: gorm.Model{ID: uint(userID)}}).
		Omit("Roles.*").
		Association("Roles").
		Append(roles)

	if err != nil {
		fmt.Println("Cannot add user role:", err)
		return false
	}

	return true
}

func RemoveUserRole(userID int, roleName string) bool {
	role := GetRolesByName(&[]string{
		roleName,
	})

	err := DB.Model(&User{Model: gorm.Model{ID: uint(userID)}}).Association("Roles").Delete(role)

	if err != nil {
		fmt.Println("Cannot remove user role:", err)
		return false
	}

	return true
}

func GetRolesByName(roleNames *[]string) *[]Role {

	roles := []Role{}

	result := DB.Debug().Where("name IN (?)", *roleNames).Find(&roles)

	if result.Error != nil {
		fmt.Println("Cannot fetch role:", result.Error)
		return nil
	}

	return &roles
}
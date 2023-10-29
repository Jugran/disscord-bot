package models

// https://foaas.com/

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	Name              *string `json:"name" gorm:"type:varchar(30)"`
	Disrespect        float32 `json:"disrespect" gorm:"default:0.5"`
	DiscordID         string  `json:"discord_id" gorm:"uniqueIndex;not null"`
	SeverityThreshold uint8   `json:"severity_threshold" gorm:"default:2"`
	Roles             []Role  `gorm:"many2many:user_roles"`
}

type Role struct {
	gorm.Model
	Name    string   `json:"name" gorm:"uniqueIndex;not null;default:null"`
	Insults []Insult `gorm:"many2many:insult_roles"`
}

type APIRole struct {
	Name string
}

// DB Model Actions

func FetchAllUsersActions() ([]User, int64) {
	var users []User
	result := DB.Preload("Roles").Preload(clause.Associations).Find(&users)

	if result.Error != nil {
		fmt.Println("Cannot fetch user data:", result.Error)
		return []User{}, 0
	}

	return users, result.RowsAffected
}

func AddUserAction(user *User, roleNames *[]string) bool {
	roles := GetRolesByName(roleNames)

	if roles != nil {
		user.Roles = *roles
	}

	result := DB.Debug().Omit("Roles.*").Create(&user)

	if result.Error != nil {
		fmt.Println("Cannot add user data:", result.Error)
		return false
	}

	return true
}

func (u *User) AddNewUser(roleNames *[]string) bool {
	roles := GetRolesByName(roleNames)

	if roles != nil {
		u.Roles = *roles
	}

	result := DB.Debug().Omit("Roles.*").Create(&u)

	if result.Error != nil {
		fmt.Println("Cannot add user data:", result.Error)
		return false
	}

	return true
}

func AddRoleAction(role *Role) bool {
	result := DB.Create(&role)

	if result.Error != nil {
		fmt.Println("Cannot add role data:", result.Error)
		return false
	}

	return true
}

func AddUserRolesAction(userID int, roleNames []string) bool {
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

func RemoveUserRoleAction(userID int, roleName string) bool {
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

	if result.Error == nil && len(roles) > 0 {
		return &roles
	}

	fmt.Println("Cannot fetch roles:", result.Error)

	roles = []Role{}
	// loop through roleNames
	for _, roleName := range *roleNames {
		role := Role{}
		result := DB.FirstOrCreate(&role, Role{Name: roleName})
		if result.Error != nil {
			continue
		}
		roles = append(roles, role)
	}

	return &roles
}

func CheckDiscordUser(discordID string) (*User, error) {
	user := User{
		DiscordID: discordID,
	}

	result := DB.Where(&User{DiscordID: discordID}).First(&user)

	if result.Error != nil {
		return &user, result.Error
	}

	return &user, nil
}

package models

import (
	"fmt"

	"gorm.io/gorm"
)

type InsultData struct {
	// old data insult data will be replaced
	Insults []Insult
}

type Insult struct {
	gorm.Model
	Text     string  `json:"text"`
	Severity uint8   `json:"severity" gorm:"default:0"`
	Roles    []*Role `gorm:"many2many:insult_roles"`
}

type APIInsult struct {
	Text string `json:"text"`
}

// TODO: make all actions receiver functions

func FindInsultsAction() ([]Insult, int64) {
	var insults []Insult

	result := DB.Preload("Roles").Find(&insults)

	if result.Error != nil {
		fmt.Println("Cannot fetch insult data:", result.Error)
		return []Insult{}, 0
	}
	return insults, result.RowsAffected

}

func AddInsultAction(insult *Insult, roleNames *[]string) bool {
	roles := GetRolesByName(roleNames)

	err := DB.Create(&insult)

	if err.Error != nil {
		fmt.Println("Cannot add insult data:", err)
		return false
	}

	err1 := DB.Model(&Insult{Model: gorm.Model{ID: insult.ID}}).
		Omit("Roles.*").
		Association("Roles").
		Append(roles)

	if err1 != nil {
		fmt.Println("Cannot add insult roles:", err1)
		return false
	}

	return true
}

func UpdateInsultAction(insult *Insult) int64 {
	result := DB.Model(&insult).Updates(insult)

	if result.Error != nil {
		fmt.Println("Cannot update insult:", result.Error)
		return 0
	}

	return result.RowsAffected
}

func DeleteInsultAction(insult *Insult) bool {
	result := DB.Delete(&insult)

	if result.Error != nil {
		fmt.Println("Cannot add insult data:", result.Error)
		return false
	}

	return true
}

func GetInsultForUser(user_id *uint) (*gorm.DB, *APIInsult) {

	var insults APIInsult

	var result *gorm.DB

	if user_id != nil && *user_id != 0 {
		// user_id -> role -> insult
		var insult string

		result = DB.Table("insult_roles").
			Joins("INNER JOIN user_roles ON user_roles.role_id = insult_roles.role_id").
			Where("user_id = ?", user_id).
			Distinct("insult_id").
			Joins("INNER JOIN insults ON insults.id = insult_roles.insult_id").
			Select("insults.text").
			Order("random()").
			Limit(1).
			Pluck("text", &insult)

		if insult != "" {
			insults.Text = insult
		} else {
			result = DB.Model(&Insult{}).Select("text").Order("random()").Limit(1).Take(&insults)
		}

		if result.Error != nil {
			return result, &APIInsult{}
		}

	} else {
		// TODO: add severity condition
		result = DB.Table("insults").Order("random()").Limit(1).Take(&insults)
	}

	return result, &insults
}

package models

import (
	"fmt"
	"math/rand"

	"gorm.io/gorm"
)

type InsultData struct {
	// old data insult data will be replaced
	Insults []Insult
}

type Insult struct {
	gorm.Model
	Text     string `json:"text"`
	Severity uint8  `json:"severity" gorm:"default:0"`
	Roles    []Role `gorm:"many2many:insult_roles"`
}

// TODO: make all actions receiver functions

func FindInsultsAction() ([]Insult, int64) {
	var insults []Insult

	result := DB.Find(&insults)

	if result.Error != nil {
		fmt.Println("Cannot fetch insult data:", result.Error)
		return []Insult{}, 0
	}
	return insults, result.RowsAffected

}

func AddInsultAction(insult *Insult, roleNames *[]string) bool {
	roles := GetRolesByName(roleNames)

	if roles != nil {
		insult.Roles = *roles
	}

	result := DB.Debug().Omit("Roles.*").Create(&insult)

	if result.Error != nil {
		fmt.Println("Cannot add insult data:", result.Error)
		return false
	}

	return true
}

func UpdateInsultAction(insult *Insult) int64 {
	// This won't work, unless pass in ID from API
	result := DB.Model(&insult).Updates(insult)

	if result.Error != nil {
		fmt.Println("Cannot update insult:", result.Error)
		return 0
	}

	return result.RowsAffected
}

func DeleteInsultAction(insult *Insult) bool {
	result := DB.Debug().Delete(&insult)

	if result.Error != nil {
		fmt.Println("Cannot add insult data:", result.Error)
		return false
	}

	return true
}

// ---------- old shit ------------

func NewInsults() InsultData {
	insultData := InsultData{
		Insults: []Insult{},
	}

	insultData.LoadInsults()
	return insultData
}

func (insultData *InsultData) LoadInsults() {
	insultData.Insults = []Insult{
		{
			Text:     "You're a jerk.",
			Severity: 1,
		},
		{
			Text:     "You're an idiot.",
			Severity: 1,
		},
	}
}

func (insData *InsultData) GetInsult() string {
	randIndex := rand.Intn(len(insData.Insults))
	insult := insData.Insults[randIndex].Text
	return insult
}

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
	Severity uint8  `json:"severity"`
}

func FindInsults() ([]Insult, int64) {
	var insults []Insult

	result := DB.Find(&insults)

	if result.Error != nil {
		fmt.Println("Cannot fetch insult data:", result.Error)
		return []Insult{}, 0
	}
	return insults, result.RowsAffected

}

func AddInsult(insult *Insult, roles []string) bool {

	var userRoles []Role

	for _, role := range roles {
		userRoles = append(userRoles, Role{
			Name: role,
		})
	}

	err := DB.Create(&insult).Association("Roles").Append(userRoles)

	if err != nil {
		fmt.Println("Cannot add insult data:", err)
		return false
	}

	return true
}

func UpdateInsult(insult *Insult) int64 {
	result := DB.Model(&insult).Updates(insult)

	if result.Error != nil {
		fmt.Println("Cannot update insult:", result.Error)
		return 0
	}

	return result.RowsAffected
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

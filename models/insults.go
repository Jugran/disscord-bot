package models

import (
	"fmt"
	"math/rand"

	"gorm.io/gorm"
)

type InsultData struct {
	Insults []Insult
}

type Insult struct {
	gorm.Model
	Text     string `json:"text"`
	Target   string `json:"target"`
	Severity uint8  `json:"severity"`
}

func FindInsults() []Insult {
	var insults []Insult

	DB.Find(&insults)

	return insults
}

func AddInsult(insult *Insult) bool {
	result := DB.Create(&insult)

	if result.Error != nil {
		fmt.Println("Cannot add insult data:", result.Error)
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
			Target:   "members",
			Severity: 1,
		},
		{
			Text:     "You're an idiot.",
			Target:   "members",
			Severity: 1,
		},
	}
}

func (insData *InsultData) GetInsult() string {
	randIndex := rand.Intn(len(insData.Insults))
	insult := insData.Insults[randIndex].Text
	return insult
}

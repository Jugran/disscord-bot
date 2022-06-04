package insults

import "math/rand"

type InsultData struct {
	Insults []insult
}

type insult struct {
	Insult   string `json:"insult"`
	Target   string `json:"target"`
	Severity uint8  `json:"severity"`
}

func NewInsults() InsultData {
	insultData := InsultData{
		Insults: []insult{},
	}

	insultData.LoadInsults()
	return insultData
}

func (insultData *InsultData) LoadInsults() {
	insultData.Insults = []insult{
		{
			Insult:   "You're a jerk.",
			Target:   "members",
			Severity: 1,
		},
		{
			Insult:   "You're an idiot.",
			Target:   "members",
			Severity: 1,
		},
	}
}

func (insData *InsultData) GetInsult() string {
	randIndex := rand.Intn(len(insData.Insults))
	insult := insData.Insults[randIndex].Insult
	return insult
}

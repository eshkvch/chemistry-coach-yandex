package catalog

import "fmt"

type Goal struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Difficulty  int    `json:"difficulty"`
	Recommended bool   `json:"recommended"`
}

type Persona struct {
	ID            string   `json:"id"`
	Initials      string   `json:"initials"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Tags          []string `json:"tags"`
	Difficulty    int      `json:"difficulty"`
	ShowDemoBadge bool     `json:"showDemoBadge,omitempty"`
}

var Goals = []Goal{
	{ID: "first-message", Title: "Первое сообщение после мэтча", Description: "Завязать разговор и не уйти в банальности.", Difficulty: 1},
	{ID: "keep-interest", Title: "Удержание интереса", Description: "Поддержать диалог, когда он остывает.", Difficulty: 2, Recommended: true},
	{ID: "suggest-call", Title: "Переход на звонок", Description: "Предложить созвон без давления.", Difficulty: 2},
	{ID: "boundaries", Title: "Разговор о границах", Description: "Обсудить ожидания и комфорт уважительно.", Difficulty: 3},
}

var Personas = []Persona{
	{
		ID: "calm-careful", Initials: "МВ", Title: "Мягкая и вдумчивая",
		Description: "Отвечает медленно, ценит глубину. Молчание — тоже ответ.",
		Tags:        []string{"Терпеливая", "Вдумчивая", "Чуткая"}, Difficulty: 1,
	},
	{
		ID: "ironic-fast", Initials: "ИБ", Title: "Ироничная и быстрая",
		Description: "Отвечает коротко. Скучает быстро. Ценит конкретику и юмор.",
		Tags:        []string{"Ироничная", "Быстрая", "Требовательная"}, Difficulty: 2, ShowDemoBadge: true,
	},
	{
		ID: "busy", Initials: "ПЗ", Title: "Прямая и занятая",
		Description: "Пишет между делами. Ценит краткость. Не терпит воды.",
		Tags:        []string{"Прямая", "Занятая", "Лаконичная"}, Difficulty: 3,
	},
}

func GetGoal(id string) (*Goal, bool) {
	for i := range Goals {
		if Goals[i].ID == id {
			return &Goals[i], true
		}
	}
	return nil, false
}

func GetPersona(id string) (*Persona, bool) {
	for i := range Personas {
		if Personas[i].ID == id {
			return &Personas[i], true
		}
	}
	return nil, false
}

func GoalsWithRecommendation(recommendedID string) []Goal {
	out := make([]Goal, len(Goals))
	for i, g := range Goals {
		out[i] = g
		out[i].Recommended = g.ID == recommendedID
	}
	return out
}

func RecommendedGoalForSkill(skill string) string {
	switch skill {
	case "clarity":
		return "first-message"
	case "confidence":
		return "keep-interest"
	case "respect":
		return "boundaries"
	case "balance":
		return "suggest-call"
	default:
		return "keep-interest"
	}
}

func SkillLabel(skill string) string {
	switch skill {
	case "clarity":
		return "ясность"
	case "confidence":
		return "уверенность"
	case "respect":
		return "уважение"
	case "balance":
		return "баланс"
	default:
		return skill
	}
}

func BuildSystemPrompt(goalID, personaID string) string {
	goal, okG := GetGoal(goalID)
	persona, okP := GetPersona(personaID)
	if !okG || !okP {
		return "Ты — собеседница в тренажёре романтической коммуникации 18+."
	}
	return fmt.Sprintf(`Ты — AI-собеседница в тренажёре романтической коммуникации для взрослых (18+).

Персона: %s — %s
Теги: %v
Стиль: %s

Цель пользователя в этой сессии: %s — %s

Правила:
- Оставайся в роли персоны, отвечай коротко и естественно на русском.
- Не выходи из роли, не давай советов пользователю напрямую.
- Учитывай цель сессии в своих реакциях.
- Если пользователь нарушает границы — мягко обозначь дискомфорт в характере персоны.
- Не генерируй откровенно сексуальный контент; фокус на флирте, интересе и уважении.`,
		persona.Title, persona.Description, persona.Tags, persona.Description,
		goal.Title, goal.Description)
}

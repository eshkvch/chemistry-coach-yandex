package usecase

import (
	"chemistry-coach/internal/catalog"
	"chemistry-coach/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func newUserID() string    { return "usr_" + uuid.New().String()[:8] }
func newSessionID() string { return "ses_" + uuid.New().String()[:8] }
func newMessageID() string { return "msg_" + uuid.New().String()[:8] }
func newDebriefID() string { return "deb_" + uuid.New().String()[:8] }

func weakestSkill(avg map[string]float64) string {
	if len(avg) == 0 {
		return domain.SkillClarity
	}
	skills := []string{domain.SkillClarity, domain.SkillConfidence, domain.SkillRespect, domain.SkillBalance}
	weakest := skills[0]
	min := 11.0
	for _, s := range skills {
		v, ok := avg[s]
		if !ok {
			continue
		}
		if v < min {
			min = v
			weakest = s
		}
	}
	return weakest
}

func toChatHistory(msgs []domain.Message) []domain.ChatMessage {
	out := make([]domain.ChatMessage, 0, len(msgs))
	for _, m := range msgs {
		if m.Sender == domain.SenderUser || m.Sender == domain.SenderPersona {
			out = append(out, domain.ChatMessage{Role: m.Sender, Text: m.Text})
		}
	}
	return out
}

func formatSessionDate(t time.Time) (date, timeStr string) {
	months := []string{"", "янв", "фев", "мар", "апр", "мая", "июн", "июл", "авг", "сен", "окт", "ноя", "дек"}
	return fmt.Sprintf("%d %s", t.Day(), months[t.Month()]), t.Format("15:04")
}

func sessionScoresFromDebrief(scores domain.DebriefScores) domain.SkillScores {
	return domain.SkillScores{
		Clarity: scores.Clarity.Value, Confidence: scores.Confidence.Value,
		Respect: scores.Respect.Value, Balance: scores.Balance.Value,
	}
}

func goalTitle(id string) string {
	if g, ok := catalog.GetGoal(id); ok {
		return g.Title
	}
	return id
}

func personaTitle(id string) string {
	if p, ok := catalog.GetPersona(id); ok {
		return p.Title
	}
	return id
}

func personaInitials(id string) string {
	if p, ok := catalog.GetPersona(id); ok {
		return p.Initials
	}
	return "??"
}

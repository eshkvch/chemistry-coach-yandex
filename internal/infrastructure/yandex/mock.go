package yandex

import (
	"chemistry-coach/internal/domain"
	"context"
	"fmt"
)

// MockClient returns deterministic responses when Yandex API is not configured.
type MockClient struct{}

func NewMockClient() *MockClient { return &MockClient{} }

func (m *MockClient) GenerateOpening(_ context.Context, _ string) (string, error) {
	return "Привет 🙂 чем сегодня занимался?", nil
}

func (m *MockClient) ProcessMessage(_ context.Context, _ string, history []domain.ChatMessage, userText string) (*domain.ChatTurnResult, error) {
	reply := "Интересно, расскажи подробнее."
	if len(history) > 2 {
		reply = "Какую?"
	}
	return &domain.ChatTurnResult{
		PersonaText: reply,
		Score: domain.MessageScore{
			Clarity: 7, Confidence: 7, Respect: 9, Balance: 6,
			Status: "success", StatusLabel: "Отлично: конкретика, открытость, повод для продолжения",
		},
		Consent: domain.ConsentRisk{
			Detected: stringsContainsPressure(userText),
		},
	}, nil
}

func (m *MockClient) SuggestReply(_ context.Context, _ string, _ []domain.ChatMessage, draft string) (string, error) {
	return fmt.Sprintf("«Гражданская оборона» — слушал у отца на проигрывателе. %s", draft), nil
}

func (m *MockClient) SummarizeSession(_ context.Context, _ string, _ []domain.ChatMessage) (*domain.SessionSummary, error) {
	return &domain.SessionSummary{
		Scores: domain.DebriefScores{
			Clarity:    domain.ScoreDetail{Value: 6, Comment: "Часто отвечал размыто, теряешь конкретику."},
			Confidence: domain.ScoreDetail{Value: 7, Comment: "Хорошо держишь темп, но иногда оправдываешься."},
			Respect:    domain.ScoreDetail{Value: 9, Comment: "Границы соблюдены, тон спокойный."},
			Balance:    domain.ScoreDetail{Value: 5, Comment: "Слишком много вопросов в лоб, мало реакций."},
		},
		Strengths:  []string{"Уверенно держишь темп диалога", "Соблюдаешь границы и не давишь"},
		Weaknesses: []string{"Размытые ответы вместо конкретики", "Перегружаешь диалог вопросами"},
		RiskFlags: []domain.RiskFlag{
			{
				Severity: "medium", Quote: "Ну такую, советскую.",
				Explanation: "Размытый ответ обесценивает её вопрос.",
				Suggestion:  "«Гражданская оборона», слушал у отца на проигрывателе.",
			},
		},
		ImprovedReplies: []domain.ImprovedReply{
			{Original: "Ну такую, советскую.", Improved: "«Гражданская оборона». Слушал у отца.", Reason: "Конкретика цепляет."},
		},
		TipForNext: "Меньше прямых вопросов в лоб. Сначала — короткая реакция на её реплику.",
		HasRisk:    true,
	}, nil
}

func (m *MockClient) GenerateDailyExercise(_ context.Context, focusSkill string) (*domain.DailyExercise, error) {
	return &domain.DailyExercise{
		Title:       fmt.Sprintf("Дай конкретику вместо общих фраз (%s)", focusSkill),
		Description: "В следующей сессии хотя бы один раз ответь конкретной деталью из жизни.",
		Criterion:   "минимум одно сообщение со скором «Ясность» 8+",
	}, nil
}

func stringsContainsPressure(s string) bool {
	pressure := []string{"должна", "обязана", "надо тебе", "сейчас же"}
	for _, p := range pressure {
		if len(s) > 0 && containsFold(s, p) {
			return true
		}
	}
	return false
}

func containsFold(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(len(s) > 0 && findSub(s, sub)))
}

func findSub(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if equalFoldAt(s, sub, i) {
			return true
		}
	}
	return false
}

func equalFoldAt(s, sub string, i int) bool {
	for j := 0; j < len(sub); j++ {
		a, b := s[i+j], sub[j]
		if a >= 'A' && a <= 'Z' {
			a += 'a' - 'A'
		}
		if b >= 'A' && b <= 'Z' {
			b += 'a' - 'A'
		}
		if a != b {
			return false
		}
	}
	return true
}

var _ domain.LLMService = (*MockClient)(nil)

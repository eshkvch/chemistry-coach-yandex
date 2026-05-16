package yandex

import (
	"bytes"
	"chemistry-coach/internal/config"
	"chemistry-coach/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	apiKey   string
	folderID string
	model    string
	baseURL  string
	http     *http.Client
}

func NewClient(cfg config.YandexConfig) *Client {
	return &Client{
		apiKey:   cfg.APIKey,
		folderID: cfg.FolderID,
		model:    cfg.Model,
		baseURL:  strings.TrimSuffix(cfg.BaseURL, "/"),
		http:     &http.Client{Timeout: 90 * time.Second},
	}
}

func (c *Client) modelURI() string {
	return fmt.Sprintf("gpt://%s/%s", c.folderID, c.model)
}

type completionRequest struct {
	ModelURI          string              `json:"modelUri"`
	CompletionOptions completionOptions   `json:"completionOptions"`
	Messages          []completionMessage `json:"messages"`
	JSONSchema        *jsonSchemaWrapper  `json:"jsonSchema,omitempty"`
}

type completionOptions struct {
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"maxTokens"`
}

type completionMessage struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type jsonSchemaWrapper struct {
	Schema map[string]interface{} `json:"schema"`
}

type completionResponse struct {
	Result struct {
		Alternatives []struct {
			Message struct {
				Role string `json:"role"`
				Text string `json:"text"`
			} `json:"message"`
			Status string `json:"status"`
		} `json:"alternatives"`
	} `json:"result"`
}

type ResponseRequest struct {
	Model           string  `json:"model"`
	Temperature     float64 `json:"temperature"`
	Instructions    string  `json:"instructions"`
	Input           string  `json:"input"`
	MaxOutputTokens int     `json:"max_output_tokens"`
}

type ResponseData struct {
	Output []struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

func (r ResponseData) GetOutputText() string {
	for _, output := range r.Output {
		if len(output.Content) > 0 {
			return output.Content[0].Text
		}
	}
	return ""
}

func (c *Client) complete(ctx context.Context, messages []completionMessage, schema map[string]interface{}, temperature float64, maxTokens int) (string, error) {
	if c.apiKey == "" || c.folderID == "" {
		return "", fmt.Errorf("yandex ai not configured")
	}

	fmt.Println("messages:", messages)
	fmt.Println("schema", schema)

	mes, err := json.Marshal(messages)
	if err != nil {
		fmt.Println("error marshall message: %w", err)
	}

	reqData := ResponseRequest{
		Model:           c.modelURI(),
		Temperature:     0.3,
		Instructions:    "",
		Input:           string(mes),
		MaxOutputTokens: 500,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		log.Fatal(err)
	}

	// if schema != nil {
	// 	reqBody.JSONSchema = &jsonSchemaWrapper{Schema: schema}
	// }

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Api-Key "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Project", c.folderID)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("yandex api %d: %s", resp.StatusCode, string(body))
	}

	var response ResponseData
	outTest := ""
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Error parsing response: %v", err)
	} else {
		outTest = response.GetOutputText()
		fmt.Println("outTest:", outTest)
	}
	return strings.TrimSpace(outTest), nil
}

func (c *Client) GenerateOpening(ctx context.Context, systemPrompt string) (string, error) {
	messages := []completionMessage{
		{Role: "system", Text: systemPrompt},
		{Role: "user", Text: "Начни диалог первой репликой от лица персоны. Одно короткое сообщение на русском, без пояснений."},
	}
	return c.complete(ctx, messages, nil, 0.7, 300)
}

func (c *Client) ProcessMessage(ctx context.Context, systemPrompt string, history []domain.ChatMessage, userText string) (*domain.ChatTurnResult, error) {
	var hist strings.Builder
	for _, m := range history {
		role := "Пользователь"
		if m.Role == domain.SenderPersona {
			role = "Персона"
		}
		fmt.Fprintf(&hist, "%s: %s\n", role, m.Text)
	}
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"personaReply":       map[string]interface{}{"type": "string"},
			"clarity":            map[string]interface{}{"type": "integer", "minimum": 0, "maximum": 10},
			"confidence":         map[string]interface{}{"type": "integer", "minimum": 0, "maximum": 10},
			"respect":            map[string]interface{}{"type": "integer", "minimum": 0, "maximum": 10},
			"balance":            map[string]interface{}{"type": "integer", "minimum": 0, "maximum": 10},
			"statusLabel":        map[string]interface{}{"type": "string"},
			"status":             map[string]interface{}{"type": "string", "enum": []string{"success", "warning", "danger"}},
			"consentDetected":    map[string]interface{}{"type": "boolean"},
			"consentSeverity":    map[string]interface{}{"type": "string", "enum": []string{"low", "medium", "high", "none"}},
			"consentExplanation": map[string]interface{}{"type": "string"},
			"consentSuggestion":  map[string]interface{}{"type": "string"},
		},
		"required": []string{
			"personaReply", "clarity", "confidence", "respect", "balance",
			"statusLabel", "status", "consentDetected",
		},
	}
	prompt := fmt.Sprintf(`История диалога:
%s
Новое сообщение пользователя: %s

Выполни function calling логику:
1) score_message — оцени сообщение по 4 шкалам 0-10
2) detect_consent_risk — проверь давление/нарушение согласия
3) сгенерируй короткий ответ персоны personaReply

Верни только JSON по схеме.`, hist.String(), userText)

	messages := []completionMessage{
		{Role: "system", Text: systemPrompt + "\n\nТы также аналитик. Отвечай строго JSON."},
		{Role: "user", Text: prompt},
	}
	text, err := c.complete(ctx, messages, schema, 0.4, 1200)
	if err != nil {
		return nil, err
	}
	return parseChatTurn(text)
}

func (c *Client) SuggestReply(ctx context.Context, systemPrompt string, history []domain.ChatMessage, draft string) (string, error) {
	var hist strings.Builder
	for _, m := range history {
		fmt.Fprintf(&hist, "%s: %s\n", m.Role, m.Text)
	}
	messages := []completionMessage{
		{Role: "system", Text: systemPrompt},
		{Role: "user", Text: fmt.Sprintf("Контекст:\n%s\nЧерновик пользователя: %s\nУлучши формулировку с учётом цели сессии. Верни только готовую реплику на русском.", hist.String(), draft)},
	}
	return c.complete(ctx, messages, nil, 0.6, 400)
}

func (c *Client) SummarizeSession(ctx context.Context, systemPrompt string, history []domain.ChatMessage) (*domain.SessionSummary, error) {
	var hist strings.Builder
	for _, m := range history {
		fmt.Fprintf(&hist, "%s: %s\n", m.Role, m.Text)
	}
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"scores": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"clarity":    scoreDetailSchema(),
					"confidence": scoreDetailSchema(),
					"respect":    scoreDetailSchema(),
					"balance":    scoreDetailSchema(),
				},
			},
			"strengths":       map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
			"weaknesses":      map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
			"riskFlags":       map[string]interface{}{"type": "array"},
			"improvedReplies": map[string]interface{}{"type": "array"},
			"tipForNext":      map[string]interface{}{"type": "string"},
			"hasRisk":         map[string]interface{}{"type": "boolean"},
		},
		"required": []string{"scores", "strengths", "weaknesses", "tipForNext", "hasRisk"},
	}
	messages := []completionMessage{
		{Role: "system", Text: systemPrompt},
		{Role: "user", Text: fmt.Sprintf("Проанализируй сессию (summarize_session):\n%s\nВерни JSON разбора.", hist.String())},
	}
	text, err := c.complete(ctx, messages, schema, 0.3, 2000)
	if err != nil {
		return nil, err
	}
	return parseSummary(text)
}

func (c *Client) GenerateDailyExercise(ctx context.Context, focusSkill string) (*domain.DailyExercise, error) {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"title":       map[string]interface{}{"type": "string"},
			"description": map[string]interface{}{"type": "string"},
			"criterion":   map[string]interface{}{"type": "string"},
		},
		"required": []string{"title", "description", "criterion"},
	}
	messages := []completionMessage{
		{Role: "user", Text: fmt.Sprintf("Сгенерируй упражнение дня для развития навыка «%s» в тренажёре романтической коммуникации. JSON на русском.", focusSkill)},
	}
	text, err := c.complete(ctx, messages, schema, 0.5, 500)
	if err != nil {
		return nil, err
	}
	var ex domain.DailyExercise
	if err := decodeJSON(text, &ex); err != nil {
		return nil, err
	}
	return &ex, nil
}

func scoreDetailSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"value":   map[string]interface{}{"type": "integer", "minimum": 0, "maximum": 10},
			"comment": map[string]interface{}{"type": "string"},
		},
	}
}

func parseChatTurn(text string) (*domain.ChatTurnResult, error) {
	var raw struct {
		PersonaReply       string `json:"personaReply"`
		Clarity            int    `json:"clarity"`
		Confidence         int    `json:"confidence"`
		Respect            int    `json:"respect"`
		Balance            int    `json:"balance"`
		StatusLabel        string `json:"statusLabel"`
		Status             string `json:"status"`
		ConsentDetected    bool   `json:"consentDetected"`
		ConsentSeverity    string `json:"consentSeverity"`
		ConsentExplanation string `json:"consentExplanation"`
		ConsentSuggestion  string `json:"consentSuggestion"`
	}
	if err := decodeJSON(text, &raw); err != nil {
		return nil, err
	}
	clamp := func(v int) int {
		if v < 0 {
			return 0
		}
		if v > 10 {
			return 10
		}
		return v
	}
	result := &domain.ChatTurnResult{
		PersonaText: raw.PersonaReply,
		Score: domain.MessageScore{
			Clarity: clamp(raw.Clarity), Confidence: clamp(raw.Confidence),
			Respect: clamp(raw.Respect), Balance: clamp(raw.Balance),
			StatusLabel: raw.StatusLabel, Status: raw.Status,
		},
		Consent: domain.ConsentRisk{Detected: raw.ConsentDetected},
	}
	if raw.ConsentDetected && raw.ConsentSeverity != "" && raw.ConsentSeverity != "none" {
		sev := raw.ConsentSeverity
		result.Consent.Severity = &sev
		if raw.ConsentExplanation != "" {
			e := raw.ConsentExplanation
			result.Consent.Explanation = &e
		}
		if raw.ConsentSuggestion != "" {
			s := raw.ConsentSuggestion
			result.Consent.Suggestion = &s
		}
	}
	return result, nil
}

func parseSummary(text string) (*domain.SessionSummary, error) {
	var s domain.SessionSummary
	if err := decodeJSON(text, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func decodeJSON(text string, v interface{}) error {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "```") {
		lines := strings.Split(text, "\n")
		if len(lines) >= 2 {
			text = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}
	return json.Unmarshal([]byte(text), v)
}

var _ domain.LLMService = (*Client)(nil)

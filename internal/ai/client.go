package ai

/*
import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client struct {
	ApiKey string
	Folder string
}

func (c *Client) Generate(prompt string) (string, error) {

	body := map[string]interface{}{
		"modelUri": "gpt://" + c.Folder + "/yandexgpt/latest",
		"completionOptions": map[string]interface{}{
			"temperature": 0.7,
			"maxTokens":   500,
		},
		"messages": []map[string]string{
			{
				"role": "system",
				"text": prompt,
			},
		},
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(
		"POST",
		"https://llm.api.cloud.yandex.net/foundationModels/v1/completion",
		bytes.NewBuffer(jsonBody),
	)

	req.Header.Set("Authorization", "Api-Key "+c.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	return "response", nil
}
*/

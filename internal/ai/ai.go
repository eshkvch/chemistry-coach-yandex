package ai

/*
import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func AskYandex(prompt string) (string, error) {

	apiKey := os.Getenv("YANDEX_API_KEY")
	folder := os.Getenv("YANDEX_FOLDER_ID")

	body := map[string]interface{}{
		"modelUri": "gpt://" + folder + "/yandexgpt/latest",
		"completionOptions": map[string]interface{}{
			"stream":      false,
			"temperature": 0.7,
			"maxTokens":   500,
		},
		"messages": []map[string]string{
			{
				"role": "user",
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

	req.Header.Set("Authorization", "Api-Key "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	return string(responseBody), nil
}
*/

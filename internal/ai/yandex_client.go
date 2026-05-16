package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type Client struct {
	ApiKey string
	Folder string
}

func NewClient() *Client {
	return &Client{
		ApiKey: os.Getenv("YANDEX_API_KEY"),
		Folder: os.Getenv("YANDEX_FOLDER_ID"),
	}
}

func (c *Client) Generate(prompt string) (string, error) {

	body := map[string]interface{}{
		"modelUri": "gpt://" + c.Folder + "/aliceai-llm/latest",
		"completionOptions": map[string]interface{}{
			"temperature": 0.8,
			"maxTokens":   500,
		},
		// "messages": []map[string]string{
		// 	{
		// 		"role": "user",
		// 		"text": prompt,
		// 	},
		// },
		"messages": []map[string]string{
			{
				"role": "user",
				"text": "Привет",
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

	b, _ := io.ReadAll(resp.Body)

	return string(b), nil
}

//

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	YANDEX_CLOUD_FOLDER := os.Getenv("YANDEX_CLOUD_FOLDER")
	YANDEX_CLOUD_API_KEY := os.Getenv("YANDEX_CLOUD_API_KEY")
	YANDEX_CLOUD_MODEL := os.Getenv("YANDEX_CLOUD_MODEL")

	reqData := ResponseRequest{
		Model:           fmt.Sprintf("gpt://%s/%s", YANDEX_CLOUD_FOLDER, YANDEX_CLOUD_MODEL),
		Temperature:     0.3,
		Instructions:    "",
		Input:           "какая погода в москве?",
		MaxOutputTokens: 500,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", "https://ai.api.cloud.yandex.net/v1/responses", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Api-Key "+YANDEX_CLOUD_API_KEY)
	req.Header.Set("OpenAI-Project", YANDEX_CLOUD_FOLDER)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %s\n", string(body))

	if resp.StatusCode == 200 {
		var response ResponseData
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Printf("Error parsing response: %v", err)
		} else {
			fmt.Println(response.GetOutputText())
		}
	}
}

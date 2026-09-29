package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIRequest struct {
	Model    string      `json:"model"`
	Messages []AIMessage `json:"messages"`
}

type AIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GetAiResponse(message string) (string, error) {

	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		return "", fmt.Errorf("API_KEY is not set")
	}

	reqBody := AIRequest{
		Model: "openai/gpt-oss-120b",

		Messages: []AIMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant inside the GoTalk chat app. Keep replies clear and concise.",
			},
			{
				Role:    "user",
				Content: message,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// Important: get actual Groq error
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return "", fmt.Errorf(
			"Groq API error: %s - %s",
			resp.Status,
			string(body),
		)
	}

	var aiResp AIResponse

	err = json.NewDecoder(resp.Body).Decode(&aiResp)
	if err != nil {
		return "", err
	}

	if len(aiResp.Choices) == 0 {
		return "", fmt.Errorf("empty AI response")
	}

	return aiResp.Choices[0].Message.Content, nil
}

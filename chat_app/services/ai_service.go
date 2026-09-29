package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)



type AIRequest struct{
	Model string `json:"model"`
	Input string `json:"input"`
}

type AIResponse struct {
    Choices []struct {
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
}

func GetAiResponse(message string) (string,error){
	apiKey:=os.Getenv("API_KEY")
	if apiKey==""{
		return "",fmt.Errorf("API Key is not Set")
	}
	reqBody := AIRequest{
		Model: "llama-3.3-70b-versatile",
		Input: message,
	}

	JsonData,err:=json.Marshal(reqBody)
	if err!=nil{
		return "",err
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(JsonData),
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

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI API returned status: %s", resp.Status)
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
package parser

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "regexp"
    "strings"

    "resumeparser/models"
)

const nvidiaURL = "https://integrate.api.nvidia.com/v1/chat/completions"
const model = "meta/llama-3.1-8b-instruct"  

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type RequestBody struct {
    Model       string    `json:"model"`
    Messages    []Message `json:"messages"`
    MaxTokens   int       `json:"max_tokens"`
    Temperature float64   `json:"temperature"`
    Stream      bool      `json:"stream"`
}

type Choice struct {
    Message Message `json:"message"`
}

type APIResponse struct {
    Choices []Choice `json:"choices"`
}

var fenceRegex = regexp.MustCompile("(?s)```[a-z]*\\n?(.*?)```")

func ParseResume(text string) (*models.Resume, error) {
    apiKey := os.Getenv("NVIDIA_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("NVIDIA_API_KEY not set")
    }

prompt := fmt.Sprintf(`
Parse this resume and return ONLY valid JSON.

Schema:
{
 "name": "string",
 "email": "string",
 "phone": "string",
 "location": "string",
 "summary": "string",
 "skills": ["string"],
 "experience": [
   {
     "company": "string",
     "title": "string",
     "dates": "string",
     "highlights": ["string"]
   }
 ],
 "education": [
   {
     "institution": "string",
     "degree": "string",
     "year": "string"
   }
 ],
 "certifications": ["string"]
}

Return empty array [] if none.

Resume:
%s
`, text)

    body := RequestBody{
        Model: model,
        Messages: []Message{
            {
                Role:    "system",
                Content: "You are a resume parser. Return ONLY valid JSON, no markdown, no explanation.",
            },
            {
                Role:    "user",
                Content: prompt,
            },
        },
        MaxTokens:   2000,
        Temperature: 0.1,  // low temp = more consistent JSON output
        Stream:      false,
    }

    bodyBytes, _ := json.Marshal(body)

    req, err := http.NewRequest("POST", nvidiaURL, bytes.NewReader(bodyBytes))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiKey)  // NVIDIA uses Bearer token

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        b, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(b))
    }

    var apiResp APIResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    if len(apiResp.Choices) == 0 {
        return nil, fmt.Errorf("empty response from API")
    }

    raw := strings.TrimSpace(apiResp.Choices[0].Message.Content)

fmt.Println("------------ AI RAW RESPONSE ------------")
fmt.Println(raw)
fmt.Println("-----------------------------------------")
    // strip markdown fences if present
    if match := fenceRegex.FindStringSubmatch(raw); len(match) > 1 {
        raw = strings.TrimSpace(match[1])
    }

    var resume models.Resume
    if err := json.Unmarshal([]byte(raw), &resume); err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %w", err)
    }

    return &resume, nil
}
func CheckConnection() error {
    apiKey := os.Getenv("NVIDIA_API_KEY")
    if apiKey == "" {
        return fmt.Errorf("NVIDIA_API_KEY not set")
    }

    // send a tiny test request to verify the key works
    body := RequestBody{
        Model: model,
        Messages: []Message{
            {Role: "system", Content: "You are a helpful assistant."},
            {Role: "user", Content: "Say OK"},
        },
        MaxTokens:   5,
        Temperature: 0.1,
        Stream:      false,
    }

    bodyBytes, _ := json.Marshal(body)
    req, err := http.NewRequest("POST", nvidiaURL, bytes.NewReader(bodyBytes))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiKey)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return fmt.Errorf("connection failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        b, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("API error %d: %s", resp.StatusCode, string(b))
    }

    return nil
}
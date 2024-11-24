package api

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "llm/config"
    "net/http"
)

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type RequestPayload struct {
    Messages []Message `json:"messages"`
    Model    string    `json:"model"`
}

type ResponsePayload struct {
    Choices []struct {
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
}

func MakeGroqRequest(prompt string) (string, error) {
    payload := RequestPayload{
        Messages: []Message{
            {Role: "user", Content: prompt},
        },
        Model: "llama3-8b-8192", // Model specified in your curl example
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return "", fmt.Errorf("failed to marshal JSON: %v", err)
    }

    req, err := http.NewRequest("POST", config.ApiUrl, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("failed to create HTTP request: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+config.ApiKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("failed to make API request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("API returned status code %d", resp.StatusCode)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response body: %v", err)
    }

    var responsePayload ResponsePayload
    if err := json.Unmarshal(body, &responsePayload); err != nil {
        return "", fmt.Errorf("failed to unmarshal JSON response: %v", err)
    }

    if len(responsePayload.Choices) > 0 {
        return responsePayload.Choices[0].Message.Content, nil
    }

    return "", fmt.Errorf("no response from API")
}
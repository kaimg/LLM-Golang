package db

import (
    "fmt"
)

func SavePrompt(userID interface{}, prompt string, response string) error {
    _, err := DB.Exec(
        "INSERT INTO prompts (user_id, prompt, response) VALUES ($1, $2, $3)",
        userID, prompt, response,
    )
    if err != nil {
        return fmt.Errorf("failed to save prompt: %v", err)
    }
    return nil
}

func GetPromptsByUser(userID int) ([]string, error) {
    rows, err := DB.Query(
        "SELECT prompt FROM prompts WHERE user_id = $1 ORDER BY created_at DESC",
        userID,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to fetch prompts: %v", err)
    }
    defer rows.Close()

    var prompts []string
    for rows.Next() {
        var prompt string
        if err := rows.Scan(&prompt); err != nil {
            return nil, fmt.Errorf("failed to scan prompt: %v", err)
        }
        prompts = append(prompts, prompt)
    }

    return prompts, nil
}

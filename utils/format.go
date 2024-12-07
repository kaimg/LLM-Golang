package utils

import (
	"regexp"
	"strings"
	"time"
)

func ParseDateAndTime(input string) (string) {
    t, err := time.Parse(time.RFC3339, input)
    if err != nil {
        return "Undefined"
    }
    return t.Format("15:04 01/02/2006")
}

// formatMarkdown converts a Markdown-style response to basic HTML.
func FormatMarkdown(text string) string {
    // Convert headers (e.g., # Header) to <h1>, <h2>, etc.
    text = regexp.MustCompile(`(?m)^### (.+)$`).ReplaceAllString(text, "<h3>$1</h3>")
    text = regexp.MustCompile(`(?m)^## (.+)$`).ReplaceAllString(text, "<h2>$1</h2>")
    text = regexp.MustCompile(`(?m)^# (.+)$`).ReplaceAllString(text, "<h1>$1</h1>")

    // Convert bold (**text**) to <b>text</b>
    text = regexp.MustCompile(`\*\*(.+?)\*\*`).ReplaceAllString(text, "<b>$1</b>")

    // Convert italic (*text*) to <i>text</i>
    text = regexp.MustCompile(`\*(.+?)\*`).ReplaceAllString(text, "<i>$1</i>")

	// Convert code (```text```) to <code>text<code>
	text = regexp.MustCompile("(?s)```(.*?)```").ReplaceAllString(text, "<code>$1</code>")
	
    // Convert newline characters to <br> for line breaks
    text = strings.ReplaceAll(text, "\n", "<br>")

    return text
}
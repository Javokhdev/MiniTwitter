package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	genai "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// AIResponse structure
type AIResponse struct {
	Candidates []struct {
		Content string `json:"content"`
	} `json:"candidates"`
}

// Helper function to retrieve the Gemini API key
func getGeminiAPIKey() (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", errors.New("API key is not set in environment variables")
	}
	return apiKey, nil
}

// sanitizeContent removes special characters, emojis, and non-ASCII characters
func sanitizeContent(content string) string {
	var sanitizedContent []rune

	for _, r := range content {
		if r <= unicode.MaxASCII && unicode.IsPrint(r) && !unicode.IsSymbol(r) {
			sanitizedContent = append(sanitizedContent, r)
		}
	}

	return strings.TrimSpace(string(sanitizedContent))
}

func GetTagsFromAI(ctx context.Context, content string) ([]string, error) {
	// Sanitize the content to ensure it only contains valid characters
	sanitizedContent := sanitizeContent(content)

	// Retrieve API key
	apiKey, err := getGeminiAPIKey()
	if err != nil {
		return nil, err
	}

	// Create a new GenAI client
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating GenAI client: %w", err)
	}
	defer client.Close()

	// Configure the model with safety settings
	model := client.GenerativeModel("gemini-1.5-flash")
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockOnlyHigh,
		},
	}

	// Construct the prompt
	prompt := fmt.Sprintf("Extract relevant tags and hashtags for the following content: \"%s\"", sanitizedContent)

	// Generate the response
	timeoutCtx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()
	resp, err := model.GenerateContent(timeoutCtx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("error generating AI response: %w", err)
	}

	// Check for empty candidates or parts
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return []string{"default_tag"}, nil
	}

	// Extract the first part of the content
	contentPart := resp.Candidates[0].Content.Parts[0]
	temp, err := json.Marshal(contentPart)
	if err != nil {
		return nil, fmt.Errorf("error marshaling response content: %w", err)
	}

	// Convert the JSON to a plain string
	unescapedData, err := strconv.Unquote(string(temp))
	if err != nil {
		return nil, fmt.Errorf("error unquoting response content: %w", err)
	}

	// Parse the tags from the unescaped string
	tags := strings.Split(unescapedData, ",")
	for i := range tags {
		tags[i] = strings.TrimSpace(tags[i])
	}

	// Ensure at least one valid tag is present
	if len(tags) == 0 {
		tags = []string{"default_tag"}
	}

	return tags, nil
}

package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "os"
	"strings"

	// "github.com/joho/godotenv"
)


type Candidate struct {
	Content struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"content"`
}

type Response struct {
	Candidates []Candidate `json:"candidates"`
}

type RequestBody struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

func AskFromGemeni(content string, defaultTags []string) ([]string, error) {
	// .env faylni yuklash
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

	// API kalitni .env fayldan olish
	// GEMINI_API_KEY := os.Getenv("GEMINI_API_KEY")
	// if GEMINI_API_KEY == "" {
	// 	log.Fatal("GEMINI_API_KEY is not set in .env file")
	// }

	GEMINI_API_KEY := ""

	// API URL
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", GEMINI_API_KEY)

	// Prompt yaratish
	promt := fmt.Sprintf("Given the content: '%s', which of the following tags are most relevant? %v. Just write only tags", content, defaultTags)

	// So‘rov body tayyorlash
	requestBody := RequestBody{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: promt},
				},
			},
		},
	}

	// JSONga o‘girish
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}

	// POST so‘rov yuborish
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Javobni tahlil qilish
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	// Teglarni ajratish
	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		aiResponse := response.Candidates[0].Content.Parts[0].Text
		matchedTags := []string{}
		for _, tag := range defaultTags {
			if strings.Contains(aiResponse, tag[1:]) {
				matchedTags = append(matchedTags, tag)
			}
		}
		return matchedTags, nil
	}

	return []string{}, fmt.Errorf("no valid response from AI")
}

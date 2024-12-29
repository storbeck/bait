package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Read API keys
var (
	openaiAPIKey     = os.Getenv("OPENAI_API_KEY")
	elevenlabsAPIKey = os.Getenv("ELEVENLABS_API_KEY")
	callbackNumber   = "1-800-555-0123"
)

// Helper function for HTTP POST requests
func postRequest(url string, body []byte, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

type APIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Helper function to generate OpenAI content
func generateOpenAIContent(prompt string) (string, error) {
	req := map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]interface{}{
			{"role": "user", "content": prompt},
		},
		"max_tokens": 1500,
	}

	reqBody, _ := json.Marshal(req)
	headers := map[string]string{
		"Authorization": "Bearer " + openaiAPIKey,
		"Content-Type":  "application/json",
	}

	response, err := postRequest("https://api.openai.com/v1/chat/completions", reqBody, headers)
	if err != nil {
		return "", err
	}

	var result APIResponse
	if err := json.Unmarshal(response, &result); err != nil {
		return "", err
	}

	if result.Error.Message != "" {
		return "", fmt.Errorf("OpenAI API error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	return result.Choices[0].Message.Content, nil
}

// Helper function to generate voicemail
func generateVoicemail(script string) error {
	voiceID := "nPczCjzI2devNBz1zQrb" // Brian voice
	req := map[string]interface{}{
		"text":     script,
		"model_id": "eleven_multilingual_v2",
		"voice_settings": map[string]interface{}{
			"stability":        0.8,
			"similarity_boost": 0.65,
		},
	}

	reqBody, _ := json.Marshal(req)
	headers := map[string]string{
		"xi-api-key":   elevenlabsAPIKey,
		"Content-Type": "application/json",
	}

	response, err := postRequest(fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voiceID), reqBody, headers)
	if err != nil {
		return err
	}

	if len(response) < 1024 {
		return fmt.Errorf("ElevenLabs API error: %s", string(response))
	}

	return ioutil.WriteFile("voicemail.mp3", response, 0644)
}

func main() {
	if openaiAPIKey == "" || elevenlabsAPIKey == "" {
		fmt.Println("Error: Missing API keys.")
		return
	}

	// Generate company details
	fmt.Println("Generating IT support department details...")
	companyDetails, err := generateOpenAIContent(`Generate a fictional but legitimate-sounding IT support department with the following details:
	- Company name (use a well-known tech or financial company)
	- Security operations center location
	- Department phone number (fictional but properly formatted)
	- IT Support ticket reference number (fictional but properly formatted)
	
	Format as plain text, one detail per line.`)
	if err != nil {
		fmt.Println("Error generating company details:", err)
		return
	}

	fmt.Println("Generated Details:", companyDetails)
	if !confirmContinue() {
		main()
		return
	}

	// Generate script
	for {
		fmt.Println("Generating IT security alert script...")
		script, err := generateOpenAIContent(fmt.Sprintf(`Create an urgent IT security voicemail script using this company: %s

		The script should:
		- Be 30-45 seconds long
		- Sound professional and authoritative
		- Mention detecting suspicious login activity "just minutes ago"
		- Express urgency about potential account compromise
		- Use this exact callback number: %s
		- Request immediate callback to verify identity and prevent account lockout
		- Include a reference/ticket number
		- End with clear urgency to call back immediately

		Format as a simple reading script without any dialog markers, brackets, or annotations.
		Keep the tone professional but convey urgency.`, companyDetails, callbackNumber))
		if err != nil {
			fmt.Println("Error generating script:", err)
			return
		}

		fmt.Println("\nGenerated Script:")
		fmt.Println(script)

		if confirmContinue() {
			// Generate voicemail
			fmt.Println("Generating voicemail...")
			if err := generateVoicemail(script); err != nil {
				fmt.Println("Error generating voicemail:", err)
				return
			}
			fmt.Println("Voicemail saved as 'voicemail.mp3'.")
			break
		}
	}
}

func confirmContinue() bool {
	fmt.Println("\nWould you like to proceed? (yes/no)")
	var response string
	fmt.Scanln(&response)
	return response == "yes"
}

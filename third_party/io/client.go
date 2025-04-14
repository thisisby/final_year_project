package io

import (
	"backend/internal/http/data_transfers"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const ionetEndpoint = "https://api.intelligence.io.solutions/api/v1/chat/completions"

type ChatCompletionRequest struct {
	Model               string    `json:"model"`
	Messages            []Message `json:"messages"`
	Temperature         float64   `json:"temperature"`
	Stream              bool      `json:"stream"`
	MaxCompletionTokens int       `json:"max_completion_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Client struct {
	ApiKey string
}

type WorkoutPlan struct {
	WorkoutName string     `json:"workout_name"`
	Description string     `json:"description"`
	Exercises   []Exercise `json:"exercises"`
}

type Exercise struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		ApiKey: apiKey,
	}
}

func (c *Client) GenerateWorkout(generateRequest data_transfers.WorkoutGenerateRequest) (WorkoutPlan, error) {

	content := fmt.Sprintf(`
{
	"role": "system",
	"content": "You are a certified professional gym trainer. Your task is to generate a personalized and detailed workout plan based on the user's preferences.

Rules to follow strictly:
1. Respond with a single **valid JSON object only**. No markdown, no text outside the JSON.
2. The JSON must include:
   - **workout_name** (based on the user's goal and body areas, max 30 characters),
   - **description**: A detailed plan overview including purpose, progression logic, rest days, and how often to train per week. Mention when to increase weights or take deloads.
   - **exercises**: A list of 4–5 exercises, each with:
     - name: string
     - description: a step-by-step guide (including sets, reps, tempo if needed, rest time, key form tips)
     - link: a YouTube tutorial (must be relevant and accurate)

3. The workout should be tailored precisely for:
   - Level: %s
   - Goal: %s
   - Body areas: %s
   - Gender: %s
   - Age: %s
   - Additional details: %s

Output example format (do not reuse the values directly, just structure similarly):

{
  "workout_name": "Upper Body Hypertrophy",
  "description": "This upper body workout targets chest, back, shoulders, and arms. Designed for intermediate lifters focused on muscle growth. Follow a 4-day split with alternating push/pull days. Increase weights by ~5%% weekly if all sets are completed with proper form. Every 4th week is a deload (reduce weights by 15%% and sets by 1). Rest 60–90 seconds between hypertrophy sets and 2–3 mins on compound lifts. Ensure proper warm-up and cool-down.",
  "exercises": [
    {
      "name": "Incline Barbell Press",
      "description": "4 sets of 8–10 reps at moderate weight.\nRest 90 seconds.\nKeep elbows at ~45 degrees. Focus on slow controlled negatives (3s down).",
      "link": "https://www.youtube.com/watch?v=SrqOu55lrYU"
    },
    ...
  ]
}
"
}`,
		generateRequest.Level,
		generateRequest.Goal,
		generateRequest.BodyAreas,
		generateRequest.Gender,
		generateRequest.Age,
		generateRequest.Details,
	)

	requestBody := ChatCompletionRequest{
		Model: "meta-llama/Llama-3.3-70B-Instruct",
		Messages: []Message{
			{
				Role:    "system",
				Content: content,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Error marshaling request: %v\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", ionetEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API returned error: %s\n", string(body))
		os.Exit(1)
	}

	var response ChatCompletionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		os.Exit(1)
	}

	// Check if we have choices
	if len(response.Choices) == 0 {
		fmt.Println("No completion choices returned")
		os.Exit(1)
	}

	var plan WorkoutPlan
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &plan)
	if err != nil {
		fmt.Printf("Error unmarshaling workout plan: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Workout Plan: %+v\n", plan)

	return plan, nil
}

package structures

// GeminiAIRequest structure for Gemini AI API
type GeminiAIRequest struct {
	Contents []Content `json:"contents"`
}

// Content structure for Gemini AI API
type Content struct {
	Parts []Part `json:"parts"`
}

// Part structure for Gemini AI API
type Part struct {
	Text string `json:"text"`
}

// GeminiAIResponse structure for Gemini AI API response
type GeminiAIResponse struct {
	Candidates     []Candidate    `json:"candidates"`
	PromptFeedback PromptFeedback `json:"promptFeedback"`
}

type Candidate struct {
	Content       Content        `json:"content"`
	FinishReason  string         `json:"finishReason"`
	Index         int            `json:"index"`
	SafetyRatings []SafetyRating `json:"safetyRatings"`
}

type SafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

type PromptFeedback struct {
	SafetyRatings []SafetyRating `json:"safetyRatings"`
}

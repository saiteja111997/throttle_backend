package structures

type UserAction struct {
	TextContent string `json:"text_content"`
	Type        string `json:"type"`
	Id          int    `json:"id"`
}

type RawErrorResponse struct {
	ErrorInfo   Errors       `json:"error_info"`
	UserActions []UserAction `json:"user_actions"`
}
type Title struct {
	Title string `json:"title"`
}

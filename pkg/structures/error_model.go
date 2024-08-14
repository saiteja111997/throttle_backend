package structures

type UserAction struct {
	TextContent string `json:"text_content"`
	Type        string `json:"type"`
	Id          int    `json:"id"`
}

type GetUnresolvedJourneys struct {
	Title string `json:"title"`
	Id    string `json:"id"`
}

type RawErrorResponse struct {
	ErrorInfo   Errors       `json:"error_info"`
	UserActions []UserAction `json:"user_actions"`
}

type Title struct {
	Title string `json:"title"`
}

type DashboardData struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	DocFilePath string `json:"doc_file_path"`
}

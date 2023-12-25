package structures

type FulFillmentBody struct {
	Tag string `json:"tag"`
}

type IntentParameterValue struct {
	OriginalValue string `json:"originalValue"`
	ResolvedValue string `json:"resolvedValue"`
}

type IntentBody struct {
	LastMatchedIntent string                          `json:"lastMatchedIntent"`
	DisplayName       string                          `json:"displayName"`
	Confidence        float64                         `json:"confidence"`
	Parameters        map[string]IntentParameterValue `json:"parameters"`
}

type SessionBody struct {
	Session    string                 `json:"session"`
	Parameters map[string]interface{} `json:"parameters"`
}

// webhookRequest is used to unmarshal a WebhookRequest JSON object. Note that
// not all members need to be defined--just those that you need to process.
// As an alternative, you could use the types provided by
// the Dialogflow protocol buffers:
// https://godoc.org/google.golang.org/genproto/googleapis/cloud/dialogflow/v2#WebhookRequest
type WebhookRequest struct {
	DetectIntentResponseId string          `json:"detectIntentResponseId"`
	Language               string          `json:"languageCode"`
	FulFillmentInfo        FulFillmentBody `json:"fulfillmentInfo"`
	IntentInfo             IntentBody      `json:"intentInfo"`
	SessionInfo            SessionBody     `json:"sessionInfo"`
	Text                   string          `json:"text"`
}

// webhookResponse is used to marshal a WebhookResponse JSON object. Note that
// not all members need to be defined--just those that you need to process.
// As an alternative, you could use the types provided by
// the Dialogflow protocol buffers:
// https://godoc.org/google.golang.org/genproto/googleapis/cloud/dialogflow/v2#WebhookResponse

type TextObj struct {
	Text []string `json:"text"`
}

type MessageObject struct {
	Text TextObj `json:"text"`
}

type Message struct {
	Messages []MessageObject `json:"messages"`
}

type WebhookResponse struct {
	FulfillmentResponse Message `json:"fulfillmentResponse"`
}

// type SummaryUpdateMessageBody struct {
// 	Message string `json:"message"`
// }

package models

type WebhookPayload struct {
	Action		string		`json:"action"`
	Ref 		string		`json:"ref"`		// Not empty Ref field is a signal for repository push event
	Issue		Issue		`json:"issue"`
	Label		Label		`json:"label"`
	Repository	Repos		`json:"repository"`
}

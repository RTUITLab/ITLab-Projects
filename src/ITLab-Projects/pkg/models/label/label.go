package label

type Label struct {
	ID          int    `json:"id,omitempty"`
	NodeID      string `json:"node_id,omitempty"`
	URL         string `json:"url,omitempty"`
	Name        string `json:"name"`
	Type		string	`json:"type"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}
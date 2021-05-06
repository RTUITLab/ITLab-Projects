package label

type Label struct {
	ID          int    	`json:"id,omitempty"`
	NodeID      string 	`json:"node_id,omitempty"`
	URL         string 	`json:"url,omitempty"`
	Type		string	`json:"type"`
	Description string 	`json:"description,omitempty"`
	CompactLabel		`bson:",inline"`
}

type CompactLabel struct {
	Name        string 	`json:"name"`
	Color       string 	`json:"color,omitempty"`
}
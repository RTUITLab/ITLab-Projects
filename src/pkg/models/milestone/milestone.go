package milestone

import (
	. "github.com/ITLab-Projects/pkg/models/user"
)

type Milestone struct {
	ID        			uint64     	`json:"id"`
	Number    			uint64     	`json:"number"`
	State     			string     	`json:"state"`
	Title     			string     	`json:"title"`
	Description 		string 		`json:"description"`
	Creator      		User    	`json:"creator"`
}
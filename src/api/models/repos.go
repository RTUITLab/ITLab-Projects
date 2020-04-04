package models

type Repos struct {
	ID           uint64			    `json:"id"`
	Name		 string             `json:"name"`
	HTMLUrl      string             `json:"html_url"`
	Description  string             `json:"description"`
	CreatedAt	 string				`json:"created_at"`
	UpdatedAt	 string				`json:"updated_at"`
	Language     string				`json:"language"`
	Archived	 bool				`json:"archived"`
}

type Issue struct {
	ID           uint64 			`json:"id"`
	Number		 uint64				`json:"number"`
	Title		 string				`json:"title"`
	User		 User				`json:"user"`
	State        string				`json:"state"`
	Assignees    []Assignee		    `json:"assignees"`
	Milestone	 Milestone	        `json:"milestone"`
	CreatedAt	 string				`json:"created_at"`
	UpdatedAt	 string				`json:"updated_at"`
	ClosedAt	 string				`json:"closed_at"`
}

type User struct {
	ID           uint64			    `json:"id"`
	Login		 string				`json:"login"`
	AvatarURL    string 			`json:"avatar_url"`
	URL		 	 string				`json:"html_url"`
}

type Assignee struct {
	ID           uint64			    `json:"id"`
	Login		 string				`json:"login"`
	AvatarURL    string 			`json:"avatar_url"`
	URL	 	 	 string				`json:"html_url"`
}

type Milestone struct {

}
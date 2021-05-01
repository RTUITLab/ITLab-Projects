package user

type User struct {
	ID           		uint64		`json:"id"`
	Login		 		string		`json:"login"`
	AvatarURL    		string 		`json:"avatar_url"`
	URL		 	 		string		`json:"html_url"`
}

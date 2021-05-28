package err

type Message struct {
	Message string `json:"message"`
}

type Err struct {
	Err string `json:"error"`
	Message
}
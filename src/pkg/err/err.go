package err

type Err struct {
	Err error `json:"error"`
	Message string `json:"string"`
}
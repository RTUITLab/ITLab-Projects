package responce

import (
	"errors"
	"net/http"

	e "github.com/ITLab-Projects/pkg/err"
	"github.com/ITLab-Projects/pkg/statuscode"
	"gopkg.in/square/go-jose.v2/json"
)

var (
	ErrMessageIsNil 	= errors.New("Message is nil")
)

type Messager interface {
	Message() *e.Message
}

type Statuser interface {
	StatusCode() int
}

type BodyEncoder interface {
	Encode(w http.ResponseWriter) error
}

type HTTPResponce interface {
	Statuser
	BodyEncoder
}

type HTTPErrResponce interface {
	Responcer
	WriteHeader(w http.ResponseWriter)
	WriteMessage(w http.ResponseWriter) error
}

type Responce struct {
	Status *statuscode.StatusCode
}

type Responcer interface {
	Messager
	Statuser
}

func (r *Responce) Message() *e.Message {
	if r.Status.Err != nil {
		return &e.Message{
			Message: r.Status.Err.Error(),
		}
	}

	return nil
}

func (r *Responce) StatusCode() int {
	return r.Status.Status
}

func (r *Responce) WriteHeader(w http.ResponseWriter) {
	w.WriteHeader(r.StatusCode())
}

func (r *Responce) WriteMessage(w http.ResponseWriter) error {
	w.Header().Add("Content-Type", "application/json")
	if r.Message() == nil {
		return ErrMessageIsNil
	}

	return json.NewEncoder(w).Encode(r.Message())
}

// Return err from responce
// if err is nil return responce with nil err
// and 200 status code
func FromErr(err error) (*Responce) {
	status, ok := err.(*statuscode.StatusCode)
	if !ok {
		return &Responce{
			Status: &statuscode.StatusCode{
				Err: nil,
				Status: http.StatusOK,
			},
		}
	}

	return &Responce{
		Status: status,
	}
}


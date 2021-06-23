package responce

import (
	log "github.com/sirupsen/logrus"
	"context"
	"errors"
	"fmt"
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

type Header interface {
	Headers(ctx context.Context, w http.ResponseWriter)
}

type HTTPResponce interface {
	Statuser
	BodyEncoder
	Header
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
// 
// if err is not nil but not status err return 500 with msg internal status error
func FromErr(err error) (*Responce) {
	if err == nil {
		return &Responce{
			Status: &statuscode.StatusCode{
				Err: nil,
				Status: http.StatusOK,
			},
		}
	}
	status, ok := err.(*statuscode.StatusCode)
	if !ok {
		log.WithFields(
			log.Fields{
				"pkg": "responce",
				"err": err,
			},
		).Error("Unhandled error, err not implemetns responce")
		return &Responce{
			Status: &statuscode.StatusCode{
				Err: fmt.Errorf("Unexpected err"),
				Status: http.StatusInternalServerError,
			},
		}
	}

	return &Responce{
		Status: status,
	}
}


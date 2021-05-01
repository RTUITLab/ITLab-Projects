package clientwrapper

import (
	"net/http"
)

type WrapReq func(*http.Request)

func (wr WrapReq) Wrap(w WrapReq) WrapReq {
	return func(r *http.Request) {
		wr(r)
		w(r)
	}
}
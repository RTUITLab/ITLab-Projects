package clientwrapper
// ClientWraper is package that allow to wrap default client Do()
// You can wrap a changes with req

import (
	"net/http"
)

// ClientWithWrap allow to wrap your request to change them
type ClientWithWrap struct {
	client 		*http.Client
	wrapReq 	WrapReq
}

// Wrap wrap a req ot a new func that change them
func (cw *ClientWithWrap) Wrap(wrapFunc WrapReq) {
	if cw.wrapReq == nil {
		cw.wrapReq = func(r *http.Request) {
			wrapFunc(r)
		}
	} else {
		cw.wrapReq = cw.wrapReq.Wrap(wrapFunc)
	}
}

// Do is simmalr like the *http.Client.Do(*http.Request) but with wrappers
func (cw *ClientWithWrap) Do(req *http.Request) (*http.Response, error) {
	if cw.wrapReq != nil {
		cw.wrapReq(req)
	}

	return cw.client.Do(req)
}

// Return new ClientWithWrap
func New(client *http.Client) *ClientWithWrap {
	return &ClientWithWrap{
		client: client,
	}
}
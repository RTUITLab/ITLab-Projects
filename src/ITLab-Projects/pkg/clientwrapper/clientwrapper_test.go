package clientwrapper_test

import (
	"net/http"
	"testing"

	"github.com/ITLab-Projects/pkg/clientwrapper"
)

func TestFunc_Wrap(t *testing.T) {
	cw := clientwrapper.New(&http.Client{})
	
	cw.Wrap(func(r *http.Request) {
		r.Header.Add("Some_Head_1", "1")
		t.Log("Wrap header 1")
	})

	cw.Wrap(func(r *http.Request) {
		r.Header.Add("Some_Head_2", "2")
		t.Log("Wrap header 2")
	})

	cw.Wrap(func(r *http.Request) {
		r.Header.Add("Some_Head_3", "3")
		t.Log("Wrap header 3")
	})

	req, err := http.NewRequest("GET", "google.com", nil)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	cw.Do(req)

	if req.Header.Get("Some_Head_1") != "1" {
		t.Log("Assert error")
		t.FailNow()
	}

	if req.Header.Get("Some_Head_2") != "2" {
		t.Log("Assert error")
		t.FailNow()
	}

	if req.Header.Get("Some_Head_3") != "3" {
		t.Log("Assert error")
		t.FailNow()
	}
}

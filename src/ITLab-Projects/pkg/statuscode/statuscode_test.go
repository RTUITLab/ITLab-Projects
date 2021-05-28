package statuscode_test

import (
	"errors"
	
	"net/http"
	"testing"

	"github.com/ITLab-Projects/pkg/statuscode"
)

func TestFunc_GetStatus(t *testing.T) {
	someErr := errors.New("Some error")
	wrapped := statuscode.WrapStatusError(
		someErr,
		http.StatusNotFound,
	)
	t.Log(wrapped)

	if status, ok := statuscode.GetStatus(wrapped); ok {
		t.Log(status)
		if status != http.StatusNotFound {
			t.Log("Assert error")
			t.FailNow()
		}
	} else {
		t.Log("Assert error")
		t.FailNow()
	}

	err := statuscode.GetError(
		wrapped,
	)

	if err != someErr {
		t.Log("Assert error")
		t.FailNow()
	}
}
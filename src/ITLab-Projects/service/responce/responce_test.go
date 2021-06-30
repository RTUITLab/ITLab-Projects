package responce_test

import (
	"fmt"
	"testing"

	"github.com/ITLab-Projects/pkg/statuscode"
	"github.com/ITLab-Projects/service/responce"
)

type SomeResponce struct {
	responce.Responce
}

func TestFunc_FromErr(t *testing.T) {
	err := statuscode.WrapStatusError(
		fmt.Errorf("Some_err"),
		101,
	)

	resp := responce.FromErr(err)
	if resp == nil {
		panic("Assert error")
	}

	if resp.StatusCode() != 101 {
		panic("Assert error")
	}

	if resp.Message().Message != "Some_err" {
		panic("Assert error")
	}
}

func TestFunc_FromErr_Nil(t *testing.T) {
	resp := responce.FromErr(nil)

	if resp.StatusCode() != 200 {
		t.Log("Assert error")
		t.FailNow()
	}

	if resp.Message() != nil {
		t.Log("Assert error")
		t.FailNow()
	}
}

func TestFunc_FromErr_OtherError(t *testing.T) {
	resp := responce.FromErr(fmt.Errorf("some_err"))

	if resp.StatusCode() != 500 {
		t.Log("Assert error")
		t.Log(resp.StatusCode())
		t.FailNow()
	}

	if resp.Message() == nil {
		t.Log("Assert error")
		t.FailNow()
	}
}

func TestFunc_Cast(t *testing.T) {
	resp := &SomeResponce{
		responce.Responce{
			Status: &statuscode.StatusCode{
				Err:    nil,
				Status: 200,
			},
		},
	}

	f := func(r responce.Responcer) {
		t.Log(r.StatusCode())
		if r.Message() != nil {
			t.FailNow()
		}
	}

	f(resp)
}

func TestFunc_NilErrButOtherStatusCode(t *testing.T) {
	err := statuscode.WrapStatusError(
		nil,
		201,
	)

	resp := responce.FromErr(err)

	if resp.StatusCode() != 201 {
		t.Log("Asser error")
		t.FailNow()
	}

	if resp.Message() != nil {
		t.Log("Assert error")
		t.FailNow()
	}
}

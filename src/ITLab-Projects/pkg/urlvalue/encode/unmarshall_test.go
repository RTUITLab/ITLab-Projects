package encode_test

import (
	"fmt"
	"testing"

	"github.com/ITLab-Projects/pkg/urlvalue/encode"
)

func TestFunc_Unmarshall(t *testing.T) {
	type B struct {
		Number int `query:"number,int"`
	}
	type A struct {
		Start *int   `query:"start,int" json:",inline"`
		Count int    `query:"count,int"`
		Name  string `query:"name,string"`
		B
	}

	a := &A{}

	if err := encode.UrlQueryUnmarshall(
		a,
		map[string][]string{
			"start":  {"10"},
			"count":  {"12"},
			"name":   {"orbit"},
			"number": {"27"},
		},
	); err != nil {
		t.Log(err)
	}

	fmt.Printf("%+v", a)
}
